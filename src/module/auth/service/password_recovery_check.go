package service

import (
	"context"
	"errors"
	"hris/module/shared/primitive"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type InternalPasswordRecoveryCheckIn struct {
	Token string
	UID   string `json:"uid"`
}

type InternalPasswordRecoveryCheckOut struct {
	primitive.CommonResult
}

func ValidateInternalPasswordRecoveryCheckIn(in InternalPasswordRecoveryCheckIn) *primitive.RequestValidationError {
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

func (s *InternalAuthService) PasswordRecoveryCheck(ctx context.Context, in InternalPasswordRecoveryCheckIn) (out InternalPasswordRecoveryCheckOut) {
	// validate the request
	if err := ValidateInternalPasswordRecoveryCheckIn(in); err != nil {
		out.SetResponse(http.StatusBadRequest, "request validation failed", err)
		return
	}

	// check if the given token has not expired yet
	existingToken, err := s.AuthRepo.RedisGetPasswordRecoveryToken(ctx, in.UID)
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
