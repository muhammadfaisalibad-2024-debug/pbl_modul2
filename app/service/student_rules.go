package service

import (
	"errors"

	"api-students/app/model"
)

// Sentinel validation errors
var (
	ErrNIMRequired  = errors.New("NIM is required")
	ErrNameRequired = errors.New("Name is required")
	ErrEmptyPatch   = errors.New("no fields to update")
)

// ValidateCreate validates a POST (create) request.
func ValidateCreate(req *model.CreateStudentRequest) error {
	if req.NIM == "" {
		return ErrNIMRequired
	}
	if req.Name == "" {
		return ErrNameRequired
	}
	return nil
}

// ValidateReplace validates a PUT (full replace) request.
func ValidateReplace(req *model.ReplaceStudentRequest) error {
	if req.NIM == "" {
		return ErrNIMRequired
	}
	if req.Name == "" {
		return ErrNameRequired
	}
	return nil
}

// IsEmptyPatch checks whether a PatchStudentRequest has at least one non-nil field.
func IsEmptyPatch(req *model.PatchStudentRequest) bool {
	return req.NIM == nil && req.Name == nil && req.Grade == nil && req.IsActive == nil
}

// CountTotalPages computes total pages given total records and page size.
func CountTotalPages(total, limit int) int {
	if total == 0 || limit <= 0 {
		return 0
	}
	return (total + limit - 1) / limit
}
