interface SortHeaderProps {
  label: string;
  field: string;
  currentSort?: string;
  currentOrder?: 'asc' | 'desc';
  onSort: (field: string) => void;
}

export function SortHeader({ label, field, currentSort, currentOrder, onSort }: SortHeaderProps) {
  const isActive = currentSort === field;

  return (
    <th
      className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 select-none"
      onClick={() => onSort(field)}
    >
      <div className="flex items-center gap-1">
        {label}
        <span className="text-gray-400">
          {isActive ? (
            currentOrder === 'asc' ? '↑' : '↓'
          ) : (
            '↕'
          )}
        </span>
      </div>
    </th>
  );
}
