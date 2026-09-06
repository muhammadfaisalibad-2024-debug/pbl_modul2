package service

import (
	"testing"

	"api-students/app/model"
)

// --- ValidateCreate ---

func TestValidateCreate_Success(t *testing.T) {
	req := &model.CreateStudentRequest{
		NIM:  "12345",
		Name: "Faisal",
	}
	if err := ValidateCreate(req); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestValidateCreate_MissingNIM(t *testing.T) {
	req := &model.CreateStudentRequest{
		Name: "Faisal",
	}
	if err := ValidateCreate(req); err != ErrNIMRequired {
		t.Errorf("expected ErrNIMRequired, got %v", err)
	}
}

func TestValidateCreate_MissingName(t *testing.T) {
	req := &model.CreateStudentRequest{
		NIM: "12345",
	}
	if err := ValidateCreate(req); err != ErrNameRequired {
		t.Errorf("expected ErrNameRequired, got %v", err)
	}
}

// --- ValidateReplace ---

func TestValidateReplace_Success(t *testing.T) {
	req := &model.ReplaceStudentRequest{
		NIM:  "99999",
		Name: "Budi",
	}
	if err := ValidateReplace(req); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestValidateReplace_MissingNIM(t *testing.T) {
	req := &model.ReplaceStudentRequest{
		Name: "Budi",
	}
	if err := ValidateReplace(req); err != ErrNIMRequired {
		t.Errorf("expected ErrNIMRequired, got %v", err)
	}
}

func TestValidateReplace_MissingName(t *testing.T) {
	req := &model.ReplaceStudentRequest{
		NIM: "99999",
	}
	if err := ValidateReplace(req); err != ErrNameRequired {
		t.Errorf("expected ErrNameRequired, got %v", err)
	}
}

// --- IsEmptyPatch ---

func TestIsEmptyPatch_AllNil(t *testing.T) {
	req := &model.PatchStudentRequest{}
	if !IsEmptyPatch(req) {
		t.Error("expected IsEmptyPatch to return true for all-nil struct")
	}
}

func TestIsEmptyPatch_WithNIM(t *testing.T) {
	nim := "55555"
	req := &model.PatchStudentRequest{NIM: &nim}
	if IsEmptyPatch(req) {
		t.Error("expected IsEmptyPatch to return false when NIM is set")
	}
}

func TestIsEmptyPatch_WithGrade(t *testing.T) {
	grade := 3.75
	req := &model.PatchStudentRequest{Grade: &grade}
	if IsEmptyPatch(req) {
		t.Error("expected IsEmptyPatch to return false when Grade is set")
	}
}

// --- CountTotalPages ---

func TestCountTotalPages_Zero(t *testing.T) {
	if CountTotalPages(0, 10) != 0 {
		t.Error("expected 0 total pages when total is 0")
	}
}

func TestCountTotalPages_Exact(t *testing.T) {
	if CountTotalPages(20, 10) != 2 {
		t.Errorf("expected 2 total pages, got %d", CountTotalPages(20, 10))
	}
}

func TestCountTotalPages_Remainder(t *testing.T) {
	if CountTotalPages(21, 10) != 3 {
		t.Errorf("expected 3 total pages, got %d", CountTotalPages(21, 10))
	}
}
