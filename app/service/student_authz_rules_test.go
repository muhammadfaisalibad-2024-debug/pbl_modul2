package service

import (
	"testing"

	"api-students/app/model"
	"api-students/helper"
)

func TestCanAccessStudentOwnerMayUpdate(t *testing.T) {
	perms := helper.NewPermissionSet(map[string][]string{"user": {}})
	if !CanAccessStudent(model.AuthUser{UserID: 5, Role: "user"}, 5, perms, "student:update:any") {
		t.Fatal("student owner should be allowed to update")
	}
}

func TestCanAccessStudentNonOwnerWithoutPermissionDenied(t *testing.T) {
	perms := helper.NewPermissionSet(map[string][]string{"user": {}})
	if CanAccessStudent(model.AuthUser{UserID: 5, Role: "user"}, 4, perms, "student:update:any") {
		t.Fatal("non-owner without update-any permission should be denied")
	}
}

func TestCanAccessStudentAdminWithPermissionAllowed(t *testing.T) {
	perms := helper.NewPermissionSet(map[string][]string{"admin": {"student:update:any"}})
	if !CanAccessStudent(model.AuthUser{UserID: 5, Role: "admin"}, 4, perms, "student:update:any") {
		t.Fatal("admin with update-any permission should be allowed")
	}
}
