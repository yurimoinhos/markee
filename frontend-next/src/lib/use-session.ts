'use client';

import { useEffect, useState } from 'react';

import { APP_BASE_PATH } from '@/lib/config';
import type { SessionResponse, SessionUser } from '@/lib/auth';

type UseSessionResult = {
  loading: boolean;
  authenticated: boolean;
  user: SessionUser | null;
};

export function useSession(): UseSessionResult {
  const [loading, setLoading] = useState(true);
  const [authenticated, setAuthenticated] = useState(false);
  const [user, setUser] = useState<SessionUser | null>(null);

  useEffect(() => {
    let mounted = true;

    async function load(): Promise<void> {
      try {
        const response = await fetch(`${APP_BASE_PATH}/api/session`, { cache: 'no-store' });
        const payload = (await response.json().catch(() => null)) as SessionResponse | null;
        if (!mounted || !payload?.authenticated || !payload.user) {
          return;
        }
        setAuthenticated(true);
        setUser(payload.user);
      } finally {
        if (mounted) {
          setLoading(false);
        }
      }
    }

    void load();
    return () => {
      mounted = false;
    };
  }, []);

  return { loading, authenticated, user };
}

