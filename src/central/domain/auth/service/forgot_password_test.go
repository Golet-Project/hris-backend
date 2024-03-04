package service_test

import (
	"context"
	"fmt"
	"hroost/central/domain/auth/model"
	"hroost/shared/primitive"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	authService "hroost/central/domain/auth/service"
)

type MockForgotPasswordDb struct {
	mock.Mock
}

func (m *MockForgotPasswordDb) GetLoginCredential(ctx context.Context, email string) (model.GetLoginCredentialOut, *primitive.RepoError) {
	ret := m.Called(ctx, email)

	var repoError *primitive.RepoError
	if err := ret.Get(1); err != nil {
		repoError = err.(*primitive.RepoError)
	}

	return ret.Get(0).(model.GetLoginCredentialOut), repoError
}

type MockForgotPasswordMemory struct {
	mock.Mock
}

func (m *MockForgotPasswordMemory) GetPasswordRecoveryToken(ctx context.Context, userId string) (token string, repoErr *primitive.RepoError) {
	ret := m.Called(ctx, userId)

	if err := ret.Get(1); err != nil {
		repoErr = err.(*primitive.RepoError)
	}

	return ret.String(0), repoErr
}

func (m *MockForgotPasswordMemory) SetPasswordRecoveryToken(ctx context.Context, userId string, token string) (repoError *primitive.RepoError) {
	ret := m.Called(ctx, userId, token)

	if err := ret.Get(0); err != nil {
		repoError = err.(*primitive.RepoError)
	}

	return repoError
}

type ForgotPasswordTestSuite struct {
	suite.Suite

	db     *MockForgotPasswordDb
	memory *MockForgotPasswordMemory

	service authService.ForgotPassword
}

func (s *ForgotPasswordTestSuite) SetupSubTest() {
	db := new(MockForgotPasswordDb)
	memory := new(MockForgotPasswordMemory)

	s.db = db
	s.memory = memory

	s.service = authService.ForgotPassword{
		Db:     db,
		Memory: memory,
	}
}

func TestForgotPasswordTestSuite(t *testing.T) {
	suite.Run(t, new(ForgotPasswordTestSuite))
}

func (s *ForgotPasswordTestSuite) TestExec() {
	validPayload := authService.ForgotPasswordIn{
		Email: "mail@example.com",
		AppID: primitive.CentralAppID,
	}

	s.Run("invalid request payload", func() {
		// arrange
		mockPayload := validPayload
		mockPayload.Email = ""
		ctx := context.Background()

		// action
		out := s.service.Exec(ctx, mockPayload)

		// assert
		s.Assert().Equal(http.StatusBadRequest, out.GetCode())
		s.Assert().Equal("request validation failed", out.GetMessage())
		s.db.AssertNumberOfCalls(s.T(), "GetLoginCredential", 0)
		s.memory.AssertNumberOfCalls(s.T(), "SetPasswordRecoveryToken", 0)
	})

	s.Run("handle user not found when getting the credential", func() {
		// arrange
		ctx := context.Background()

		// mock
		s.db.On("GetLoginCredential", ctx, validPayload.Email).Return(model.GetLoginCredentialOut{}, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeDataNotFound,
			Err:   fmt.Errorf("mock not found error"),
		})

		// action
		out := s.service.Exec(ctx, validPayload)

		// assert
		s.Assert().Equal(http.StatusNotFound, out.GetCode())
		s.Assert().Equal("user not found", out.GetMessage())
		s.db.AssertExpectations(s.T())
		s.db.AssertNumberOfCalls(s.T(), "GetLoginCredential", 1)
		s.memory.AssertNumberOfCalls(s.T(), "SetPasswordRecoveryToken", 0)
	})

	s.Run("handle error 500 when getting the credential", func() {
		// arrange
		ctx := context.Background()

		// mock
		s.db.On("GetLoginCredential", ctx, validPayload.Email).Return(model.GetLoginCredentialOut{}, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   fmt.Errorf("mock not server error"),
		})

		// action
		out := s.service.Exec(ctx, validPayload)

		// assert
		s.Assert().Equal(http.StatusInternalServerError, out.GetCode())
		s.Assert().Equal("internal server error", out.GetMessage())
		s.db.AssertExpectations(s.T())
		s.db.AssertNumberOfCalls(s.T(), "GetLoginCredential", 1)
		s.memory.AssertNumberOfCalls(s.T(), "SetPasswordRecoveryToken", 0)
	})

	s.Run("should return 200 when password is empty (login with OAuth)", func() {
		// arrange
		ctx := context.Background()

		// mock
		s.db.On("GetLoginCredential", ctx, validPayload.Email).Return(model.GetLoginCredentialOut{
			UserUID:  "a84c2c59-748c-48d0-b628-4a73b1c3a8d7",
			Email:    "mail@example.com",
			Password: primitive.String{String: "", Valid: false},
		}, nil)

		// action
		out := s.service.Exec(ctx, validPayload)

		// assert
		s.Assert().Equal(http.StatusOK, out.GetCode())
		s.db.AssertExpectations(s.T())
		s.db.AssertNumberOfCalls(s.T(), "GetLoginCredential", 1)
		s.memory.AssertNumberOfCalls(s.T(), "SetPasswordRecoveryToken", 0)
	})

	s.Run("handle error 500 when checking password recovery token", func() {
		// arrange
		ctx := context.Background()

		// mock
		s.db.On("GetLoginCredential", ctx, validPayload.Email).Return(model.GetLoginCredentialOut{}, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   fmt.Errorf("mock server error"),
		})
		s.memory.On("GetPasswordRecoveryToken", ctx)

		// action
		out := s.service.Exec(ctx, validPayload)

		// assert
		s.Assert().Equal(http.StatusInternalServerError, out.GetCode())
		s.Assert().Equal("internal server error", out.GetMessage())
		s.db.AssertExpectations(s.T())
		s.db.AssertNumberOfCalls(s.T(), "GetLoginCredential", 1)
		s.memory.AssertNumberOfCalls(s.T(), "SetPasswordRecoveryToken", 0)
	})

	s.Run("should return reponse 400 when existing password recovery token is already exists", func() {
		// arrange
		ctx := context.Background()

		// mock
		s.db.On("GetLoginCredential", ctx, validPayload.Email).Return(model.GetLoginCredentialOut{
			UserUID:  "a84c2c59-748c-48d0-b628-4a73b1c3a8d7",
			Email:    "mail@example.com",
			Password: primitive.String{String: "mock password", Valid: true},
		}, nil)
		s.memory.On("GetPasswordRecoveryToken", ctx, "a84c2c59-748c-48d0-b628-4a73b1c3a8d7").Return("mock existing token", nil)

		// action
		out := s.service.Exec(ctx, validPayload)

		// assert
		s.Assert().Equal(http.StatusBadRequest, out.GetCode())
		s.Assert().Equal("password recovery link has already been sent to your email", out.GetMessage())
		s.db.AssertExpectations(s.T())
		s.db.AssertNumberOfCalls(s.T(), "GetLoginCredential", 1)
		s.memory.AssertExpectations(s.T())
		s.memory.AssertNumberOfCalls(s.T(), "GetPasswordRecoveryToken", 1)
		s.memory.AssertNumberOfCalls(s.T(), "SetPasswordRecoveryToken", 0)
	})

	s.Run("should set password recovery token when password recovery token doesn't exists", func() {
		s.T().SkipNow()
		// arrange
		ctx := context.Background()

		// mock
		s.db.On("GetLoginCredential", ctx, validPayload.Email).Return(model.GetLoginCredentialOut{
			UserUID:  "a84c2c59-748c-48d0-b628-4a73b1c3a8d7",
			Email:    "mail@example.com",
			Password: primitive.String{String: "mock password", Valid: true},
		}, nil)
		s.memory.
			On("GetPasswordRecoveryToken", ctx, "a84c2c59-748c-48d0-b628-4a73b1c3a8d7").Return("", &primitive.RepoError{
			Issue: primitive.RepoErrorCodeDataNotFound,
			Err:   fmt.Errorf("mock not found error"),
		}).
			On("SetPasswordRecoveryToken", ctx, "a84c2c59-748c-48d0-b628-4a73b1c3a8d7", mock.Anything).Return(nil)

		// action
		out := s.service.Exec(ctx, validPayload)

		// assert
		s.Assert().Equal(http.StatusOK, out.GetCode())
		s.Assert().Equal("password recovery link has already been sent to your email", out.GetMessage())
		s.db.AssertExpectations(s.T())
		s.db.AssertNumberOfCalls(s.T(), "GetLoginCredential", 1)
		s.memory.AssertExpectations(s.T())
		s.memory.AssertNumberOfCalls(s.T(), "GetPasswordRecoveryToken", 1)
		s.memory.AssertNumberOfCalls(s.T(), "SetPasswordRecoveryToken", 1)
	})
}

func (s *ForgotPasswordTestSuite) TestValidateForgotPasswordPayload() {
	validPayload := authService.ForgotPasswordIn{
		Email: "mail@example.com",
		AppID: primitive.CentralAppID,
	}

	s.Run("email required", func() {
		mockPayload := validPayload
		mockPayload.Email = ""
		var expectedError *primitive.RequestValidationError

		err := s.service.ValidateForgotPasswordPayload(mockPayload)

		if s.Assert().NotNil(err) {
			s.Assert().ErrorAs(err, &expectedError)
			s.Assert().Greater(len(err.Issues), 0)

			expectedIssue := primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeRequired,
				Field:   "email",
				Message: "email is required",
			}
			var containCorrectIssue = false
			for _, issue := range err.Issues {
				if assert.ObjectsAreEqual(expectedIssue, issue) {
					containCorrectIssue = true
					break
				}
			}
			s.Assert().True(containCorrectIssue)
		}
	})

	s.Run("email is valid", func() {
		mockPayload := validPayload
		mockPayload.Email = "invalid email"
		var expectedError *primitive.RequestValidationError

		err := s.service.ValidateForgotPasswordPayload(mockPayload)

		if s.Assert().NotNil(err) {
			s.Assert().ErrorAs(err, &expectedError)
			s.Assert().Greater(len(err.Issues), 0)

			expectedIssue := primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "email",
				Message: "invalid email address",
			}
			var containCorrectIssue = false
			for _, issue := range err.Issues {
				if assert.ObjectsAreEqual(expectedIssue, issue) {
					containCorrectIssue = true
					break
				}
			}
			s.Assert().True(containCorrectIssue)
		}
	})

	s.Run("App ID required", func() {
		mockPayload := validPayload
		mockPayload.AppID = ""
		var expectedError *primitive.RequestValidationError

		err := s.service.ValidateForgotPasswordPayload(mockPayload)

		if s.Assert().NotNil(err) {
			s.Assert().ErrorAs(err, &expectedError)
			s.Assert().Greater(len(err.Issues), 0)

			expectedIssue := primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeRequired,
				Field:   "X-App-ID",
				Message: "X-App-ID header is required",
			}
			var containCorrectIssue = false
			for _, issue := range err.Issues {
				if assert.ObjectsAreEqual(expectedIssue, issue) {
					containCorrectIssue = true
					break
				}
			}
			s.Assert().True(containCorrectIssue)
		}
	})

	s.Run("App ID is invalid", func() {
		mockPayload := validPayload
		mockPayload.AppID = primitive.MobileAppID
		var expectedError *primitive.RequestValidationError

		err := s.service.ValidateForgotPasswordPayload(mockPayload)

		if s.Assert().NotNil(err) {
			s.Assert().ErrorAs(err, &expectedError)
			s.Assert().Greater(len(err.Issues), 0)

			expectedIssue := primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeProhibitedValue,
				Field:   "X-App-ID",
				Message: "X-App-ID header has a prohibited value",
			}
			var containCorrectIssue = false
			for _, issue := range err.Issues {
				if assert.ObjectsAreEqual(expectedIssue, issue) {
					containCorrectIssue = true
					break
				}
			}
			s.Assert().True(containCorrectIssue)
		}
	})

	s.Run("valid payload", func() {
		err := s.service.ValidateForgotPasswordPayload(validPayload)

		s.Assert().Nil(err)
	})
}
