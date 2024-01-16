package service

import (
	"context"
	"errors"
	"hroost/module/shared/jwt"
	"hroost/module/shared/primitive"
	"hroost/module/shared/validator"
	"net/http"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type BasicAuthLoginIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type employee struct {
	Email          string           `json:"email"`
	FullName       string           `json:"full_name"`
	Gender         primitive.Gender `json:"gender"`
	BirthDate      primitive.Date   `json:"birth_date"`
	ProfilePicture primitive.String `json:"profile_picture"`
	Address        primitive.String `json:"address"`
	JoinDate       primitive.Date   `json:"join_date"`
}

type BasicAuthLoginOut struct {
	primitive.CommonResult

	AccessToken string   `json:"access_token"`
	Employee    employee `json:"employee"`
}

func ValidateBasicAuthLoginBody(body BasicAuthLoginIn) *primitive.RequestValidationError {
	var allIssues []primitive.RequestValidationIssue

	// validate email
	if len(body.Email) == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeTooShort,
			Field:   "email",
			Message: "email is required",
		})
	} else if issues := validator.IsEmailValid(body.Email); len(issues) > 0 {
		allIssues = append(allIssues, issues...)
	}

	// validate password
	if len(body.Password) == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeTooShort,
			Field:   "password",
			Message: "password is required",
		})
	} else if issues := validator.IsPasswordValid(body.Password); len(issues) > 0 {
		allIssues = append(allIssues, issues...)
	}

	if len(allIssues) > 0 {
		return &primitive.RequestValidationError{
			Issues: allIssues,
		}
	}

	return nil
}

func (s *Service) BasicAuthLogin(ctx context.Context, body BasicAuthLoginIn) (out BasicAuthLoginOut) {
	// validate request body
	if err := ValidateBasicAuthLoginBody(body); err != nil {
		out.SetResponse(http.StatusBadRequest, "request validation failed", err)
		return
	}

	// get the credentials
	loginCredential, err := s.db.GetLoginCredential(ctx, body.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			out.SetResponse(http.StatusNotFound, "user not found")
			return
		} else {
			out.SetResponse(http.StatusInternalServerError, "internal server error", err)
			return
		}
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(loginCredential.Password.String), []byte(body.Password)); err != nil {
		out.SetResponse(http.StatusUnauthorized, "invalid password")
		return
	}

	// if success, get the user detail at the domain
	userDetail, err := s.db.GetEmployeeDetail(ctx, loginCredential.Domain, loginCredential.UserUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			out.SetResponse(http.StatusNotFound, "employee not found")
			return
		} else {
			out.SetResponse(http.StatusInternalServerError, "internal server error", err)
			return
		}
	}

	// generate jwt token
	out.AccessToken = jwt.GenerateAccessToken(loginCredential.UserUID)
	out.Employee = employee{
		Email:          userDetail.Email,
		FullName:       userDetail.FullName,
		Gender:         userDetail.Gender,
		BirthDate:      userDetail.BirthDate,
		ProfilePicture: userDetail.ProfilePicture,
		Address:        userDetail.Address,
		JoinDate:       userDetail.JoinDate,
	}

	out.SetResponse(http.StatusOK, "login success")
	return
}
