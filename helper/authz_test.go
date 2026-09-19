package helper

import "testing"

func TestPermissionSetFailClosed(t *testing.T) {
	perms := NewPermissionSet(map[string][]string{"admin": {"user:list"}})

	if !perms.Can("admin", "user:list") {
		t.Fatal("permission yang diberikan harus diizinkan")
	}
	if perms.Can("unknown", "user:list") {
		t.Fatal("role tidak dikenal harus ditolak")
	}
	if perms.Can("admin", "unknown") {
		t.Fatal("permission tidak dikenal harus ditolak")
	}
	var nilPerms *PermissionSet
	if nilPerms.Can("admin", "user:list") {
		t.Fatal("PermissionSet nil harus ditolak")
	}
}
