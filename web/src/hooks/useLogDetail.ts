import { useState, useEffect, useCallback } from 'react';
import type { SessionLog } from '../types/log';
import { fetchLogById } from '../api/logs';

export interface UseLogDetailResult {
  log: SessionLog | null;
  loading: boolean;
  error: string | null;
  refetch: () => void;
}

export function useLogDetail(id: string | null): UseLogDetailResult {
  const [log, setLog] = useState<SessionLog | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadLog = useCallback(async () => {
    if (!id) {
      setLog(null);
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await fetchLogById(id);
      setLog(response);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load log');
    } finally {
      setLoading(false);
    }
  }, [id]);

  useEffect(() => {
    loadLog();
  }, [loadLog]);

  return { log, loading, error, refetch: loadLog };
}
