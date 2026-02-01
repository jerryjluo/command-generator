import { Routes, Route, useNavigate, useParams } from 'react-router-dom';
import { Layout } from './components/Layout';
import { LogTable } from './components/LogTable';
import { LogDetail } from './components/LogDetail';
import { useLogs } from './hooks/useLogs';
import { useLogDetail } from './hooks/useLogDetail';
import { useFilters } from './hooks/useFilters';

function LogListPage() {
  const navigate = useNavigate();
  const { filters, setFilter, clearFilters, setSort } = useFilters();
  const { logs, total, loading, error } = useLogs(filters);

  const handleSelectLog = (id: string) => {
    navigate(`/logs/${id}`);
  };

  return (
    <LogTable
      logs={logs}
      total={total}
      loading={loading}
      error={error}
      filters={filters}
      onFilterChange={setFilter}
      onClearFilters={clearFilters}
      onSort={setSort}
      onSelectLog={handleSelectLog}
    />
  );
}

function LogDetailPage() {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const { log, loading, error } = useLogDetail(id || null);

  const handleBack = () => {
    navigate('/');
  };

  return (
    <LogDetail
      log={log!}
      loading={loading}
      error={error}
      onBack={handleBack}
    />
  );
}

function App() {
  return (
    <Layout>
      <Routes>
        <Route path="/" element={<LogListPage />} />
        <Route path="/logs/:id" element={<LogDetailPage />} />
      </Routes>
    </Layout>
  );
}

export default App;
