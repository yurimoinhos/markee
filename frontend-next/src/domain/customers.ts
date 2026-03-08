import type {
  CreateCustomerInput,
  Customer,
  CustomerDetail,
  UpdateCustomerInput,
} from '@/generated/models';
import { apiRequest } from '@/lib/http';

export async function listCustomers(): Promise<Customer[]> {
  return apiRequest<Customer[]>('/api/bff/customers', { cache: 'no-store' });
}

export async function createCustomer(payload: CreateCustomerInput): Promise<void> {
  await apiRequest('/api/bff/customers', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export async function updateCustomer(id: string, payload: UpdateCustomerInput): Promise<void> {
  await apiRequest(`/api/bff/customers/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload),
  });
}

export async function getCustomerDetail(id: string): Promise<CustomerDetail> {
  return apiRequest<CustomerDetail>(`/api/bff/customers/${id}`, { cache: 'no-store' });
}
