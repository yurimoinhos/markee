import type { AutomationRun, RunResult } from '@/generated/models';
import { apiRequest } from '@/lib/http';

export async function listAutomationRuns(): Promise<AutomationRun[]> {
  return apiRequest<AutomationRun[]>('/api/bff/automation/runs', { cache: 'no-store' });
}

export async function runAutomationNow(): Promise<RunResult> {
  return apiRequest<RunResult>('/api/bff/automation/run', {
    method: 'POST',
    body: JSON.stringify({}),
  });
}
