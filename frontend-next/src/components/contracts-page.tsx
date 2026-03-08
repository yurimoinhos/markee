'use client';

import { FormEvent, useEffect, useState } from 'react';

import type {
  Contract,
  CreateContractInput,
  GenerateContractInput,
  SendSignatureInput,
} from '@/generated/models';
import { AppShell } from '@/components/app-shell';
import { DataTable } from '@/components/data-table';
import {
  createContract as createContractRequest,
  generateContractVersion,
  listContracts,
  sendContractSignature,
} from '@/domain/contracts';
import { Feedback } from '@/components/feedback';
import { SectionCard } from '@/components/section-card';
import { formatMoney } from '@/lib/format';
import { hasPermission } from '@/lib/auth';
import { useSession } from '@/lib/use-session';

export function ContractsPage() {
  const session = useSession();
  const permissions = session.user?.permissions ?? [];
  const canReadContracts = hasPermission(permissions, 'contracts.read');
  const canWriteContracts = hasPermission(permissions, 'contracts.write');

  const [items, setItems] = useState<Contract[]>([]);
  const [success, setSuccess] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  async function load(): Promise<void> {
    if (!canReadContracts) {
      return;
    }
    try {
      setItems(await listContracts());
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao listar contratos');
    }
  }

  useEffect(() => {
    if (session.loading || !canReadContracts) {
      return;
    }
    void load();
  }, [session.loading, canReadContracts]);

  async function createContract(event: FormEvent<HTMLFormElement>): Promise<void> {
    event.preventDefault();
    if (!canWriteContracts) {
      return;
    }
    setError(null);
    const form = new FormData(event.currentTarget);
    const payload: CreateContractInput = {
      customer_id: String(form.get('customer_id') ?? ''),
      title: String(form.get('title') ?? ''),
      contract_type: String(form.get('contract_type') ?? 'service'),
      billing_type: String(form.get('billing_type') ?? 'monthly'),
      amount_cents: Number(form.get('amount_cents') ?? 0),
      auto_renew: Boolean(form.get('auto_renew')),
    };

    try {
      await createContractRequest(payload);
      setSuccess('Contrato criado com sucesso.');
      event.currentTarget.reset();
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao criar contrato');
    }
  }

  async function generateVersion(event: FormEvent<HTMLFormElement>): Promise<void> {
    event.preventDefault();
    if (!canWriteContracts) {
      return;
    }
    setError(null);
    const form = new FormData(event.currentTarget);
    const id = String(form.get('contract_id') ?? '').trim();
    const payload: GenerateContractInput = {
      template_name: String(form.get('template_name') ?? 'standard'),
      editable_content: String(form.get('editable_content') ?? ''),
    };

    try {
      await generateContractVersion(id, payload);
      setSuccess('Versão de contrato gerada.');
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao gerar contrato');
    }
  }

  async function sendSignature(event: FormEvent<HTMLFormElement>): Promise<void> {
    event.preventDefault();
    if (!canWriteContracts) {
      return;
    }
    setError(null);
    const form = new FormData(event.currentTarget);
    const id = String(form.get('contract_id') ?? '').trim();
    const payload: SendSignatureInput = {
      signer_name: String(form.get('signer_name') ?? ''),
      signer_email: String(form.get('signer_email') ?? ''),
    };

    try {
      await sendContractSignature(id, payload);
      setSuccess('Envio para assinatura concluído.');
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao enviar assinatura');
    }
  }

  return (
    <AppShell title="Contratos">
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

      {canWriteContracts ? (
        <div className="grid gap-4 xl:grid-cols-3">
          <SectionCard title="Criar contrato">
            <form className="grid gap-2" onSubmit={(event) => void createContract(event)}>
              <input name="customer_id" placeholder="Customer ID" required />
              <input name="title" placeholder="Título" required />
              <input defaultValue="service" name="contract_type" placeholder="Tipo" required />
              <input defaultValue="monthly" name="billing_type" placeholder="Billing" required />
              <input name="amount_cents" placeholder="Valor (centavos)" type="number" required />
              <label className="inline-flex items-center gap-2 text-sm text-slate-700">
                <input name="auto_renew" type="checkbox" className="h-4 w-4" /> Auto renovação
              </label>
              <button className="rounded-xl bg-primary px-4 py-2 font-semibold text-white" type="submit">
                Criar
              </button>
            </form>
          </SectionCard>

          <SectionCard title="Gerar versão">
            <form className="grid gap-2" onSubmit={(event) => void generateVersion(event)}>
              <input name="contract_id" placeholder="Contract ID" required />
              <input defaultValue="standard" name="template_name" placeholder="Template" required />
              <textarea name="editable_content" placeholder="Conteúdo editável" rows={5} required />
              <button className="rounded-xl bg-accent px-4 py-2 font-semibold text-white" type="submit">
                Gerar versão
              </button>
            </form>
          </SectionCard>

          <SectionCard title="Enviar assinatura">
            <form className="grid gap-2" onSubmit={(event) => void sendSignature(event)}>
              <input name="contract_id" placeholder="Contract ID" required />
              <input name="signer_name" placeholder="Nome do signatário" required />
              <input name="signer_email" placeholder="Email do signatário" required type="email" />
              <button className="rounded-xl bg-success px-4 py-2 font-semibold text-white" type="submit">
                Enviar assinatura
              </button>
            </form>
          </SectionCard>
        </div>
      ) : null}

      {canReadContracts ? (
        <div className="mt-4">
          <DataTable
            rows={items as unknown as Record<string, unknown>[]}
            columns={[
              { key: 'id', label: 'ID' },
              { key: 'title', label: 'Título' },
              {
                key: 'amount_cents',
                label: 'Valor',
                render: (row) => formatMoney(Number((row as Contract).amount_cents ?? 0)),
              },
              { key: 'billing_type', label: 'Billing' },
            ]}
          />
        </div>
      ) : (
        <SectionCard title="Contratos">
          <p className="text-sm text-slate-600">
            Você não possui permissão de leitura para contratos.
          </p>
        </SectionCard>
      )}
    </AppShell>
  );
}
