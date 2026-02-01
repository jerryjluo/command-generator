import { CodeBlock } from '../../common/CodeBlock';

interface SystemPromptTabProps {
  content: string;
}

export function SystemPromptTab({ content }: SystemPromptTabProps) {
  if (!content) {
    return <p className="text-gray-500 italic">No system prompt available</p>;
  }

  return (
    <div>
      <p className="text-sm text-gray-600 mb-4">
        The system prompt sent to Claude to configure its behavior.
      </p>
      <CodeBlock code={content} title="System Prompt" />
    </div>
  );
}
