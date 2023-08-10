package service

import (
	"context"
	"errors"
	"hris/module/shared/mailer"
	"hris/module/shared/primitive"
	"hris/module/shared/random"
	"hris/module/shared/validator"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type InternalForgotPasswordIn struct {
	Email string `json:"email"`
	AppID primitive.AppID
}

type InternalForgotPasswordOut struct {
	primitive.CommonResult
}

func ValidateInternalForgotPasswordPayload(body InternalForgotPasswordIn) *primitive.RequestValidationError {
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

	if len(allIssues) > 0 {
		return &primitive.RequestValidationError{
			Issues: allIssues,
		}
	}

	// validate app id
	if len(body.AppID) == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeTooShort,
			Field:   "X-App-ID",
			Message: "x-app-id header is required",
		})
	} else {
		if body.AppID != primitive.InternalAppID && body.AppID != primitive.MobileAppID && body.AppID != primitive.WebAppID {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeProhibitedValue,
				Field:   "x-app-id",
				Message: "x-app-id header is invalid",
			})
		}
	}

	return nil
}

func (s *InternalAuthService) ForgotPassword(ctx context.Context, in InternalForgotPasswordIn) (out InternalForgotPasswordOut) {
	// validate request body
	if err := ValidateInternalForgotPasswordPayload(in); err != nil {
		out.SetResponse(http.StatusBadRequest, "request validation failed", err)
		return
	}

	// check email
	admin, err := s.AuthRepo.InternalGetLoginCredential(ctx, in.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			out.SetResponse(http.StatusNotFound, "user not found")
			return
		} else {
			out.SetResponse(http.StatusInternalServerError, "internal server error", err)
			return
		}
	}

	// if password is not set (login using o auth) return fake response immediately
	if !admin.Password.Valid || admin.Password.String == "" {
		out.SetResponse(http.StatusOK, "password recovery link has been sent to your email")
		return
	}

	// check token is exists
	existingToken, err := s.AuthRepo.RedisGetPasswordRecoveryToken(ctx, admin.UserUID)
	if err != nil && !errors.Is(err, redis.Nil) {
		out.SetResponse(http.StatusInternalServerError, "internal server error", err)
		return
	}
	if err == nil && existingToken != "" {
		out.SetResponse(http.StatusBadRequest, "password recovery link has already been sent to your email")
		return
	}

	// make password revocery token
	token, err := random.Base64String(48)
	if err != nil {
		out.SetResponse(http.StatusInternalServerError, "internal server error", err)
		return
	}

	// store into redis with expire time
	err = s.AuthRepo.RedisSetPasswordRecoveryToken(ctx, admin.UserUID, token)

	// send password recovery link via email
	emailRecoveryLink := os.Getenv("INTERNAL_WEB_BASE_URL") + "/password-recovery?token=" + token + "&uid=" + admin.UserUID + "&cid=" + in.AppID.String()

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
