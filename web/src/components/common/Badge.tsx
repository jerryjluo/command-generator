interface BadgeProps {
  variant: 'success' | 'error' | 'warning' | 'info';
  children: React.ReactNode;
}

const variantStyles = {
  success: 'bg-green-100 text-green-800',
  error: 'bg-red-100 text-red-800',
  warning: 'bg-yellow-100 text-yellow-800',
  info: 'bg-blue-100 text-blue-800',
};

export function Badge({ variant, children }: BadgeProps) {
  return (
    <span className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${variantStyles[variant]}`}>
      {children}
    </span>
  );
}

export function StatusBadge({ status }: { status: string }) {
  const variant = status === 'accepted' ? 'success' : status === 'rejected' ? 'error' : 'warning';
  return <Badge variant={variant}>{status}</Badge>;
}

export function ModelBadge({ model }: { model: string }) {
  return <Badge variant="info">{model}</Badge>;
}
