import type {
  Charge,
  CreateChargeInput,
  PaymentLinkResponse,
  PaymentQRResponse,
} from '@/generated/models';
import { apiRequest } from '@/lib/http';

export async function listCharges(): Promise<Charge[]> {
  return apiRequest<Charge[]>('/api/bff/charges', { cache: 'no-store' });
}

export async function createCharge(payload: CreateChargeInput): Promise<void> {
  await apiRequest('/api/bff/charges', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export async function generatePayLink(id: string): Promise<PaymentLinkResponse> {
  return apiRequest<PaymentLinkResponse>(`/api/bff/charges/${id}/pay-link`, {
    method: 'POST',
    body: JSON.stringify({}),
  });
}

export async function generatePayQR(id: string): Promise<PaymentQRResponse> {
  return apiRequest<PaymentQRResponse>(`/api/bff/charges/${id}/pay-qr`, {
    method: 'POST',
    body: JSON.stringify({}),
  });
}
