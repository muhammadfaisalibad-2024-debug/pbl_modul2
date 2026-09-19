package service

import (
	"testing"

	"api-students/app/model"
	"api-students/helper"
)

func testPermissions() *helper.PermissionSet {
	return helper.NewPermissionSet(map[string][]string{
		"admin": {"user:read:any", "student:read:any"},
		"user":  {},
	})
}

func TestCanAccessUser(t *testing.T) {
	perms := testPermissions()
	if !CanAccessUser(model.AuthUser{UserID: 1, Role: "user"}, 1, perms, "user:read:any") {
		t.Fatal("pemilik harus dapat mengakses data sendiri")
	}
	if !CanAccessUser(model.AuthUser{UserID: 1, Role: "admin"}, 2, perms, "user:read:any") {
		t.Fatal("role dengan permission harus dapat mengakses user lain")
	}
	if CanAccessUser(model.AuthUser{UserID: 1, Role: "user"}, 2, perms, "user:read:any") {
		t.Fatal("non-pemilik tanpa permission harus ditolak")
	}
}

func TestCanAccessStudent(t *testing.T) {
	perms := testPermissions()
	if !CanAccessStudent(model.AuthUser{UserID: 1, Role: "user"}, 1, perms, "student:read:any") {
		t.Fatal("pemilik harus dapat mengakses student sendiri")
	}
	if !CanAccessStudent(model.AuthUser{UserID: 1, Role: "admin"}, 2, perms, "student:read:any") {
		t.Fatal("role dengan permission harus dapat mengakses student lain")
	}
	if CanAccessStudent(model.AuthUser{UserID: 1, Role: "user"}, 2, perms, "student:read:any") {
		t.Fatal("non-pemilik tanpa permission harus ditolak")
	}
}

func TestValidateAssignRole(t *testing.T) {
	perms := testPermissions()
	tests := []struct {
		name    string
		current model.AuthUser
		target  int
		role    string
		valid   bool
	}{
		{"valid", model.AuthUser{UserID: 1}, 2, "user", true},
		{"empty", model.AuthUser{UserID: 1}, 2, " ", false},
		{"unknown", model.AuthUser{UserID: 1}, 2, "owner", false},
		{"self", model.AuthUser{UserID: 1}, 1, "user", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateAssignRole(tt.current, tt.target, model.AssignRoleRequest{Role: tt.role}, perms)
			if (len(errs) == 0) != tt.valid {
				t.Fatalf("valid=%v, errors=%v", tt.valid, errs)
			}
		})
	}
}
