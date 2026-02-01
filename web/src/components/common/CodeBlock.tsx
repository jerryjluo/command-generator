import { useState } from 'react';

interface CodeBlockProps {
  code: string;
  language?: string;
  title?: string;
}

export function CodeBlock({ code, title }: CodeBlockProps) {
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(code);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {
      // Clipboard API not available
    }
  };

  return (
    <div className="relative group">
      {title && (
        <div className="bg-gray-800 text-gray-300 px-4 py-2 text-sm font-medium rounded-t-lg">
          {title}
        </div>
      )}
      <div className={`relative ${title ? '' : 'rounded-lg'}`}>
        <pre className={`bg-gray-900 text-gray-100 p-4 overflow-x-auto text-sm ${title ? 'rounded-b-lg' : 'rounded-lg'}`}>
          <code>{code}</code>
        </pre>
        <button
          onClick={handleCopy}
          className="absolute top-2 right-2 bg-gray-700 hover:bg-gray-600 text-gray-300 px-2 py-1 rounded text-xs opacity-0 group-hover:opacity-100 transition-opacity"
        >
          {copied ? 'Copied!' : 'Copy'}
        </button>
      </div>
    </div>
  );
}
