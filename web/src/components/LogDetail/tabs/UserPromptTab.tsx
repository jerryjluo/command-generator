import { CodeBlock } from '../../common/CodeBlock';

interface UserPromptTabProps {
  content: string;
}

export function UserPromptTab({ content }: UserPromptTabProps) {
  if (!content) {
    return <p className="text-gray-500 italic">No user prompt available</p>;
  }

  return (
    <div>
      <p className="text-sm text-gray-600 mb-4">
        The complete user prompt including query, context, and instructions.
      </p>
      <CodeBlock code={content} title="User Prompt" />
    </div>
  );
}
