package service_test

import (
	"errors"
	authService "hroost/central/domain/auth/service"
	"hroost/shared/primitive"
	"testing"
)

func TestValidateBasicAuthLoginBody(t *testing.T) {
	service := authService.BasicAuthLogin{}

	t.Run("payload is valid (1)", func(t *testing.T) {
		payload := authService.BasicAuthLoginIn{
			Email:    "email@email.com",
			Password: "Password@321",
		}

		err := service.ValidateBasicAuthLoginBody(payload)

		if err != nil {
			t.Errorf("Expect: nil\nGot: %T", err)
		}
	})

	validPayload := authService.BasicAuthLoginIn{
		Email:    "email@email.com",
		Password: "Password@321",
	}
	t.Run("email is missing", func(t *testing.T) {
		mock := validPayload
		mock.Email = ""
		var requestValidationError *primitive.RequestValidationError

		err := service.ValidateBasicAuthLoginBody(mock)

		if err == nil {
			t.Errorf("Expect: error\nGot: nil")
		}

		if !errors.As(err, &requestValidationError) {
			t.Errorf("Expect: %T\nGot: %T", err, requestValidationError)
		}
	})

	t.Run("email is invalid", func(t *testing.T) {
		mock := validPayload
		mock.Email = "email"
		var requestValidationError *primitive.RequestValidationError

		got := service.ValidateBasicAuthLoginBody(mock)

		if got == nil {
			t.Errorf("Expect: error\nGot: nil")
		}

		if !errors.As(got, &requestValidationError) {
			t.Errorf("Expect: %T\nGot: %T", got, requestValidationError)
		}
	})

	t.Run("password is required", func(t *testing.T) {
		mock := validPayload
		mock.Password = ""
		var requestValidationError *primitive.RequestValidationError

		got := service.ValidateBasicAuthLoginBody(mock)

		if got == nil {
			t.Errorf("Expect: error\nGot: nil")
		}

		if !errors.As(got, &requestValidationError) {
			t.Errorf("Expect: %T\nGot: %T", got, requestValidationError)
		}

	})

	t.Run("password is too short", func(t *testing.T) {
		mock := validPayload
		mock.Password = "pass"
		var requestValidationError *primitive.RequestValidationError

		got := service.ValidateBasicAuthLoginBody(mock)

		if got == nil {
			t.Errorf("Expect: error\nGot: nil")
		}

		if !errors.As(got, &requestValidationError) {
			t.Errorf("Expect: %T\nGot: %T", got, requestValidationError)
		}
	})

	t.Run("password contains no lowercase characters", func(t *testing.T) {
		mock := validPayload
		mock.Password = "PASSWORD123"
		var requestValidationError *primitive.RequestValidationError

		got := service.ValidateBasicAuthLoginBody(mock)

		if got == nil {
			t.Errorf("Expect: error\nGot: nil")
		}

		if !errors.As(got, &requestValidationError) {
			t.Errorf("Expect: %T\nGot: %T", got, requestValidationError)
		}
	})

	t.Run("password contains no uppercase characters", func(t *testing.T) {
		mock := validPayload
		mock.Password = "password321"
		var requestValidationError *primitive.RequestValidationError

		got := service.ValidateBasicAuthLoginBody(mock)

		if got == nil {
			t.Errorf("Expect: error\nGot: nil")
		}

		if !errors.As(got, &requestValidationError) {
			t.Errorf("Expect: %T\nGot: %T", got, requestValidationError)
		}
	})

	t.Run("password contains no number", func(t *testing.T) {
		mock := validPayload
		mock.Password = "Passwords"
		var requestValidationError *primitive.RequestValidationError

		got := service.ValidateBasicAuthLoginBody(mock)

		if got == nil {
			t.Errorf("Expect: error\nGot: nil")
		}

		if !errors.As(got, &requestValidationError) {
			t.Errorf("Expect: %T\nGot: %T", got, requestValidationError)
		}
	})
}
