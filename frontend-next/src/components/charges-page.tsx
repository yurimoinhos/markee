'use client';

import { FormEvent, useEffect, useState } from 'react';

import type {
  Charge,
  CreateChargeInput,
  PaymentLinkResponse,
  PaymentQRResponse,
} from '@/generated/models';
import { AppShell } from '@/components/app-shell';
import { DataTable } from '@/components/data-table';
import {
  createCharge as createChargeRequest,
  generatePayLink,
  generatePayQR,
  listCharges,
} from '@/domain/charges';
import { Feedback } from '@/components/feedback';
import { SectionCard } from '@/components/section-card';
import { formatMoney } from '@/lib/format';
import { hasPermission } from '@/lib/auth';
import { useSession } from '@/lib/use-session';

export function ChargesPage() {
  const session = useSession();
  const permissions = session.user?.permissions ?? [];
  const canReadCharges = hasPermission(permissions, 'charges.read');
  const canWriteCharges = hasPermission(permissions, 'charges.write');

  const [items, setItems] = useState<Charge[]>([]);
  const [actionResult, setActionResult] = useState<string>('');
  const [success, setSuccess] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  async function load(): Promise<void> {
    if (!canReadCharges) {
      return;
    }
    try {
      setItems(await listCharges());
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao listar cobranças');
    }
  }

  useEffect(() => {
    if (session.loading || !canReadCharges) {
      return;
    }
    void load();
  }, [session.loading, canReadCharges]);

  async function createCharge(event: FormEvent<HTMLFormElement>): Promise<void> {
    event.preventDefault();
    if (!canWriteCharges) {
      return;
    }
    setError(null);
    const form = new FormData(event.currentTarget);
    const payload: CreateChargeInput = {
      customer_id: String(form.get('customer_id') ?? ''),
      contract_id: String(form.get('contract_id') ?? '') || undefined,
      charge_type: String(form.get('charge_type') ?? 'monthly'),
      amount_cents: Number(form.get('amount_cents') ?? 0),
      payment_method: String(form.get('payment_method') ?? 'pix'),
      description: String(form.get('description') ?? '') || undefined,
    };

    try {
      await createChargeRequest(payload);
      setSuccess('Cobrança criada com sucesso.');
      event.currentTarget.reset();
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao criar cobrança');
    }
  }

  async function requestPayLink(event: FormEvent<HTMLFormElement>): Promise<void> {
    event.preventDefault();
    if (!canWriteCharges) {
      return;
    }
    setError(null);
    const id = String(new FormData(event.currentTarget).get('charge_id') ?? '').trim();

    try {
      const result = await generatePayLink(id);
      setActionResult(JSON.stringify(result, null, 2));
      setSuccess('Link de pagamento gerado.');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao gerar link');
    }
  }

  async function requestPayQR(event: FormEvent<HTMLFormElement>): Promise<void> {
    event.preventDefault();
    if (!canWriteCharges) {
      return;
    }
    setError(null);
    const id = String(new FormData(event.currentTarget).get('charge_id') ?? '').trim();

    try {
      const result = await generatePayQR(id);
      setActionResult(JSON.stringify(result, null, 2));
      setSuccess('QR Code gerado.');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao gerar QR');
    }
  }

  return (
    <AppShell title="Cobranças">
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

      {canWriteCharges ? (
        <div className="grid gap-4 xl:grid-cols-3">
          <SectionCard title="Nova cobrança">
            <form className="grid gap-2" onSubmit={(event) => void createCharge(event)}>
              <input name="customer_id" placeholder="Customer ID" required />
              <input name="contract_id" placeholder="Contract ID" />
              <input defaultValue="monthly" name="charge_type" placeholder="Tipo" />
              <input name="amount_cents" placeholder="Valor (centavos)" type="number" required />
              <input defaultValue="pix" name="payment_method" placeholder="Pagamento" />
              <input name="description" placeholder="Descrição" />
              <button className="rounded-xl bg-primary px-4 py-2 font-semibold text-white" type="submit">
                Criar cobrança
              </button>
            </form>
          </SectionCard>

          <SectionCard title="Gerar link de pagamento">
            <form className="grid gap-2" onSubmit={(event) => void requestPayLink(event)}>
              <input name="charge_id" placeholder="Charge ID" required />
              <button className="rounded-xl bg-accent px-4 py-2 font-semibold text-white" type="submit">
                Gerar link
              </button>
            </form>
          </SectionCard>

          <SectionCard title="Gerar QR de pagamento">
            <form className="grid gap-2" onSubmit={(event) => void requestPayQR(event)}>
              <input name="charge_id" placeholder="Charge ID" required />
              <button className="rounded-xl bg-success px-4 py-2 font-semibold text-white" type="submit">
                Gerar QR
              </button>
            </form>
          </SectionCard>
        </div>
      ) : null}

      <section className="mt-4 rounded-2xl border border-slate-200 bg-slate-900 p-4 shadow-sm">
        <h3 className="mb-2 text-sm font-semibold text-slate-100">Resultado da ação</h3>
        <pre className="overflow-x-auto text-xs text-emerald-200">{actionResult || '{ }'}</pre>
      </section>

      {canReadCharges ? (
        <div className="mt-4">
          <DataTable
            rows={items as unknown as Record<string, unknown>[]}
            columns={[
              { key: 'id', label: 'ID' },
              { key: 'customer_id', label: 'Customer' },
              {
                key: 'amount_cents',
                label: 'Valor',
                render: (row) => formatMoney(Number((row as Charge).amount_cents ?? 0)),
              },
              { key: 'payment_method', label: 'Método' },
            ]}
          />
        </div>
      ) : (
        <SectionCard title="Cobranças">
          <p className="text-sm text-slate-600">
            Você não possui permissão de leitura para cobranças.
          </p>
        </SectionCard>
      )}
    </AppShell>
  );
}
