import Link from 'next/link';

export const revalidate = 900;

export default function HomePage() {
  return (
    <main className="mx-auto min-h-screen max-w-5xl px-4 py-16 sm:px-6">
      <section className="rounded-3xl border border-slate-200 bg-white/90 p-8 shadow-xl backdrop-blur">
        <p className="mb-2 text-sm font-semibold uppercase tracking-wider text-teal-700">AggiPay Platform</p>
        <h1 className="text-3xl font-bold text-slate-900 sm:text-4xl">
          Gestão financeira, contratos e automações em um único painel
        </h1>
        <p className="mt-4 max-w-2xl text-slate-600">
          Interface otimizada para operações de cobrança, pagamentos e projetos. Renderização híbrida
          com navegação SPA-like e integração completa com a API existente.
        </p>
        <div className="mt-6 flex flex-wrap gap-3">
          <Link
            href="/login"
            className="rounded-xl bg-primary px-4 py-2 font-semibold text-white shadow hover:opacity-95"
          >
            Entrar
          </Link>
          <Link
            href="/register"
            className="rounded-xl border border-teal-300 bg-teal-50 px-4 py-2 font-semibold text-teal-800"
          >
            Registrar
          </Link>
          <Link
            href="/dashboard"
            className="rounded-xl border border-slate-300 bg-white px-4 py-2 font-semibold text-slate-800"
          >
            Ir para dashboard
          </Link>
        </div>
      </section>
    </main>
  );
}
