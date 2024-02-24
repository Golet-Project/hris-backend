package service

import (
	"context"
	"errors"
	"hroost/mobile/domain/auth/model"
	"hroost/shared/primitive"
	"hroost/shared/utils"
	"net/http"

	"github.com/google/uuid"
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

type ChangePasswordDb interface {
	ChangePassword(ctx context.Context, data model.ChangePasswordIn) (rowsAffected int64, err *primitive.RepoError)
}

type ChangePasswordMemory interface {
	GetPasswordRecoveryToken(ctx context.Context, userId string) (token string, err *primitive.RepoError)
	DeletePasswordRecoveryToken(ctx context.Context, userId string) (err *primitive.RepoError)
}

type ChangePassword struct {
	Db     ChangePasswordDb
	Memory ChangePasswordMemory
}

func (s *ChangePassword) Exec(ctx context.Context, in ChangePasswordIn) (out ChangePasswordOut) {
	// validate request body
	if err := s.ValidateChangePasswordRequest(in); err != nil {
		out.SetResponse(http.StatusBadRequest, "request validation failed", err)
		return
	}

	// check if the given token has not expired yet
	existingToken, repoError := s.Memory.GetPasswordRecoveryToken(ctx, in.UID)
	if repoError != nil {
		switch repoError.Issue {
		case primitive.RepoErrorCodeDataNotFound:
			out.SetResponse(http.StatusBadRequest, "password recovery token has expired")
			return
		default:
			out.SetResponse(http.StatusInternalServerError, "internal server error", repoError)
			return
		}
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
	rowsAffected, err := s.Db.ChangePassword(ctx, model.ChangePasswordIn{
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
	err = s.Memory.DeletePasswordRecoveryToken(ctx, in.UID)
	if err != nil && !errors.Is(err, redis.Nil) {
		out.SetResponse(http.StatusInternalServerError, "internal server error", err)
		return
	}

	out.SetResponse(http.StatusOK, "password has been successfully changed")

	return
}

func (s *ChangePassword) ValidateChangePasswordRequest(in ChangePasswordIn) *primitive.RequestValidationError {
	var allIssues []primitive.RequestValidationIssue

	// validate token
	if len(in.Token) == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeTooShort,
			Field:   "token",
			Message: "token is required",
		})
	}

	// vaildate password
	if len(in.Password) == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeTooShort,
			Field:   "password",
			Message: "password is required",
		})
	} else {
		if issues := utils.IsPasswordValid(in.Password); len(issues) > 0 {
			allIssues = append(allIssues, issues...)
		}
	}

	// validate uid
	if len(in.UID) == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeTooShort,
			Field:   "uid",
			Message: "uid is required",
		})
	} else {
		_, err := uuid.Parse(in.UID)
		if err != nil {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "uid",
				Message: "uid is not valid",
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
