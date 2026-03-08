import type { CashFlowPoint, Dashboard, DefaultMetrics } from '@/generated/models';
import { apiRequest } from '@/lib/http';

export async function getDashboard(): Promise<Dashboard> {
  return apiRequest<Dashboard>('/api/bff/finance/dashboard', { cache: 'no-store' });
}

export async function getCashflow(): Promise<CashFlowPoint[]> {
  return apiRequest<CashFlowPoint[]>('/api/bff/finance/cashflow', { cache: 'no-store' });
}

export async function getDefaultMetrics(): Promise<DefaultMetrics> {
  return apiRequest<DefaultMetrics>('/api/bff/finance/defaults', { cache: 'no-store' });
}
