import { useState, useEffect } from 'react';
import type { FilterParams } from '../../types/log';

interface FiltersProps {
  filters: FilterParams;
  onFilterChange: (key: keyof FilterParams, value: string | undefined) => void;
  onClear: () => void;
}

export function Filters({ filters, onFilterChange, onClear }: FiltersProps) {
  const [searchValue, setSearchValue] = useState(filters.search || '');

  // Debounce search input
  useEffect(() => {
    const timer = setTimeout(() => {
      onFilterChange('search', searchValue || undefined);
    }, 300);
    return () => clearTimeout(timer);
  }, [searchValue, onFilterChange]);

  // Sync local state with filters
  useEffect(() => {
    setSearchValue(filters.search || '');
  }, [filters.search]);

  const hasFilters = filters.status || filters.model || filters.search || filters.from || filters.to;

  return (
    <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-200 mb-4">
      <div className="flex flex-wrap gap-4 items-end">
        {/* Search */}
        <div className="flex-1 min-w-[200px]">
          <label className="block text-sm font-medium text-gray-700 mb-1">Search</label>
          <input
            type="text"
            value={searchValue}
            onChange={(e) => setSearchValue(e.target.value)}
            placeholder="Search queries, commands..."
            className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
          />
        </div>

        {/* Status filter */}
        <div className="w-32">
          <label className="block text-sm font-medium text-gray-700 mb-1">Status</label>
          <select
            value={filters.status || ''}
            onChange={(e) => onFilterChange('status', e.target.value || undefined)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="">All</option>
            <option value="accepted">Accepted</option>
            <option value="rejected">Rejected</option>
            <option value="quit">Quit</option>
          </select>
        </div>

        {/* Model filter */}
        <div className="w-32">
          <label className="block text-sm font-medium text-gray-700 mb-1">Model</label>
          <select
            value={filters.model || ''}
            onChange={(e) => onFilterChange('model', e.target.value || undefined)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="">All</option>
            <option value="opus">Opus</option>
            <option value="sonnet">Sonnet</option>
            <option value="haiku">Haiku</option>
          </select>
        </div>

        {/* Date range */}
        <div className="w-40">
          <label className="block text-sm font-medium text-gray-700 mb-1">From</label>
          <input
            type="date"
            value={filters.from?.split('T')[0] || ''}
            onChange={(e) => onFilterChange('from', e.target.value ? `${e.target.value}T00:00:00Z` : undefined)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
          />
        </div>

        <div className="w-40">
          <label className="block text-sm font-medium text-gray-700 mb-1">To</label>
          <input
            type="date"
            value={filters.to?.split('T')[0] || ''}
            onChange={(e) => onFilterChange('to', e.target.value ? `${e.target.value}T23:59:59Z` : undefined)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
          />
        </div>

        {/* Clear button */}
        {hasFilters && (
          <button
            onClick={onClear}
            className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-md"
          >
            Clear
          </button>
        )}
      </div>
    </div>
  );
}
