import type { LogListResponse, SessionLog, FilterParams } from '../types/log';

const API_BASE = '/api/v1';

export async function fetchLogs(params: FilterParams = {}): Promise<LogListResponse> {
  const searchParams = new URLSearchParams();

  if (params.status) searchParams.set('status', params.status);
  if (params.model) searchParams.set('model', params.model);
  if (params.search) searchParams.set('search', params.search);
  if (params.from) searchParams.set('from', params.from);
  if (params.to) searchParams.set('to', params.to);
  if (params.sort) searchParams.set('sort', params.sort);
  if (params.order) searchParams.set('order', params.order);
  if (params.limit) searchParams.set('limit', params.limit.toString());
  if (params.offset) searchParams.set('offset', params.offset.toString());

  const url = `${API_BASE}/logs${searchParams.toString() ? '?' + searchParams.toString() : ''}`;
  const response = await fetch(url);

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to fetch logs');
  }

  return response.json();
}

export async function fetchLogById(id: string): Promise<SessionLog> {
  const response = await fetch(`${API_BASE}/logs/${id}`);

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to fetch log');
  }

  return response.json();
}
