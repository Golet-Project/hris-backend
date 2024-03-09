package service_test

import (
	"context"
	"fmt"
	"hroost/mobile/domain/attendance/service"
	"hroost/shared/primitive"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockCheckoutDb struct {
	mock.Mock
}

func (m *MockCheckoutDb) GetDomainByUid(ctx context.Context, uid string) (domain string, repoError *primitive.RepoError) {
	ret := m.Called(ctx, uid)

	if err := ret.Get(1); err != nil {
		repoError = err.(*primitive.RepoError)
	}

	return ret.String(0), repoError
}

func (m *MockCheckoutDb) CheckTodayAttendanceById(ctx context.Context, domain string, uid string, timezone primitive.Timezone) (exists bool, repoError *primitive.RepoError) {
	ret := m.Called(ctx, domain, uid, timezone)

	if err := ret.Get(1); err != nil {
		repoError = err.(*primitive.RepoError)
	}

	return ret.Bool(0), repoError
}

func (m *MockCheckoutDb) Checkout(ctx context.Context, domain string, uid string) (rowsAffected int64, repoError *primitive.RepoError) {
	ret := m.Called(ctx, domain, uid)

	if err := ret.Get(1); err != nil {
		repoError = err.(*primitive.RepoError)
	}

	return ret.Get(0).(int64), repoError
}

type CheckoutTestSuite struct {
	suite.Suite

	db *MockCheckoutDb

	validPayload service.CheckoutIn

	service service.Checkout
}

func (t *CheckoutTestSuite) SetupSubTest() {
	db := new(MockCheckoutDb)

	t.db = db

	t.validPayload = service.CheckoutIn{
		UID:      "a84c2c59-748c-48d0-b628-4a73b1c3a8d7",
		Timezone: primitive.WIB,
	}

	t.service = service.Checkout{
		Db: db,
	}
}

func TestCheckoutTestSuite(t *testing.T) {
	suite.Run(t, new(CheckoutTestSuite))
}

func (t *CheckoutTestSuite) TestExec_InvalidRequestPayload() {
	t.Run("should return error 400", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.UID = ""

		// action
		out := t.service.Exec(context.Background(), mockPayload)

		// assert
		t.Assert().Equal(http.StatusBadRequest, out.GetCode())
		t.Assert().Equal("request validation failed", out.GetMessage())
	})
}

func (t *CheckoutTestSuite) TestExec_GetDomain() {
	t.Run("server error", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.db.On("GetDomainByUid", ctx, t.validPayload.UID).Return("", &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   fmt.Errorf("mock server error"),
		})

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Assert().Equal(http.StatusInternalServerError, out.GetCode())
		t.Assert().Equal("internal server error", out.GetMessage())
		t.db.AssertExpectations(t.T())
		t.db.AssertNumberOfCalls(t.T(), "Checkout", 0)
	})

	t.Run("user not found", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.db.On("GetDomainByUid", ctx, t.validPayload.UID).Return("", &primitive.RepoError{
			Issue: primitive.RepoErrorCodeDataNotFound,
			Err:   fmt.Errorf("mock not found error"),
		})

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Assert().Equal(http.StatusNotFound, out.GetCode())
		t.Assert().Equal("employee not found", out.GetMessage())
		t.db.AssertExpectations(t.T())
		t.db.AssertNumberOfCalls(t.T(), "Checkout", 0)
	})
}

func (t *CheckoutTestSuite) TestExec_CheckTodayAttendance() {
	t.Run("server error", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.db.On("GetDomainByUid", ctx, t.validPayload.UID).Return("mock_domain", nil)
		t.db.On("CheckTodayAttendanceById", ctx, "mock_domain", t.validPayload.UID, t.validPayload.Timezone).Return(false, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   fmt.Errorf("mock server error"),
		})

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Assert().Equal(http.StatusInternalServerError, out.GetCode())
		t.Assert().Equal("internal server error", out.GetMessage())
		t.db.AssertExpectations(t.T())
		t.db.AssertNumberOfCalls(t.T(), "Checkout", 0)
	})

	t.Run("not yet check in", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.db.On("GetDomainByUid", ctx, t.validPayload.UID).Return("mock_domain", nil)
		t.db.On("CheckTodayAttendanceById", ctx, "mock_domain", t.validPayload.UID, t.validPayload.Timezone).Return(false, nil)

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Assert().Equal(http.StatusBadRequest, out.GetCode())
		t.Assert().Equal("you haven't checkin yet", out.GetMessage())
		t.db.AssertExpectations(t.T())
		t.db.AssertNumberOfCalls(t.T(), "Checkout", 0)
	})
}

func (t *CheckoutTestSuite) TestExec_Checkout() {
	t.Run("can checkout", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.db.On("GetDomainByUid", ctx, t.validPayload.UID).Return("mock_domain", nil)
		t.db.On("CheckTodayAttendanceById", ctx, "mock_domain", t.validPayload.UID, t.validPayload.Timezone).Return(true, nil)
		t.db.On("Checkout", ctx, "mock_domain", t.validPayload.UID).Return(int64(1), nil)

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Assert().Equal(http.StatusCreated, out.GetCode())
		t.Assert().Equal("success", out.GetMessage())
		t.db.AssertExpectations(t.T())
		t.db.AssertNumberOfCalls(t.T(), "Checkout", 1)
		t.db.AssertCalled(t.T(), "Checkout", ctx, "mock_domain", t.validPayload.UID)
	})

	t.Run("server error", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.db.On("GetDomainByUid", ctx, t.validPayload.UID).Return("mock_domain", nil)
		t.db.On("CheckTodayAttendanceById", ctx, "mock_domain", t.validPayload.UID, t.validPayload.Timezone).Return(true, nil)
		t.db.On("Checkout", ctx, "mock_domain", t.validPayload.UID).Return(int64(0), &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   fmt.Errorf("mock server error"),
		})

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Assert().Equal(http.StatusInternalServerError, out.GetCode())
		t.Assert().Equal("internal server error", out.GetMessage())
		t.db.AssertExpectations(t.T())
		t.db.AssertNumberOfCalls(t.T(), "Checkout", 1)
		t.db.AssertCalled(t.T(), "Checkout", ctx, "mock_domain", t.validPayload.UID)
	})

	t.Run("no rows affected", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.db.On("GetDomainByUid", ctx, t.validPayload.UID).Return("mock_domain", nil)
		t.db.On("CheckTodayAttendanceById", ctx, "mock_domain", t.validPayload.UID, t.validPayload.Timezone).Return(true, nil)
		t.db.On("Checkout", ctx, "mock_domain", t.validPayload.UID).Return(int64(0), nil)

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Assert().Equal(http.StatusNotFound, out.GetCode())
		t.Assert().Equal("you already checkout", out.GetMessage())
		t.db.AssertExpectations(t.T())
		t.db.AssertNumberOfCalls(t.T(), "Checkout", 1)
		t.db.AssertCalled(t.T(), "Checkout", ctx, "mock_domain", t.validPayload.UID)
	})
}

func (t *CheckoutTestSuite) TestValidateRequestPayload_UID() {
	t.Run("UID is empty", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.UID = ""

		// action
		err := t.service.ValidateRequestPayload(mockPayload)

		// assert
		if t.Assert().NotNil(err) {
			var expectedError *primitive.RequestValidationError
			t.Assert().ErrorAs(err, &expectedError)
			t.Assert().Greater(len(err.Issues), 0)

			expectedIssue := primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeRequired,
				Field:   "uid",
				Message: "must not be empty",
			}
			t.Assert().Contains(err.Issues, expectedIssue)
		}
	})

	t.Run("invalid UUID", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.UID = "invalid UUID"

		// action
		err := t.service.ValidateRequestPayload(mockPayload)

		// assert
		if t.Assert().NotNil(err) {
			var expectedError *primitive.RequestValidationError
			t.Assert().ErrorAs(err, &expectedError)
			t.Assert().Greater(len(err.Issues), 0)

			expectedIssue := primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "uid",
				Message: "must be a valid UUID",
			}
			t.Assert().Contains(err.Issues, expectedIssue)
		}
	})

	t.Run("must be valid UUIDV4", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.UID = "115a0ee6-da39-11ee-a506-0242ac120002"

		// action
		err := t.service.ValidateRequestPayload(mockPayload)

		// assert
		if t.Assert().NotNil(err) {
			var expectedError *primitive.RequestValidationError
			t.Assert().ErrorAs(err, &expectedError)
			t.Assert().Greater(len(err.Issues), 0)

			expectedIssue := primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "uid",
				Message: "must be a valid UUIDV4",
			}
			t.Assert().Contains(err.Issues, expectedIssue)
		}
	})
}

func (t *CheckoutTestSuite) TestValidateRequestPayload_Timezone() {
	t.Run("empty timezone", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.Timezone = 0

		// action
		err := t.service.ValidateRequestPayload(mockPayload)

		// assert
		if t.Assert().NotNil(err) {
			var expectedError *primitive.RequestValidationError
			t.Assert().ErrorAs(err, &expectedError)
			t.Assert().Greater(len(err.Issues), 0)

			expectedIssue := primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeRequired,
				Field:   "timezone",
				Message: "must not be empty",
			}
			t.Assert().Contains(err.Issues, expectedIssue)
		}
	})

	t.Run("invalid timezone", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.Timezone = 10

		// action
		err := t.service.ValidateRequestPayload(mockPayload)

		// assert
		if t.Assert().NotNil(err) {
			var expectedError *primitive.RequestValidationError
			t.Assert().ErrorAs(err, &expectedError)
			t.Assert().Greater(len(err.Issues), 0)

			expectedIssue := primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "timezone",
				Message: "must be a valid timezone",
			}
			t.Assert().Contains(err.Issues, expectedIssue)
		}
	})
}

func (t *CheckoutTestSuite) TestValidateRequestPayload_ValidPayload() {
	t.Run("should return nil error", func() {
		validPayloads := []service.CheckoutIn{
			{
				UID:      t.validPayload.UID,
				Timezone: t.validPayload.Timezone,
			},
			{
				UID:      t.validPayload.UID,
				Timezone: primitive.WIB,
			},
			{
				UID:      t.validPayload.UID,
				Timezone: primitive.WIT,
			},
			{
				UID:      t.validPayload.UID,
				Timezone: primitive.WITA,
			},
		}

		for _, validPayload := range validPayloads {
			// action
			err := t.service.ValidateRequestPayload(validPayload)

			// assert
			t.Assert().Nil(err)
		}
	})
}
