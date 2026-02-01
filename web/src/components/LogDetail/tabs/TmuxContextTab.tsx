import type { TmuxInfo } from '../../../types/log';
import { CodeBlock } from '../../common/CodeBlock';

interface TmuxContextTabProps {
  terminalContext: string;
  tmuxInfo: TmuxInfo;
}

export function TmuxContextTab({ terminalContext, tmuxInfo }: TmuxContextTabProps) {
  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-sm font-medium text-gray-700 mb-2">Tmux Session Info</h3>
        <div className="bg-gray-50 rounded-lg p-4">
          {tmuxInfo.in_tmux ? (
            <dl className="grid grid-cols-2 gap-4 text-sm">
              <div>
                <dt className="text-gray-500">Session</dt>
                <dd className="font-medium text-gray-900">{tmuxInfo.session || '-'}</dd>
              </div>
              <div>
                <dt className="text-gray-500">Window</dt>
                <dd className="font-medium text-gray-900">{tmuxInfo.window || '-'}</dd>
              </div>
              <div>
                <dt className="text-gray-500">Pane</dt>
                <dd className="font-medium text-gray-900">{tmuxInfo.pane || '-'}</dd>
              </div>
            </dl>
          ) : (
            <p className="text-gray-500 italic">Not running in tmux</p>
          )}
        </div>
      </div>

      <div>
        <h3 className="text-sm font-medium text-gray-700 mb-2">Terminal Scrollback</h3>
        {terminalContext ? (
          <CodeBlock code={terminalContext} title="Recent Terminal Output" />
        ) : (
          <p className="text-gray-500 italic">No terminal context captured</p>
        )}
      </div>
    </div>
  );
}
