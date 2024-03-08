package service_test

import (
	"context"
	"fmt"
	"hroost/mobile/domain/attendance/model"
	"hroost/mobile/domain/attendance/service"
	"hroost/shared/primitive"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockAddAttendanceDb struct {
	mock.Mock
}

func (m *MockAddAttendanceDb) GetDomainByUid(ctx context.Context, uid string) (domain string, repoErr *primitive.RepoError) {
	ret := m.Called(ctx, uid)

	if err := ret.Get(1); err != nil {
		repoErr = err.(*primitive.RepoError)
	}

	return ret.String(0), repoErr
}

func (m *MockAddAttendanceDb) EmployeeExistsById(ctx context.Context, domain string, uid string) (exists bool, repoErr *primitive.RepoError) {
	ret := m.Called(ctx, domain, uid)

	if err := ret.Get(1); err != nil {
		repoErr = err.(*primitive.RepoError)
	}

	return ret.Bool(0), repoErr
}

func (m *MockAddAttendanceDb) CheckTodayAttendanceById(ctx context.Context, domain string, uid string, timezone primitive.Timezone) (exists bool, repoErr *primitive.RepoError) {
	ret := m.Called(ctx, domain, uid, timezone)

	if err := ret.Get(1); err != nil {
		repoErr = err.(*primitive.RepoError)
	}

	return ret.Bool(0), repoErr
}

func (m *MockAddAttendanceDb) AddAttendance(ctx context.Context, domain string, data model.AddAttendanceIn) (repoErr *primitive.RepoError) {
	ret := m.Called(ctx, domain, data)

	if err := ret.Get(0); err != nil {
		repoErr = err.(*primitive.RepoError)
	}

	return repoErr
}

type AddAttendanceTestSuite struct {
	suite.Suite

	db *MockAddAttendanceDb

	validPayload service.AddAttendanceIn

	service service.AddAttendance
}

func (t *AddAttendanceTestSuite) SetupSubTest() {
	db := new(MockAddAttendanceDb)

	t.db = db

	t.validPayload = service.AddAttendanceIn{
		UID:      "a84c2c59-748c-48d0-b628-4a73b1c3a8d7",
		Timezone: primitive.WIB,
		Coordinate: primitive.Coordinate{
			Latitude:  -7.782870711329031,
			Longitude: 110.36707035197001,
		},
	}

	t.service = service.AddAttendance{
		Db: db,
	}
}

func TestAddAttendanceTestSuite(t *testing.T) {
	suite.Run(t, new(AddAttendanceTestSuite))
}

func (t *AddAttendanceTestSuite) TestExec_InvalidPayload() {
	t.Run("should return 400", func() {
		// arrange
		ctx := context.Background()

		// mock
		mockPayload := t.validPayload
		mockPayload.Coordinate.Latitude = -91

		// action
		out := t.service.Exec(ctx, mockPayload)

		// assert
		t.Assert().Equal(http.StatusBadRequest, out.GetCode())
		t.Assert().Equal("request validation failed", out.GetMessage())
		t.db.AssertNumberOfCalls(t.T(), "GetDomainByUid", 0)
		t.db.AssertNumberOfCalls(t.T(), "AddAttendance", 0)
		t.db.AssertNumberOfCalls(t.T(), "EmployeeExistsBYId", 0)
		t.db.AssertNumberOfCalls(t.T(), "CheckTodayAttendanceById", 0)
	})
}

func (t *AddAttendanceTestSuite) TestExec_GetDomain() {
	t.Run("should return 404 when user not found", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.db.On("GetDomainByUid", ctx, t.validPayload.UID).Return("", &primitive.RepoError{
			Issue: primitive.RepoErrorCodeDataNotFound,
			Err:   fmt.Errorf("mock not found"),
		})

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Assert().Equal(http.StatusNotFound, out.GetCode())
		t.Assert().Equal("user not found", out.GetMessage())
		t.db.AssertNumberOfCalls(t.T(), "AddAttendance", 0)
	})

	t.Run("should return 500 when server error", func() {
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
		t.db.AssertNumberOfCalls(t.T(), "AddAttendance", 0)
	})
}

func (t *AddAttendanceTestSuite) TestExec_CheckEmployeeExists() {
	t.Run("not found", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.db.On("GetDomainByUid", ctx, t.validPayload.UID).Return("mock_tenant", nil)
		t.db.On("EmployeeExistsById", ctx, "mock_tenant", t.validPayload.UID).Return(false, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeDataNotFound,
			Err:   fmt.Errorf("mock not found"),
		})

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Assert().Equal(http.StatusNotFound, out.GetCode())
		t.Assert().Equal("employee not found", out.GetMessage())
		t.db.AssertExpectations(t.T())
		t.db.AssertNumberOfCalls(t.T(), "GetDomainByUid", 1)
		t.db.AssertNumberOfCalls(t.T(), "AddAttendance", 0)
	})

	t.Run("server error", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.db.On("GetDomainByUid", ctx, t.validPayload.UID).Return("mock_tenant", nil)
		t.db.On("EmployeeExistsById", ctx, "mock_tenant", t.validPayload.UID).Return(false, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   fmt.Errorf("mock server error"),
		})

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Assert().Equal(http.StatusInternalServerError, out.GetCode())
		t.Assert().Equal("internal server error", out.GetMessage())
		t.db.AssertExpectations(t.T())
		t.db.AssertNumberOfCalls(t.T(), "GetDomainByUid", 1)
		t.db.AssertNumberOfCalls(t.T(), "AddAttendance", 0)
	})

	t.Run("employee doesn't exists", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.db.On("GetDomainByUid", ctx, t.validPayload.UID).Return("mock_tenant", nil)
		t.db.On("EmployeeExistsById", ctx, "mock_tenant", t.validPayload.UID).Return(false, nil)

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Assert().Equal(http.StatusNotFound, out.GetCode())
		t.Assert().Equal("employee not found", out.GetMessage())
		t.db.AssertExpectations(t.T())
		t.db.AssertNumberOfCalls(t.T(), "GetDomainByUid", 1)
		t.db.AssertNumberOfCalls(t.T(), "AddAttendance", 0)
	})
}

func (t *AddAttendanceTestSuite) TestExec_CheckTodayAttendance() {
	t.Run("server error", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.db.On("GetDomainByUid", ctx, t.validPayload.UID).Return("mock_tenant", nil)
		t.db.On("EmployeeExistsById", ctx, "mock_tenant", t.validPayload.UID).Return(true, nil)
		t.db.On("CheckTodayAttendanceById", ctx, "mock_tenant", t.validPayload.UID, t.validPayload.Timezone).Return(false, &primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   fmt.Errorf("mock server error"),
		})

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Assert().Equal(http.StatusInternalServerError, out.GetCode())
		t.Assert().Equal("internal server error", out.GetMessage())
		t.db.AssertExpectations(t.T())
		t.db.AssertNumberOfCalls(t.T(), "CheckTodayAttendanceById", 1)
		t.db.AssertNumberOfCalls(t.T(), "AddAttendance", 0)
	})

	t.Run("attendance already exists", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.db.On("GetDomainByUid", ctx, t.validPayload.UID).Return("mock_tenant", nil)
		t.db.On("EmployeeExistsById", ctx, "mock_tenant", t.validPayload.UID).Return(true, nil)
		t.db.On("CheckTodayAttendanceById", ctx, "mock_tenant", t.validPayload.UID, t.validPayload.Timezone).Return(true, nil)

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Assert().Equal(http.StatusConflict, out.GetCode())
		t.Assert().Equal("attendance already exist", out.GetMessage())
		t.db.AssertExpectations(t.T())
		t.db.AssertNumberOfCalls(t.T(), "CheckTodayAttendanceById", 1)
		t.db.AssertNumberOfCalls(t.T(), "AddAttendance", 0)
	})
}

func (t *AddAttendanceTestSuite) TestExec_AddAttendance() {
	t.Run("server error", func() {
		// arrange
		ctx := context.Background()

		// mock
		t.db.On("GetDomainByUid", ctx, t.validPayload.UID).Return("mock_tenant", nil)
		t.db.On("EmployeeExistsById", ctx, "mock_tenant", t.validPayload.UID).Return(true, nil)
		t.db.On("CheckTodayAttendanceById", ctx, "mock_tenant", t.validPayload.UID, t.validPayload.Timezone).Return(false, nil)
		t.db.On("AddAttendance", ctx, "mock_tenant", model.AddAttendanceIn{
			EmployeeUID: t.validPayload.UID,
			Timezone:    t.validPayload.Timezone,
			Coordinate:  t.validPayload.Coordinate,
		}).Return(&primitive.RepoError{
			Issue: primitive.RepoErrorCodeServerError,
			Err:   fmt.Errorf("mock server error"),
		})

		// action
		out := t.service.Exec(ctx, t.validPayload)

		// assert
		t.Assert().Equal(http.StatusInternalServerError, out.GetCode())
		t.Assert().Equal("internal server error", out.GetMessage())
		t.db.AssertExpectations(t.T())
		t.db.AssertNumberOfCalls(t.T(), "AddAttendance", 1)
	})
	t.Run("can add attendance", func() {
		// arrange
		validPayloads := []service.AddAttendanceIn{
			{
				UID:        t.validPayload.UID,
				Timezone:   t.validPayload.Timezone,
				Coordinate: t.validPayload.Coordinate,
			},
			{
				UID:      t.validPayload.UID,
				Timezone: t.validPayload.Timezone,
				Coordinate: primitive.Coordinate{
					Latitude:  90,
					Longitude: 180,
				},
			},
			{
				UID:      t.validPayload.UID,
				Timezone: t.validPayload.Timezone,
				Coordinate: primitive.Coordinate{
					Latitude:  -90,
					Longitude: 180,
				},
			},
			{
				UID:      t.validPayload.UID,
				Timezone: t.validPayload.Timezone,
				Coordinate: primitive.Coordinate{
					Latitude:  90,
					Longitude: -180,
				},
			},
			{
				UID:      t.validPayload.UID,
				Timezone: t.validPayload.Timezone,
				Coordinate: primitive.Coordinate{
					Latitude:  -90,
					Longitude: -180,
				},
			},
			{
				UID:      t.validPayload.UID,
				Timezone: primitive.WIB,
				Coordinate: primitive.Coordinate{
					Latitude:  -90,
					Longitude: 180,
				},
			},
			{
				UID:      t.validPayload.UID,
				Timezone: primitive.WIT,
				Coordinate: primitive.Coordinate{
					Latitude:  -90,
					Longitude: 180,
				},
			},
			{
				UID:      t.validPayload.UID,
				Timezone: primitive.WITA,
				Coordinate: primitive.Coordinate{
					Latitude:  -90,
					Longitude: 180,
				},
			},
		}
		ctx := context.Background()

		for _, payload := range validPayloads {
			// mock
			t.db.On("GetDomainByUid", ctx, payload.UID).Return("mock_tenant", nil)
			t.db.On("EmployeeExistsById", ctx, "mock_tenant", payload.UID).Return(true, nil)
			t.db.On("CheckTodayAttendanceById", ctx, "mock_tenant", payload.UID, payload.Timezone).Return(false, nil)
			t.db.On("AddAttendance", ctx, "mock_tenant", model.AddAttendanceIn{
				EmployeeUID: payload.UID,
				Timezone:    payload.Timezone,
				Coordinate:  payload.Coordinate,
			}).Return(nil)

			// action
			out := t.service.Exec(ctx, payload)

			// assert
			t.Assert().Equal(http.StatusCreated, out.GetCode())
			t.Assert().Equal("success", out.GetMessage())
			t.db.AssertExpectations(t.T())
		}
	})
}

func (t *AddAttendanceTestSuite) TestValidateAddAttendanceRequest_ValidPayload() {
	t.Run("should return no error", func() {
		validPayloads := []service.AddAttendanceIn{
			{
				UID:        t.validPayload.UID,
				Timezone:   t.validPayload.Timezone,
				Coordinate: t.validPayload.Coordinate,
			},
			{
				UID:      t.validPayload.UID,
				Timezone: t.validPayload.Timezone,
				Coordinate: primitive.Coordinate{
					Latitude:  90,
					Longitude: 180,
				},
			},
			{
				UID:      t.validPayload.UID,
				Timezone: t.validPayload.Timezone,
				Coordinate: primitive.Coordinate{
					Latitude:  -90,
					Longitude: 180,
				},
			},
			{
				UID:      t.validPayload.UID,
				Timezone: t.validPayload.Timezone,
				Coordinate: primitive.Coordinate{
					Latitude:  90,
					Longitude: -180,
				},
			},
			{
				UID:      t.validPayload.UID,
				Timezone: t.validPayload.Timezone,
				Coordinate: primitive.Coordinate{
					Latitude:  -90,
					Longitude: -180,
				},
			},
			{
				UID:      t.validPayload.UID,
				Timezone: primitive.WIB,
				Coordinate: primitive.Coordinate{
					Latitude:  -90,
					Longitude: 180,
				},
			},
			{
				UID:      t.validPayload.UID,
				Timezone: primitive.WIT,
				Coordinate: primitive.Coordinate{
					Latitude:  -90,
					Longitude: 180,
				},
			},
			{
				UID:      t.validPayload.UID,
				Timezone: primitive.WITA,
				Coordinate: primitive.Coordinate{
					Latitude:  -90,
					Longitude: 180,
				},
			},
		}

		// action
		for _, payload := range validPayloads {
			err := t.service.ValidateAddAttendanceRequest(payload)

			// assert
			t.Assert().Nil(err)
		}
	})
}

func (t *AddAttendanceTestSuite) TestValidateAddAttendanceRequest_UID() {
	t.Run("uid required", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.UID = ""

		// action
		err := t.service.ValidateAddAttendanceRequest(mockPayload)

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

	t.Run("uid is not a uuid", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.UID = "not a uuid"

		// action
		err := t.service.ValidateAddAttendanceRequest(mockPayload)

		// assert
		if t.Assert().NotNil(err) {
			var expectedError *primitive.RequestValidationError
			t.Assert().ErrorAs(err, &expectedError)
			t.Assert().Greater(len(err.Issues), 0)

			expectedIssue := primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "uid",
				Message: "uid is not a valid UUID",
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

	t.Run("uid is not a valid UUIDV4", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.UID = "115a0ee6-da39-11ee-a506-0242ac120002"

		// action
		err := t.service.ValidateAddAttendanceRequest(mockPayload)

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

func (t *AddAttendanceTestSuite) TestValidateAddAttendanceRequest_Timeozne() {
	t.Run("timezone required", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.Timezone = 0

		// action
		err := t.service.ValidateAddAttendanceRequest(mockPayload)

		// assert
		if t.Assert().NotNil(err) {
			var expectedError *primitive.RequestValidationError
			t.Assert().ErrorAs(err, &expectedError)
			t.Assert().Greater(len(err.Issues), 0)

			var correctIssue = primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeRequired,
				Field:   "timezone",
				Message: "timezone header is required",
			}
			var containCorrectIssue = false
			for _, issue := range err.Issues {
				if assert.ObjectsAreEqual(correctIssue, issue) {
					containCorrectIssue = true
					break
				}
			}
			t.Assert().True(containCorrectIssue)
		}
	})

	t.Run("timezone invalid", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.Timezone = 10

		// action
		err := t.service.ValidateAddAttendanceRequest(mockPayload)

		// assert
		if t.Assert().NotNil(err) {
			var expectedError *primitive.RequestValidationError
			t.Assert().ErrorAs(err, &expectedError)
			t.Assert().Greater(len(err.Issues), 0)

			var correctIssue = primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "timezone",
				Message: "timezone header has an invalid value",
			}
			var containCorrectIssue = false
			for _, issue := range err.Issues {
				if assert.ObjectsAreEqual(correctIssue, issue) {
					containCorrectIssue = true
					break
				}
			}
			t.Assert().True(containCorrectIssue)
		}
	})
}

func (t *AddAttendanceTestSuite) TestValidateAddAttendanceRequest_Coordinate() {
	t.Run("coordinate required", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.Coordinate = primitive.Coordinate{}

		// action
		err := t.service.ValidateAddAttendanceRequest(mockPayload)

		// assert
		if t.Assert().NotNil(err) {
			var expectedError *primitive.RequestValidationError
			t.Assert().ErrorAs(err, &expectedError)
			t.Assert().Greater(len(err.Issues), 0)

			var correctIssue = primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeRequired,
				Field:   "coordinate",
				Message: "must not empty",
			}
			var containCorrectIssue = false
			for _, issue := range err.Issues {
				if assert.ObjectsAreEqual(correctIssue, issue) {
					containCorrectIssue = true
					break
				}
			}
			t.Assert().True(containCorrectIssue)
		}
	})
}

func (t *AddAttendanceTestSuite) TestValidateAddAttendanceRequest_Coordinate_Latitude() {
	t.Run("must be greater than or equal -90", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.Coordinate.Latitude = -91

		// action
		err := t.service.ValidateAddAttendanceRequest(mockPayload)

		// assert
		if t.Assert().NotNil(err) {
			var expectedError *primitive.RequestValidationError
			t.Assert().ErrorAs(err, &expectedError)
			t.Assert().Greater(len(err.Issues), 0)

			var correctIssue = primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "coordinate.latitude",
				Message: "must be between -90 and 90",
			}
			var containCorrectIssue = false
			for _, issue := range err.Issues {
				if assert.ObjectsAreEqual(correctIssue, issue) {
					containCorrectIssue = true
					break
				}
			}
			t.Assert().True(containCorrectIssue)
		}
	})

	t.Run("must be less than or equal 90", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.Coordinate.Latitude = 91

		// action
		err := t.service.ValidateAddAttendanceRequest(mockPayload)

		// assert
		if t.Assert().NotNil(err) {
			var expectedError *primitive.RequestValidationError
			t.Assert().ErrorAs(err, &expectedError)
			t.Assert().Greater(len(err.Issues), 0)

			var correctIssue = primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "coordinate.latitude",
				Message: "must be between -90 and 90",
			}
			var containCorrectIssue = false
			for _, issue := range err.Issues {
				if assert.ObjectsAreEqual(correctIssue, issue) {
					containCorrectIssue = true
					break
				}
			}
			t.Assert().True(containCorrectIssue)
		}
	})
}

func (t *AddAttendanceTestSuite) TestValidateAddAttendanceRequest_Coordinate_Longitude() {
	t.Run("greater than or equal -180", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.Coordinate.Longitude = -181

		// action
		err := t.service.ValidateAddAttendanceRequest(mockPayload)

		// assert
		if t.Assert().NotNil(err) {
			var expectedError *primitive.RequestValidationError
			t.Assert().ErrorAs(err, &expectedError)
			t.Assert().Greater(len(err.Issues), 0)

			var correctIssue = primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "coordinate.longitude",
				Message: "must be between -180 and 180",
			}
			var containCorrectIssue = false
			for _, issue := range err.Issues {
				if assert.ObjectsAreEqual(correctIssue, issue) {
					containCorrectIssue = true
					break
				}
			}
			t.Assert().True(containCorrectIssue)
		}
	})

	t.Run("less than or equal 180", func() {
		// mock
		mockPayload := t.validPayload
		mockPayload.Coordinate.Longitude = 181

		// action
		err := t.service.ValidateAddAttendanceRequest(mockPayload)

		// assert
		if t.Assert().NotNil(err) {
			var expectedError *primitive.RequestValidationError
			t.Assert().ErrorAs(err, &expectedError)
			t.Assert().Greater(len(err.Issues), 0)

			var correctIssue = primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "coordinate.longitude",
				Message: "must be between -180 and 180",
			}
			var containCorrectIssue = false
			for _, issue := range err.Issues {
				if assert.ObjectsAreEqual(correctIssue, issue) {
					containCorrectIssue = true
					break
				}
			}
			t.Assert().True(containCorrectIssue)
		}
	})
}
