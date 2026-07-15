package authz

import "testing"

func TestDefaultPolicySeparatesPlatformAdministrationFromChangeApproval(t *testing.T) {
	p := DefaultPolicy()
	platformAdmin := Actor{ID: "platform-admin", Roles: []Role{RolePlatformAdmin}}
	if p.Allows(platformAdmin, PermissionApproveChange) {
		t.Fatal("platform administrator must not automatically receive change approval authority")
	}
}

func TestNetworkSecurityCanApproveAndAudit(t *testing.T) {
	p := DefaultPolicy()
	actor := Actor{ID: "security-01", Roles: []Role{RoleNetworkSecurity}}
	if !p.Allows(actor, PermissionApproveChange) || !p.Allows(actor, PermissionViewAudit) {
		t.Fatal("network security role should hold the expected governed permissions")
	}
}
