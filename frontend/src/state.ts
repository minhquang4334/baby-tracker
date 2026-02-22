import type { Child, SleepLog, FeedingLog } from './types/models';

type Listener<T> = (val: T) => void;

function signal<T>(initial: T) {
  let value = initial;
  const listeners: Listener<T>[] = [];
  return {
    get: () => value,
    set: (v: T) => {
      value = v;
      listeners.forEach(fn => fn(v));
    },
    subscribe: (fn: Listener<T>) => {
      listeners.push(fn);
      return () => {
        const i = listeners.indexOf(fn);
        if (i >= 0) listeners.splice(i, 1);
      };
    },
  };
}

export const state = {
  child: signal<Child | null>(null),
  activeSleep: signal<SleepLog | null>(null),
  activeFeeding: signal<FeedingLog | null>(null),
  lastBottleML: signal<number>(120),
};
