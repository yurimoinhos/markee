import type {
  Contract,
  CreateContractInput,
  GenerateContractInput,
  SendSignatureInput,
} from '@/generated/models';
import { apiRequest } from '@/lib/http';

export async function listContracts(): Promise<Contract[]> {
  return apiRequest<Contract[]>('/api/bff/contracts', { cache: 'no-store' });
}

export async function createContract(payload: CreateContractInput): Promise<void> {
  await apiRequest('/api/bff/contracts', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export async function generateContractVersion(id: string, payload: GenerateContractInput): Promise<void> {
  await apiRequest(`/api/bff/contracts/${id}/generate`, {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export async function sendContractSignature(id: string, payload: SendSignatureInput): Promise<void> {
  await apiRequest(`/api/bff/contracts/${id}/send-signature`, {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}
