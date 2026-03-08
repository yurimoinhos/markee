export type SessionUser = {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  role: string;
  roles: string[];
  permissions: string[];
};

export type SessionResponse = {
  authenticated: boolean;
  user?: SessionUser;
};

export function hasPermission(permissions: string[], permission: string): boolean {
  return permissions.includes(permission);
}

export function hasAnyPermission(permissions: string[], requiredPermissions: string[]): boolean {
  return requiredPermissions.some((permission) => hasPermission(permissions, permission));
}

export function hasAccess(
  role: string,
  permissions: string[],
  requiredRoles?: string[],
  requiredPermissions?: string[]
): boolean {
  if (requiredRoles && requiredRoles.length > 0 && !requiredRoles.includes(role)) {
    return false;
  }

  if (requiredPermissions && requiredPermissions.length > 0) {
    return hasAnyPermission(permissions, requiredPermissions);
  }

  return true;
}
