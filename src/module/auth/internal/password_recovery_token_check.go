package internal

import (
	"context"
	"errors"
	"hris/module/shared/primitive"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type PasswordRecoveryTokenCheckIn struct {
	Token string
	UID   string `json:"uid"`
}

type PasswordRecoveryTokenCheckOut struct {
	primitive.CommonResult
}

// ValidatePasswordRecoveryTokenCheckIn validate the request body
func ValidatePasswordRecoveryTokenCheckIn(in PasswordRecoveryTokenCheckIn) *primitive.RequestValidationError {
	var allIssues []primitive.RequestValidationIssue

	// validate token
	if in.Token == "" {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeTooShort,
			Field:   "token",
			Message: "token is required",
		})
	}

	// validate uid
	if in.UID == "" {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeTooShort,
			Field:   "uid",
			Message: "uid is required",
		})
	}

	if len(allIssues) > 0 {
		return &primitive.RequestValidationError{
			Issues: allIssues,
		}
	}

	return nil
}

// PasswordRecoveryTokenCheck check if the user has the valid password recovery token
func (i *Internal) PasswordRecoveryTokenCheck(ctx context.Context, in PasswordRecoveryTokenCheckIn) (out PasswordRecoveryTokenCheckOut) {
	// validate the request
	if err := ValidatePasswordRecoveryTokenCheckIn(in); err != nil {
		out.SetResponse(http.StatusBadRequest, "request validation failed", err)
		return
	}

	// check if the given token has not expired yet
	existingToken, err := i.redis.GetPasswordRecoveryToken(ctx, in.UID)
	if err != nil && !errors.Is(err, redis.Nil) {
		out.SetResponse(http.StatusInternalServerError, "internal server error", err)
		return
	}
	if errors.Is(err, redis.Nil) || existingToken == "" {
		out.SetResponse(http.StatusBadRequest, "password recovery token has expired")
		return
	}

	// check if the given token is correct
	if existingToken != in.Token {
		out.SetResponse(http.StatusBadRequest, "password recovery token has expired")
		return
	}

	out.SetResponse(http.StatusNoContent, "OK")
	return
}
