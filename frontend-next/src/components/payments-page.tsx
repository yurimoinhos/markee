'use client';

import { FormEvent, useEffect, useState } from 'react';

import type {
  ConfirmPaymentInput,
  EvidenceInput,
  Payment,
} from '@/generated/models';
import { AppShell } from '@/components/app-shell';
import { DataTable } from '@/components/data-table';
import {
  addEvidence as addEvidenceRequest,
  confirmPayment as confirmPaymentRequest,
  listPayments,
} from '@/domain/payments';
import { Feedback } from '@/components/feedback';
import { SectionCard } from '@/components/section-card';
import { formatMoney } from '@/lib/format';
import { hasPermission } from '@/lib/auth';
import { useSession } from '@/lib/use-session';

export function PaymentsPage() {
  const session = useSession();
  const permissions = session.user?.permissions ?? [];
  const canReadPayments = hasPermission(permissions, 'payments.read');
  const canWritePayments = hasPermission(permissions, 'payments.write');

  const [items, setItems] = useState<Payment[]>([]);
  const [success, setSuccess] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  async function load(): Promise<void> {
    if (!canReadPayments) {
      return;
    }
    try {
      setItems(await listPayments());
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao listar pagamentos');
    }
  }

  useEffect(() => {
    if (session.loading || !canReadPayments) {
      return;
    }
    void load();
  }, [session.loading, canReadPayments]);

  async function confirmPayment(event: FormEvent<HTMLFormElement>): Promise<void> {
    event.preventDefault();
    if (!canWritePayments) {
      return;
    }
    setError(null);
    const form = new FormData(event.currentTarget);
    const payload: ConfirmPaymentInput = {
      charge_id: String(form.get('charge_id') ?? ''),
      amount_cents: Number(form.get('amount_cents') ?? 0),
      method: String(form.get('method') ?? 'pix'),
      tx_hash: String(form.get('tx_hash') ?? '') || undefined,
    };

    try {
      await confirmPaymentRequest(payload);
      setSuccess('Pagamento confirmado com sucesso.');
      event.currentTarget.reset();
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao confirmar pagamento');
    }
  }

  async function addEvidence(event: FormEvent<HTMLFormElement>): Promise<void> {
    event.preventDefault();
    if (!canWritePayments) {
      return;
    }
    setError(null);
    const form = new FormData(event.currentTarget);
    const paymentID = String(form.get('payment_id') ?? '').trim();
    const payload: EvidenceInput = {
      file_url: String(form.get('file_url') ?? '') || undefined,
      note: String(form.get('note') ?? '') || undefined,
      tx_hash: String(form.get('tx_hash') ?? '') || undefined,
    };

    try {
      await addEvidenceRequest(paymentID, payload);
      setSuccess('Evidência registrada com sucesso.');
      event.currentTarget.reset();
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao adicionar evidência');
    }
  }

  return (
    <AppShell title="Pagamentos">
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

      {canWritePayments ? (
        <div className="grid gap-4 xl:grid-cols-2">
          <SectionCard title="Confirmar pagamento">
            <form className="grid gap-2" onSubmit={(event) => void confirmPayment(event)}>
              <input name="charge_id" placeholder="Charge ID" required />
              <input name="amount_cents" placeholder="Valor (centavos)" type="number" required />
              <input defaultValue="pix" name="method" placeholder="Método" />
              <input name="tx_hash" placeholder="Hash da transação" />
              <button className="rounded-xl bg-primary px-4 py-2 font-semibold text-white" type="submit">
                Confirmar pagamento
              </button>
            </form>
          </SectionCard>

          <SectionCard title="Adicionar evidência">
            <form className="grid gap-2" onSubmit={(event) => void addEvidence(event)}>
              <input name="payment_id" placeholder="Payment ID" required />
              <input name="file_url" placeholder="URL do arquivo" />
              <input name="note" placeholder="Observação" />
              <input name="tx_hash" placeholder="Hash da transação" />
              <button className="rounded-xl bg-accent px-4 py-2 font-semibold text-white" type="submit">
                Adicionar evidência
              </button>
            </form>
          </SectionCard>
        </div>
      ) : null}

      {canReadPayments ? (
        <div className="mt-4">
          <DataTable
            rows={items as unknown as Record<string, unknown>[]}
            columns={[
              { key: 'id', label: 'ID' },
              { key: 'charge_id', label: 'Charge' },
              {
                key: 'amount_cents',
                label: 'Valor',
                render: (row) => formatMoney(Number((row as Payment).amount_cents ?? 0)),
              },
              { key: 'method', label: 'Método' },
            ]}
          />
        </div>
      ) : (
        <SectionCard title="Pagamentos">
          <p className="text-sm text-slate-600">
            Você não possui permissão de leitura para pagamentos.
          </p>
        </SectionCard>
      )}
    </AppShell>
  );
}
