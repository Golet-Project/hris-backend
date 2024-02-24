package model

import "hroost/shared/primitive"

type FindHomePageIn struct {
	UID      string
	Timezone primitive.Timezone
}

type FindHomePageOut struct {
	TodayAttendance
}

type TodayAttendance struct {
	Timezone     primitive.Timezone
	CheckinTime  primitive.Time
	CheckoutTime primitive.Time
	ApprovedAt   primitive.Time
}
