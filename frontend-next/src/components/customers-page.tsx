'use client';

import { FormEvent, useEffect, useState } from 'react';

import type {
  CreateCustomerInput,
  Customer,
  CustomerDetail,
  UpdateCustomerInput,
} from '@/generated/models';
import { AppShell } from '@/components/app-shell';
import { DataTable } from '@/components/data-table';
import { Feedback } from '@/components/feedback';
import { SectionCard } from '@/components/section-card';
import {
  createCustomer as createCustomerRequest,
  getCustomerDetail,
  listCustomers,
  updateCustomer as updateCustomerRequest,
} from '@/domain/customers';
import { hasPermission } from '@/lib/auth';
import { useSession } from '@/lib/use-session';

export function CustomersPage() {
  const session = useSession();
  const permissions = session.user?.permissions ?? [];
  const canReadCustomers = hasPermission(permissions, 'customers.read');
  const canWriteCustomers = hasPermission(permissions, 'customers.write');

  const [items, setItems] = useState<Customer[]>([]);
  const [detail, setDetail] = useState<CustomerDetail | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function load(): Promise<void> {
    if (!canReadCustomers) {
      return;
    }
    setLoading(true);
    setError(null);
    try {
      setItems(await listCustomers());
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao carregar clientes');
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    if (session.loading || !canReadCustomers) {
      return;
    }
    void load();
  }, [session.loading, canReadCustomers]);

  async function createCustomer(event: FormEvent<HTMLFormElement>): Promise<void> {
    event.preventDefault();
    if (!canWriteCustomers) {
      return;
    }
    setError(null);
    const form = new FormData(event.currentTarget);
    const payload: CreateCustomerInput = {
      name: String(form.get('name') ?? ''),
      cpf_cnpj: String(form.get('cpf_cnpj') ?? ''),
      email: String(form.get('email') ?? ''),
      preferred_payment_method: String(form.get('preferred_payment_method') ?? 'pix'),
      company: String(form.get('company') ?? '') || undefined,
      phone: String(form.get('phone') ?? '') || undefined,
      address: String(form.get('address') ?? '') || undefined,
    };

    try {
      await createCustomerRequest(payload);
      setSuccess('Cliente criado com sucesso.');
      event.currentTarget.reset();
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao criar cliente');
    }
  }

  async function updateCustomer(event: FormEvent<HTMLFormElement>): Promise<void> {
    event.preventDefault();
    if (!canWriteCustomers) {
      return;
    }
    setError(null);
    const form = new FormData(event.currentTarget);
    const id = String(form.get('id') ?? '').trim();
    if (!id) {
      setError('Informe o ID do cliente para atualizar.');
      return;
    }

    const payload: UpdateCustomerInput = {
      name: String(form.get('name') ?? '') || undefined,
      email: String(form.get('email') ?? '') || undefined,
      phone: String(form.get('phone') ?? '') || undefined,
      preferred_payment_method: String(form.get('preferred_payment_method') ?? '') || undefined,
    };

    try {
      await updateCustomerRequest(id, payload);
      setSuccess('Cliente atualizado com sucesso.');
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao atualizar cliente');
    }
  }

  async function loadDetail(id: string): Promise<void> {
    if (!canReadCustomers) {
      return;
    }
    setError(null);
    try {
      setDetail(await getCustomerDetail(id));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao carregar detalhe');
    }
  }

  return (
    <AppShell title="Clientes">
      <Feedback success={success} error={error} />

      <div className="mb-4 flex justify-end">
        <button
          type="button"
          onClick={() => void load()}
          className="rounded-xl border border-slate-300 bg-white px-4 py-2 text-sm font-semibold"
        >
          {loading ? 'Atualizando...' : 'Atualizar'}
        </button>
      </div>

      {canWriteCustomers ? (
        <div className="grid gap-4 xl:grid-cols-2">
          <SectionCard title="Novo cliente">
            <form className="grid gap-2" onSubmit={(event) => void createCustomer(event)}>
              <input name="name" placeholder="Nome" required />
              <input name="cpf_cnpj" placeholder="CPF/CNPJ" required />
              <input name="email" placeholder="Email" required type="email" />
              <input name="company" placeholder="Empresa" />
              <input name="phone" placeholder="Telefone" />
              <input name="address" placeholder="Endereço" />
              <input name="preferred_payment_method" placeholder="Forma de pagamento" defaultValue="pix" />
              <button className="rounded-xl bg-primary px-4 py-2 font-semibold text-white" type="submit">
                Criar cliente
              </button>
            </form>
          </SectionCard>

          <SectionCard title="Atualizar cliente">
            <form className="grid gap-2" onSubmit={(event) => void updateCustomer(event)}>
              <input name="id" placeholder="ID do cliente" required />
              <input name="name" placeholder="Nome" />
              <input name="email" placeholder="Email" type="email" />
              <input name="phone" placeholder="Telefone" />
              <input name="preferred_payment_method" placeholder="Forma de pagamento" />
              <button className="rounded-xl bg-accent px-4 py-2 font-semibold text-white" type="submit">
                Atualizar cliente
              </button>
            </form>
          </SectionCard>
        </div>
      ) : null}

      {detail ? (
        <section className="mt-4 rounded-2xl border border-slate-200 bg-white p-4 shadow-sm">
          <h3 className="mb-2 text-base font-semibold text-slate-900">Detalhe selecionado</h3>
          <p className="text-sm text-slate-700">
            {detail.customer.name} ({detail.customer.email})
          </p>
          <div className="mt-2 flex flex-wrap gap-4 text-sm text-slate-600">
            <span>Contratos: {detail.financial_summary.contracts_count}</span>
            <span>Pendentes: {detail.financial_summary.pending_charges}</span>
            <span>Pagas: {detail.financial_summary.paid_charges}</span>
          </div>
        </section>
      ) : null}

      {canReadCustomers ? (
        <div className="mt-4">
          <DataTable
            rows={items as unknown as Record<string, unknown>[]}
            columns={[
              { key: 'id', label: 'ID' },
              { key: 'name', label: 'Nome' },
              { key: 'email', label: 'Email' },
              { key: 'preferred_payment_method', label: 'Pagamento' },
              {
                key: 'actions',
                label: 'Ações',
                render: (row) => (
                  <button
                    type="button"
                    className="rounded-lg border border-slate-300 px-2 py-1 text-xs"
                    onClick={() => void loadDetail(String((row as Customer).id ?? ''))}
                  >
                    Detalhar
                  </button>
                ),
              },
            ]}
          />
        </div>
      ) : (
        <SectionCard title="Clientes">
          <p className="text-sm text-slate-600">
            Você não possui permissão de leitura para clientes.
          </p>
        </SectionCard>
      )}
    </AppShell>
  );
}
