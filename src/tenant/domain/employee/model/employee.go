package model

import (
	"hroost/shared/primitive"
	"time"
)

type CreateEmployeeIn struct {
	Domain         string
	Email          string
	Password       string
	FirstName      string
	LastName       string
	Gender         primitive.Gender
	BirthDate      string
	Address        string
	ProvinceId     string
	RegencyId      string
	DistrictId     string
	VillageId      string
	JoinDate       string
	EmployeeStatus primitive.EmployeeStatus
}

type FindAllEmployeeOut struct {
	UID            string
	Email          string
	FullName       string
	BirthDate      time.Time
	Gender         primitive.Gender
	EmployeeStatus primitive.EmployeeStatus
	JoinDate       time.Time
	EndDate        primitive.Date
}
