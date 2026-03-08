import { cookies } from 'next/headers';
import { NextResponse } from 'next/server';

import { BACKEND_API_URL, SESSION_COOKIE_NAME } from '@/lib/config';

type Params = { params: Promise<{ path?: string[] }> };

async function proxy(request: Request, params: Params['params']): Promise<Response> {
  const { path = [] } = await params;
  const cookieStore = await cookies();
  const token = cookieStore.get(SESSION_COOKIE_NAME)?.value;

  if (!token) {
    return NextResponse.json(
      { error: 'Não autorizado', errorDescription: 'sessão ausente' },
      { status: 401 }
    );
  }

  const url = new URL(request.url);
  const query = url.search ? url.search : '';
  const upstreamURL = `${BACKEND_API_URL}/${path.join('/')}${query}`;

  const headers = new Headers();
  headers.set('Authorization', `Bearer ${token}`);

  const contentType = request.headers.get('content-type');
  if (contentType) {
    headers.set('Content-Type', contentType);
  }

  const method = request.method.toUpperCase();
  const canHaveBody = !['GET', 'HEAD'].includes(method);
  let body: string | undefined;
  if (canHaveBody) {
    const text = await request.text();
    body = text === '' ? undefined : text;
  }

  const upstream = await fetch(upstreamURL, {
    method,
    headers,
    body,
    cache: 'no-store',
  });

  const text = await upstream.text();
  const responseHeaders = new Headers();
  const upstreamType = upstream.headers.get('content-type') ?? 'application/json';
  responseHeaders.set('Content-Type', upstreamType);

  return new Response(text, {
    status: upstream.status,
    headers: responseHeaders,
  });
}

export async function GET(request: Request, ctx: Params): Promise<Response> {
  return proxy(request, ctx.params);
}

export async function POST(request: Request, ctx: Params): Promise<Response> {
  return proxy(request, ctx.params);
}

export async function PUT(request: Request, ctx: Params): Promise<Response> {
  return proxy(request, ctx.params);
}

export async function PATCH(request: Request, ctx: Params): Promise<Response> {
  return proxy(request, ctx.params);
}

export async function DELETE(request: Request, ctx: Params): Promise<Response> {
  return proxy(request, ctx.params);
}
