import { APP_BASE_PATH } from '@/lib/config';
import { extractProblemMessage } from '@/lib/problem';

export class HttpError extends Error {
  constructor(
    public readonly status: number,
    public readonly body: unknown,
    message: string
  ) {
    super(message);
  }
}

async function parseBody(res: Response): Promise<unknown> {
  const text = await res.text();
  if (!text) {
    return null;
  }
  try {
    return JSON.parse(text);
  } catch {
    return text;
  }
}

export async function apiRequest<T>(
  path: string,
  init?: RequestInit
): Promise<T> {
  const response = await fetch(`${APP_BASE_PATH}${path}`, {
    credentials: 'include',
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...(init?.headers ?? {}),
    },
  });

  const body = await parseBody(response);
  if (!response.ok) {
    throw new HttpError(
      response.status,
      body,
      extractProblemMessage(body, `Erro HTTP ${response.status}`)
    );
  }

  return body as T;
}

export async function apiFormRequest<T>(
  path: string,
  form: URLSearchParams,
  init?: RequestInit
): Promise<T> {
  const response = await fetch(`${APP_BASE_PATH}${path}`, {
    method: 'POST',
    credentials: 'include',
    ...init,
    body: form.toString(),
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
      ...(init?.headers ?? {}),
    },
  });

  const body = await parseBody(response);
  if (!response.ok) {
    throw new HttpError(
      response.status,
      body,
      extractProblemMessage(body, `Erro HTTP ${response.status}`)
    );
  }

  return body as T;
}
