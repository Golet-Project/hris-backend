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

type MockChangePasswordDb struct {
	mock.Mock
}

func (m *MockChangePasswordDb) ChangePassword(ctx context.Context, param model.ChangePasswordIn) (rowsAffected int64, repoError *primitive.RepoError) {
	ret := m.Called(ctx, param)

	if err := ret.Get(1); err != nil {
		repoError = err.(*primitive.RepoError)
	}

	return ret.Get(0).(int64), repoError
}

func (m *MockChangePasswordDb) DeletePasswordRecoveryToken(ctx context.Context, userId string) (repoError *primitive.RepoError) {
	ret := m.Called(ctx, userId)

	if err := ret.Get(0); err != nil {
		repoError = err.(*primitive.RepoError)
	}

	return
}

type MockChangePasswordMemory struct {
	mock.Mock
}

func (m *MockChangePasswordMemory) GetPasswordRecoveryToken(ctx context.Context, userId string) (token string, repoError *primitive.RepoError) {
	ret := m.Called(ctx, userId)

	if err := ret.Get(1); err != nil {
		repoError = err.(*primitive.RepoError)
	}

	return ret.Get(0).(string), repoError
}

func (m *MockChangePasswordMemory) DeletePasswordRecoveryToken(ctx context.Context, userId string) (repoError *primitive.RepoError) {
	ret := m.Called(ctx, userId)

	if err := ret.Get(0); err != nil {
		repoError = err.(*primitive.RepoError)
	}

	return
}

type ChangePasswordTestSuite struct {
	suite.Suite

	db     *MockChangePasswordDb
	memory *MockChangePasswordMemory

	service authService.ChangePassword
}

func (s *ChangePasswordTestSuite) SetupSubTest() {
	db := new(MockChangePasswordDb)
	memory := new(MockChangePasswordMemory)

	s.db = db
	s.memory = memory

	s.service = authService.ChangePassword{
		Db:     db,
		Memory: memory,
	}
}

func TestChangePasswordSuite(t *testing.T) {
	suite.Run(t, new(ChangePasswordTestSuite))
}

func (s *ChangePasswordTestSuite) TestExec() {
	validPayload := authService.ChangePasswordIn{
		Token:    "change_password_token",
		UID:      "a84c2c59-748c-48d0-b628-4a73b1c3a8d7",
		Password: "NewPassword123",
	}

	s.Run("invalid payload", func() {
		mockPayload := validPayload
		mockPayload.Token = ""

		out := s.service.Exec(context.Background(), mockPayload)

		s.Assert().Equal(http.StatusBadRequest, out.GetCode())
		s.memory.AssertNumberOfCalls(s.T(), "GetPasswordRecoveryToken", 0)
		s.db.AssertNumberOfCalls(s.T(), "ChangePassword", 0)
	})

	s.Run("token doesn't exists", func() {
		// arrange
		ctx := context.Background()

		// mock
		s.memory.On("GetPasswordRecoveryToken", ctx, validPayload.UID).Return("", &primitive.RepoError{
			Issue: primitive.RepoErrorCodeDataNotFound,
			Err:   fmt.Errorf("token not found"),
		})

		// action
		out := s.service.Exec(ctx, validPayload)

		// assert
		s.Assert().Equal(http.StatusBadRequest, out.GetCode())
		s.Assert().Equal("password recovery token has expired", out.GetMessage())
		s.memory.AssertExpectations(s.T())
		s.memory.AssertNumberOfCalls(s.T(), "GetPasswordRecoveryToken", 1)
		s.memory.AssertExpectations(s.T())
		s.db.AssertNumberOfCalls(s.T(), "ChangePassword", 0)
	})

	s.Run("error when retrieving token", func() {
		// arrange
		ctx := context.Background()

		// mock
		s.memory.On("GetPasswordRecoveryToken", ctx, validPayload.UID).Return("", &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   fmt.Errorf("server error"),
		})

		// action
		out := s.service.Exec(ctx, validPayload)

		// assert
		s.Assert().Equal(http.StatusInternalServerError, out.GetCode())
		s.Assert().Equal("internal server error", out.GetMessage())
		s.memory.AssertExpectations(s.T())
		s.memory.AssertNumberOfCalls(s.T(), "GetPasswordRecoveryToken", 1)
		s.db.AssertNumberOfCalls(s.T(), "ChangePassword", 0)
	})

	s.Run("token mismatch", func() {
		// arrange
		ctx := context.Background()

		// mock
		s.memory.On("GetPasswordRecoveryToken", ctx, validPayload.UID).Return("mismatch_token", nil)

		// action
		out := s.service.Exec(ctx, validPayload)

		// assert
		s.Assert().Equal(http.StatusBadRequest, out.GetCode())
		s.Assert().Equal("password recovery token has expired", out.GetMessage())
		s.memory.AssertExpectations(s.T())
		s.memory.AssertNumberOfCalls(s.T(), "GetPasswordRecoveryToken", 1)
		s.db.AssertNumberOfCalls(s.T(), "ChangePassword", 0)
		s.memory.AssertExpectations(s.T())
	})

	s.Run("can change password", func() {
		// arrange
		ctx := context.Background()
		mockHashPassword := "$2a$10$4MX0foX8163XXVEWIpR9munu6UnyGoN61086iRsfJ7qH6rQ.PwsD"

		// mock
		s.service.GenerateFromPassword = func(password []byte, cost int) (hash []byte, err error) {
			return []byte(mockHashPassword), nil
		}
		s.memory.On("GetPasswordRecoveryToken", ctx, validPayload.UID).Return(validPayload.Token, nil)
		s.db.On("ChangePassword", ctx, model.ChangePasswordIn{
			UID:      validPayload.UID,
			Password: mockHashPassword,
		}).Return(int64(1), nil)
		s.memory.On("DeletePasswordRecoveryToken", ctx, validPayload.UID).Return(nil)

		// action
		out := s.service.Exec(ctx, validPayload)

		// assert
		s.Assert().Equal(http.StatusOK, out.GetCode())
		s.memory.AssertNumberOfCalls(s.T(), "DeletePasswordRecoveryToken", 1)
		s.memory.AssertCalled(s.T(), "DeletePasswordRecoveryToken", ctx, validPayload.UID)
		s.memory.AssertExpectations(s.T())
		s.db.AssertExpectations(s.T())
	})

	s.Run("should handle error 500 when invalidating password recovery token", func() {
		// arrange
		ctx := context.Background()
		mockHashPassword := "$2a$10$4MX0foX8163XXVEWIpR9munu6UnyGoN61086iRsfJ7qH6rQ.PwsD"

		// mock
		s.service.GenerateFromPassword = func(password []byte, cost int) (hash []byte, err error) {
			return []byte(mockHashPassword), nil
		}
		s.memory.
			On("GetPasswordRecoveryToken", ctx, validPayload.UID).Return(validPayload.Token, nil).
			On("DeletePasswordRecoveryToken", ctx, validPayload.UID).Return(&primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   fmt.Errorf("mock server error"),
		})

		s.db.On("ChangePassword", ctx, model.ChangePasswordIn{
			UID:      validPayload.UID,
			Password: mockHashPassword,
		}).Return(int64(1), nil)

		// action
		out := s.service.Exec(ctx, validPayload)

		// assert
		s.Assert().Equal(http.StatusInternalServerError, out.GetCode())
		s.memory.AssertNumberOfCalls(s.T(), "DeletePasswordRecoveryToken", 1)
		s.memory.AssertCalled(s.T(), "DeletePasswordRecoveryToken", ctx, validPayload.UID)
		s.memory.AssertExpectations(s.T())
		s.db.AssertExpectations(s.T())
	})

	s.Run("should not return an error if the token is not found when invalidating the token", func() {
		// arrange
		ctx := context.Background()
		mockHashPassword := "$2a$10$4MX0foX8163XXVEWIpR9munu6UnyGoN61086iRsfJ7qH6rQ.PwsD"

		// mock
		s.service.GenerateFromPassword = func(password []byte, cost int) (hash []byte, err error) {
			return []byte(mockHashPassword), nil
		}
		s.memory.
			On("GetPasswordRecoveryToken", ctx, validPayload.UID).Return(validPayload.Token, nil).
			On("DeletePasswordRecoveryToken", ctx, validPayload.UID).Return(&primitive.RepoError{
			Issue: primitive.RepoErrorCodeDataNotFound,
			Err:   fmt.Errorf("mock not found error"),
		})

		s.db.On("ChangePassword", ctx, model.ChangePasswordIn{
			UID:      validPayload.UID,
			Password: mockHashPassword,
		}).Return(int64(1), nil)

		// action
		out := s.service.Exec(ctx, validPayload)

		// assert
		s.Assert().Equal(http.StatusOK, out.GetCode())
		s.memory.AssertNumberOfCalls(s.T(), "DeletePasswordRecoveryToken", 1)
		s.memory.AssertCalled(s.T(), "DeletePasswordRecoveryToken", ctx, validPayload.UID)
		s.memory.AssertExpectations(s.T())
		s.db.AssertExpectations(s.T())
	})
}

func (s *ChangePasswordTestSuite) TestValidateChangePasswordRequest() {
	validPayload := authService.ChangePasswordIn{
		Token:    "change_password_token",
		UID:      "a84c2c59-748c-48d0-b628-4a73b1c3a8d7",
		Password: "NewPassword123",
	}

	s.Run("token required", func() {
		mockPayload := validPayload
		mockPayload.Token = ""
		var expectError *primitive.RequestValidationError

		err := s.service.ValidateChangePasswordRequest(mockPayload)

		if s.Assert().NotNil(err) {
			s.Assert().ErrorAs(err, &expectError)

			expectedIssue := primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeRequired,
				Field:   "token",
				Message: "token is required",
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

	s.Run("UID required", func() {
		mockPayload := validPayload
		mockPayload.UID = ""
		var expectError *primitive.RequestValidationError

		err := s.service.ValidateChangePasswordRequest(mockPayload)

		if s.Assert().NotNil(err) {
			s.Assert().ErrorAs(err, &expectError)
			expectedIssue := primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeRequired,
				Field:   "uid",
				Message: "uid is required",
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

	s.Run("UID must a valid uuid", func() {
		mockPayload := validPayload
		mockPayload.UID = "invalid uuid"
		var expectError *primitive.RequestValidationError

		err := s.service.ValidateChangePasswordRequest(mockPayload)

		if s.Assert().NotNil(err) {
			s.Assert().ErrorAs(err, &expectError)

			expectedIssue := primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "uid",
				Message: "uid is invalid",
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

	s.Run("UID must a valid uuid v4", func() {
		mockPayload := validPayload
		mockPayload.UID = "c20c0de6-d8b5-11ee-a506-0242ac120002"
		var expectError *primitive.RequestValidationError

		err := s.service.ValidateChangePasswordRequest(mockPayload)

		if s.Assert().NotNil(err) {
			s.Assert().ErrorAs(err, &expectError)

			expectedIssue := primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "uid",
				Message: "uid must be a valid uuid v4",
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

	s.Run("password required", func() {
		mockPayload := validPayload
		mockPayload.Password = ""
		var expectError *primitive.RequestValidationError

		err := s.service.ValidateChangePasswordRequest(mockPayload)

		if s.Assert().NotNil(err) {
			s.Assert().ErrorAs(err, &expectError)

			expectedIssue := primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeRequired,
				Field:   "password",
				Message: "password is required",
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

	s.Run("password must less than 25 characters", func() {
		mockPayload := validPayload
		mockPayload.Password = "Password1234567890Password1234567890"
		var expectError *primitive.RequestValidationError

		err := s.service.ValidateChangePasswordRequest(mockPayload)

		if s.Assert().NotNil(err) {
			s.Assert().ErrorAs(err, &expectError)

			expectedIssue := primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeTooLong,
				Field:   "password",
				Message: "password must be less than 25 characters",
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

	s.Run("password is too short", func() {
		mockPayload := validPayload
		mockPayload.Password = "pass"
		var expectError *primitive.RequestValidationError

		err := s.service.ValidateChangePasswordRequest(mockPayload)

		if s.Assert().NotNil(err) {
			s.Assert().ErrorAs(err, &expectError)

			expectedIssue := primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeTooShort,
				Field:   "password",
				Message: "password must be at least 8 characters long",
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

	s.Run("password contains no lowercase characters", func() {
		mockPayload := validPayload
		mockPayload.Password = "PASSWORD123"
		var expectError *primitive.RequestValidationError

		err := s.service.ValidateChangePasswordRequest(mockPayload)

		if s.Assert().NotNil(err) {
			s.Assert().ErrorAs(err, &expectError)

			expectedIssue := primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "password",
				Message: "must contain at least one lowercase character",
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

	s.Run("password contains no uppercase characters", func() {
		mockPyload := validPayload
		mockPyload.Password = "password321"
		var expectError *primitive.RequestValidationError

		err := s.service.ValidateChangePasswordRequest(mockPyload)

		if s.Assert().NotNil(err) {
			s.Assert().ErrorAs(err, &expectError)

			expectedIssue := primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "password",
				Message: "must contain at least one uppercase character",
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

	s.Run("password contains no number", func() {
		mockPayload := validPayload
		mockPayload.Password = "Passwords"
		var expectError *primitive.RequestValidationError

		err := s.service.ValidateChangePasswordRequest(mockPayload)

		if s.Assert().NotNil(err) {
			s.Assert().ErrorAs(err, &expectError)

			expectedIssue := primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "password",
				Message: "must contain at least one number",
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
}
