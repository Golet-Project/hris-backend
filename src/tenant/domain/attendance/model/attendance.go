package model

import "hroost/shared/primitive"

type FindAllAttendanceOut struct {
	UID          string
	FullName     string
	CheckinTime  primitive.Time
	CheckoutTime primitive.Time
	ApprovedAt   primitive.Time
	ApprovedBy   primitive.String
}
