package service_test

import (
	"context"
	"fmt"
	"hroost/central/domain/auth/model"
	authService "hroost/central/domain/auth/service"
	"hroost/shared/primitive"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockBasicAuthLoginDb struct {
	mock.Mock
}

func (m *MockBasicAuthLoginDb) GetLoginCredential(ctx context.Context, email string) (model.GetLoginCredentialOut, *primitive.RepoError) {
	ret := m.Called(ctx, email)

	var errRet *primitive.RepoError
	if err := ret.Get(1); err != nil {
		errRet = err.(*primitive.RepoError)
	}

	return ret.Get(0).(model.GetLoginCredentialOut), errRet
}

type BasicAuthLoginTestSuite struct {
	suite.Suite
	service authService.BasicAuthLogin
	db      *MockBasicAuthLoginDb
}

func TestBasicAuthLoginSuite(t *testing.T) {
	suite.Run(t, new(BasicAuthLoginTestSuite))
}

func (s *BasicAuthLoginTestSuite) SetupSubTest() {
	fmt.Println("SETUP SUB TEST")
	db := new(MockBasicAuthLoginDb)

	s.db = db
	s.service = authService.BasicAuthLogin{
		Db: db,
	}
}

func (s *BasicAuthLoginTestSuite) TestBasicAuthLogin_Exec() {
	validPayload := authService.BasicAuthLoginIn{
		Email:    "mail@example.com",
		Password: "SuperSecretPassword123",
	}
	const correctPasswordHash = "$2a$10$qVaBCtt3TDvH.Nz3s4B73eTuTZhLaMI8AbnnWlBsBtdJgTzT82Paa"

	s.Run("invalid request body", func() {
		mockPayload := validPayload
		mockPayload.Password = "weak"

		out := s.service.Exec(context.Background(), mockPayload)

		s.db.AssertNumberOfCalls(s.T(), "GetLoginCredential", 0)
		s.Assert().Equal(http.StatusBadRequest, out.GetCode())
	})

	s.Run("get login credential", func() {
		s.Run("user not found", func() {
			// arrange
			mockPayload := validPayload
			mockPayload.Email = "notfound@example.com"

			// mock
			s.db.On("GetLoginCredential", context.Background(), validPayload.Email).Return(model.GetLoginCredentialOut{}, &primitive.RepoError{
				Issue: primitive.RepoErrorCodeDataNotFound,
			})

			// action
			out := s.service.Exec(context.Background(), validPayload)

			// assert
			s.db.AssertExpectations(s.T())
			s.Assert().Equal(http.StatusNotFound, out.GetCode())
		})

		s.Run("internal server error", func() {
			// mock
			s.db.On("GetLoginCredential", context.Background(), validPayload.Email).Return(model.GetLoginCredentialOut{}, &primitive.RepoError{
				Issue: primitive.RepoErrorCodeServerError,
			})

			// action
			out := s.service.Exec(context.Background(), validPayload)

			// assert
			s.db.AssertExpectations(s.T())
			s.Assert().Equal(http.StatusInternalServerError, out.GetCode())
			s.Assert().Equal("internal server error", out.GetMessage())
		})
	})

	s.Run("compare password", func() {
		s.Run("invalid password", func() {
			// mock
			mockPayload := validPayload
			mockPayload.Password = "InvalidPassword123"

			s.db.On("GetLoginCredential", context.Background(), mockPayload.Email).Return(model.GetLoginCredentialOut{
				UserUID:  "a84c2c59-748c-48d0-b628-4a73b1c3a8d7",
				Email:    mockPayload.Email,
				Password: primitive.String{String: correctPasswordHash, Valid: true},
			}, nil)

			// action
			out := s.service.Exec(context.Background(), mockPayload)

			// assert
			s.db.AssertNumberOfCalls(s.T(), "GetLoginCredential", 1)
			s.Assert().Equal(http.StatusUnauthorized, out.GetCode())
		})
	})

	s.Run("valid payload", func() {
		s.db.On("GetLoginCredential", context.Background(), validPayload.Email).Return(model.GetLoginCredentialOut{
			UserUID:  "a84c2c59-748c-48d0-b628-4a73b1c3a8d7",
			Email:    "mail@example.com",
			Password: primitive.String{String: correctPasswordHash, Valid: true},
		}, nil)

		out := s.service.Exec(context.Background(), validPayload)

		s.Assert().NotEmpty(out.AccessToken)
		s.Assert().Equal(http.StatusOK, out.GetCode())
		s.db.AssertNumberOfCalls(s.T(), "GetLoginCredential", 1)
	})
}

func (s *BasicAuthLoginTestSuite) TestValidateBasicAuthLoginBody() {
	s.Run("payload is valid (1)", func() {
		payload := authService.BasicAuthLoginIn{
			Email:    "email@email.com",
			Password: "Password@321",
		}

		err := s.service.ValidateBasicAuthLoginBody(payload)

		s.Assert().Nil(err)
	})

	validPayload := authService.BasicAuthLoginIn{
		Email:    "email@email.com",
		Password: "Password@321",
	}

	s.Run("email is missing", func() {
		mock := validPayload
		mock.Email = ""
		var requestValidationError *primitive.RequestValidationError

		err := s.service.ValidateBasicAuthLoginBody(mock)

		s.Assert().NotNil(err)
		s.Assert().ErrorAs(err, &requestValidationError)
	})

	s.Run("email is invalid", func() {
		mock := validPayload
		mock.Email = "email"
		var requestValidationError *primitive.RequestValidationError

		got := s.service.ValidateBasicAuthLoginBody(mock)

		s.Assert().NotNil(got)
		s.Assert().ErrorAs(got, &requestValidationError)
	})

	s.Run("password is required", func() {
		mock := validPayload
		mock.Password = ""
		var requestValidationError *primitive.RequestValidationError

		got := s.service.ValidateBasicAuthLoginBody(mock)

		s.Assert().NotNil(got)
		s.Assert().ErrorAs(got, &requestValidationError)
	})

	s.Run("password is too short", func() {
		mock := validPayload
		mock.Password = "pass"
		var requestValidationError *primitive.RequestValidationError

		got := s.service.ValidateBasicAuthLoginBody(mock)

		s.Assert().NotNil(got)
		s.Assert().ErrorAs(got, &requestValidationError)
	})

	s.Run("password contains no lowercase characters", func() {
		mock := validPayload
		mock.Password = "PASSWORD123"
		var requestValidationError *primitive.RequestValidationError

		got := s.service.ValidateBasicAuthLoginBody(mock)

		s.Assert().NotNil(got)
		s.Assert().ErrorAs(got, &requestValidationError)
	})

	s.Run("password contains no uppercase characters", func() {
		mock := validPayload
		mock.Password = "password321"
		var requestValidationError *primitive.RequestValidationError

		got := s.service.ValidateBasicAuthLoginBody(mock)

		s.Assert().NotNil(got)
		s.Assert().ErrorAs(got, &requestValidationError)
	})

	s.Run("password contains no number", func() {
		mock := validPayload
		mock.Password = "Passwords"
		var requestValidationError *primitive.RequestValidationError

		got := s.service.ValidateBasicAuthLoginBody(mock)

		s.Assert().NotNil(got)
		s.Assert().ErrorAs(got, &requestValidationError)
	})
}
