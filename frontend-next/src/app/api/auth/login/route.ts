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
  const email = String(form.get('email') ?? '').trim();
  const password = String(form.get('password') ?? '').trim();

  const body = new URLSearchParams();
  body.set('email', email);
  body.set('password', password);

  const upstream = await fetch(`${BACKEND_API_URL}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    body: body.toString(),
    cache: 'no-store',
  });

  const text = await upstream.text();
  let payload: TokenResponse | ProblemDetail | string | null = null;
  if (text) {
    try {
      payload = JSON.parse(text) as TokenResponse | ProblemDetail;
    } catch {
      payload = text;
    }
  }

  if (!upstream.ok) {
    return NextResponse.json(payload ?? { error: 'Falha no login' }, { status: upstream.status });
  }

  const token = typeof payload === 'object' && payload && 'access_token' in payload
    ? String(payload.access_token)
    : '';

  if (!token) {
    return NextResponse.json(
      { error: 'Falha no login', errorDescription: 'token ausente na resposta do backend' },
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
