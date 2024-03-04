package service

import (
	"context"
	"hroost/shared/primitive"
	"net/http"

	"github.com/google/uuid"
)

type PasswordRecoveryTokenCheckIn struct {
	Token string
	UID   string `json:"uid"`
}

type PasswordRecoveryTokenCheckOut struct {
	primitive.CommonResult
}

type PasswordRecoveryTokenCheckMemory interface {
	GetPasswordRecoveryToken(ctx context.Context, userId string) (token string, err *primitive.RepoError)
}

type PasswordRecoveryTokenCheck struct {
	Memory PasswordRecoveryTokenCheckMemory
}

// PasswordRecoveryTokenCheck check if the user has the valid password recovery token
func (s *PasswordRecoveryTokenCheck) Exec(ctx context.Context, in PasswordRecoveryTokenCheckIn) (out PasswordRecoveryTokenCheckOut) {
	// validate the request
	if err := s.ValidatePasswordRecoveryTokenCheckIn(in); err != nil {
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

	out.SetResponse(http.StatusNoContent, "OK")
	return
}

// ValidatePasswordRecoveryTokenCheckIn validate the request body
func (s *PasswordRecoveryTokenCheck) ValidatePasswordRecoveryTokenCheckIn(in PasswordRecoveryTokenCheckIn) *primitive.RequestValidationError {
	var allIssues []primitive.RequestValidationIssue

	// validate token
	if in.Token == "" {
		allIssues = append(allIssues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeRequired,
			Field:   "token",
			Message: "token is required",
		})
	}

	// validate uid
	if in.UID == "" {
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
				Message: "uid is not a valid uuid",
			})
		}

		if parsed.Version() != 4 {
			allIssues = append(allIssues, primitive.RequestValidationIssue{
				Code:    primitive.RequestValidationCodeInvalidValue,
				Field:   "uid",
				Message: "uid is not a valid UUIDV4",
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
