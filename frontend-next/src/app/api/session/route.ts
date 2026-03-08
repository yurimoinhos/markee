import { cookies } from 'next/headers';
import { NextResponse } from 'next/server';

import { BACKEND_API_URL, SESSION_COOKIE_NAME } from '@/lib/config';

export async function GET(): Promise<Response> {
  const cookieStore = await cookies();
  const token = cookieStore.get(SESSION_COOKIE_NAME)?.value;

  if (!token) {
    return NextResponse.json({ authenticated: false });
  }

  const upstream = await fetch(`${BACKEND_API_URL}/auth/me`, {
    method: 'GET',
    headers: {
      Authorization: `Bearer ${token}`,
    },
    cache: 'no-store',
  });

  const text = await upstream.text();
  let payload: Record<string, unknown> | null = null;
  if (text) {
    try {
      payload = JSON.parse(text) as Record<string, unknown>;
    } catch {
      payload = null;
    }
  }

  if (!upstream.ok || !payload) {
    if (upstream.status === 401 || upstream.status === 403) {
      cookieStore.delete(SESSION_COOKIE_NAME);
    }
    return NextResponse.json({ authenticated: false });
  }

  const permissions = Array.isArray(payload.permissions)
    ? payload.permissions.filter((entry): entry is string => typeof entry === 'string')
    : [];
  const roles = Array.isArray(payload.roles)
    ? payload.roles.filter((entry): entry is string => typeof entry === 'string')
    : [];

  return NextResponse.json({
    authenticated: true,
    user: {
      id: String(payload.id ?? ''),
      email: String(payload.email ?? ''),
      firstName: String(payload.firstName ?? ''),
      lastName: String(payload.lastName ?? ''),
      role: String(payload.role ?? ''),
      roles,
      permissions,
    },
  });
}
