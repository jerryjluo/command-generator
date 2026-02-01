import { useState, useCallback } from 'react';
import type { FilterParams } from '../types/log';

export interface UseFiltersResult {
  filters: FilterParams;
  setFilter: (key: keyof FilterParams, value: string | number | undefined) => void;
  clearFilters: () => void;
  setSort: (field: string) => void;
}

const defaultFilters: FilterParams = {
  sort: 'timestamp',
  order: 'desc',
};

export function useFilters(): UseFiltersResult {
  const [filters, setFilters] = useState<FilterParams>(defaultFilters);

  const setFilter = useCallback((key: keyof FilterParams, value: string | number | undefined) => {
    setFilters(prev => {
      if (value === undefined || value === '') {
        const { [key]: _, ...rest } = prev;
        return rest;
      }
      return { ...prev, [key]: value };
    });
  }, []);

  const clearFilters = useCallback(() => {
    setFilters(defaultFilters);
  }, []);

  const setSort = useCallback((field: string) => {
    setFilters(prev => {
      if (prev.sort === field) {
        // Toggle order if same field
        return { ...prev, order: prev.order === 'asc' ? 'desc' : 'asc' };
      }
      // New field, default to desc
      return { ...prev, sort: field, order: 'desc' };
    });
  }, []);

  return { filters, setFilter, clearFilters, setSort };
}
