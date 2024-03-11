package service_test

import (
	"context"
	"encoding/json"
	"fmt"
	"hroost/mobile/domain/attendance/model"
	"hroost/mobile/domain/attendance/service"
	"hroost/shared/primitive"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockAttendanceHistoryDb struct {
	mock.Mock
}

func (m *MockAttendanceHistoryDb) GetDomainByUid(ctx context.Context, uid string) (domain string, err *primitive.RepoError) {
	ret := m.Called(ctx, uid)

	var r0 string
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(string)
	}

	var r1 *primitive.RepoError
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(*primitive.RepoError)
	}
	return r0, r1
}

func (m *MockAttendanceHistoryDb) FindAttendanceHistory(ctx context.Context, domain string, in model.FindAttendanceHistoryIn) (length int64, attendances []model.FindAttendanceHistoryOut, err *primitive.RepoError) {
	ret := m.Called(ctx, domain, in)

	var r1 []model.FindAttendanceHistoryOut
	if ret.Get(1) != nil {
		r1 = ret.Get(1).([]model.FindAttendanceHistoryOut)
	}

	var r2 *primitive.RepoError
	if ret.Get(2) != nil {
		r2 = ret.Get(2).(*primitive.RepoError)
	}
	return ret.Get(0).(int64), r1, r2
}

type AttendanceHistoryTestSuite struct {
	suite.Suite

	db *MockAttendanceHistoryDb

	validPayload service.AttendanceHistoryIn

	service service.AttendanceHistory
}

func (t *AttendanceHistoryTestSuite) SetupSubTest() {
	db := new(MockAttendanceHistoryDb)

	t.db = db
	t.validPayload = service.AttendanceHistoryIn{
		EmployeeId: "a84c2c59-748c-48d0-b628-4a73b1c3a8d7",
		StartDate:  "2024-03-01",
		EndDate:    "2024-03-31",
	}

	t.service = service.AttendanceHistory{
		Db: db,

		GetNow: t.getNow,
	}
}

func (t *AttendanceHistoryTestSuite) getNow() time.Time {
	t.T().Helper()

	return time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
}

func TestAttendanceHistoryTestSuite(t *testing.T) {
	suite.Run(t, new(AttendanceHistoryTestSuite))
}

func (t *AttendanceHistoryTestSuite) TestExec_InvalidPayload() {
	t.Run("one of date filter is empty", func() {
		// arrange
		ctx := context.Background()

		// mock
		mockPayload := t.validPayload
		mockPayload.StartDate = ""

		// action
		out := t.service.Exec(ctx, mockPayload)

		// assert
		t.Equal(http.StatusBadRequest, out.GetCode())
		t.Equal("request validation failed", out.GetMessage())
		t.db.AssertNumberOfCalls(t.T(), "GetDomainByUid", 0)
		t.db.AssertNumberOfCalls(t.T(), "FindAttendanceHistory", 0)
	})
}

func (t *AttendanceHistoryTestSuite) TestExec_GetDomainByUid() {
	t.Run("server error", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.db.On("GetDomainByUid", ctx, t.validPayload.EmployeeId).Return("", &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   fmt.Errorf("mock server error"),
		})

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Equal(http.StatusInternalServerError, out.GetCode())
		t.Equal("internal server error", out.GetMessage())
		t.db.AssertExpectations(t.T())
		t.db.AssertNumberOfCalls(t.T(), "GetDomainByUid", 1)
		t.db.AssertNumberOfCalls(t.T(), "FindAttendanceHistory", 0)
	})

	t.Run("user not found", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.db.On("GetDomainByUid", ctx, t.validPayload.EmployeeId).Return("", &primitive.RepoError{
			Issue: primitive.RepoErrorCodeDataNotFound,
			Err:   fmt.Errorf("mock not found error"),
		})

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Equal(http.StatusNotFound, out.GetCode())
		t.Equal("user not found", out.GetMessage())
		t.db.AssertExpectations(t.T())
		t.db.AssertNumberOfCalls(t.T(), "GetDomainByUid", 1)
		t.db.AssertNumberOfCalls(t.T(), "FindAttendanceHistory", 0)
	})
}

func (t *AttendanceHistoryTestSuite) TestExec_FindAttendanceHistory() {
	t.Run("server error", func() {
		// arrange
		ctx := context.Background()
		domain := "mock_domain"

		// mock
		t.db.On("GetDomainByUid", ctx, t.validPayload.EmployeeId).Return(domain, nil)
		t.db.On("FindAttendanceHistory", ctx, domain, model.FindAttendanceHistoryIn{
			StartDate: t.validPayload.StartDate,
			EndDate:   t.validPayload.EndDate,
		}).Return(int64(0), []model.FindAttendanceHistoryOut{}, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   fmt.Errorf("mock server error"),
		})

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Equal(http.StatusInternalServerError, out.GetCode())
		t.Equal("internal server error", out.GetMessage())
		t.db.AssertExpectations(t.T())
		t.db.AssertNumberOfCalls(t.T(), "GetDomainByUid", 1)
		t.db.AssertNumberOfCalls(t.T(), "FindAttendanceHistory", 1)
	})

	t.Run("return empty array if not found", func() {
		// arrange
		ctx := context.Background()
		domain := "mock_domain"

		// mock
		t.db.On("GetDomainByUid", ctx, t.validPayload.EmployeeId).Return(domain, nil)
		t.db.On("FindAttendanceHistory", ctx, domain, model.FindAttendanceHistoryIn{
			StartDate: t.validPayload.StartDate,
			EndDate:   t.validPayload.EndDate,
		}).Return(int64(0), []model.FindAttendanceHistoryOut{}, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeDataNotFound,
			Err:   fmt.Errorf("mock not found error"),
		})

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Equal(http.StatusOK, out.GetCode())
		t.Equal(int64(0), out.Length)
		t.Equal("success", out.GetMessage())
		t.Empty(out.Attendances)
		t.db.AssertExpectations(t.T())
		t.db.AssertNumberOfCalls(t.T(), "GetDomainByUid", 1)
		t.db.AssertNumberOfCalls(t.T(), "FindAttendanceHistory", 1)
	})

	t.Run("can get attendance history with default date filter if date filter is not provided", func() {
		// arrange
		ctx := context.Background()
		domain := "mock_domain"
		now := t.getNow()
		startDate := now.AddDate(0, 0, -now.Day()+1).Format("2006-01-02")
		endDate := now.AddDate(0, 1, -now.Day()).Format("2006-01-02")

		// mock
		mockPayload := t.validPayload
		mockPayload.StartDate = ""
		mockPayload.EndDate = ""

		t.db.On("GetDomainByUid", ctx, t.validPayload.EmployeeId).Return(domain, nil)
		t.db.On("FindAttendanceHistory", ctx, domain, model.FindAttendanceHistoryIn{
			StartDate: startDate,
			EndDate:   endDate,
		}).Return(int64(2), []model.FindAttendanceHistoryOut{
			{
				ID:           "mock_id",
				Date:         primitive.Date{String: "2024-03-01", Valid: true},
				CheckinTime:  primitive.Time{Time: time.Date(2024, 3, 1, 9, 0, 0, 0, time.UTC), Valid: true},
				CheckoutTime: primitive.Time{Time: time.Date(2024, 3, 1, 17, 0, 0, 0, time.UTC), Valid: true},
				ApprovedAt:   primitive.Time{Time: time.Time{}, Valid: false},
				Radius:       100,

				Coordinate: primitive.Coordinate{
					Latitude:  -7.782870711329031,
					Longitude: 110.36707035197001,
				},
			},
			{
				ID:           "mock_id_2",
				Date:         primitive.Date{String: "2024-03-02", Valid: true},
				CheckinTime:  primitive.Time{Time: time.Date(2024, 3, 2, 9, 0, 0, 0, time.UTC), Valid: true},
				CheckoutTime: primitive.Time{Time: time.Date(2024, 3, 2, 17, 0, 0, 0, time.UTC), Valid: true},
				ApprovedAt:   primitive.Time{Time: time.Time{}, Valid: false},
				Radius:       100,

				Coordinate: primitive.Coordinate{
					Latitude:  -7.782870711329031,
					Longitude: 110.36707035197001,
				},
			},
			{
				ID:           "mock_id_3",
				Date:         primitive.Date{String: "2024-03-03", Valid: true},
				CheckinTime:  primitive.Time{Time: time.Date(2024, 3, 3, 9, 0, 0, 0, time.UTC), Valid: true},
				CheckoutTime: primitive.Time{Time: time.Date(2024, 3, 3, 17, 0, 0, 0, time.UTC), Valid: true},
				ApprovedAt:   primitive.Time{Time: time.Date(2024, 3, 3, 10, 0, 0, 0, time.UTC), Valid: true},
				Radius:       100,

				Coordinate: primitive.Coordinate{
					Latitude:  -7.782870711329031,
					Longitude: 110.36707035197001,
				},
			},
		}, nil)

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		type Attendance struct {
			ID               string `json:"id"`
			Date             string `json:"date"`
			CheckinTime      string `json:"checkin_time"`
			CheckoutTime     string `json:"checkout_time"`
			ApprovedAt       string `json:"approved_at"`
			AttendanceRadius int64  `json:"attendance_radius"` // radius in meters
		}
		expectedResponse := struct {
			Attendances []Attendance `json:"attendances"`
		}{
			Attendances: []Attendance{
				{
					ID:               "mock_id",
					Date:             "2024-03-01",
					CheckinTime:      time.Date(2024, 3, 1, 9, 0, 0, 0, time.UTC).Format(primitive.UtcRFC3339),
					CheckoutTime:     time.Date(2024, 3, 1, 17, 0, 0, 0, time.UTC).Format(primitive.UtcRFC3339),
					ApprovedAt:       "",
					AttendanceRadius: 100,
				},
				{
					ID:               "mock_id_2",
					Date:             "2024-03-02",
					CheckinTime:      time.Date(2024, 3, 2, 9, 0, 0, 0, time.UTC).Format(primitive.UtcRFC3339),
					CheckoutTime:     time.Date(2024, 3, 2, 17, 0, 0, 0, time.UTC).Format(primitive.UtcRFC3339),
					ApprovedAt:       "",
					AttendanceRadius: 100,
				},
				{
					ID:               "mock_id_3",
					Date:             "2024-03-03",
					CheckinTime:      time.Date(2024, 3, 3, 9, 0, 0, 0, time.UTC).Format(primitive.UtcRFC3339),
					CheckoutTime:     time.Date(2024, 3, 3, 17, 0, 0, 0, time.UTC).Format(primitive.UtcRFC3339),
					ApprovedAt:       time.Date(2024, 3, 3, 10, 0, 0, 0, time.UTC).Format(primitive.UtcRFC3339),
					AttendanceRadius: 100,
				},
			},
		}
		expectedJson, _ := json.Marshal(expectedResponse)
		respnseJson, err := json.Marshal(out)
		if !t.NoError(err) {
			return
		}

		t.Equal(http.StatusOK, out.GetCode())
		t.Equal(int64(2), out.Length)
		t.JSONEq(string(expectedJson), string(respnseJson))
		t.Equal("success", out.GetMessage())
		t.db.AssertExpectations(t.T())
		t.db.AssertNumberOfCalls(t.T(), "GetDomainByUid", 1)
		t.db.AssertNumberOfCalls(t.T(), "FindAttendanceHistory", 1)
	})

	t.Run("can get attendance history", func() {
		// arrange
		ctx := context.Background()
		domain := "mock_domain"

		// mock
		t.db.On("GetDomainByUid", ctx, t.validPayload.EmployeeId).Return(domain, nil)
		t.db.On("FindAttendanceHistory", ctx, domain, model.FindAttendanceHistoryIn{
			StartDate: t.validPayload.StartDate,
			EndDate:   t.validPayload.EndDate,
		}).Return(int64(2), []model.FindAttendanceHistoryOut{
			{
				ID:           "mock_id",
				Date:         primitive.Date{String: "2024-03-01", Valid: true},
				CheckinTime:  primitive.Time{Time: time.Date(2024, 3, 1, 9, 0, 0, 0, time.UTC), Valid: true},
				CheckoutTime: primitive.Time{Time: time.Date(2024, 3, 1, 17, 0, 0, 0, time.UTC), Valid: true},
				ApprovedAt:   primitive.Time{Time: time.Time{}, Valid: false},
				Radius:       100,

				Coordinate: primitive.Coordinate{
					Latitude:  -7.782870711329031,
					Longitude: 110.36707035197001,
				},
			},
			{
				ID:           "mock_id_2",
				Date:         primitive.Date{String: "2024-03-02", Valid: true},
				CheckinTime:  primitive.Time{Time: time.Date(2024, 3, 2, 9, 0, 0, 0, time.UTC), Valid: true},
				CheckoutTime: primitive.Time{Time: time.Date(2024, 3, 2, 17, 0, 0, 0, time.UTC), Valid: true},
				ApprovedAt:   primitive.Time{Time: time.Time{}, Valid: false},
				Radius:       100,

				Coordinate: primitive.Coordinate{
					Latitude:  -7.782870711329031,
					Longitude: 110.36707035197001,
				},
			},
			{
				ID:           "mock_id_3",
				Date:         primitive.Date{String: "2024-03-03", Valid: true},
				CheckinTime:  primitive.Time{Time: time.Date(2024, 3, 3, 9, 0, 0, 0, time.UTC), Valid: true},
				CheckoutTime: primitive.Time{Time: time.Date(2024, 3, 3, 17, 0, 0, 0, time.UTC), Valid: true},
				ApprovedAt:   primitive.Time{Time: time.Date(2024, 3, 3, 10, 0, 0, 0, time.UTC), Valid: true},
				Radius:       100,

				Coordinate: primitive.Coordinate{
					Latitude:  -7.782870711329031,
					Longitude: 110.36707035197001,
				},
			},
		}, nil)

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		type Attendance struct {
			ID               string `json:"id"`
			Date             string `json:"date"`
			CheckinTime      string `json:"checkin_time"`
			CheckoutTime     string `json:"checkout_time"`
			ApprovedAt       string `json:"approved_at"`
			AttendanceRadius int64  `json:"attendance_radius"` // radius in meters
		}
		expectedResponse := struct {
			Attendances []Attendance `json:"attendances"`
		}{
			Attendances: []Attendance{
				{
					ID:               "mock_id",
					Date:             "2024-03-01",
					CheckinTime:      time.Date(2024, 3, 1, 9, 0, 0, 0, time.UTC).Format(primitive.UtcRFC3339),
					CheckoutTime:     time.Date(2024, 3, 1, 17, 0, 0, 0, time.UTC).Format(primitive.UtcRFC3339),
					ApprovedAt:       "",
					AttendanceRadius: 100,
				},
				{
					ID:               "mock_id_2",
					Date:             "2024-03-02",
					CheckinTime:      time.Date(2024, 3, 2, 9, 0, 0, 0, time.UTC).Format(primitive.UtcRFC3339),
					CheckoutTime:     time.Date(2024, 3, 2, 17, 0, 0, 0, time.UTC).Format(primitive.UtcRFC3339),
					ApprovedAt:       "",
					AttendanceRadius: 100,
				},
				{
					ID:               "mock_id_3",
					Date:             "2024-03-03",
					CheckinTime:      time.Date(2024, 3, 3, 9, 0, 0, 0, time.UTC).Format(primitive.UtcRFC3339),
					CheckoutTime:     time.Date(2024, 3, 3, 17, 0, 0, 0, time.UTC).Format(primitive.UtcRFC3339),
					ApprovedAt:       time.Date(2024, 3, 3, 10, 0, 0, 0, time.UTC).Format(primitive.UtcRFC3339),
					AttendanceRadius: 100,
				},
			},
		}
		expectedJson, _ := json.Marshal(expectedResponse)
		respnseJson, err := json.Marshal(out)
		if !t.NoError(err) {
			return
		}

		t.Equal(http.StatusOK, out.GetCode())
		t.Equal(int64(2), out.Length)
		t.JSONEq(string(expectedJson), string(respnseJson))
		t.Equal("success", out.GetMessage())
		t.db.AssertExpectations(t.T())
		t.db.AssertNumberOfCalls(t.T(), "GetDomainByUid", 1)
		t.db.AssertNumberOfCalls(t.T(), "FindAttendanceHistory", 1)
	})
}

func (t *AttendanceHistoryTestSuite) TestValidateRequest_EmployeeId() {
	t.Run("empty", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.EmployeeId = ""

		// action
		err := t.service.ValidateRequest(mockPayload)

		// assert
		if t.NotNil(err) {
			var expectedError *primitive.RequestValidationError
			t.ErrorAs(err, &expectedError)
			t.Greater(len(err.Issues), 0)

			var correctIssue = primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeRequired,
				Field:   "employee_id",
				Message: "must not be empty",
			}
			t.Contains(err.Issues, correctIssue)
		}
	})

	t.Run("is not UUID", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.EmployeeId = "not-uuid"

		// action
		err := t.service.ValidateRequest(mockPayload)

		// assert
		if t.NotNil(err) {
			var expectedError *primitive.RequestValidationError
			t.ErrorAs(err, &expectedError)
			t.Greater(len(err.Issues), 0)

			var correctIssue = primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "employee_id",
				Message: "must be a valid UUID",
			}
			t.Contains(err.Issues, correctIssue)
		}
	})

	t.Run("not UUIDV4", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.EmployeeId = "115a0ee6-da39-11ee-a506-0242ac120002"

		// action
		err := t.service.ValidateRequest(mockPayload)

		// assert
		if t.NotNil(err) {
			var expectedError *primitive.RequestValidationError
			t.ErrorAs(err, &expectedError)
			t.Greater(len(err.Issues), 0)

			var correctIssue = primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "employee_id",
				Message: "must be a valid UUIDV4",
			}
			t.Contains(err.Issues, correctIssue)
		}
	})
}

func (t *AttendanceHistoryTestSuite) TestValidateRequest_StartDate() {
	t.Run("can't empty when end_date is not empty", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.StartDate = ""

		// action
		err := t.service.ValidateRequest(mockPayload)

		// assert
		if t.NotNil(err) {
			var expectedError *primitive.RequestValidationError
			t.ErrorAs(err, &expectedError)
			t.Greater(len(err.Issues), 0)

			var correctIssue = primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeRequired,
				Field:   "start_date",
				Message: "must not be empty if end_date is not empty",
			}
			t.Contains(err.Issues, correctIssue)
		}
	})

	t.Run("can empty when end_date is empty", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.StartDate = ""
		mockPayload.EndDate = ""

		// action
		err := t.service.ValidateRequest(mockPayload)

		// assert
		t.Nil(err)
	})

	t.Run("invalid date", func() {
		// mock
		mockPayloads := []service.AttendanceHistoryIn{
			{
				EmployeeId: t.validPayload.EmployeeId,
				StartDate:  "not date",
				EndDate:    t.validPayload.EndDate,
			},
			{
				EmployeeId: t.validPayload.EmployeeId,
				StartDate:  "2024-01-32",
				EndDate:    t.validPayload.EndDate,
			},
			{
				EmployeeId: t.validPayload.EmployeeId,
				StartDate:  "2023-02-29",
				EndDate:    t.validPayload.EndDate,
			},
		}

		for _, mockPayload := range mockPayloads {
			// action
			err := t.service.ValidateRequest(mockPayload)

			// asset
			if t.NotNil(err) {
				var expectedError *primitive.RequestValidationError
				t.ErrorAs(err, &expectedError)
				t.Greater(len(err.Issues), 0)

				var correctIssue = primitive.RequestValidationIssue{
					Code:    primitive.RequestValidationCodeInvalidValue,
					Field:   "start_date",
					Message: "must be in the YYYY-MM-DD format",
				}
				t.Contains(err.Issues, correctIssue)
			}
		}
	})

	t.Run("invalid format", func() {
		// mock
		mockPayloads := []service.AttendanceHistoryIn{
			{
				EmployeeId: t.validPayload.EmployeeId,
				StartDate:  "2024-03",
				EndDate:    t.validPayload.EndDate,
			},
			{
				EmployeeId: t.validPayload.EmployeeId,
				StartDate:  "31-01-2024",
				EndDate:    t.validPayload.EndDate,
			},
			{
				EmployeeId: t.validPayload.EmployeeId,
				StartDate:  "01 March 2024",
				EndDate:    t.validPayload.EndDate,
			},
			{
				EmployeeId: t.validPayload.EmployeeId,
				StartDate:  "01-31-2024",
				EndDate:    t.validPayload.EndDate,
			},
		}

		for _, mockPayload := range mockPayloads {
			// action
			err := t.service.ValidateRequest(mockPayload)

			// assert
			if t.NotNil(err) {
				var expectedError *primitive.RequestValidationError
				t.ErrorAs(err, &expectedError)
				t.Greater(len(err.Issues), 0)

				var correctIssue = primitive.RequestValidationIssue{
					Code:    primitive.RequestValidationCodeInvalidValue,
					Field:   "start_date",
					Message: "must be in the YYYY-MM-DD format",
				}
				t.Contains(err.Issues, correctIssue)
			}
		}
	})
}

func (t *AttendanceHistoryTestSuite) TestValidateRequest_EndDate() {
	t.Run("can't empty when start_date is not empty", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.EndDate = ""

		// action
		err := t.service.ValidateRequest(mockPayload)

		// assert
		if t.NotNil(err) {
			var expectedError *primitive.RequestValidationError
			t.ErrorAs(err, &expectedError)
			t.Greater(len(err.Issues), 0)

			var correctIssue = primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeRequired,
				Field:   "end_date",
				Message: "must not be empty if start_date is not empty",
			}
			t.Contains(err.Issues, correctIssue)
		}
	})

	t.Run("can empty when start_date is empty", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.StartDate = ""
		mockPayload.EndDate = ""

		// action
		err := t.service.ValidateRequest(mockPayload)

		// assert
		t.Nil(err)
	})

	t.Run("invalid date", func() {
		// mock
		mockPayloads := []service.AttendanceHistoryIn{
			{
				EmployeeId: t.validPayload.EmployeeId,
				StartDate:  t.validPayload.StartDate,
				EndDate:    "not date",
			},
			{
				EmployeeId: t.validPayload.EmployeeId,
				StartDate:  t.validPayload.StartDate,
				EndDate:    "2024-01-32",
			},
			{
				EmployeeId: t.validPayload.EmployeeId,
				StartDate:  t.validPayload.StartDate,
				EndDate:    "2023-02-29",
			},
		}

		for _, mockPayload := range mockPayloads {
			// action
			err := t.service.ValidateRequest(mockPayload)

			// asset
			if t.NotNil(err) {
				var expectedError *primitive.RequestValidationError
				t.ErrorAs(err, &expectedError)
				t.Greater(len(err.Issues), 0)

				var correctIssue = primitive.RequestValidationIssue{
					Code:    primitive.RequestValidationCodeInvalidValue,
					Field:   "end_date",
					Message: "must be in the YYYY-MM-DD format",
				}
				t.Contains(err.Issues, correctIssue)
			}
		}
	})

	t.Run("invalid format", func() {
		// mock
		mockPayloads := []service.AttendanceHistoryIn{
			{
				EmployeeId: t.validPayload.EmployeeId,
				StartDate:  t.validPayload.StartDate,
				EndDate:    "2024-03",
			},
			{
				EmployeeId: t.validPayload.EmployeeId,
				StartDate:  t.validPayload.StartDate,
				EndDate:    "31-01-2024",
			},
			{
				EmployeeId: t.validPayload.EmployeeId,
				StartDate:  t.validPayload.StartDate,
				EndDate:    "01 March 2024",
			},
			{
				EmployeeId: t.validPayload.EmployeeId,
				StartDate:  t.validPayload.StartDate,
				EndDate:    "01-31-2024",
			},
		}

		for _, mockPayload := range mockPayloads {
			// action
			err := t.service.ValidateRequest(mockPayload)

			// assert
			if t.NotNil(err) {
				var expectedError *primitive.RequestValidationError
				t.ErrorAs(err, &expectedError)
				t.Greater(len(err.Issues), 0)

				var correctIssue = primitive.RequestValidationIssue{
					Code:    primitive.RequestValidationCodeInvalidValue,
					Field:   "end_date",
					Message: "must be in the YYYY-MM-DD format",
				}
				t.Contains(err.Issues, correctIssue)
			}
		}
	})
}
