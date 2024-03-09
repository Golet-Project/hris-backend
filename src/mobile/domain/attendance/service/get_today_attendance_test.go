package service_test

import (
	"context"
	"encoding/json"
	"fmt"
	"hroost/mobile/domain/attendance/model"
	"hroost/mobile/domain/attendance/service"
	"hroost/shared/entities"
	"hroost/shared/primitive"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockGetTodayAttendanceDb struct {
	mock.Mock
}

func (m *MockGetTodayAttendanceDb) GetDomainByUid(ctx context.Context, uid string) (domain string, err *primitive.RepoError) {
	ret := m.Called(ctx, uid)

	var r1 *primitive.RepoError
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(*primitive.RepoError)
	}

	return ret.Get(0).(string), r1
}

func (m *MockAddAttendanceDb) GetTodayAttendance(ctx context.Context, domain string, in model.GetTodayAttendanceIn) (out model.GetTodayAttendanceOut, err *primitive.RepoError) {
	ret := m.Called(ctx, domain, in)

	var r1 *primitive.RepoError
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(*primitive.RepoError)
	}

	return ret.Get(0).(model.GetTodayAttendanceOut), r1
}

type GetTodayAttendanceTestSuite struct {
	suite.Suite

	db *MockAddAttendanceDb

	now          time.Time
	validPayload service.GetTodayAttendanceIn

	service service.GetTodayAttendance
}

func (t *GetTodayAttendanceTestSuite) SetupSubTest() {
	db := new(MockAddAttendanceDb)

	t.db = db
	t.validPayload = service.GetTodayAttendanceIn{
		EmployeeUID: "a84c2c59-748c-48d0-b628-4a73b1c3a8d7",
		Timezone:    primitive.WIB,
	}
	t.now = time.Date(2024, 3, 9, 12, 0, 0, 0, time.UTC)

	t.service = service.GetTodayAttendance{
		Db: db,
		GetNow: func() time.Time {
			return t.now
		},
	}
}

func TestGetTodayAttendanceTestSuite(t *testing.T) {
	suite.Run(t, new(GetTodayAttendanceTestSuite))
}

func (t *GetTodayAttendanceTestSuite) TestExec_InvalidPayload() {
	t.Run("should return error 400", func() {
		// arrange
		ctx := context.Background()

		// mock
		mockPayload := t.validPayload
		mockPayload.EmployeeUID = ""

		// action
		out := t.service.Exec(ctx, mockPayload)

		// assert
		t.Equal(http.StatusBadRequest, out.GetCode())
		t.Equal("request validation failed", out.GetMessage())
		t.db.AssertNumberOfCalls(t.T(), "GetDomainByUid", 0)
		t.db.AssertNumberOfCalls(t.T(), "GetTodayAttendance", 0)
	})
}

func (t *GetTodayAttendanceTestSuite) TestExec_GetDomainByUid() {
	t.Run("server error", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.db.On("GetDomainByUid", ctx, t.validPayload.EmployeeUID).Return("", &primitive.RepoError{
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
	})

	t.Run("employee not found", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.db.On("GetDomainByUid", ctx, t.validPayload.EmployeeUID).Return("", &primitive.RepoError{
			Issue: primitive.RepoErrorCodeDataNotFound,
			Err:   fmt.Errorf("mock employee not found"),
		})

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Equal(http.StatusNotFound, out.GetCode())
		t.Equal("employee not found", out.GetMessage())
		t.db.AssertExpectations(t.T())
		t.db.AssertNumberOfCalls(t.T(), "GetDomainByUid", 1)
	})
}

func (t *GetTodayAttendanceTestSuite) TestExec_GetTodayAttendance() {
	t.Run("server error", func() {
		// arrange
		ctx := context.Background()
		domain := "mock-domain"

		// mock
		t.db.On("GetDomainByUid", ctx, t.validPayload.EmployeeUID).Return(domain, nil)
		t.db.On("GetTodayAttendance", ctx, domain, model.GetTodayAttendanceIn{
			EmployeeUID: t.validPayload.EmployeeUID,
			Timezone:    t.validPayload.Timezone,
		}).Return(model.GetTodayAttendanceOut{}, &primitive.RepoError{
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
		t.db.AssertNumberOfCalls(t.T(), "GetTodayAttendance", 1)
	})

	t.Run("should not return error when data not found", func() {
		// arrange
		ctx := context.Background()
		domain := "mock-domain"

		// mock
		t.db.On("GetDomainByUid", ctx, t.validPayload.EmployeeUID).Return(domain, nil)
		t.db.On("GetTodayAttendance", ctx, domain, model.GetTodayAttendanceIn{
			EmployeeUID: t.validPayload.EmployeeUID,
			Timezone:    t.validPayload.Timezone,
		}).Return(model.GetTodayAttendanceOut{}, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeDataNotFound,
			Err:   fmt.Errorf("mock data not found"),
		})

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Equal(http.StatusOK, out.GetCode())
		t.Equal("success", out.GetMessage())
		t.db.AssertExpectations(t.T())
		t.db.AssertNumberOfCalls(t.T(), "GetDomainByUid", 1)
		t.db.AssertNumberOfCalls(t.T(), "GetTodayAttendance", 1)
	})

	t.Run("can get today attendance and has a correct response", func() {
		// arrange
		ctx := context.Background()
		domain := "mock-domain"

		// mock
		t.db.On("GetDomainByUid", ctx, t.validPayload.EmployeeUID).Return(domain, nil)
		t.db.On("GetTodayAttendance", ctx, domain, model.GetTodayAttendanceIn{
			EmployeeUID: t.validPayload.EmployeeUID,
			Timezone:    t.validPayload.Timezone,
		}).Return(model.GetTodayAttendanceOut{
			Timezone:         primitive.WIB,
			AttendanceRadius: primitive.Int64{Int64: 100, Valid: true},

			CheckinTime: primitive.Time{
				Time:  time.Date(2024, 3, 9, 9, 0, 0, 0, time.UTC),
				Valid: true,
			},
			CheckoutTime: primitive.Time{
				Time:  time.Date(2024, 3, 9, 17, 0, 0, 0, time.UTC),
				Valid: true,
			},
			ApprovedAt: primitive.Time{
				Time:  time.Time{},
				Valid: false,
			},

			StartWorkingHour: primitive.Time{
				Time:  time.Date(0, 0, 0, 9, 0, 0, 0, time.UTC),
				Valid: true,
			},
			EndWorkingHour: primitive.Time{
				Time:  time.Date(0, 0, 0, 17, 0, 0, 0, time.UTC),
				Valid: true,
			},

			Company: entities.Company{
				Coordinate: primitive.Coordinate{
					Latitude:  -7.782870711329031,
					Longitude: 110.36707035197001,
				},
				Address: primitive.String{
					String: "Jl. Kaliurang KM 5,3, Sleman, Yogyakarta",
					Valid:  true,
				},
			},
		}, nil)

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		var s = struct {
			Timezone         primitive.Timezone `json:"timezone"`
			CurrentTime      string             `json:"current_time"`
			CheckinTime      string             `json:"checkin_time"`
			CheckoutTime     string             `json:"checkout_time"`
			ApprovedAt       string             `json:"approved_at"`
			StartWorkingHour string             `json:"start_working_hour"`
			EndWorkingHour   string             `json:"end_working_hour"`
			AttendanceRadius int64              `json:"attendance_radius"` // radius in meters

			Company struct {
				Coordinate primitive.Coordinate `json:"coordinate,omitempty"`
				Address    primitive.String     `json:"address,omitempty"`
			} `json:"company"`
		}{
			Timezone:         primitive.WIB,
			CurrentTime:      t.now.Format(primitive.UtcRFC3339),
			CheckinTime:      time.Date(2024, 3, 9, 9, 0, 0, 0, time.UTC).Format(primitive.UtcRFC3339),
			CheckoutTime:     time.Date(2024, 3, 9, 17, 0, 0, 0, time.UTC).Format(primitive.UtcRFC3339),
			ApprovedAt:       "",
			StartWorkingHour: time.Date(0, 0, 0, 9, 0, 0, 0, time.UTC).Format("15:04"),
			EndWorkingHour:   time.Date(0, 0, 0, 17, 0, 0, 0, time.UTC).Format("15:04"),
			AttendanceRadius: 100,
			Company: entities.Company{
				Coordinate: primitive.Coordinate{
					Latitude:  -7.782870711329031,
					Longitude: 110.36707035197001,
				},
				Address: primitive.String{
					String: "Jl. Kaliurang KM 5,3, Sleman, Yogyakarta",
					Valid:  true,
				},
			},
		}

		expectedJson, _ := json.Marshal(s)

		t.Equal(http.StatusOK, out.GetCode())
		t.Equal("success", out.GetMessage())

		b, err := json.Marshal(out)
		if err != nil {
			t.Fail(err.Error())
		}
		t.Assert().JSONEq(string(expectedJson), string(b))
	})
}

func (t *GetTodayAttendanceTestSuite) TestValidateGetTodayAttendanceIn_EmployeeUID() {
	t.Run("EmployeeUID is empty", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.EmployeeUID = ""

		// action
		err := t.service.ValidateGetTodayAttendanceIn(mockPayload)

		// assert
		if t.NotNil(err) {
			var expectedError *primitive.RequestValidationError
			t.ErrorAs(err, &expectedError)
			t.Greater(len(err.Issues), 0)

			var expectedIssue = primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeRequired,
				Field:   "employee_uid",
				Message: "must not be empty",
			}
			t.Contains(err.Issues, expectedIssue)
		}
	})

	t.Run("EmployeeUID is not a valid UUID", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.EmployeeUID = "not-a-uuid"

		// action
		err := t.service.ValidateGetTodayAttendanceIn(mockPayload)

		// assert
		if t.NotNil(err) {
			var expectedError *primitive.RequestValidationError
			t.ErrorAs(err, &expectedError)
			t.Greater(len(err.Issues), 0)

			var expectedIssue = primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "employee_uid",
				Message: "must be a valid UUID",
			}
			t.Contains(err.Issues, expectedIssue)
		}
	})

	t.Run("EmployeeUID is not a valid UUIV4", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.EmployeeUID = "115a0ee6-da39-11ee-a506-0242ac120002"

		// action
		err := t.service.ValidateGetTodayAttendanceIn(mockPayload)

		// assert
		if t.NotNil(err) {
			var expectedError *primitive.RequestValidationError
			t.ErrorAs(err, &expectedError)
			t.Greater(len(err.Issues), 0)

			var expectedIssue = primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "employee_uid",
				Message: "must be a valid UUIDV4",
			}
			t.Contains(err.Issues, expectedIssue)
		}
	})
}

func (t *GetTodayAttendanceTestSuite) TestValidateGetTodayAttendanceIn_Timeozone() {
	t.Run("Timezone is empty", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.Timezone = 0

		// action
		err := t.service.ValidateGetTodayAttendanceIn(mockPayload)

		// assert
		if t.NotNil(err) {
			var expectedError *primitive.RequestValidationError
			t.ErrorAs(err, &expectedError)
			t.Greater(len(err.Issues), 0)

			var expectedIssue = primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeRequired,
				Field:   "timezone",
				Message: "must not be empty",
			}
			t.Contains(err.Issues, expectedIssue)
		}
	})

	t.Run("invalid timezone", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.Timezone = 10

		// action
		err := t.service.ValidateGetTodayAttendanceIn(mockPayload)

		// assert
		if t.NotNil(err) {
			var expectedError *primitive.RequestValidationError
			t.ErrorAs(err, &expectedError)
			t.Greater(len(err.Issues), 0)

			var expectedIssue = primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "timezone",
				Message: "must be a valid timezone header",
			}
			t.Contains(err.Issues, expectedIssue)
		}
	})
}
