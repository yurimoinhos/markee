'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { FormEvent, useState } from 'react';

import { APP_BASE_PATH } from '@/lib/config';
import { extractProblemMessage } from '@/lib/problem';

export function LoginPage() {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function onSubmit(event: FormEvent<HTMLFormElement>): Promise<void> {
    event.preventDefault();
    setLoading(true);
    setError(null);

    const body = new URLSearchParams();
    body.set('email', email);
    body.set('password', password);

    try {
      const response = await fetch(`${APP_BASE_PATH}/api/auth/login`, {
        method: 'POST',
        body,
      });

      if (!response.ok) {
        const payload = await response.json().catch(() => null);
        throw new Error(extractProblemMessage(payload, 'Falha ao autenticar'));
      }

      router.push('/dashboard');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao autenticar');
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="grid min-h-screen place-items-center px-4 py-10">
      <section className="w-full max-w-md rounded-3xl border border-slate-200 bg-white p-6 shadow-xl">
        <p className="mb-2 text-xs font-semibold uppercase tracking-[0.2em] text-teal-700">AggiPay</p>
        <h1 className="text-2xl font-bold text-slate-900">Acessar console</h1>
        <p className="mt-1 text-sm text-slate-600">Use email e senha para entrar.</p>

        {error ? (
          <div className="mt-4 rounded-xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-danger">
            {error}
          </div>
        ) : null}

        <form className="mt-5 space-y-3" onSubmit={onSubmit}>
          <label className="block text-sm font-medium text-slate-700">
            Email
            <input
              required
              type="email"
              value={email}
              onChange={(event) => setEmail(event.target.value)}
              placeholder="voce@empresa.com"
              className="mt-1"
            />
          </label>
          <label className="block text-sm font-medium text-slate-700">
            Senha
            <input
              required
              type="password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              placeholder="********"
              className="mt-1"
            />
          </label>

          <button
            type="submit"
            disabled={loading}
            className="w-full rounded-xl bg-primary px-4 py-2 font-semibold text-white disabled:cursor-not-allowed disabled:opacity-70"
          >
            {loading ? 'Entrando...' : 'Entrar'}
          </button>
        </form>

        <p className="mt-4 text-center text-sm text-slate-600">
          Ainda não tem conta?{' '}
          <Link href="/register" className="font-semibold text-teal-700 hover:underline">
            Criar conta
          </Link>
        </p>
      </section>
    </main>
  );
}
