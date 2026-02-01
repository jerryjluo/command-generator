import { useState } from 'react';

interface TimeAgoProps {
  timestamp: string;
}

function formatTimeAgo(timestamp: string): string {
  const date = new Date(timestamp);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffSec = Math.floor(diffMs / 1000);
  const diffMin = Math.floor(diffSec / 60);
  const diffHour = Math.floor(diffMin / 60);
  const diffDay = Math.floor(diffHour / 24);

  if (diffSec < 60) return 'just now';
  if (diffMin < 60) return `${diffMin}m ago`;
  if (diffHour < 24) return `${diffHour}h ago`;
  if (diffDay < 7) return `${diffDay}d ago`;

  return date.toLocaleDateString();
}

function formatAbsolute(timestamp: string): string {
  const date = new Date(timestamp);
  return date.toLocaleString();
}

export function TimeAgo({ timestamp }: TimeAgoProps) {
  const [showAbsolute, setShowAbsolute] = useState(false);

  return (
    <span
      className="cursor-pointer hover:underline"
      onClick={() => setShowAbsolute(!showAbsolute)}
      title={showAbsolute ? 'Click for relative time' : 'Click for absolute time'}
    >
      {showAbsolute ? formatAbsolute(timestamp) : formatTimeAgo(timestamp)}
    </span>
  );
}
