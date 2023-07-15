package primitive

import (
	"fmt"
	"strings"
)

// RequestValidationCode provides a typed string for validation error codes.
type RequestValidationCode string

const (
	RequestValidationCodeRequired        RequestValidationCode = "field_required"
	RequestValidationCodeTooShort        RequestValidationCode = "too_short"
	RequestValidationCodeTooLong         RequestValidationCode = "too_long"
	RequestValidationCodeProhibitedValue RequestValidationCode = "prohibited_value"
	RequestValidationCodeInvalidValue    RequestValidationCode = "invalid_value"
)

// RequestValidationIssue contains a specific validation issue for each field and rules.
// It should be embedded as array inside the RequestValidationError struct.
type RequestValidationIssue struct {
	// Code specifies the error code. You must not enter a custom code, instead
	// add another entry for the RequestValidationCode type.
	//
	// This should be aligned with the documentation on how the consumers (or users)
	// handle their validation errors from us.
	Code RequestValidationCode `json:"code"`
	// Field specifies the field that the error happened. If the field is on a nested object,
	// you can separate it using a dot. For example:
	//
	// 		{ "customer": { "name": "string" } }
	// becomes
	// 		customer.name
	Field string `json:"field"`
	// Message must contain helpful message that helps the user create proper request.
	// This should also be simple, and should not repeat what's on the code and field value.
	// For example: "maximum of 50 characters", "must be numeric", "must not empty".
	Message string `json:"message"`
}

type RequestValidationError struct {
	Issues []RequestValidationIssue `json:"issue"`
}

func (r RequestValidationError) Error() string {
	var s strings.Builder
	for _, issue := range r.Issues {
		s.WriteString(fmt.Sprintf("%s for %s: %s", issue.Code, issue.Field, issue.Message))
	}

	return s.String()
}
