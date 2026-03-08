'use client';

import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import { useMemo, useState, type ReactNode } from 'react';

import { APP_BASE_PATH } from '@/lib/config';
import { hasAccess } from '@/lib/auth';
import { useSession } from '@/lib/use-session';

type AppShellProps = {
  title: string;
  children: ReactNode;
};

type NavItem = {
  href: string;
  label: string;
  requiredRoles?: string[];
  requiredPermissions?: string[];
};

const navItems: NavItem[] = [
  { href: '/dashboard', label: 'Dashboard', requiredPermissions: ['dashboard.read', 'finance.read'] },
  { href: '/customers', label: 'Clientes', requiredPermissions: ['customers.read'] },
  { href: '/contracts', label: 'Contratos', requiredPermissions: ['contracts.read'] },
  { href: '/charges', label: 'Cobranças', requiredPermissions: ['charges.read'] },
  { href: '/payments', label: 'Pagamentos', requiredPermissions: ['payments.read'] },
  { href: '/projects', label: 'Projetos', requiredPermissions: ['projects.read'] },
  { href: '/automation', label: 'Automação', requiredPermissions: ['automation.read'] },
];

export function AppShell({ title, children }: AppShellProps) {
  const pathname = usePathname();
  const router = useRouter();
  const [menuOpen, setMenuOpen] = useState(false);
  const session = useSession();
  const sessionRole = session.user?.role ?? '';
  const sessionPermissions = session.user?.permissions ?? [];
  const sessionLoaded = !session.loading;

  const items = useMemo(
    () =>
      navItems
        .filter((item) =>
          hasAccess(
            sessionRole,
            sessionPermissions,
            item.requiredRoles,
            item.requiredPermissions
          )
        )
        .map((item) => ({
          ...item,
          active: pathname === `${APP_BASE_PATH}${item.href}`,
        })),
    [pathname, sessionRole, sessionPermissions]
  );

  async function logout(): Promise<void> {
    await fetch(`${APP_BASE_PATH}/api/auth/logout`, { method: 'POST' });
    router.push('/login');
  }

  const nav = (
    <nav className="space-y-1">
      <h2 className="mb-4 text-xl font-bold text-slate-100">AggiPay Console</h2>
      {sessionLoaded ? (
        <p className="mb-3 rounded-lg border border-slate-700 bg-slate-800 px-2 py-1 text-xs uppercase tracking-wide text-slate-300">
          Perfil: {sessionRole || 'indefinido'}
        </p>
      ) : null}
      {items.map((item) => (
        <Link
          key={item.href}
          href={item.href}
          className={`block rounded-xl px-3 py-2 text-sm transition ${
            item.active
              ? 'bg-teal-700/70 text-white'
              : 'text-slate-300 hover:bg-slate-700/60 hover:text-white'
          }`}
          onClick={() => setMenuOpen(false)}
        >
          {item.label}
        </Link>
      ))}
      {sessionLoaded && items.length === 0 ? (
        <p className="rounded-xl bg-slate-800 px-3 py-2 text-xs text-slate-300">
          Sem permissões para exibir menus.
        </p>
      ) : null}
    </nav>
  );

  return (
    <div className="min-h-screen bg-surface">
      <div className="grid min-h-screen grid-cols-1 lg:grid-cols-[280px_1fr]">
        <aside className="hidden bg-slate-900 p-5 lg:block">{nav}</aside>

        <div className="flex min-h-screen flex-col">
          <header className="sticky top-0 z-20 flex items-center justify-between border-b border-slate-200 bg-white/95 px-4 py-3 backdrop-blur lg:px-6">
            <div className="flex items-center gap-2">
              <button
                className="rounded-lg border border-slate-300 px-3 py-1 text-sm lg:hidden"
                onClick={() => setMenuOpen((v) => !v)}
                type="button"
              >
                Menu
              </button>
              <h1 className="text-lg font-semibold text-slate-900">{title}</h1>
            </div>
            <button
              className="rounded-lg bg-danger px-3 py-2 text-sm font-semibold text-white"
              onClick={logout}
              type="button"
            >
              Sair
            </button>
          </header>

          {menuOpen && (
            <div className="border-b border-slate-300 bg-slate-900 p-4 lg:hidden">{nav}</div>
          )}

          <main className="flex-1 p-4 lg:p-6">{children}</main>
        </div>
      </div>
    </div>
  );
}
