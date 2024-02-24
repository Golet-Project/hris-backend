package model

import "hroost/shared/primitive"

type GetEmployeeDetailOut struct {
	Email          string
	FullName       string
	Gender         primitive.Gender
	BirthDate      primitive.Date
	ProfilePicture primitive.String
	Address        primitive.String
	JoinDate       primitive.Date
}
