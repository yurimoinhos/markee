import { cookies } from 'next/headers';
import { NextResponse } from 'next/server';

import type { ProblemDetail, TokenResponse } from '@/generated/models';
import {
  BACKEND_API_URL,
  SESSION_COOKIE_NAME,
  SESSION_COOKIE_TTL_SECONDS,
} from '@/lib/config';

export async function POST(request: Request): Promise<Response> {
  const form = await request.formData();
  const firstName = String(form.get('firstName') ?? '').trim();
  const lastName = String(form.get('lastName') ?? '').trim();
  const email = String(form.get('email') ?? '').trim();
  const phoneNumber = String(form.get('phoneNumber') ?? '').trim();
  const password = String(form.get('password') ?? '').trim();

  const registerBody = new URLSearchParams();
  registerBody.set('firstName', firstName);
  registerBody.set('lastName', lastName);
  registerBody.set('email', email);
  if (phoneNumber) {
    registerBody.set('phoneNumber', phoneNumber);
  }
  registerBody.set('password', password);

  const registerResponse = await fetch(`${BACKEND_API_URL}/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    body: registerBody.toString(),
    cache: 'no-store',
  });

  const registerText = await registerResponse.text();
  let registerPayload: Record<string, unknown> | ProblemDetail | string | null = null;
  if (registerText) {
    try {
      registerPayload = JSON.parse(registerText) as Record<string, unknown> | ProblemDetail;
    } catch {
      registerPayload = registerText;
    }
  }

  if (!registerResponse.ok) {
    return NextResponse.json(
      registerPayload ?? { error: 'Falha no registro' },
      { status: registerResponse.status }
    );
  }

  const loginBody = new URLSearchParams();
  loginBody.set('email', email);
  loginBody.set('password', password);

  const loginResponse = await fetch(`${BACKEND_API_URL}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    body: loginBody.toString(),
    cache: 'no-store',
  });

  const loginText = await loginResponse.text();
  let loginPayload: TokenResponse | ProblemDetail | string | null = null;
  if (loginText) {
    try {
      loginPayload = JSON.parse(loginText) as TokenResponse | ProblemDetail;
    } catch {
      loginPayload = loginText;
    }
  }

  if (!loginResponse.ok) {
    return NextResponse.json(
      loginPayload ?? { error: 'Registro realizado, mas o login automático falhou' },
      { status: loginResponse.status }
    );
  }

  const token =
    typeof loginPayload === 'object' && loginPayload && 'access_token' in loginPayload
      ? String(loginPayload.access_token)
      : '';

  if (!token) {
    return NextResponse.json(
      {
        error: 'Registro realizado, mas o login automático falhou',
        errorDescription: 'token ausente na resposta do backend',
      },
      { status: 502 }
    );
  }

  const cookieStore = await cookies();
  cookieStore.set(SESSION_COOKIE_NAME, token, {
    httpOnly: true,
    sameSite: 'lax',
    secure: process.env.NODE_ENV === 'production',
    path: '/',
    maxAge: SESSION_COOKIE_TTL_SECONDS,
  });

  return NextResponse.json({ ok: true, token_type: 'Bearer' });
}

