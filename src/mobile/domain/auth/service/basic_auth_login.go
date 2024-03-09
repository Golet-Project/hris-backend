package service

import (
	"context"
	"hroost/mobile/lib/jwt"
	"hroost/shared/primitive"
	"hroost/shared/utils"
	"net/http"

	"hroost/mobile/domain/auth/model"
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

type BasicAuthLoginBcrypt interface {
	CompareHashAndPassword(hashedPassword []byte, password []byte) error
}

type BasicAuthLoginDb interface {
	GetLoginCredential(ctx context.Context, email string) (credential model.GetLoginCredentialOut, err *primitive.RepoError)
	GetEmployeeDetail(ctx context.Context, domain string, userId string) (employee model.GetEmployeeDetailOut, err *primitive.RepoError)
}

type BasicAuthLogin struct {
	Db     BasicAuthLoginDb
	Bcrypt BasicAuthLoginBcrypt
}

func (s *BasicAuthLogin) Exec(ctx context.Context, body BasicAuthLoginIn) (out BasicAuthLoginOut) {
	// validate request body
	if err := s.ValidateBasicAuthLoginBody(body); err != nil {
		out.SetResponse(http.StatusBadRequest, "request validation failed", err)
		return
	}

	// get the credentials
	loginCredential, repoError := s.Db.GetLoginCredential(ctx, body.Email)
	if repoError != nil {
		switch repoError.Issue {
		case primitive.RepoErrorCodeDataNotFound:
			out.SetResponse(http.StatusNotFound, "user not found")
			return
		default:
			out.SetResponse(http.StatusInternalServerError, "internal server error", repoError)
			return
		}
	}

	// compare password
	if err := s.Bcrypt.CompareHashAndPassword([]byte(loginCredential.Password.String), []byte(body.Password)); err != nil {
		out.SetResponse(http.StatusUnauthorized, "invalid password")
		return
	}

	// if success, get the user detail at the domain
	userDetail, repoError := s.Db.GetEmployeeDetail(ctx, loginCredential.Domain, loginCredential.UserUID)
	if repoError != nil {
		switch repoError.Issue {
		case primitive.RepoErrorCodeDataNotFound:
			out.SetResponse(http.StatusNotFound, "employee not found")
			return
		default:
			out.SetResponse(http.StatusInternalServerError, "internal server error", repoError)
			return
		}
	}

	// TODO: testing
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

func (s *BasicAuthLogin) ValidateBasicAuthLoginBody(body BasicAuthLoginIn) *primitive.RequestValidationError {
	var allIssues []primitive.RequestValidationIssue

	// validate email
	if len(body.Email) == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeTooShort,
			Field:   "email",
			Message: "email is required",
		})
	} else if issues := utils.IsEmailValid(body.Email); len(issues) > 0 {
		allIssues = append(allIssues, issues...)
	}

	// validate password
	if len(body.Password) == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeTooShort,
			Field:   "password",
			Message: "password is required",
		})
	} else if issues := utils.IsPasswordValid(body.Password); len(issues) > 0 {
		allIssues = append(allIssues, issues...)
	}

	if len(allIssues) > 0 {
		return &primitive.RequestValidationError{
			Issues: allIssues,
		}
	}

	return nil
}
