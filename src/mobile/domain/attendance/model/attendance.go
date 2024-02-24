package model

import (
	"hroost/shared/entities"
	"hroost/shared/primitive"
)

type AddAttendanceIn struct {
	EmployeeUID string
	Timezone    primitive.Timezone
	Coordinate  primitive.Coordinate
}

type GetTodayAttendanceIn struct {
	EmployeeUID string
	Timezone    primitive.Timezone
}

type GetTodayAttendanceOut struct {
	Timezone         primitive.Timezone
	AttendanceRadius primitive.Int64
	CheckinTime      primitive.Time
	CheckoutTime     primitive.Time
	ApprovedAt       primitive.Time

	StartWorkingHour primitive.Time
	EndWorkingHour   primitive.Time

	Company entities.Company
}
