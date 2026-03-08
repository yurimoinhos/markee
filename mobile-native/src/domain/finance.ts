import type { Dashboard } from '../generated/api';
import { request } from '../lib/http';

export async function getDashboard(): Promise<Dashboard> {
  return request<Dashboard>('/finance/dashboard');
}
