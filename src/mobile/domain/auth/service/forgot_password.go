package service

import (
	"context"
	"errors"
	"hroost/shared/lib/mailer"
	"hroost/shared/primitive"
	"hroost/shared/utils"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type ForgotPasswordIn struct {
	Email string `json:"email"`
}

type ForgotPasswordOut struct {
	primitive.CommonResult
}

func ValidateForgotPasswordRequest(req ForgotPasswordIn) *primitive.RequestValidationError {
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

func (s *Service) ForgotPassword(ctx context.Context, req ForgotPasswordIn) (out ForgotPasswordOut) {
	// validate the request
	if err := ValidateForgotPasswordRequest(req); err != nil {
		out.SetResponse(http.StatusBadRequest, "request validation failed", err)
		return
	}

	// get the login data
	user, err := s.db.GetLoginCredential(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			out.SetResponse(http.StatusNotFound, "user not found")
			return
		} else {
			out.SetResponse(http.StatusInternalServerError, "internal server error", err)
			return
		}
	}

	// if password is not set (login using o auth) return fake response imediately
	if !user.Password.Valid || user.Password.String == "" {
		out.SetResponse(http.StatusOK, "password recovery link has been sent to your email")
		return
	}

	// check if token exists
	token, err := s.memory.GetPasswordRecoveryToken(ctx, user.UserUID)
	if err != nil && !errors.Is(err, redis.Nil) {
		out.SetResponse(http.StatusInternalServerError, "internal server error", err)
	}
	if err == nil && token != "" {
		out.SetResponse(http.StatusOK, "password recovery link has been sent to your email")
		return
	}

	// make password recovery token
	token, err = utils.Base64String(48)
	if err != nil {
		out.SetResponse(http.StatusInternalServerError, "internal server error", err)
		return
	}

	// save the token to redis
	err = s.memory.SetPasswordRecoveryToken(ctx, user.UserUID, token)
	if err != nil {
		out.SetResponse(http.StatusInternalServerError, "internal server error", err)
		return
	}

	// send password recovery link via email
	emailRecoveryLink := os.Getenv("INTERNAL_WEB_BASE_URL") + "/auth/password-recovery?token=" + token + "&uid=" + user.UserUID + "&cid=" + primitive.MobileAppID.String()

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
