package model

import "hroost/shared/primitive"

type GetEmployeeDetailOut struct {
	UID            string
	FullName       string
	Email          string
	Gender         primitive.Gender
	BirthDate      primitive.Date
	ProfilePicture primitive.String
	Address        primitive.String
	JoinDate       primitive.Date
}
