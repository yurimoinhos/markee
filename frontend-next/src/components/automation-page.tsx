'use client';

import { useEffect, useState } from 'react';

import type { AutomationRun, RunResult } from '@/generated/models';
import { AppShell } from '@/components/app-shell';
import { DataTable } from '@/components/data-table';
import { listAutomationRuns, runAutomationNow } from '@/domain/automation';
import { Feedback } from '@/components/feedback';
import { SectionCard } from '@/components/section-card';
import { hasPermission } from '@/lib/auth';
import { useSession } from '@/lib/use-session';

export function AutomationPage() {
  const session = useSession();
  const permissions = session.user?.permissions ?? [];
  const canReadAutomation = hasPermission(permissions, 'automation.read');
  const canRunAutomation = hasPermission(permissions, 'automation.write');

  const [items, setItems] = useState<AutomationRun[]>([]);
  const [runResult, setRunResult] = useState<RunResult | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  async function load(): Promise<void> {
    if (!canReadAutomation) {
      return;
    }
    try {
      setItems(await listAutomationRuns());
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao listar automações');
    }
  }

  useEffect(() => {
    if (session.loading || !canReadAutomation) {
      return;
    }
    void load();
  }, [session.loading, canReadAutomation]);

  async function runAutomation(): Promise<void> {
    if (!canRunAutomation) {
      return;
    }
    setError(null);
    try {
      const res = await runAutomationNow();
      setRunResult(res);
      setSuccess('Automação executada com sucesso.');
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao executar automação');
    }
  }

  return (
    <AppShell title="Automação">
      <Feedback success={success} error={error} />

      <div className="mb-4 flex justify-end">
        <button
          type="button"
          onClick={() => void load()}
          className="rounded-xl border border-slate-300 bg-white px-4 py-2 text-sm font-semibold"
        >
          Atualizar
        </button>
      </div>

      {canRunAutomation ? (
        <SectionCard title="Execução manual">
          <p className="text-sm text-slate-600">
            Dispara lembretes, recorrências e relatório mensal de forma imediata.
          </p>
          <button
            className="mt-3 rounded-xl bg-primary px-4 py-2 font-semibold text-white"
            type="button"
            onClick={() => void runAutomation()}
          >
            Executar agora
          </button>

          <pre className="mt-4 overflow-x-auto rounded-xl bg-slate-900 p-3 text-xs text-emerald-200">
            {runResult ? JSON.stringify(runResult, null, 2) : '{ }'}
          </pre>
        </SectionCard>
      ) : null}

      {canReadAutomation ? (
        <div className="mt-4">
          <DataTable
            rows={items as unknown as Record<string, unknown>[]}
            columns={[
              { key: 'id', label: 'ID' },
              { key: 'created_at', label: 'Criado em' },
              { key: 'status', label: 'Status' },
            ]}
          />
        </div>
      ) : (
        <SectionCard title="Automação">
          <p className="text-sm text-slate-600">
            Você não possui permissão de leitura para automação.
          </p>
        </SectionCard>
      )}
    </AppShell>
  );
}
