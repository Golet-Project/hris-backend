package service_test

import (
	"context"
	"fmt"
	"hroost/shared/primitive"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	authService "hroost/central/domain/auth/service"
)

type MockPasswordRecoveryTokenCheckMemory struct {
	mock.Mock
}

func (m *MockPasswordRecoveryTokenCheckMemory) GetPasswordRecoveryToken(ctx context.Context, userId string) (token string, repoError *primitive.RepoError) {
	ret := m.Called(ctx, userId)

	if err := ret.Get(1); err != nil {
		repoError = err.(*primitive.RepoError)
	}

	fmt.Println("RET", ret)

	return ret.String(0), repoError
}

type PasswordRecoveryTokenCheckTestSuite struct {
	suite.Suite

	memory *MockPasswordRecoveryTokenCheckMemory

	validPayload authService.PasswordRecoveryTokenCheckIn

	service authService.PasswordRecoveryTokenCheck
}

func (t *PasswordRecoveryTokenCheckTestSuite) SetupSubTest() {
	memory := new(MockPasswordRecoveryTokenCheckMemory)

	t.memory = memory

	t.validPayload = authService.PasswordRecoveryTokenCheckIn{
		Token: "example_token",
		UID:   "95bbfe2e-7801-46d0-8ee5-c0b839d596c8",
	}

	t.service = authService.PasswordRecoveryTokenCheck{
		Memory: memory,
	}
}

func TestPasswordRecoveryTokenTestSuite(t *testing.T) {
	suite.Run(t, new(PasswordRecoveryTokenCheckTestSuite))
}

func (t *PasswordRecoveryTokenCheckTestSuite) TestExec_InvalidPayload() {
	t.Run("should return correct response", func() {
		// arrange
		mockPayload := t.validPayload
		mockPayload.UID = ""
		ctx := context.Background()

		// action
		err := t.service.Exec(ctx, mockPayload)

		// assert
		t.Assert().Equal(http.StatusBadRequest, err.GetCode())
		t.Assert().Equal("request validation failed", err.GetMessage())
		t.memory.AssertNumberOfCalls(t.T(), "GetPasswordRecoveryToken", 0)
	})
}

func (t *PasswordRecoveryTokenCheckTestSuite) TestExec_GetRecoveryToken() {
	t.Run("correct token", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.memory.On("GetPasswordRecoveryToken", ctx, t.validPayload.UID).Return(t.validPayload.Token, nil)

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Assert().Equal(http.StatusNoContent, out.GetCode())
		t.Assert().Equal("OK", out.GetMessage())
		t.memory.AssertExpectations(t.T())
		t.memory.AssertNumberOfCalls(t.T(), "GetPasswordRecoveryToken", 1)
	})

	t.Run("token not found", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.memory.On("GetPasswordRecoveryToken", ctx, t.validPayload.UID).Return("", &primitive.RepoError{
			Issue: primitive.RepoErrorCodeDataNotFound,
			Err:   fmt.Errorf("mock not found error"),
		})

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Assert().Equal(http.StatusBadRequest, out.GetCode())
		t.Assert().Equal("password recovery token has expired", out.GetMessage())
		t.memory.AssertExpectations(t.T())
		t.memory.AssertNumberOfCalls(t.T(), "GetPasswordRecoveryToken", 1)
	})

	t.Run("error when getting token", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.memory.On("GetPasswordRecoveryToken", ctx, t.validPayload.UID).Return("", &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   fmt.Errorf("mock server error"),
		})

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Assert().Equal(http.StatusInternalServerError, out.GetCode())
		t.Assert().Equal("internal server error", out.GetMessage())
		t.memory.AssertExpectations(t.T())
		t.memory.AssertNumberOfCalls(t.T(), "GetPasswordRecoveryToken", 1)
	})

	t.Run("token mismatch", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.memory.On("GetPasswordRecoveryToken", ctx, t.validPayload.UID).Return("wrong_token", nil)

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Assert().Equal(http.StatusBadRequest, out.GetCode())
		t.Assert().Equal("password recovery token has expired", out.GetMessage())
		t.memory.AssertExpectations(t.T())
		t.memory.AssertNumberOfCalls(t.T(), "GetPasswordRecoveryToken", 1)
	})
}

func (t *PasswordRecoveryTokenCheckTestSuite) TestValidatePasswordRecoveryTokenCheckIn_ValidPayload() {
	t.Run("should return no error", func() {
		// action
		err := t.service.ValidatePasswordRecoveryTokenCheckIn(t.validPayload)

		// assert
		t.Assert().Nil(err)
	})
}

func (t *PasswordRecoveryTokenCheckTestSuite) TestValidatePasswordRecoveryTokenCheckIn_Token() {
	t.Run("token requied", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.Token = ""

		// action
		err := t.service.ValidatePasswordRecoveryTokenCheckIn(mockPayload)

		// assert
		if t.Assert().NotNil(err) {

			var expectedError *primitive.RequestValidationError
			t.Assert().ErrorAs(err, &expectedError)
			t.Assert().Greater(len(err.Issues), 0)

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
			t.Assert().True(containCorrectIssue)
		}
	})
}

func (t *PasswordRecoveryTokenCheckTestSuite) TestValidatePasswordRecoveryTokenCheckIn_UID() {
	t.Run("uid required", func() {
		// arrange
		mockPayload := t.validPayload
		mockPayload.UID = ""

		// action
		err := t.service.ValidatePasswordRecoveryTokenCheckIn(mockPayload)

		// assert
		if t.Assert().NotNil(err) {

			var expectedError *primitive.RequestValidationError
			t.Assert().ErrorAs(err, &expectedError)
			t.Assert().Greater(len(err.Issues), 0)

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
			t.Assert().True(containCorrectIssue)
		}
	})

	t.Run("uid is not a valid uuid", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.UID = "invalid uid"

		// action
		err := t.service.ValidatePasswordRecoveryTokenCheckIn(mockPayload)

		// assert
		if t.Assert().NotNil(err) {

			var expectedError *primitive.RequestValidationError
			t.Assert().ErrorAs(err, &expectedError)
			t.Assert().Greater(len(err.Issues), 0)

			expectedIssue := primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "uid",
				Message: "uid is not a valid uuid",
			}
			var containCorrectIssue = false
			for _, issue := range err.Issues {
				if assert.ObjectsAreEqual(expectedIssue, issue) {
					containCorrectIssue = true
					break
				}
			}
			t.Assert().True(containCorrectIssue)
		}
	})

	t.Run("uid is not UUIDV4", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.UID = "115a0ee6-da39-11ee-a506-0242ac120002"

		// action
		err := t.service.ValidatePasswordRecoveryTokenCheckIn(mockPayload)

		// assert
		if t.Assert().NotNil(err) {

			var expectedError *primitive.RequestValidationError
			t.Assert().ErrorAs(err, &expectedError)
			t.Assert().Greater(len(err.Issues), 0)

			expectedIssue := primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "uid",
				Message: "uid is not a valid UUIDV4",
			}
			var containCorrectIssue = false
			for _, issue := range err.Issues {
				if assert.ObjectsAreEqual(expectedIssue, issue) {
					containCorrectIssue = true
					break
				}
			}
			t.Assert().True(containCorrectIssue)
		}
	})
}
