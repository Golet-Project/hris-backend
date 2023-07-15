package service

import (
	"context"
	"errors"
	"hris/module/auth/repo/auth"
	"hris/module/shared/primitive"
	"hris/module/shared/validator"
	"net/http"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type InternalChangePasswordIn struct {
	Token    string
	UID      string `json:"uid"`
	Password string `json:"password"`
}

type InternalChangePasswordOut struct {
	primitive.CommonResult
}

func ValidateInternalChangePasswordRequest(in InternalChangePasswordIn) *primitive.RequestValidationError {
	var allIssues []primitive.RequestValidationIssue

	// validate token
	if len(in.Token) == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeTooShort,
			Field:   "token",
			Message: "token is required",
		})
	}

	// validate password
	if in.Password == "" {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeTooShort,
			Field:   "password",
			Message: "password is required",
		})
	} else if issues := validator.IsPasswordValid(in.Password); len(issues) > 0 {
		allIssues = append(allIssues, issues...)
	}

	if len(allIssues) > 0 {
		return &primitive.RequestValidationError{
			Issues: allIssues,
		}
	}

	return nil
}

func (s *InternalAuthService) ChangePassword(ctx context.Context, in InternalChangePasswordIn) (out InternalChangePasswordOut) {
	// validate request body
	if err := ValidateInternalChangePasswordRequest(in); err != nil {
		out.SetResponse(http.StatusBadRequest, "request validation failed", err)
		return
	}

	// check if the given toke has not expired yet
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

	// hash & change password
	passByte, err := bcrypt.GenerateFromPassword([]byte(in.Password), 10)
	if err != nil {
		out.SetResponse(http.StatusInternalServerError, "internal server error", err)
		return
	}
	rowsAffected, err := s.AuthRepo.InternalChangePassword(ctx, auth.InternalChangePasswordIn{
		UID:      in.UID,
		Password: string(passByte),
	})
	if err != nil {
		out.SetResponse(http.StatusInternalServerError, "internal server error", err)
		return
	}
	if rowsAffected == 0 {
		out.SetResponse(http.StatusNotFound, "user not found")
		return
	}

	// delete delete password recovery link
	err = s.AuthRepo.RedisDeletePasswordRecoveryToken(ctx, in.UID)
	if err != nil && !errors.Is(err, redis.Nil) {
		out.SetResponse(http.StatusInternalServerError, "internal server error", err)
		return
	}

	out.SetResponse(http.StatusOK, "password has been successfully changed")
	return
}
