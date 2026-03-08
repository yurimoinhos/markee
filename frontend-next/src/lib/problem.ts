import type { ProblemDetail } from '@/generated/models';

export function extractProblemMessage(payload: unknown, fallback = 'Operação falhou'): string {
  if (!payload || typeof payload !== 'object') {
    return fallback;
  }

  const p = payload as ProblemDetail;
  if (typeof p.errorDescription === 'string' && p.errorDescription.trim() !== '') {
    return p.errorDescription;
  }
  if (typeof p.error === 'string' && p.error.trim() !== '') {
    return p.error;
  }
  return fallback;
}
