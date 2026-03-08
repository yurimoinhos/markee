'use client';

import { useEffect, useMemo, useState } from 'react';

import type { CashFlowPoint, Dashboard, DefaultMetrics } from '@/generated/models';
import { AppShell } from '@/components/app-shell';
import { Feedback } from '@/components/feedback';
import { SectionCard } from '@/components/section-card';
import { getCashflow, getDashboard, getDefaultMetrics } from '@/domain/finance';
import { formatMoney } from '@/lib/format';
import { hasAnyPermission } from '@/lib/auth';
import { useSession } from '@/lib/use-session';

export function DashboardPage() {
  const session = useSession();
  const permissions = session.user?.permissions ?? [];
  const canReadDashboard = hasAnyPermission(permissions, ['dashboard.read', 'finance.read']);

  const [dashboard, setDashboard] = useState<Dashboard | null>(null);
  const [cashflow, setCashflow] = useState<CashFlowPoint[]>([]);
  const [defaults, setDefaults] = useState<DefaultMetrics | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  async function load(): Promise<void> {
    if (!canReadDashboard) {
      return;
    }
    setLoading(true);
    setError(null);
    try {
      const [d, cf, def] = await Promise.all([
        getDashboard(),
        getCashflow(),
        getDefaultMetrics(),
      ]);
      setDashboard(d);
      setCashflow(cf);
      setDefaults(def);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao carregar dashboard');
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    if (session.loading || !canReadDashboard) {
      return;
    }
    void load();
  }, [session.loading, canReadDashboard]);

  const maxCashflow = useMemo(
    () => Math.max(1, ...cashflow.map((item) => Math.max(item.in_cents ?? 0, item.pending_cents ?? 0))),
    [cashflow]
  );

  return (
    <AppShell title="Dashboard Financeiro">
      <Feedback error={error} />

      <div className="mb-4 flex justify-end">
        <button
          type="button"
          onClick={() => void load()}
          className="rounded-xl border border-slate-300 bg-white px-4 py-2 text-sm font-semibold"
        >
          {loading ? 'Atualizando...' : 'Atualizar'}
        </button>
      </div>

      {!canReadDashboard ? (
        <SectionCard title="Dashboard">
          <p className="text-sm text-slate-600">
            Você não possui permissão de leitura para o dashboard.
          </p>
        </SectionCard>
      ) : (
        <>
      <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
        <Kpi label="Receita Mensal" value={formatMoney(dashboard?.monthly_revenue_cents)} />
        <Kpi label="MRR" value={formatMoney(dashboard?.mrr_cents)} />
        <Kpi label="Pendências" value={String(dashboard?.pending_payments ?? 0)} />
        <Kpi label="Recebidos" value={String(dashboard?.payments_received ?? 0)} />
        <Kpi label="Contratos Ativos" value={String(dashboard?.active_contracts ?? 0)} />
        <Kpi label="Vencendo" value={String(dashboard?.expiring_contracts ?? 0)} />
      </div>

      <div className="mt-4 grid gap-4 xl:grid-cols-[1fr_320px]">
        <SectionCard title="Fluxo de caixa">
          <div className="space-y-2">
            {cashflow.length === 0 ? (
              <p className="text-sm text-slate-500">Sem dados para o período.</p>
            ) : (
              cashflow.map((point) => {
                const inPct = ((point.in_cents ?? 0) / maxCashflow) * 100;
                const pendingPct = ((point.pending_cents ?? 0) / maxCashflow) * 100;

                return (
                  <div key={point.date} className="grid gap-2 sm:grid-cols-[120px_1fr] sm:items-center">
                    <span className="text-xs font-medium text-slate-600">{point.date}</span>
                    <div className="relative h-5 overflow-hidden rounded-full bg-slate-200">
                      <div
                        className="absolute inset-y-0 left-0 rounded-full bg-primary"
                        style={{ width: `${inPct}%` }}
                      />
                      <div
                        className="absolute inset-y-0 left-0 rounded-full bg-warning/80"
                        style={{ width: `${pendingPct}%` }}
                      />
                    </div>
                  </div>
                );
              })
            )}
          </div>
        </SectionCard>

        <SectionCard title="Inadimplência e crescimento">
          <ul className="space-y-2 text-sm text-slate-700">
            <li className="flex justify-between">
              <span>Em atraso</span>
              <strong>{defaults?.overdue_charges ?? 0}</strong>
            </li>
            <li className="flex justify-between">
              <span>Taxa de inadimplência</span>
              <strong>{defaults?.default_rate_percent?.toFixed(1) ?? '0.0'}%</strong>
            </li>
            <li className="flex justify-between">
              <span>Crescimento</span>
              <strong>{defaults?.growth_percent?.toFixed(1) ?? '0.0'}%</strong>
            </li>
          </ul>
        </SectionCard>
      </div>
        </>
      )}
    </AppShell>
  );
}

function Kpi({ label, value }: { label: string; value: string }) {
  return (
    <article className="rounded-2xl border border-slate-200 bg-white p-4 shadow-sm">
      <p className="text-xs uppercase tracking-wide text-slate-500">{label}</p>
      <p className="mt-1 text-2xl font-bold text-slate-900">{value}</p>
    </article>
  );
}
