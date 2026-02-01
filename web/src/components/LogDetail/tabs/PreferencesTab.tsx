import { CodeBlock } from '../../common/CodeBlock';

interface PreferencesTabProps {
  claudeMdContent: string;
}

export function PreferencesTab({ claudeMdContent }: PreferencesTabProps) {
  if (!claudeMdContent) {
    return (
      <div className="bg-gray-50 rounded-lg p-4">
        <p className="text-gray-500 italic">No user preferences configured (claude.md)</p>
      </div>
    );
  }

  return (
    <div>
      <p className="text-sm text-gray-600 mb-4">
        User preferences from ~/.config/cmd/claude.md
      </p>
      <CodeBlock code={claudeMdContent} title="claude.md" />
    </div>
  );
}
