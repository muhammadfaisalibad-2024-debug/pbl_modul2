package service

import (
	"testing"

	"api-students/app/model"
)

func TestValidateRegisterValid(t *testing.T) {
	req := model.RegisterRequest{
		Username: "faisal_01",
		Email:    "faisal@example.com",
		Password: "rahasia123",
	}

	errs := ValidateRegister(req)

	if len(errs) != 0 {
		t.Fatalf("expected no validation errors, got %v", errs)
	}
}

func TestValidateRegisterPasswordTooShort(t *testing.T) {
	req := model.RegisterRequest{
		Username: "faisal",
		Email:    "faisal@example.com",
		Password: "abc12",
	}

	errs := ValidateRegister(req)

	if errs["password"] != "minimal 8 karakter" {
		t.Fatalf(
			"expected password too short error, got %v",
			errs["password"],
		)
	}
}

func TestValidateRegisterWeakPassword(t *testing.T) {
	req := model.RegisterRequest{
		Username: "faisal",
		Email:    "faisal@example.com",
		Password: "password1",
	}

	errs := ValidateRegister(req)

	if errs["password"] != "password terlalu umum" {
		t.Fatalf(
			"expected weak password error, got %v",
			errs["password"],
		)
	}
}

func TestValidateRegisterInvalidUsername(t *testing.T) {
	req := model.RegisterRequest{
		Username: "faisal@123",
		Email:    "faisal@example.com",
		Password: "rahasia123",
	}

	errs := ValidateRegister(req)

	if errs["username"] == "" {
		t.Fatal("expected username validation error")
	}
}

func TestValidateRegisterInvalidEmail(t *testing.T) {
	req := model.RegisterRequest{
		Username: "faisal",
		Email:    "faisal",
		Password: "rahasia123",
	}

	errs := ValidateRegister(req)

	if errs["email"] != "format email tidak valid" {
		t.Fatalf(
			"expected invalid email error, got %v",
			errs["email"],
		)
	}
}

func TestValidateLoginEmpty(t *testing.T) {
	req := model.LoginRequest{
		Username: "",
		Password: "",
	}

	errs := ValidateLogin(req)

	if errs["username"] == "" {
		t.Fatal("expected username validation error")
	}

	if errs["password"] == "" {
		t.Fatal("expected password validation error")
	}
}
