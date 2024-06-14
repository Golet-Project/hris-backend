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
	Id             string
	Email          string
	FullName       string
	BirthDate      time.Time
	Gender         primitive.Gender
	EmployeeStatus primitive.EmployeeStatus
	JoinDate       time.Time
	EndDate        primitive.Date
}

type GetByIdOut struct {
	Address struct {
		Detail       primitive.String
		ProvinceId   primitive.String
		ProvinceName primitive.String
		RegencyId    primitive.String
		RegencyName  primitive.String
		DistrictId   primitive.String
		DistrictName primitive.String
		VillageId    primitive.String
		VillageName  primitive.String
	}
	Id             string
	Email          string
	FullName       string
	Gender         string
	EmployeeStatus primitive.EmployeeStatus
	BirthDate      primitive.Date
	JoinDate       primitive.Date
}
