export const APP_BASE_PATH = '/app';

export const BACKEND_API_URL =
  process.env.BACKEND_API_URL ?? 'http://127.0.0.1:8000/api/v1';

export const SESSION_COOKIE_NAME =
  process.env.SESSION_COOKIE_NAME ?? 'aggipay_token';

const ttl = Number.parseInt(process.env.SESSION_COOKIE_TTL_SECONDS ?? '86400', 10);
export const SESSION_COOKIE_TTL_SECONDS = Number.isFinite(ttl) ? ttl : 86400;
