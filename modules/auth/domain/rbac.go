package domain

import (
	"maps"
	"slices"
	"strings"
)

const (
	RoleAdmin   = "admin"
	RoleManager = "manager"
	RoleViewer  = "viewer"

	DefaultRole = RoleManager
)

const (
	PermDashboardRead = "dashboard.read"
	PermFinanceRead   = "finance.read"

	PermCustomersRead  = "customers.read"
	PermCustomersWrite = "customers.write"

	PermContractsRead  = "contracts.read"
	PermContractsWrite = "contracts.write"

	PermChargesRead  = "charges.read"
	PermChargesWrite = "charges.write"

	PermPaymentsRead  = "payments.read"
	PermPaymentsWrite = "payments.write"

	PermProjectsRead  = "projects.read"
	PermProjectsWrite = "projects.write"

	PermAutomationRead  = "automation.read"
	PermAutomationWrite = "automation.write"

	PermUsersRead  = "users.read"
	PermUsersWrite = "users.write"
)

var rolePermissions = map[string][]string{
	RoleAdmin: {
		PermDashboardRead,
		PermFinanceRead,
		PermCustomersRead,
		PermCustomersWrite,
		PermContractsRead,
		PermContractsWrite,
		PermChargesRead,
		PermChargesWrite,
		PermPaymentsRead,
		PermPaymentsWrite,
		PermProjectsRead,
		PermProjectsWrite,
		PermAutomationRead,
		PermAutomationWrite,
		PermUsersRead,
		PermUsersWrite,
	},
	RoleManager: {
		PermDashboardRead,
		PermFinanceRead,
		PermCustomersRead,
		PermCustomersWrite,
		PermContractsRead,
		PermContractsWrite,
		PermChargesRead,
		PermChargesWrite,
		PermPaymentsRead,
		PermPaymentsWrite,
		PermProjectsRead,
		PermProjectsWrite,
		PermAutomationRead,
		PermAutomationWrite,
		PermUsersRead,
	},
	RoleViewer: {
		PermDashboardRead,
		PermFinanceRead,
		PermCustomersRead,
		PermContractsRead,
		PermChargesRead,
		PermPaymentsRead,
		PermProjectsRead,
		PermAutomationRead,
		PermUsersRead,
	},
}

func RolePermissions() map[string][]string {
	out := make(map[string][]string, len(rolePermissions))
	for role, permissions := range rolePermissions {
		copied := make([]string, len(permissions))
		copy(copied, permissions)
		out[role] = copied
	}
	return out
}

func PermissionDescriptions() map[string]string {
	return map[string]string{
		PermDashboardRead:   "Visualizar dashboard",
		PermFinanceRead:     "Visualizar finanças",
		PermCustomersRead:   "Visualizar clientes",
		PermCustomersWrite:  "Criar/editar clientes",
		PermContractsRead:   "Visualizar contratos",
		PermContractsWrite:  "Criar/editar contratos",
		PermChargesRead:     "Visualizar cobranças",
		PermChargesWrite:    "Criar/editar cobranças",
		PermPaymentsRead:    "Visualizar pagamentos",
		PermPaymentsWrite:   "Registrar/editar pagamentos",
		PermProjectsRead:    "Visualizar projetos",
		PermProjectsWrite:   "Criar/editar projetos",
		PermAutomationRead:  "Visualizar automações",
		PermAutomationWrite: "Executar automações",
		PermUsersRead:       "Visualizar usuários",
		PermUsersWrite:      "Gerenciar usuários",
	}
}

func NormalizeRole(role string) string {
	normalized := strings.ToLower(strings.TrimSpace(role))
	if _, ok := rolePermissions[normalized]; ok {
		return normalized
	}
	return ""
}

func NormalizeRoles(roles []string) []string {
	if len(roles) == 0 {
		return []string{DefaultRole}
	}

	unique := make(map[string]struct{}, len(roles))
	out := make([]string, 0, len(roles))
	for _, role := range roles {
		normalized := NormalizeRole(role)
		if normalized == "" {
			continue
		}
		if _, exists := unique[normalized]; exists {
			continue
		}
		unique[normalized] = struct{}{}
		out = append(out, normalized)
	}

	if len(out) == 0 {
		return []string{DefaultRole}
	}

	slices.Sort(out)
	return out
}

func PrimaryRole(roles []string) string {
	normalized := NormalizeRoles(roles)
	priority := []string{RoleAdmin, RoleManager, RoleViewer}
	for _, candidate := range priority {
		if slices.Contains(normalized, candidate) {
			return candidate
		}
	}
	return normalized[0]
}

func PermissionsForRoles(roles []string) []string {
	normalizedRoles := NormalizeRoles(roles)
	unique := make(map[string]struct{})
	for _, role := range normalizedRoles {
		for _, permission := range rolePermissions[role] {
			unique[permission] = struct{}{}
		}
	}
	return NormalizePermissions(slices.Collect(maps.Keys(unique)))
}

func DefaultPermissions(role string) []string {
	return PermissionsForRoles([]string{role})
}

func EnsurePermissionsForRole(role string, perms []string) []string {
	return EnsurePermissionsForRoles([]string{role}, perms)
}

func NormalizePermissions(perms []string) []string {
	if len(perms) == 0 {
		return nil
	}

	unique := make(map[string]struct{}, len(perms))
	out := make([]string, 0, len(perms))
	for _, perm := range perms {
		normalized := strings.ToLower(strings.TrimSpace(perm))
		if normalized == "" {
			continue
		}
		if _, exists := unique[normalized]; exists {
			continue
		}
		unique[normalized] = struct{}{}
		out = append(out, normalized)
	}

	slices.Sort(out)
	return out
}

func EnsurePermissionsForRoles(roles []string, perms []string) []string {
	normalizedRoles := NormalizeRoles(roles)
	allowedSet := make(map[string]struct{})
	for _, permission := range PermissionsForRoles(normalizedRoles) {
		allowedSet[permission] = struct{}{}
	}

	normalizedPerms := NormalizePermissions(perms)
	if len(normalizedPerms) == 0 {
		return PermissionsForRoles(normalizedRoles)
	}

	filtered := make([]string, 0, len(normalizedPerms))
	for _, permission := range normalizedPerms {
		if _, ok := allowedSet[permission]; ok {
			filtered = append(filtered, permission)
		}
	}

	if len(filtered) == 0 {
		return PermissionsForRoles(normalizedRoles)
	}
	return NormalizePermissions(filtered)
}
