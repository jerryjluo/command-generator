import { useState, useEffect, useCallback } from 'react';
import type { LogSummary, FilterParams } from '../types/log';
import { fetchLogs } from '../api/logs';

export interface UseLogsResult {
  logs: LogSummary[];
  total: number;
  loading: boolean;
  error: string | null;
  refetch: () => void;
}

export function useLogs(filters: FilterParams): UseLogsResult {
  const [logs, setLogs] = useState<LogSummary[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadLogs = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await fetchLogs(filters);
      setLogs(response.logs);
      setTotal(response.total);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load logs');
    } finally {
      setLoading(false);
    }
  }, [JSON.stringify(filters)]);

  useEffect(() => {
    loadLogs();
  }, [loadLogs]);

  return { logs, total, loading, error, refetch: loadLogs };
}
