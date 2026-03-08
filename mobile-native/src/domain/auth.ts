import type { TokenResponse } from '../generated/api';
import { clearToken, setToken } from '../lib/http';

export async function login(email: string, password: string): Promise<void> {
  const body = new URLSearchParams();
  body.set('email', email);
  body.set('password', password);

  const response = await fetch(
    `${process.env.EXPO_PUBLIC_BACKEND_API_URL ?? 'http://127.0.0.1:8000/api/v1'}/auth/login`,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      body: body.toString(),
    }
  );

  const payload = (await response.json()) as TokenResponse;
  if (!response.ok || !payload.access_token) {
    throw new Error('Falha ao autenticar');
  }

  await setToken(payload.access_token);
}

export async function logout(): Promise<void> {
  await clearToken();
}
