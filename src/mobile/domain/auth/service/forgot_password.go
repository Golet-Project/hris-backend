package service

import (
	"context"
	"hroost/mobile/domain/auth/model"
	"hroost/shared/lib/mailer"
	"hroost/shared/primitive"
	"hroost/shared/utils"
	"net/http"
	"os"
)

type ForgotPasswordIn struct {
	Email string `json:"email"`
}

type ForgotPasswordOut struct {
	primitive.CommonResult
}

type ForgotPasswordDb interface {
	GetLoginCredential(ctx context.Context, email string) (loginCredential model.GetLoginCredentialOut, err *primitive.RepoError)
}

type ForgotPasswordMemory interface {
	GetPasswordRecoveryToken(ctx context.Context, userId string) (token string, err *primitive.RepoError)
	SetPasswordRecoveryToken(ctx context.Context, userId string, token string) (err *primitive.RepoError)
}

type ForgotPassword struct {
	Db     ForgotPasswordDb
	Memory ForgotPasswordMemory
}

func (s *ForgotPassword) Exec(ctx context.Context, req ForgotPasswordIn) (out ForgotPasswordOut) {
	// validate the request
	if err := s.ValidateForgotPasswordRequest(req); err != nil {
		out.SetResponse(http.StatusBadRequest, "request validation failed", err)
		return
	}

	// get the login data
	user, repoError := s.Db.GetLoginCredential(ctx, req.Email)
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

	// if password is not set (login using o auth) return fake response imediately
	if !user.Password.Valid || user.Password.String == "" {
		out.SetResponse(http.StatusOK, "password recovery link has been sent to your email")
		return
	}

	// check if token exists
	token, repoError := s.Memory.GetPasswordRecoveryToken(ctx, user.UserUID)
	if repoError != nil {
		switch repoError.Issue {
		case primitive.RepoErrorCodeDataNotFound:
			out.SetResponse(http.StatusNotFound, "data not found")
			return
		default:
			out.SetResponse(http.StatusInternalServerError, "internal server error", repoError)
			return
		}
	}
	if repoError == nil && token != "" {
		out.SetResponse(http.StatusOK, "password recovery link has been sent to your email")
		return
	}

	// make password recovery token
	token, err := utils.Base64String(48)
	if err != nil {
		out.SetResponse(http.StatusInternalServerError, "internal server error", err)
		return
	}

	// save the token to redis
	err = s.Memory.SetPasswordRecoveryToken(ctx, user.UserUID, token)
	if err != nil {
		out.SetResponse(http.StatusInternalServerError, "internal server error", err)
		return
	}

	// send password recovery link via email
	emailRecoveryLink := os.Getenv("INTERNAL_WEB_BASE_URL") + "/auth/password-recovery?token=" + token + "&uid=" + user.UserUID + "&cid=" + primitive.MobileAppID.String()

	// TODO: testing
	mailBuilder := mailer.NewMailer()
	mailBuilder.Subject("Email Recovery")
	mailBuilder.To([]string{
		user.Email,
	})
	err = mailBuilder.Send("Ini password recovery link, jangan diberikan ke orang lain ya\n\n" + emailRecoveryLink)
	if err != nil {
		out.SetResponse(http.StatusInternalServerError, "internal server error", err)
		return
	}

	out.SetResponse(http.StatusOK, "password recovery link has been sent to your email")
	return
}

func (s *ForgotPassword) ValidateForgotPasswordRequest(req ForgotPasswordIn) *primitive.RequestValidationError {
	var allIssues []primitive.RequestValidationIssue

	// validate email
	if len(req.Email) == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeTooShort,
			Field:   "email",
			Message: "email is required",
		})
	} else if issues := utils.IsEmailValid(req.Email); len(issues) > 0 {
		allIssues = append(allIssues, issues...)
	}

	if len(allIssues) > 0 {
		return &primitive.RequestValidationError{
			Issues: allIssues,
		}
	}

	return nil
}
