import { useState } from 'react';

interface Tab {
  id: string;
  label: string;
  content: React.ReactNode;
}

interface TabPanelProps {
  tabs: Tab[];
}

export function TabPanel({ tabs }: TabPanelProps) {
  const [activeTab, setActiveTab] = useState(tabs[0]?.id || '');

  const activeContent = tabs.find((tab) => tab.id === activeTab)?.content;

  return (
    <div>
      <div className="border-b border-gray-200">
        <nav className="flex -mb-px space-x-4" aria-label="Tabs">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`py-3 px-4 text-sm font-medium border-b-2 ${
                activeTab === tab.id
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }`}
            >
              {tab.label}
            </button>
          ))}
        </nav>
      </div>
      <div className="py-4">{activeContent}</div>
    </div>
  );
}
