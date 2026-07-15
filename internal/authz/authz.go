package authz

import (
	"errors"
	"sort"
)

type Role string

type Permission string

const (
	RoleViewer          Role = "viewer"
	RoleNetworkTech     Role = "network_technician"
	RoleNetworkAdmin    Role = "network_administrator"
	RoleNetworkSecurity Role = "network_security"
	RoleChangeApprover  Role = "change_approver"
	RoleAuditor         Role = "auditor"
	RolePlatformAdmin   Role = "platform_administrator"
)

const (
	PermissionViewDashboard Permission = "dashboard.view"
	PermissionViewInventory Permission = "inventory.view"
	PermissionCollect       Permission = "collection.execute"
	PermissionRequestChange Permission = "change.request"
	PermissionApproveChange Permission = "change.approve"
	PermissionAdminModules  Permission = "modules.administer"
	PermissionViewEvidence  Permission = "evidence.view_redacted"
	PermissionViewAudit     Permission = "audit.view"
)

type Actor struct {
	ID    string
	Roles []Role
}

type Policy struct {
	grants map[Role]map[Permission]struct{}
}

func DefaultPolicy() *Policy {
	p := &Policy{grants: make(map[Role]map[Permission]struct{})}
	p.grant(RoleViewer, PermissionViewDashboard, PermissionViewInventory)
	p.grant(RoleNetworkTech, PermissionViewDashboard, PermissionViewInventory, PermissionCollect, PermissionRequestChange, PermissionViewEvidence)
	p.grant(RoleNetworkAdmin, PermissionViewDashboard, PermissionViewInventory, PermissionCollect, PermissionRequestChange, PermissionApproveChange, PermissionViewEvidence)
	p.grant(RoleNetworkSecurity, PermissionViewDashboard, PermissionViewInventory, PermissionRequestChange, PermissionApproveChange, PermissionViewEvidence, PermissionViewAudit)
	p.grant(RoleChangeApprover, PermissionViewDashboard, PermissionViewInventory, PermissionApproveChange, PermissionViewEvidence)
	p.grant(RoleAuditor, PermissionViewDashboard, PermissionViewInventory, PermissionViewEvidence, PermissionViewAudit)
	p.grant(RolePlatformAdmin, PermissionViewDashboard, PermissionViewInventory, PermissionAdminModules, PermissionViewAudit)
	return p
}

func (p *Policy) grant(role Role, permissions ...Permission) {
	if p.grants[role] == nil {
		p.grants[role] = make(map[Permission]struct{})
	}
	for _, permission := range permissions {
		p.grants[role][permission] = struct{}{}
	}
}

func (p *Policy) Allows(actor Actor, permission Permission) bool {
	for _, role := range actor.Roles {
		if _, ok := p.grants[role][permission]; ok {
			return true
		}
	}
	return false
}

func (p *Policy) Require(actor Actor, permission Permission) error {
	if actor.ID == "" {
		return errors.New("authenticated actor is required")
	}
	if !p.Allows(actor, permission) {
		return errors.New("permission denied")
	}
	return nil
}

func (a Actor) RoleNames() []string {
	names := make([]string, 0, len(a.Roles))
	for _, role := range a.Roles {
		names = append(names, string(role))
	}
	sort.Strings(names)
	return names
}
