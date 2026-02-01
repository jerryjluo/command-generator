import { CodeBlock } from '../../common/CodeBlock';

interface UserQueryTabProps {
  content: string;
}

export function UserQueryTab({ content }: UserQueryTabProps) {
  if (!content) {
    return <p className="text-gray-500 italic">No user query available</p>;
  }

  return (
    <div>
      <p className="text-sm text-gray-600 mb-4">
        The original query from the user.
      </p>
      <CodeBlock code={content} title="User Query" />
    </div>
  );
}
