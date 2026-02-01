import { useState } from 'react';
import { Layout } from './components/Layout';
import { LogTable } from './components/LogTable';
import { LogDetail } from './components/LogDetail';
import { useLogs } from './hooks/useLogs';
import { useLogDetail } from './hooks/useLogDetail';
import { useFilters } from './hooks/useFilters';

function App() {
  const [selectedLogId, setSelectedLogId] = useState<string | null>(null);
  const { filters, setFilter, clearFilters, setSort } = useFilters();
  const { logs, total, loading: logsLoading, error: logsError } = useLogs(filters);
  const { log, loading: logLoading, error: logError } = useLogDetail(selectedLogId);

  const handleSelectLog = (id: string) => {
    setSelectedLogId(id);
  };

  const handleBack = () => {
    setSelectedLogId(null);
  };

  return (
    <Layout>
      {selectedLogId ? (
        <LogDetail
          log={log!}
          loading={logLoading}
          error={logError}
          onBack={handleBack}
        />
      ) : (
        <LogTable
          logs={logs}
          total={total}
          loading={logsLoading}
          error={logsError}
          filters={filters}
          onFilterChange={setFilter}
          onClearFilters={clearFilters}
          onSort={setSort}
          onSelectLog={handleSelectLog}
        />
      )}
    </Layout>
  );
}

export default App;
