export interface Child {
  id: string;
  name: string;
  date_of_birth: string;
  gender: string;
  photo_url?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

export interface SleepLog {
  id: string;
  child_id: string;
  start_time: string;
  end_time: string | null;
  duration_minutes: number | null;
  notes?: string;
  created_at: string;
}

export interface FeedingLog {
  id: string;
  child_id: string;
  feed_type: 'breast_left' | 'breast_right' | 'bottle';
  start_time: string;
  end_time: string | null;
  duration_minutes: number | null;
  quantity_ml: number | null;
  notes?: string;
  created_at: string;
}

export interface DiaperLog {
  id: string;
  child_id: string;
  diaper_type: 'wet' | 'dirty' | 'mixed';
  changed_at: string;
  notes?: string;
  created_at: string;
}

export interface GrowthLog {
  id: string;
  child_id: string;
  measured_on: string;
  weight_grams: number | null;
  length_mm: number | null;
  head_circumference_mm: number | null;
  notes?: string;
  created_at: string;
}

export interface DayStats {
  date: string;
  sleep_minutes: number;
  sleep_count: number;
  feeding_count: number;
  breast_feed_count: number;
  bottle_feed_count: number;
  bottle_ml_total: number;
  diaper_count: number;
  wet_count: number;
  dirty_count: number;
}

export interface DaySummary {
  date: string;
  total_sleep_minutes: number;
  sleep_count: number;
  feeding_count: number;
  diaper_count: number;
  last_weight_grams?: number;
  last_sleep_end_time?: string;
  active_sleep?: SleepLog;
  active_feeding?: FeedingLog;
}
