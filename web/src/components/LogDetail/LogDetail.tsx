import type { SessionLog } from '../../types/log';
import { StatusBadge, ModelBadge } from '../common/Badge';
import { TimeAgo } from '../common/TimeAgo';
import { TabPanel } from './TabPanel';
import { SystemPromptTab } from './tabs/SystemPromptTab';
import { UserPromptTab } from './tabs/UserPromptTab';
import { ResponseTab } from './tabs/ResponseTab';
import { TmuxContextTab } from './tabs/TmuxContextTab';
import { BuildToolsTab } from './tabs/BuildToolsTab';
import { PreferencesTab } from './tabs/PreferencesTab';

interface LogDetailProps {
  log: SessionLog;
  loading: boolean;
  error: string | null;
  onBack: () => void;
}

export function LogDetail({ log, loading, error, onBack }: LogDetailProps) {
  if (loading) {
    return (
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-8 text-center">
        <p className="text-gray-500">Loading log details...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-lg p-4">
        <p className="text-red-700">Error: {error}</p>
        <button
          onClick={onBack}
          className="mt-4 text-red-600 hover:text-red-800 font-medium"
        >
          ← Back to logs
        </button>
      </div>
    );
  }

  if (!log) {
    return null;
  }

  // Get the last iteration for display
  const lastIteration = log.iterations[log.iterations.length - 1];

  const tabs = [
    {
      id: 'response',
      label: 'Response',
      content: (
        <ResponseTab
          command={lastIteration?.model_output.command || ''}
          explanation={lastIteration?.model_output.explanation || ''}
          rawResponse={lastIteration?.model_output.raw_response}
        />
      ),
    },
    {
      id: 'system-prompt',
      label: 'System Prompt',
      content: <SystemPromptTab content={lastIteration?.model_input.system_prompt || ''} />,
    },
    {
      id: 'user-prompt',
      label: 'User Prompt',
      content: <UserPromptTab content={lastIteration?.model_input.user_prompt || ''} />,
    },
    {
      id: 'tmux-context',
      label: 'Tmux Context',
      content: (
        <TmuxContextTab
          terminalContext={log.context_sources.terminal_context}
          tmuxInfo={log.metadata.tmux_info}
        />
      ),
    },
    {
      id: 'build-tools',
      label: 'Build Tools',
      content: <BuildToolsTab userPrompt={lastIteration?.model_input.user_prompt || ''} />,
    },
    {
      id: 'preferences',
      label: 'Preferences',
      content: <PreferencesTab claudeMdContent={log.context_sources.claude_md_content} />,
    },
  ];

  return (
    <div>
      <button
        onClick={onBack}
        className="mb-4 text-blue-600 hover:text-blue-800 font-medium flex items-center gap-1"
      >
        ← Back to logs
      </button>

      <div className="bg-white rounded-lg shadow-sm border border-gray-200">
        {/* Header */}
        <div className="p-6 border-b border-gray-200">
          <div className="flex items-start justify-between gap-4">
            <div className="flex-1">
              <h2 className="text-xl font-semibold text-gray-900 mb-2">
                {log.user_query}
              </h2>
              <div className="flex items-center gap-4 text-sm text-gray-500">
                <TimeAgo timestamp={log.metadata.timestamp} />
                <span>•</span>
                <span>{log.metadata.iteration_count} iteration{log.metadata.iteration_count !== 1 ? 's' : ''}</span>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <StatusBadge status={log.metadata.final_status} />
              <ModelBadge model={log.metadata.model} />
            </div>
          </div>
        </div>

        {/* Iterations summary (if multiple) */}
        {log.iterations.length > 1 && (
          <div className="px-6 py-4 bg-gray-50 border-b border-gray-200">
            <h3 className="text-sm font-medium text-gray-700 mb-2">Iteration History</h3>
            <div className="space-y-2">
              {log.iterations.map((iter, index) => (
                <div key={index} className="flex items-center gap-4 text-sm">
                  <span className="text-gray-500">#{index + 1}</span>
                  {iter.feedback && (
                    <span className="text-gray-600 italic">Feedback: "{iter.feedback}"</span>
                  )}
                  <code className="text-xs bg-gray-200 px-2 py-1 rounded truncate max-w-md">
                    {iter.model_output.command}
                  </code>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Tabs */}
        <div className="p-6">
          <TabPanel tabs={tabs} />
        </div>
      </div>
    </div>
  );
}
