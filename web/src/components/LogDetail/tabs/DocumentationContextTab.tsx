import { CodeBlock } from '../../common/CodeBlock';

interface DocumentationContextTabProps {
  documentationContext: string;
}

export function DocumentationContextTab({ documentationContext }: DocumentationContextTabProps) {
  if (!documentationContext) {
    return (
      <div className="bg-gray-50 rounded-lg p-4">
        <p className="text-gray-500 italic">No documentation files found (README.md, CLAUDE.md, AGENTS.md)</p>
      </div>
    );
  }

  return (
    <div>
      <p className="text-sm text-gray-600 mb-4">
        Command-related sections extracted from project documentation files.
      </p>
      <CodeBlock code={documentationContext} title="Project Documentation" />
    </div>
  );
}
