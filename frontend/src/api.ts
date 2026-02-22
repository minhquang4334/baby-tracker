import type { Child, SleepLog, FeedingLog, DiaperLog, GrowthLog, DaySummary } from './types/models';

const BASE = '/api/v1';

async function req<T>(method: string, path: string, body?: unknown): Promise<T> {
  const res = await fetch(BASE + path, {
    method,
    headers: body ? { 'Content-Type': 'application/json' } : {},
    body: body ? JSON.stringify(body) : undefined,
  });
  if (res.status === 204) return undefined as T;
  const data = await res.json();
  if (!res.ok) throw new Error(data.error ?? `HTTP ${res.status}`);
  return data as T;
}

export const api = {
  // Child
  getChild: () => req<Child | null>('GET', '/child'),
  createChild: (body: Partial<Child>) => req<Child>('POST', '/child', body),
  updateChild: (body: Partial<Child>) => req<Child>('PUT', '/child', body),

  // Sleep
  getSleep: (date?: string) => req<SleepLog[]>('GET', `/sleep${date ? `?date=${date}` : ''}`),
  createSleep: (body: Partial<SleepLog>) => req<SleepLog & { stopped_feeding?: { id: string; feed_type: string; duration_minutes: number } }>('POST', '/sleep', body),
  getActiveSleep: () => req<SleepLog | null>('GET', '/sleep/active'),
  updateSleep: (id: string, body: Partial<SleepLog>) => req<SleepLog>('PUT', `/sleep/${id}`, body),
  deleteSleep: (id: string) => req<void>('DELETE', `/sleep/${id}`),

  // Feeding
  getFeeding: (date?: string) => req<FeedingLog[]>('GET', `/feeding${date ? `?date=${date}` : ''}`),
  createFeeding: (body: Partial<FeedingLog>) => req<FeedingLog & { stopped_sleep?: { id: string; duration_minutes: number } }>('POST', '/feeding', body),
  getActiveFeeding: () => req<FeedingLog | null>('GET', '/feeding/active'),
  updateFeeding: (id: string, body: Partial<FeedingLog>) => req<FeedingLog>('PUT', `/feeding/${id}`, body),
  deleteFeeding: (id: string) => req<void>('DELETE', `/feeding/${id}`),

  // Diaper
  getDiaper: (date?: string) => req<DiaperLog[]>('GET', `/diaper${date ? `?date=${date}` : ''}`),
  createDiaper: (body: Partial<DiaperLog>) => req<DiaperLog>('POST', '/diaper', body),
  updateDiaper: (id: string, body: Partial<DiaperLog>) => req<DiaperLog>('PUT', `/diaper/${id}`, body),
  deleteDiaper: (id: string) => req<void>('DELETE', `/diaper/${id}`),

  // Growth
  getGrowth: () => req<GrowthLog[]>('GET', '/growth'),
  createGrowth: (body: Partial<GrowthLog>) => req<GrowthLog>('POST', '/growth', body),
  updateGrowth: (id: string, body: Partial<GrowthLog>) => req<GrowthLog>('PUT', `/growth/${id}`, body),
  deleteGrowth: (id: string) => req<void>('DELETE', `/growth/${id}`),

  // Summary
  getSummary: (date?: string) => req<DaySummary>('GET', `/summary${date ? `?date=${date}` : ''}`),
};
