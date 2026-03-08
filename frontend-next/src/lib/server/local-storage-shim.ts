/**
 * Some server-side dependencies try to read globalThis.localStorage.
 * In non-browser runtimes this can be undefined or incompatible.
 * This shim avoids server crashes like "localStorage.getItem is not a function".
 */
type StorageLike = {
  getItem: (key: string) => string | null;
  setItem: (key: string, value: string) => void;
  removeItem: (key: string) => void;
  clear: () => void;
};

function createMemoryStorage(): StorageLike {
  const store = new Map<string, string>();
  return {
    getItem(key: string) {
      return store.has(key) ? store.get(key)! : null;
    },
    setItem(key: string, value: string) {
      store.set(key, value);
    },
    removeItem(key: string) {
      store.delete(key);
    },
    clear() {
      store.clear();
    },
  };
}

export function ensureLocalStorageShim(): void {
  if (typeof window !== 'undefined') return;

  const g = globalThis as {
    localStorage?: unknown;
  };

  const current = g.localStorage as Partial<StorageLike> | undefined;
  if (current && typeof current.getItem === 'function') return;

  const shim = createMemoryStorage();

  try {
    g.localStorage = shim;
    return;
  } catch {
    // fallback to defineProperty for read-only runtimes
  }

  try {
    Object.defineProperty(globalThis, 'localStorage', {
      configurable: true,
      enumerable: false,
      writable: true,
      value: shim,
    });
  } catch {
    // best effort; avoid throwing during bootstrap
  }
}

ensureLocalStorageShim();

export {};
