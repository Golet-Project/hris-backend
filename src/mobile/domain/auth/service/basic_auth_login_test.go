package service_test

import (
	"context"
	"errors"
	"fmt"
	"hroost/mobile/domain/auth/model"
	"hroost/mobile/domain/auth/service"
	authService "hroost/mobile/domain/auth/service"
	"hroost/shared/primitive"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockBasicAuthLoginDb struct {
	mock.Mock
}

func (m *MockBasicAuthLoginDb) GetLoginCredential(ctx context.Context, email string) (credential model.GetLoginCredentialOut, err *primitive.RepoError) {
	ret := m.Called(ctx, email)

	var r1 *primitive.RepoError
	if err := ret.Get(1); err != nil {
		r1 = err.(*primitive.RepoError)
	}

	return ret.Get(0).(model.GetLoginCredentialOut), r1
}

func (m *MockBasicAuthLoginDb) GetEmployeeDetail(ctx context.Context, domain string, userId string) (employee model.GetEmployeeDetailOut, err *primitive.RepoError) {
	ret := m.Called(ctx, domain, userId)

	var r1 *primitive.RepoError
	if err := ret.Get(1); err != nil {
		r1 = err.(*primitive.RepoError)
	}

	return ret.Get(0).(model.GetEmployeeDetailOut), r1
}

type MockBasicAuthLoginBcrypt struct {
	mock.Mock
}

func (m *MockBasicAuthLoginBcrypt) CompareHashAndPassword(hashedPassword []byte, password []byte) error {
	ret := m.Called(hashedPassword, password)

	var r0 error
	if err := ret.Get(0); err != nil {
		r0 = err.(error)
	}

	return r0
}

type BasicAuthLoginTestSuite struct {
	suite.Suite

	db     *MockBasicAuthLoginDb
	bcrypt *MockBasicAuthLoginBcrypt

	validPayload service.BasicAuthLoginIn

	service service.BasicAuthLogin
}

func (t *BasicAuthLoginTestSuite) SetupSubTest() {
	db := new(MockBasicAuthLoginDb)
	bcrypt := new(MockBasicAuthLoginBcrypt)

	t.db = db
	t.bcrypt = bcrypt

	t.validPayload = authService.BasicAuthLoginIn{
		Email:    "mail@example.com",
		Password: "Password123",
	}

	t.service = authService.BasicAuthLogin{
		Db:     t.db,
		Bcrypt: bcrypt,
	}
}

func TestBasicAuthLoginTestSuite(t *testing.T) {
	suite.Run(t, new(BasicAuthLoginTestSuite))
}

func (t *BasicAuthLoginTestSuite) TestExec_InvalidPayload() {
	t.Run("should return error 400", func() {
		// arrange
		ctx := context.Background()

		// mock
		mockPayload := t.validPayload
		mockPayload.Email = ""

		// action
		out := t.service.Exec(ctx, mockPayload)

		// assert
		t.Equal(http.StatusBadRequest, out.GetCode())
		t.Equal("request validation failed", out.GetMessage())
		t.db.AssertNumberOfCalls(t.T(), "GetLoginCredential", 0)
	})
}

func (t *BasicAuthLoginTestSuite) TestExec_GetLoginCredential() {
	t.Run("server error", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.db.On("GetLoginCredential", ctx, t.validPayload.Email).Return(model.GetLoginCredentialOut{}, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   fmt.Errorf("mock server error"),
		})

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Equal(http.StatusInternalServerError, out.GetCode())
		t.Equal("internal server error", out.GetMessage())
		t.db.AssertExpectations(t.T())
		t.db.AssertNumberOfCalls(t.T(), "GetLoginCredential", 1)
	})

	t.Run("user not found", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.db.On("GetLoginCredential", ctx, t.validPayload.Email).Return(model.GetLoginCredentialOut{}, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeDataNotFound,
			Err:   fmt.Errorf("mock not found"),
		})

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Equal(http.StatusNotFound, out.GetCode())
		t.Equal("user not found", out.GetMessage())
		t.db.AssertExpectations(t.T())
		t.db.AssertNumberOfCalls(t.T(), "GetLoginCredential", 1)
	})
}

func (t *BasicAuthLoginTestSuite) TestExec_ComparePassword() {
	t.Run("invalid password", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.db.On("GetLoginCredential", ctx, t.validPayload.Email).Return(model.GetLoginCredentialOut{
			Password: primitive.String{String: "hashedPassword", Valid: true},
		}, nil)
		t.bcrypt.On("CompareHashAndPassword", []byte("hashedPassword"), []byte(t.validPayload.Password)).Return(fmt.Errorf("mock invalid password"))

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Equal(http.StatusUnauthorized, out.GetCode())
		t.Equal("invalid password", out.GetMessage())
		t.db.AssertExpectations(t.T())
		t.db.AssertNumberOfCalls(t.T(), "GetLoginCredential", 1)
		t.bcrypt.AssertExpectations(t.T())
		t.bcrypt.AssertNumberOfCalls(t.T(), "CompareHashAndPassword", 1)
	})
}

func (t *BasicAuthLoginTestSuite) TestExec_GetEmployeeDetail() {
	t.Run("server error", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.db.On("GetLoginCredential", ctx, t.validPayload.Email).Return(model.GetLoginCredentialOut{
			Password: primitive.String{String: "hashedPassword", Valid: true},
			UserUID:  "userId",
			Domain:   "domain",
			Email:    "mail@example.com",
		}, nil)
		t.bcrypt.On("CompareHashAndPassword", []byte("hashedPassword"), []byte(t.validPayload.Password)).Return(nil)
		t.db.On("GetEmployeeDetail", ctx, "domain", "userId").Return(model.GetEmployeeDetailOut{}, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   fmt.Errorf("mock server error"),
		})

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Equal(http.StatusInternalServerError, out.GetCode())
		t.Equal("internal server error", out.GetMessage())
		t.db.AssertExpectations(t.T())
		t.db.AssertNumberOfCalls(t.T(), "GetLoginCredential", 1)
		t.bcrypt.AssertExpectations(t.T())
		t.bcrypt.AssertNumberOfCalls(t.T(), "CompareHashAndPassword", 1)
		t.db.AssertNumberOfCalls(t.T(), "GetEmployeeDetail", 1)
	})

	t.Run("employee not found", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.db.On("GetLoginCredential", ctx, t.validPayload.Email).Return(model.GetLoginCredentialOut{
			Password: primitive.String{String: "hashedPassword", Valid: true},
			UserUID:  "userId",
			Domain:   "domain",
			Email:    "mail@example.com",
		}, nil)
		t.bcrypt.On("CompareHashAndPassword", []byte("hashedPassword"), []byte(t.validPayload.Password)).Return(nil)
		t.db.On("GetEmployeeDetail", ctx, "domain", "userId").Return(model.GetEmployeeDetailOut{}, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeDataNotFound,
			Err:   fmt.Errorf("mock not found"),
		})

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Equal(http.StatusNotFound, out.GetCode())
		t.Equal("employee not found", out.GetMessage())
		t.db.AssertExpectations(t.T())
		t.db.AssertNumberOfCalls(t.T(), "GetLoginCredential", 1)
		t.bcrypt.AssertExpectations(t.T())
	})

	t.Run("can login", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.db.On("GetLoginCredential", ctx, t.validPayload.Email).Return(model.GetLoginCredentialOut{
			Password: primitive.String{String: "hashedPassword", Valid: true},
			Domain:   "domain",
			UserUID:  "userId",
		}, nil)
		t.bcrypt.On("CompareHashAndPassword", []byte("hashedPassword"), []byte(t.validPayload.Password)).Return(nil)
		t.db.On("GetEmployeeDetail", ctx, "domain", "userId").Return(model.GetEmployeeDetailOut{
			Email:     "mail@example.com",
			FullName:  "John Doe",
			Gender:    primitive.GenderMale,
			BirthDate: primitive.Date{String: "2000-01-01", Valid: true},
			ProfilePicture: primitive.String{
				String: "https://example.com/profile.jpg",
				Valid:  true,
			},
			Address: primitive.String{
				String: "Jl. Example",
				Valid:  true,
			},
			JoinDate: primitive.Date{String: "2021-01-01", Valid: true},
		}, nil)

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Equal(http.StatusOK, out.GetCode())
		t.Equal("login success", out.GetMessage())
	})
}

func TestValidateBasicAuthLoginBody(t *testing.T) {
	service := authService.BasicAuthLogin{}

	type Test struct {
		name    string
		payload authService.BasicAuthLoginIn
		want    *primitive.RequestValidationError
	}
	positiveTest := []Test{
		{
			name: "payload is valid",
			payload: authService.BasicAuthLoginIn{
				Email:    "email@email.com",
				Password: "Password321",
			},
			want: nil,
		},
		{
			name: "payload is valid",
			payload: authService.BasicAuthLoginIn{
				Email:    "email@email.com",
				Password: "Password@321",
			},
			want: nil,
		},
	}

	for _, test := range positiveTest {
		t.Run(test.name, func(t *testing.T) {
			got := service.ValidateBasicAuthLoginBody(test.payload)
			if got != test.want {
				t.Errorf("ValidateBasicAuthLoginBody() = %v, want %v", got, test.want)
			}
		})
	}

	validPayload := authService.BasicAuthLoginIn{
		Email:    "email@email.com",
		Password: "Password321",
	}
	negativeTest := []Test{
		{
			name: "email is required",
			payload: authService.BasicAuthLoginIn{
				Password: validPayload.Password,
			},
			want: &primitive.RequestValidationError{
				Issues: []primitive.RequestValidationIssue{
					{
						Code:    primitive.RequestValidationCodeTooShort,
						Field:   "email",
						Message: "email is required",
					},
				},
			},
		},
		{
			name: "email is invalid",
			payload: authService.BasicAuthLoginIn{
				Email:    "email",
				Password: validPayload.Password,
			},
			want: &primitive.RequestValidationError{
				Issues: []primitive.RequestValidationIssue{
					{
						Code:    primitive.RequestValidationCodeInvalidValue,
						Field:   "email",
						Message: "invalid email address",
					},
				},
			},
		},
		{
			name: "password is required",
			payload: authService.BasicAuthLoginIn{
				Email: validPayload.Email,
			},
			want: &primitive.RequestValidationError{
				Issues: []primitive.RequestValidationIssue{
					{
						Code:    primitive.RequestValidationCodeTooShort,
						Field:   "password",
						Message: "password is required",
					},
				},
			},
		},
		{
			name: "password is too short",
			payload: authService.BasicAuthLoginIn{
				Email:    validPayload.Email,
				Password: "pass",
			},
			want: &primitive.RequestValidationError{
				Issues: []primitive.RequestValidationIssue{
					{
						Code:    primitive.RequestValidationCodeTooShort,
						Field:   "password",
						Message: "password must be at least 8 characters long",
					},
					{
						Code:    primitive.RequestValidationCodeInvalidValue,
						Field:   "password",
						Message: "must contain at least one uppercase character",
					},
					{
						Code:    primitive.RequestValidationCodeInvalidValue,
						Field:   "password",
						Message: "must contain at least one number",
					},
				},
			},
		},
		{
			name: "password contains no lowercase characters",
			payload: authService.BasicAuthLoginIn{
				Email:    validPayload.Email,
				Password: "PASSWORD123",
			},
			want: &primitive.RequestValidationError{
				Issues: []primitive.RequestValidationIssue{
					{
						Code:    primitive.RequestValidationCodeInvalidValue,
						Field:   "password",
						Message: "must contain at least one lowercase character",
					},
				},
			},
		},
		{
			name: "password contains no uppsercase characters",
			payload: authService.BasicAuthLoginIn{
				Email:    validPayload.Email,
				Password: "password321",
			},
			want: &primitive.RequestValidationError{
				Issues: []primitive.RequestValidationIssue{
					{
						Code:    primitive.RequestValidationCodeInvalidValue,
						Field:   "password",
						Message: "must contain at least one uppercase character",
					},
				},
			},
		},
		{
			name: "password contains no number",
			payload: authService.BasicAuthLoginIn{
				Email:    validPayload.Email,
				Password: "Passwords",
			},
			want: &primitive.RequestValidationError{
				Issues: []primitive.RequestValidationIssue{
					{
						Code:    primitive.RequestValidationCodeInvalidValue,
						Field:   "password",
						Message: "must contain at least one number",
					},
				},
			},
		},
	}

	for _, test := range negativeTest {
		t.Run(test.name, func(t *testing.T) {
			var expectedErrorType *primitive.RequestValidationError
			got := service.ValidateBasicAuthLoginBody(test.payload)
			if !errors.As(got, &expectedErrorType) {
				t.Errorf("ValidateBasicAuthLoginBody()\ngot: %T\nwant: %T", got, test.want)
			}

			if len(got.Issues) != len(test.want.Issues) {
				t.Errorf("Expect %d issues, got %d issues", len(test.want.Issues), len(got.Issues))
			} else {
				for i := range got.Issues {
					if !reflect.DeepEqual(got.Issues[i], test.want.Issues[i]) {
						t.Errorf("ValidateBasicAuthLoginBody() = %v, want: %v", got, test.want)
					}
				}
			}
		})
	}
}
