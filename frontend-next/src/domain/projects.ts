import type {
  CreateMilestoneInput,
  CreateProjectInput,
  CreateWorklogInput,
  Project,
} from '@/generated/models';
import { apiRequest } from '@/lib/http';

export async function listProjects(): Promise<Project[]> {
  return apiRequest<Project[]>('/api/bff/projects', { cache: 'no-store' });
}

export async function createProject(payload: CreateProjectInput): Promise<void> {
  await apiRequest('/api/bff/projects', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export async function createMilestone(projectID: string, payload: CreateMilestoneInput): Promise<void> {
  await apiRequest(`/api/bff/projects/${projectID}/milestones`, {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export async function createWorklog(projectID: string, payload: CreateWorklogInput): Promise<void> {
  await apiRequest(`/api/bff/projects/${projectID}/worklogs`, {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}
