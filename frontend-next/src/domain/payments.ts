import type { ConfirmPaymentInput, EvidenceInput, Payment } from '@/generated/models';
import { apiRequest } from '@/lib/http';

export async function listPayments(): Promise<Payment[]> {
  return apiRequest<Payment[]>('/api/bff/payments', { cache: 'no-store' });
}

export async function confirmPayment(payload: ConfirmPaymentInput): Promise<void> {
  await apiRequest('/api/bff/payments/confirm', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export async function addEvidence(id: string, payload: EvidenceInput): Promise<void> {
  await apiRequest(`/api/bff/payments/${id}/evidence`, {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}
