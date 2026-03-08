'use client';

import { FormEvent, useEffect, useState } from 'react';

import type {
  CreateMilestoneInput,
  CreateProjectInput,
  CreateWorklogInput,
  Project,
} from '@/generated/models';
import { AppShell } from '@/components/app-shell';
import { DataTable } from '@/components/data-table';
import {
  createMilestone as createMilestoneRequest,
  createProject as createProjectRequest,
  createWorklog as createWorklogRequest,
  listProjects,
} from '@/domain/projects';
import { Feedback } from '@/components/feedback';
import { SectionCard } from '@/components/section-card';
import { hasPermission } from '@/lib/auth';
import { useSession } from '@/lib/use-session';

export function ProjectsPage() {
  const session = useSession();
  const permissions = session.user?.permissions ?? [];
  const canReadProjects = hasPermission(permissions, 'projects.read');
  const canWriteProjects = hasPermission(permissions, 'projects.write');

  const [items, setItems] = useState<Project[]>([]);
  const [success, setSuccess] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  async function load(): Promise<void> {
    if (!canReadProjects) {
      return;
    }
    try {
      setItems(await listProjects());
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao listar projetos');
    }
  }

  useEffect(() => {
    if (session.loading || !canReadProjects) {
      return;
    }
    void load();
  }, [session.loading, canReadProjects]);

  async function createProject(event: FormEvent<HTMLFormElement>): Promise<void> {
    event.preventDefault();
    if (!canWriteProjects) {
      return;
    }
    const form = new FormData(event.currentTarget);
    const payload: CreateProjectInput = {
      contract_id: String(form.get('contract_id') ?? ''),
      name: String(form.get('name') ?? ''),
    };

    try {
      await createProjectRequest(payload);
      setSuccess('Projeto criado com sucesso.');
      event.currentTarget.reset();
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao criar projeto');
    }
  }

  async function createMilestone(event: FormEvent<HTMLFormElement>): Promise<void> {
    event.preventDefault();
    if (!canWriteProjects) {
      return;
    }
    const form = new FormData(event.currentTarget);
    const projectID = String(form.get('project_id') ?? '').trim();
    const payload: CreateMilestoneInput = {
      title: String(form.get('title') ?? ''),
      amount_cents: Number(form.get('amount_cents') ?? 0) || undefined,
    };

    try {
      await createMilestoneRequest(projectID, payload);
      setSuccess('Milestone criada com sucesso.');
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao criar milestone');
    }
  }

  async function createWorklog(event: FormEvent<HTMLFormElement>): Promise<void> {
    event.preventDefault();
    if (!canWriteProjects) {
      return;
    }
    const form = new FormData(event.currentTarget);
    const projectID = String(form.get('project_id') ?? '').trim();
    const payload: CreateWorklogInput = {
      hours: Number(form.get('hours') ?? 0),
      description: String(form.get('description') ?? '') || undefined,
    };

    try {
      await createWorklogRequest(projectID, payload);
      setSuccess('Worklog criado com sucesso.');
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao criar worklog');
    }
  }

  return (
    <AppShell title="Projetos">
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

      {canWriteProjects ? (
        <div className="grid gap-4 xl:grid-cols-3">
          <SectionCard title="Novo projeto">
            <form className="grid gap-2" onSubmit={(event) => void createProject(event)}>
              <input name="contract_id" placeholder="Contract ID" required />
              <input name="name" placeholder="Nome do projeto" required />
              <button className="rounded-xl bg-primary px-4 py-2 font-semibold text-white" type="submit">
                Criar projeto
              </button>
            </form>
          </SectionCard>

          <SectionCard title="Criar milestone">
            <form className="grid gap-2" onSubmit={(event) => void createMilestone(event)}>
              <input name="project_id" placeholder="Project ID" required />
              <input name="title" placeholder="Título milestone" required />
              <input name="amount_cents" placeholder="Valor (centavos)" type="number" />
              <button className="rounded-xl bg-accent px-4 py-2 font-semibold text-white" type="submit">
                Criar milestone
              </button>
            </form>
          </SectionCard>

          <SectionCard title="Criar worklog">
            <form className="grid gap-2" onSubmit={(event) => void createWorklog(event)}>
              <input name="project_id" placeholder="Project ID" required />
              <input name="hours" placeholder="Horas" type="number" step="0.5" required />
              <input name="description" placeholder="Descrição" />
              <button className="rounded-xl bg-success px-4 py-2 font-semibold text-white" type="submit">
                Criar worklog
              </button>
            </form>
          </SectionCard>
        </div>
      ) : null}

      {canReadProjects ? (
        <div className="mt-4">
          <DataTable
            rows={items as unknown as Record<string, unknown>[]}
            columns={[
              { key: 'id', label: 'ID' },
              { key: 'name', label: 'Nome' },
              { key: 'contract_id', label: 'Contrato' },
            ]}
          />
        </div>
      ) : (
        <SectionCard title="Projetos">
          <p className="text-sm text-slate-600">
            Você não possui permissão de leitura para projetos.
          </p>
        </SectionCard>
      )}
    </AppShell>
  );
}
