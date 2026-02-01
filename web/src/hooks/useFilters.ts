import { useSearchParams } from 'react-router-dom';
import { useCallback, useMemo } from 'react';
import type { FilterParams } from '../types/log';

export interface UseFiltersResult {
  filters: FilterParams;
  setFilter: (key: keyof FilterParams, value: string | number | undefined) => void;
  clearFilters: () => void;
  setSort: (field: string) => void;
}

export function useFilters(): UseFiltersResult {
  const [searchParams, setSearchParams] = useSearchParams();

  const filters = useMemo((): FilterParams => ({
    sort: searchParams.get('sort') || 'timestamp',
    order: (searchParams.get('order') as 'asc' | 'desc') || 'desc',
    status: searchParams.get('status') || undefined,
    model: searchParams.get('model') || undefined,
    search: searchParams.get('search') || undefined,
    from: searchParams.get('from') || undefined,
    to: searchParams.get('to') || undefined,
    limit: searchParams.get('limit') ? Number(searchParams.get('limit')) : undefined,
    offset: searchParams.get('offset') ? Number(searchParams.get('offset')) : undefined,
  }), [searchParams]);

  const setFilter = useCallback((key: keyof FilterParams, value: string | number | undefined) => {
    setSearchParams(prev => {
      const params = new URLSearchParams(prev);
      if (value === undefined || value === '') {
        params.delete(key);
      } else {
        params.set(key, String(value));
      }
      // Reset offset when filters change
      if (key !== 'offset' && key !== 'limit') {
        params.delete('offset');
      }
      return params;
    }, { replace: true });
  }, [setSearchParams]);

  const clearFilters = useCallback(() => {
    setSearchParams({}, { replace: true });
  }, [setSearchParams]);

  const setSort = useCallback((field: string) => {
    setSearchParams(prev => {
      const params = new URLSearchParams(prev);
      const currentSort = params.get('sort') || 'timestamp';
      const currentOrder = params.get('order') || 'desc';

      if (currentSort === field) {
        params.set('order', currentOrder === 'asc' ? 'desc' : 'asc');
      } else {
        params.set('sort', field);
        params.set('order', 'desc');
      }
      return params;
    }, { replace: true });
  }, [setSearchParams]);

  return { filters, setFilter, clearFilters, setSort };
}
