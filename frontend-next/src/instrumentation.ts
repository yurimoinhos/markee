import { ensureLocalStorageShim } from '@/lib/server/local-storage-shim';

export async function register(): Promise<void> {
  ensureLocalStorageShim();
}
