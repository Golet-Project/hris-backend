package validator

import (
	"hroost/module/shared/primitive"
	"regexp"
)

func IsPasswordValid(pass string) []primitive.RequestValidationIssue {
	var issues []primitive.RequestValidationIssue

	if len(pass) < 8 {
		issues = append(issues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeTooShort,
			Field:   "password",
			Message: "password must be at least 8 characters long",
		})
	}

	lower, upper, num := 0, 0, 0
	for _, ch := range pass {
		if ch >= 'a' && ch <= 'z' {
			lower++
		} else if ch >= 'A' && ch <= 'Z' {
			upper++
		} else if ch >= '0' && ch <= '9' {
			num++
		}
	}

	if lower == 0 {
		issues = append(issues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeInvalidValue,
			Field:   "password",
			Message: "must contain at least one lowercase character",
		})
	}

	if upper == 0 {
		issues = append(issues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeInvalidValue,
			Field:   "password",
			Message: "must contain at least one uppercase character",
		})
	}

	if num == 0 {
		issues = append(issues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeInvalidValue,
			Field:   "password",
			Message: "must contain at least one number",
		})
	}

	return issues
}

var emailPattern, _ = regexp.Compile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`)

func IsEmailValid(email string) []primitive.RequestValidationIssue {
	var issues []primitive.RequestValidationIssue

	if !emailPattern.MatchString(email) {
		issues = append(issues, primitive.RequestValidationIssue{
			Code:    primitive.RequestValidationCodeInvalidValue,
			Field:   "email",
			Message: "invalid email address",
		})
	}

	return issues
}
