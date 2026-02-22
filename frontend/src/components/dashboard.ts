import { h } from '../utils/dom';
import { api } from '../api';
import { state } from '../state';
import { renderNav } from './nav';
import { renderQuickAdd } from './quick-add';
import { renderSleepModal } from './sleep-modal';
import { renderFeedingModal } from './feeding-modal';
import { renderDiaperModal } from './diaper-modal';
import { renderGrowthModal } from './growth-modal';
import { showToast } from './toast';
import { calcAge, formatDuration, timeAgo, formatTime, formatElapsed, elapsedSeconds, nowISO, todayISO } from '../utils/date';
import type { DaySummary } from '../types/models';

export function renderDashboard(): HTMLElement {
  const screen = h('div', { class: 'screen dashboard' });

  let summaryData: DaySummary | null = null;
  let timerInterval: ReturnType<typeof setInterval> | null = null;
  let timerBanner: HTMLElement | null = null;

  const refresh = async () => {
    try {
      const [summary, activeSleep, activeFeeding] = await Promise.all([
        api.getSummary(todayISO()),
        api.getActiveSleep(),
        api.getActiveFeeding(),
      ]);
      summaryData = summary;
      state.activeSleep.set(activeSleep ?? null);
      state.activeFeeding.set(activeFeeding ?? null);
      renderContent();
    } catch (e: any) {
      showToast('Failed to load data', 'error');
    }
  };

  const renderContent = () => {
    screen.innerHTML = '';
    if (timerInterval) clearInterval(timerInterval);

    const child = state.child.get();
    if (!child) return;

    // Header
    const avatar = h('div', { class: 'baby-avatar' });
    if (child.photo_url) {
      avatar.appendChild(h('img', { src: child.photo_url, alt: child.name }));
    } else {
      avatar.textContent = child.name[0].toUpperCase();
    }

    screen.appendChild(h('div', { class: 'dashboard-header' },
      h('div', { class: 'dashboard-baby-info' },
        avatar,
        h('div', {},
          h('div', { class: 'baby-name' }, child.name),
          h('div', { class: 'baby-age' }, calcAge(child.date_of_birth)),
        ),
      ),
    ));

    // Active timer banner
    const activeSleep = state.activeSleep.get();
    const activeFeeding = state.activeFeeding.get();

    if (activeSleep) {
      timerBanner = renderTimerBanner(
        '😴 Sleeping',
        activeSleep.start_time,
        '#8B5CF6',
        async () => {
          try {
            await api.updateSleep(activeSleep.id, { end_time: nowISO() });
            state.activeSleep.set(null);
            showToast('Sleep logged');
            refresh();
          } catch (e: any) {
            showToast(e.message, 'error');
          }
        },
      );
      screen.appendChild(timerBanner);
    } else if (activeFeeding) {
      const label = activeFeeding.feed_type === 'breast_left' ? '◀️ Left Breast' : '▶️ Right Breast';
      timerBanner = renderTimerBanner(
        `🍼 ${label}`,
        activeFeeding.start_time,
        '#EC4899',
        async () => {
          try {
            await api.updateFeeding(activeFeeding.id, { end_time: nowISO() });
            state.activeFeeding.set(null);
            showToast('Feed logged');
            refresh();
          } catch (e: any) {
            showToast(e.message, 'error');
          }
        },
      );
      screen.appendChild(timerBanner);
    }

    // Summary cards
    const s = summaryData;
    if (s) {
      const weightStr = s.last_weight_grams
        ? `${(s.last_weight_grams / 1000).toFixed(2)} kg`
        : '—';

      const openModal = (renderFn: (refresh: () => void) => HTMLElement) => {
        document.getElementById('app')!.appendChild(renderFn(refresh));
      };

      screen.appendChild(h('div', { class: 'summary-grid' },
        h('button', {
          class: 'summary-card summary-card-sleep',
          onClick: () => openModal(renderSleepModal),
        },
          h('div', { class: 'summary-card-icon' }, '😴'),
          h('div', { class: 'summary-card-value' }, formatDuration(s.total_sleep_minutes) || '0m'),
          h('div', { class: 'summary-card-label' }, `Sleep · ${s.sleep_count} sessions`),
        ),
        h('button', {
          class: 'summary-card summary-card-feeding',
          onClick: () => openModal(renderFeedingModal),
        },
          h('div', { class: 'summary-card-icon' }, '🍼'),
          h('div', { class: 'summary-card-value' }, String(s.feeding_count)),
          h('div', { class: 'summary-card-label' }, 'Feedings today'),
        ),
        h('button', {
          class: 'summary-card summary-card-diaper',
          onClick: () => openModal(renderDiaperModal),
        },
          h('div', { class: 'summary-card-icon' }, '🚼'),
          h('div', { class: 'summary-card-value' }, String(s.diaper_count)),
          h('div', { class: 'summary-card-label' }, 'Diaper changes'),
        ),
        h('button', {
          class: 'summary-card summary-card-growth',
          onClick: () => openModal(renderGrowthModal),
        },
          h('div', { class: 'summary-card-icon' }, '📏'),
          h('div', { class: 'summary-card-value', style: 'font-size: 18px' }, weightStr),
          h('div', { class: 'summary-card-label' }, 'Last weight'),
        ),
      ));
    }

    // Recent activity
    loadRecentActivity(screen, refresh);

    // Nav + FAB
    screen.appendChild(renderNav());
    screen.appendChild(renderQuickAdd(refresh));
  };

  refresh();
  return screen;
}

function renderTimerBanner(
  title: string,
  startTime: string,
  color: string,
  onStop: () => void,
): HTMLElement {
  const timeEl = h('div', { class: 'timer-banner-time' }, formatElapsed(elapsedSeconds(startTime)));
  const interval = setInterval(() => {
    timeEl.textContent = formatElapsed(elapsedSeconds(startTime));
  }, 1000);

  const banner = h('div', { class: 'timer-banner', style: `--banner-color: ${color}` },
    h('div', { class: 'timer-banner-info' },
      h('div', { class: 'timer-banner-title' }, title),
      timeEl,
    ),
    h('button', {
      class: 'timer-banner-stop',
      onClick: () => {
        clearInterval(interval);
        onStop();
      },
    }, 'Stop'),
  );

  return banner;
}

async function loadRecentActivity(screen: HTMLElement, refresh: () => void): Promise<void> {
  try {
    const today = todayISO();
    const [sleep, feeding, diaper] = await Promise.all([
      api.getSleep(today),
      api.getFeeding(today),
      api.getDiaper(today),
    ]);

    type Item = { time: string; el: HTMLElement };
    const items: Item[] = [];

    for (const s of sleep) {
      const detail = s.end_time
        ? `${formatTime(s.start_time)} → ${formatTime(s.end_time)} (${s.duration_minutes}m)`
        : `${formatTime(s.start_time)} — in progress`;
      items.push({
        time: s.start_time,
        el: h('div', { class: 'activity-item' },
          h('div', { class: 'activity-dot activity-dot-sleep' }, '😴'),
          h('div', { class: 'activity-info' },
            h('div', { class: 'activity-title' }, 'Sleep'),
            h('div', { class: 'activity-detail' }, detail),
          ),
          h('div', { class: 'activity-time' }, timeAgo(s.start_time)),
        ),
      });
    }

    for (const f of feeding) {
      const typeLabel: Record<string, string> = { breast_left: '◀ Left', breast_right: '▶ Right', bottle: '🍼 Bottle' };
      const detail = f.quantity_ml ? `${f.quantity_ml}ml` : (f.duration_minutes ? `${f.duration_minutes}m` : 'In progress');
      items.push({
        time: f.start_time,
        el: h('div', { class: 'activity-item' },
          h('div', { class: 'activity-dot activity-dot-feeding' }, '🍼'),
          h('div', { class: 'activity-info' },
            h('div', { class: 'activity-title' }, `Feed — ${typeLabel[f.feed_type] ?? f.feed_type}`),
            h('div', { class: 'activity-detail' }, detail),
          ),
          h('div', { class: 'activity-time' }, timeAgo(f.start_time)),
        ),
      });
    }

    for (const d of diaper) {
      const label: Record<string, string> = { wet: 'Wet 💧', dirty: 'Dirty 💩', mixed: 'Mixed 🔄' };
      items.push({
        time: d.changed_at,
        el: h('div', { class: 'activity-item' },
          h('div', { class: 'activity-dot activity-dot-diaper' }, '🚼'),
          h('div', { class: 'activity-info' },
            h('div', { class: 'activity-title' }, 'Diaper'),
            h('div', { class: 'activity-detail' }, label[d.diaper_type] ?? d.diaper_type),
          ),
          h('div', { class: 'activity-time' }, timeAgo(d.changed_at)),
        ),
      });
    }

    items.sort((a, b) => b.time.localeCompare(a.time));
    const recent = items.slice(0, 8);

    const titleEl = h('div', { class: 'section-title' }, "Today's Activity");
    screen.insertBefore(titleEl, screen.children[screen.children.length - 2]);

    if (recent.length === 0) {
      const empty = h('div', { class: 'empty-state', style: 'margin: 0 16px' },
        h('div', { class: 'empty-state-icon' }, '📝'),
        h('div', { class: 'empty-state-text' }, 'No activity logged today.\nTap + to get started.'),
      );
      screen.insertBefore(empty, screen.children[screen.children.length - 2]);
    } else {
      const list = h('div', { class: 'activity-list' });
      for (const item of recent) list.appendChild(item.el);
      screen.insertBefore(list, screen.children[screen.children.length - 2]);
    }
  } catch {
    // Silently fail for activity list
  }
}
