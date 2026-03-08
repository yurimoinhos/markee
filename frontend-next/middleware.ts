import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

const APP_BASE_PATH = '/app';
const SESSION_COOKIE_NAME = process.env.SESSION_COOKIE_NAME ?? 'aggipay_token';

const PUBLIC_PATHS = new Set([
  `${APP_BASE_PATH}`,
  `${APP_BASE_PATH}/login`,
  `${APP_BASE_PATH}/register`,
]);

export function middleware(request: NextRequest): NextResponse {
  const { pathname } = request.nextUrl;

  if (!pathname.startsWith(APP_BASE_PATH)) {
    return NextResponse.next();
  }

  if (pathname.startsWith(`${APP_BASE_PATH}/api`) || pathname.startsWith(`${APP_BASE_PATH}/_next`)) {
    return NextResponse.next();
  }

  const token = request.cookies.get(SESSION_COOKIE_NAME)?.value;

  if (!token && !PUBLIC_PATHS.has(pathname)) {
    const url = request.nextUrl.clone();
    url.pathname = `${APP_BASE_PATH}/login`;
    return NextResponse.redirect(url);
  }

  if (token && (pathname === `${APP_BASE_PATH}/login` || pathname === `${APP_BASE_PATH}/register`)) {
    const url = request.nextUrl.clone();
    url.pathname = `${APP_BASE_PATH}/dashboard`;
    return NextResponse.redirect(url);
  }

  return NextResponse.next();
}

export const config = {
  matcher: ['/app/:path*'],
};
