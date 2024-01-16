package service

import (
	"context"
	"errors"
	"hroost/domain/central/auth/db"
	"hroost/module/shared/primitive"
	"hroost/module/shared/validator"
	"net/http"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type ChangePasswordIn struct {
	Token    string
	UID      string `json:"uid"`
	Password string `json:"password"`
}

type ChangePasswordOut struct {
	primitive.CommonResult
}

// ValidateChangePasswordRequest validate the request body
func ValidateChangePasswordRequest(in ChangePasswordIn) *primitive.RequestValidationError {
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

func (s *Service) ChangePassword(ctx context.Context, in ChangePasswordIn) (out ChangePasswordOut) {
	// validate request body
	if err := ValidateChangePasswordRequest(in); err != nil {
		out.SetResponse(http.StatusBadRequest, "request validation failed", err)
		return
	}

	// check if the given toke has not expired yet
	existingToken, err := s.memory.GetPasswordRecoveryToken(ctx, in.UID)
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
	rowsAffected, err := s.db.ChangePassword(ctx, db.ChangePasswordIn{
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
	err = s.memory.DeletePasswordRecoveryToken(ctx, in.UID)

	if err != nil && !errors.Is(err, redis.Nil) {
		out.SetResponse(http.StatusInternalServerError, "internal server error", err)
		return
	}

	out.SetResponse(http.StatusOK, "password has been successfully changed")
	return
}
