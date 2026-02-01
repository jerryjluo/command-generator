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
        The complete prompt sent to Claude including context, build tools, and user query.
      </p>
      <CodeBlock code={content} title="User Prompt" />
    </div>
  );
}
