import type { LogSummary, FilterParams } from '../../types/log';
import { StatusBadge, ModelBadge } from '../common/Badge';
import { TimeAgo } from '../common/TimeAgo';
import { Filters } from './Filters';
import { SortHeader } from './SortHeader';

interface LogTableProps {
  logs: LogSummary[];
  total: number;
  loading: boolean;
  error: string | null;
  filters: FilterParams;
  onFilterChange: (key: keyof FilterParams, value: string | number | undefined) => void;
  onClearFilters: () => void;
  onSort: (field: string) => void;
  onSelectLog: (id: string) => void;
}

export function LogTable({
  logs,
  total,
  loading,
  error,
  filters,
  onFilterChange,
  onClearFilters,
  onSort,
  onSelectLog,
}: LogTableProps) {
  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-red-700">
        Error: {error}
      </div>
    );
  }

  return (
    <div>
      <Filters
        filters={filters}
        onFilterChange={onFilterChange}
        onClear={onClearFilters}
      />

      <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
        <div className="px-4 py-3 border-b border-gray-200 flex justify-between items-center">
          <span className="text-sm text-gray-600">
            {loading ? 'Loading...' : `${total} log${total !== 1 ? 's' : ''}`}
          </span>
        </div>

        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <SortHeader
                  label="Query"
                  field="query"
                  currentSort={filters.sort}
                  currentOrder={filters.order}
                  onSort={onSort}
                />
                <SortHeader
                  label="Status"
                  field="status"
                  currentSort={filters.sort}
                  currentOrder={filters.order}
                  onSort={onSort}
                />
                <SortHeader
                  label="Model"
                  field="model"
                  currentSort={filters.sort}
                  currentOrder={filters.order}
                  onSort={onSort}
                />
                <SortHeader
                  label="Time"
                  field="timestamp"
                  currentSort={filters.sort}
                  currentOrder={filters.order}
                  onSort={onSort}
                />
                <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Command
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {loading ? (
                <tr>
                  <td colSpan={5} className="px-4 py-8 text-center text-gray-500">
                    Loading...
                  </td>
                </tr>
              ) : logs.length === 0 ? (
                <tr>
                  <td colSpan={5} className="px-4 py-8 text-center text-gray-500">
                    No logs found
                  </td>
                </tr>
              ) : (
                logs.map((log) => (
                  <tr
                    key={log.id}
                    className="hover:bg-gray-50 cursor-pointer"
                    onClick={() => onSelectLog(log.id)}
                  >
                    <td className="px-4 py-4">
                      <div className="text-sm font-medium text-gray-900 max-w-md truncate">
                        {log.user_query}
                      </div>
                    </td>
                    <td className="px-4 py-4">
                      <StatusBadge status={log.final_status} />
                    </td>
                    <td className="px-4 py-4">
                      <ModelBadge model={log.model} />
                    </td>
                    <td className="px-4 py-4 text-sm text-gray-500">
                      <TimeAgo timestamp={log.timestamp} />
                    </td>
                    <td className="px-4 py-4">
                      <code className="text-sm text-gray-600 bg-gray-100 px-2 py-1 rounded max-w-xs block truncate">
                        {log.command_preview || '-'}
                      </code>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
