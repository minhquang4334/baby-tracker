const TZ = 'Asia/Ho_Chi_Minh'; // GMT+7
const TZ_OFFSET_MS = 7 * 60 * 60 * 1000;

/** Return a Date adjusted to GMT+7 wall-clock time (for slicing into local strings). */
function toHCMC(d: Date): Date {
  return new Date(d.getTime() + TZ_OFFSET_MS);
}

export function formatTime(iso: string): string {
  return new Date(iso).toLocaleTimeString('vi-VN', { timeZone: TZ, hour: '2-digit', minute: '2-digit', hour12: false });
}

export function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString('vi-VN', { timeZone: TZ, month: 'short', day: 'numeric' });
}

export function formatDateFull(iso: string): string {
  return new Date(iso).toLocaleDateString('vi-VN', { timeZone: TZ, weekday: 'long', month: 'long', day: 'numeric' });
}

export function timeAgo(iso: string): string {
  const diff = Math.floor((Date.now() - new Date(iso).getTime()) / 1000);
  if (diff < 60) return 'just now';
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`;
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`;
  return `${Math.floor(diff / 86400)}d ago`;
}

export function formatDuration(minutes: number): string {
  if (minutes < 60) return `${minutes}m`;
  const h = Math.floor(minutes / 60);
  const m = minutes % 60;
  return m > 0 ? `${h}h ${m}m` : `${h}h`;
}

export function elapsedSeconds(startIso: string): number {
  return Math.floor((Date.now() - new Date(startIso).getTime()) / 1000);
}

export function formatElapsed(seconds: number): string {
  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  const s = seconds % 60;
  const pad = (n: number) => String(n).padStart(2, '0');
  if (h > 0) return `${h}:${pad(m)}:${pad(s)}`;
  return `${pad(m)}:${pad(s)}`;
}

export function calcAge(dob: string): string {
  // dob is YYYY-MM-DD, compare against today in GMT+7
  const birth = new Date(dob + 'T00:00:00+07:00').getTime();
  const diffMs = Date.now() - birth;
  const days = Math.floor(diffMs / (1000 * 60 * 60 * 24));
  if (days < 7) return `${days} day${days !== 1 ? 's' : ''} old`;
  const weeks = Math.floor(days / 7);
  if (weeks < 8) return `${weeks} week${weeks !== 1 ? 's' : ''} old`;
  const months = Math.floor(days / 30.44);
  if (months < 24) return `${months} month${months !== 1 ? 's' : ''} old`;
  const years = Math.floor(months / 12);
  return `${years} year${years !== 1 ? 's' : ''} old`;
}

/** Current date as YYYY-MM-DD in GMT+7. */
export function todayISO(): string {
  return toHCMC(new Date()).toISOString().slice(0, 10);
}

/** Current instant as RFC3339 with +07:00 offset. */
export function nowISO(): string {
  const local = toHCMC(new Date());
  return local.toISOString().slice(0, 19) + '+07:00';
}

/** Convert a datetime-local input value (YYYY-MM-DDTHH:MM) to RFC3339 +07:00. */
export function localInputToISO(value: string): string {
  // value is like "2025-06-01T23:30" — treat as GMT+7
  return value.length === 16 ? value + ':00+07:00' : value + '+07:00';
}

/** Current time formatted as a datetime-local input value (YYYY-MM-DDTHH:MM) in GMT+7. */
export function nowForInput(): string {
  return toHCMC(new Date()).toISOString().slice(0, 16);
}
