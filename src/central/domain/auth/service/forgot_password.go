package service

import (
	"context"
	"hroost/central/domain/auth/model"
	"hroost/shared/lib/mailer"
	"hroost/shared/primitive"
	"hroost/shared/utils"
	"net/http"
	"os"
)

type ForgotPasswordIn struct {
	Email string `json:"email"`
	AppID primitive.AppID
}

type ForgotPasswordOut struct {
	primitive.CommonResult
}

type ForgotPasswordDb interface {
	GetLoginCredential(ctx context.Context, email string) (credential model.GetLoginCredentialOut, err *primitive.RepoError)
}

type ForgotPasswordMemory interface {
	GetPasswordRecoveryToken(ctx context.Context, userId string) (token string, err *primitive.RepoError)
	SetPasswordRecoveryToken(ctx context.Context, userId string, token string) (err *primitive.RepoError)
}

type ForgotPasswordMailer interface {
	Subject(string)
}

type ForgotPassword struct {
	Db     ForgotPasswordDb
	Memory ForgotPasswordMemory
}

// ForgotPassword send password recovery link to user email
func (s *ForgotPassword) Exec(ctx context.Context, in ForgotPasswordIn) (out ForgotPasswordOut) {
	// validate request body
	if err := s.ValidateForgotPasswordPayload(in); err != nil {
		out.SetResponse(http.StatusBadRequest, "request validation failed", err)
		return
	}

	// check email
	admin, repoError := s.Db.GetLoginCredential(ctx, in.Email)
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

	// if password is not set (login using o auth) return fake response immediately
	if !admin.Password.Valid || admin.Password.String == "" {
		out.SetResponse(http.StatusOK, "password recovery link has already been sent to your email")
		return
	}

	// check token is exists
	existingToken, repoError := s.Memory.GetPasswordRecoveryToken(ctx, admin.UserUID)
	if repoError != nil {
		if repoError.Issue != primitive.RepoErrorCodeDataNotFound {
			out.SetResponse(http.StatusInternalServerError, "internal server error", repoError)
			return
		}
	}
	if existingToken != "" {
		out.SetResponse(http.StatusBadRequest, "password recovery link has already been sent to your email")
		return
	}

	// make password revocery token
	token, err := utils.Base64String(48)
	if err != nil {
		out.SetResponse(http.StatusInternalServerError, "internal server error", err)
		return
	}

	// store into redis with expire time
	repoError = s.Memory.SetPasswordRecoveryToken(ctx, admin.UserUID, token)
	if repoError != nil {
		out.SetResponse(http.StatusInternalServerError, "internal server error", repoError)
		return
	}

	// TODO: test
	// send password recovery link via email
	emailRecoveryLink := os.Getenv("INTERNAL_WEB_BASE_URL") + "/auth/password-recovery?token=" + token + "&uid=" + admin.UserUID + "&cid=" + in.AppID.String()

	// TODO: test
	mailBuilder := mailer.NewMailer()
	mailBuilder.Subject("Email Recovery")
	mailBuilder.To([]string{
		admin.Email,
	})
	err = mailBuilder.Send("Ini password recovery, jangan diberikan ke orang lain ya\n\n" + emailRecoveryLink)
	if err != nil {
		out.SetResponse(http.StatusInternalServerError, "internal server error", err)
		return
	}

	out.SetResponse(http.StatusOK, "password recovery link has been sent to your email")
	return
}

func (s *ForgotPassword) ValidateForgotPasswordPayload(body ForgotPasswordIn) *primitive.RequestValidationError {
	var allIssues []primitive.RequestValidationIssue

	// validate email
	if len(body.Email) == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeRequired,
			Field:   "email",
			Message: "email is required",
		})
	} else {
		emailIssues := utils.IsEmailValid(body.Email)
		if len(emailIssues) > 0 {
			allIssues = append(allIssues, emailIssues...)
		}
	}

	// validate app id
	if len(body.AppID) == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeRequired,
			Field:   "X-App-ID",
			Message: "X-App-ID header is required",
		})
	} else {
		if body.AppID != primitive.CentralAppID {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeProhibitedValue,
				Field:   "X-App-ID",
				Message: "X-App-ID header has a prohibited value",
			})
		}
	}

	if len(allIssues) > 0 {
		return &primitive.RequestValidationError{
			Issues: allIssues,
		}
	}

	return nil
}
