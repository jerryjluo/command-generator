import { CodeBlock } from '../../common/CodeBlock';

interface BuildToolsTabProps {
  userPrompt: string;
}

export function BuildToolsTab({ userPrompt }: BuildToolsTabProps) {
  // Extract build tools section from user prompt
  // Match the format: "Available build tools and commands in current directory:\n---\n...content...\n---"
  const buildToolsMatch = userPrompt.match(/Available build tools and commands in current directory:\n---\n([\s\S]*?)\n---/i);
  const buildToolsContent = buildToolsMatch?.[1]?.trim();

  if (!buildToolsContent) {
    return (
      <div className="bg-gray-50 rounded-lg p-4">
        <p className="text-gray-500 italic">No build tools detected in the current directory</p>
      </div>
    );
  }

  return (
    <div>
      <p className="text-sm text-gray-600 mb-4">
        Build tools and commands detected in the working directory.
      </p>
      <CodeBlock code={buildToolsContent} title="Detected Build Tools" />
    </div>
  );
}
