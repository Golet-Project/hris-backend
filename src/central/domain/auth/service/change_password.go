package service

import (
	"context"
	"hroost/central/domain/auth/model"
	"hroost/shared/primitive"
	"hroost/shared/utils"
	"net/http"

	"github.com/google/uuid"
)

type ChangePasswordIn struct {
	Token    string
	UID      string `json:"uid"`
	Password string `json:"password"`
}

type ChangePasswordOut struct {
	primitive.CommonResult
}

type ChangePasswordMemory interface {
	GetPasswordRecoveryToken(ctx context.Context, userId string) (token string, err *primitive.RepoError)
	DeletePasswordRecoveryToken(ctx context.Context, userId string) (err *primitive.RepoError)
}

type ChangePasswordDb interface {
	ChangePassword(ctx context.Context, param model.ChangePasswordIn) (rowsAffected int64, err *primitive.RepoError)
}

type ChangePasswordBcrypt interface {
	GenerateFromPassword([]byte, int) ([]byte, error)
}

type ChangePassword struct {
	Memory               ChangePasswordMemory
	Db                   ChangePasswordDb
	GenerateFromPassword func(password []byte, cost int) (hash []byte, err error)
}

func (s *ChangePassword) Exec(ctx context.Context, in ChangePasswordIn) (out ChangePasswordOut) {
	// validate request body
	if err := s.ValidateChangePasswordRequest(in); err != nil {
		out.SetResponse(http.StatusBadRequest, "request validation failed", err)
		return
	}

	// check if the given toke has not expired yet
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
	passByte, err := s.GenerateFromPassword([]byte(in.Password), 10)
	if err != nil {
		out.SetResponse(http.StatusInternalServerError, "internal server error", err)
		return
	}

	rowsAffected, repoError := s.Db.ChangePassword(ctx, model.ChangePasswordIn{
		UID:      in.UID,
		Password: string(passByte),
	})
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
	if rowsAffected == 0 {
		out.SetResponse(http.StatusNotFound, "user not found")
		return
	}

	// delete delete password recovery link
	repoError = s.Memory.DeletePasswordRecoveryToken(ctx, in.UID)
	if repoError != nil {
		if primitive.RepoErrorCodeDataNotFound != repoError.Issue {
			out.SetResponse(http.StatusInternalServerError, "internal server error", repoError)
			return
		}
	}

	out.SetResponse(http.StatusOK, "password has been successfully changed")
	return
}

// ValidateChangePasswordRequest validate the request body
func (s *ChangePassword) ValidateChangePasswordRequest(in ChangePasswordIn) *primitive.RequestValidationError {
	var allIssues []primitive.RequestValidationIssue

	// validate token
	if len(in.Token) == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeRequired,
			Field:   "token",
			Message: "token is required",
		})
	}

	// validate uid
	if len(in.UID) == 0 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeRequired,
			Field:   "uid",
			Message: "uid is required",
		})
	} else {
		parsed, err := uuid.Parse(in.UID)
		if err != nil {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "uid",
				Message: "uid is invalid",
			})
		}
		if parsed.Version() != 4 {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "uid",
				Message: "uid must be a valid uuid v4",
			})
		}
	}

	// validate password
	if in.Password == "" {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeRequired,
			Field:   "password",
			Message: "password is required",
		})
	} else if issues := utils.IsPasswordValid(in.Password); len(issues) > 0 {
		allIssues = append(allIssues, issues...)
	}

	if len(in.Password) > 25 {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeTooLong,
			Field:   "password",
			Message: "password must be less than 25 characters",
		})
	}

	if len(allIssues) > 0 {
		return &primitive.RequestValidationError{
			Issues: allIssues,
		}
	}

	return nil
}
