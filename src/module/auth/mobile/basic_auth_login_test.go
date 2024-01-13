package mobile_test

import (
	"errors"
	"hroost/module/auth/mobile"
	"hroost/module/shared/primitive"
	"reflect"
	"testing"
)

func TestValidateBasicAuthLoginBody(t *testing.T) {
	type Test struct {
		name    string
		payload mobile.BasicAuthLoginIn
		want    *primitive.RequestValidationError
	}
	positiveTest := []Test{
		{
			name: "payload is valid",
			payload: mobile.BasicAuthLoginIn{
				Email:    "email@email.com",
				Password: "Password321",
			},
			want: nil,
		},
		{
			name: "payload is valid",
			payload: mobile.BasicAuthLoginIn{
				Email:    "email@email.com",
				Password: "Password@321",
			},
			want: nil,
		},
	}

	for _, test := range positiveTest {
		t.Run(test.name, func(t *testing.T) {
			got := mobile.ValidateBasicAuthLoginBody(test.payload)
			if got != test.want {
				t.Errorf("ValidateBasicAuthLoginBody() = %v, want %v", got, test.want)
			}
		})
	}

	validPayload := mobile.BasicAuthLoginIn{
		Email:    "email@email.com",
		Password: "Password321",
	}
	negativeTest := []Test{
		{
			name: "email is required",
			payload: mobile.BasicAuthLoginIn{
				Password: validPayload.Password,
			},
			want: &primitive.RequestValidationError{
				Issues: []primitive.RequestValidationIssue{
					{
						Code:    primitive.RequestValidationCodeTooShort,
						Field:   "email",
						Message: "email is required",
					},
				},
			},
		},
		{
			name: "email is invalid",
			payload: mobile.BasicAuthLoginIn{
				Email:    "email",
				Password: validPayload.Password,
			},
			want: &primitive.RequestValidationError{
				Issues: []primitive.RequestValidationIssue{
					{
						Code:    primitive.RequestValidationCodeInvalidValue,
						Field:   "email",
						Message: "invalid email address",
					},
				},
			},
		},
		{
			name: "password is required",
			payload: mobile.BasicAuthLoginIn{
				Email: validPayload.Email,
			},
			want: &primitive.RequestValidationError{
				Issues: []primitive.RequestValidationIssue{
					{
						Code:    primitive.RequestValidationCodeTooShort,
						Field:   "password",
						Message: "password is required",
					},
				},
			},
		},
		{
			name: "password is too short",
			payload: mobile.BasicAuthLoginIn{
				Email:    validPayload.Email,
				Password: "pass",
			},
			want: &primitive.RequestValidationError{
				Issues: []primitive.RequestValidationIssue{
					{
						Code:    primitive.RequestValidationCodeTooShort,
						Field:   "password",
						Message: "password must be at least 8 characters long",
					},
					{
						Code:    primitive.RequestValidationCodeInvalidValue,
						Field:   "password",
						Message: "must contain at least one uppercase character",
					},
					{
						Code:    primitive.RequestValidationCodeInvalidValue,
						Field:   "password",
						Message: "must contain at least one number",
					},
				},
			},
		},
		{
			name: "password contains no lowercase characters",
			payload: mobile.BasicAuthLoginIn{
				Email:    validPayload.Email,
				Password: "PASSWORD123",
			},
			want: &primitive.RequestValidationError{
				Issues: []primitive.RequestValidationIssue{
					{
						Code:    primitive.RequestValidationCodeInvalidValue,
						Field:   "password",
						Message: "must contain at least one lowercase character",
					},
				},
			},
		},
		{
			name: "password contains no uppsercase characters",
			payload: mobile.BasicAuthLoginIn{
				Email:    validPayload.Email,
				Password: "password321",
			},
			want: &primitive.RequestValidationError{
				Issues: []primitive.RequestValidationIssue{
					{
						Code:    primitive.RequestValidationCodeInvalidValue,
						Field:   "password",
						Message: "must contain at least one uppercase character",
					},
				},
			},
		},
		{
			name: "password contains no number",
			payload: mobile.BasicAuthLoginIn{
				Email:    validPayload.Email,
				Password: "Passwords",
			},
			want: &primitive.RequestValidationError{
				Issues: []primitive.RequestValidationIssue{
					{
						Code:    primitive.RequestValidationCodeInvalidValue,
						Field:   "password",
						Message: "must contain at least one number",
					},
				},
			},
		},
	}

	for _, test := range negativeTest {
		t.Run(test.name, func(t *testing.T) {
			var expectedErrorType *primitive.RequestValidationError
			got := mobile.ValidateBasicAuthLoginBody(test.payload)
			if !errors.As(got, &expectedErrorType) {
				t.Errorf("ValidateBasicAuthLoginBody()\ngot: %T\nwant: %T", got, test.want)
			}

			if len(got.Issues) != len(test.want.Issues) {
				t.Errorf("Expect %d issues, got %d issues", len(test.want.Issues), len(got.Issues))
			} else {
				for i := range got.Issues {
					if !reflect.DeepEqual(got.Issues[i], test.want.Issues[i]) {
						t.Errorf("ValidateBasicAuthLoginBody() = %v, want: %v", got, test.want)
					}
				}
			}
		})
	}
}
