import AsyncStorage from '@react-native-async-storage/async-storage';

import type { ProblemDetail } from '../generated/api';
import { MOBILE_API_BASE, TOKEN_STORAGE_KEY } from './config';

function messageFromProblem(value: unknown): string {
  if (!value || typeof value !== 'object') return 'Operação falhou';
  const p = value as ProblemDetail;
  return p.errorDescription || p.error || 'Operação falhou';
}

export async function setToken(token: string): Promise<void> {
  await AsyncStorage.setItem(TOKEN_STORAGE_KEY, token);
}

export async function clearToken(): Promise<void> {
  await AsyncStorage.removeItem(TOKEN_STORAGE_KEY);
}

export async function getToken(): Promise<string | null> {
  return AsyncStorage.getItem(TOKEN_STORAGE_KEY);
}

export async function request<T>(
  path: string,
  init?: RequestInit,
  auth = true
): Promise<T> {
  const headers = new Headers(init?.headers ?? {});
  if (!headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }

  if (auth) {
    const token = await getToken();
    if (token) {
      headers.set('Authorization', `Bearer ${token}`);
    }
  }

  const response = await fetch(`${MOBILE_API_BASE}${path}`, {
    ...init,
    headers,
  });

  const text = await response.text();
  const payload = text ? JSON.parse(text) : null;

  if (!response.ok) {
    throw new Error(messageFromProblem(payload));
  }

  return payload as T;
}
