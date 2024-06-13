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
		Detail       string
		ProvinceId   string
		ProvinceName string
		RegencyId    string
		RegencyName  string
		DistrictId   string
		DistrictName string
		VillageId    string
		VillageName  string
	}
	Id             string
	Email          string
	FullName       string
	Gender         string
	EmployeeStatus primitive.EmployeeStatus
	BirthDate      primitive.Date
	JoinDate       primitive.Date
}
