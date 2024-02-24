package model

import "hroost/shared/primitive"

type GetLoginCredentialOut struct {
	UserUID  string
	Email    string
	Password primitive.String
	Domain   string
}
