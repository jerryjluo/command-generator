import { CodeBlock } from '../../common/CodeBlock';

interface ResponseTabProps {
  command: string;
  explanation: string;
  rawResponse?: string;
}

export function ResponseTab({ command, explanation, rawResponse }: ResponseTabProps) {
  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-sm font-medium text-gray-700 mb-2">Generated Command</h3>
        <CodeBlock code={command || 'No command generated'} />
      </div>

      <div>
        <h3 className="text-sm font-medium text-gray-700 mb-2">Explanation</h3>
        <div className="bg-gray-50 rounded-lg p-4">
          <p className="text-gray-700 whitespace-pre-wrap">{explanation || 'No explanation provided'}</p>
        </div>
      </div>

      {rawResponse && (
        <div>
          <h3 className="text-sm font-medium text-gray-700 mb-2">Raw Response</h3>
          <CodeBlock code={rawResponse} title="Raw JSON Response" />
        </div>
      )}
    </div>
  );
}
