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
	AppID primitive.AppID
}

type ForgotPasswordOut struct {
	primitive.CommonResult
}

// ValidateForgotPasswordPayload validate the request body
func ValidateForgotPasswordPayload(body ForgotPasswordIn) *primitive.RequestValidationError {
	var allIssues []primitive.RequestValidationIssue

	// validate email
	if len(body.Email) == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeTooShort,
			Field:   "email",
			Message: "email is required",
		})
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
		if body.AppID != primitive.CentralAppID && body.AppID != primitive.MobileAppID && body.AppID != primitive.TenantAppID {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeProhibitedValue,
				Field:   "x-app-id",
				Message: "x-app-id header is invalid",
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

// ForgotPassword send password recovery link to user email
func (s *Service) ForgotPassword(ctx context.Context, in ForgotPasswordIn) (out ForgotPasswordOut) {
	// validate request body
	if err := ValidateForgotPasswordPayload(in); err != nil {
		out.SetResponse(http.StatusBadRequest, "request validation failed", err)
		return
	}

	// check email
	admin, err := s.db.GetLoginCredential(ctx, in.Email)
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
	existingToken, err := s.memory.GetPasswordRecoveryToken(ctx, admin.UserUID)
	if err != nil && !errors.Is(err, redis.Nil) {
		out.SetResponse(http.StatusInternalServerError, "internal server error", err)
		return
	}
	if err == nil && existingToken != "" {
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
	err = s.memory.SetPasswordRecoveryToken(ctx, admin.UserUID, token)
	if err != nil {
		out.SetResponse(http.StatusInternalServerError, "internal server error", err)
		return
	}

	// send password recovery link via email
	emailRecoveryLink := os.Getenv("INTERNAL_WEB_BASE_URL") + "/auth/password-recovery?token=" + token + "&uid=" + admin.UserUID + "&cid=" + in.AppID.String()

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
