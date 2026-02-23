import { h } from '../utils/dom';
import { api } from '../api';
import { renderNav } from './nav';
import { renderQuickAdd } from './quick-add';
import { showToast } from './toast';
import { renderBottleEditModal, renderBreastFeedEditModal } from './feeding-modal';
import { renderDiaperEditModal } from './diaper-modal';
import { renderSleepEditModal } from './sleep-edit-modal';
import { formatTime, formatDateFull, formatDuration, todayISO } from '../utils/date';
import type { SleepLog, FeedingLog, DiaperLog } from '../types/models';

type Filter = 'all' | 'sleep' | 'feeding' | 'diaper';

export function renderHistory(): HTMLElement {
  const screen = h('div', { class: 'screen history' });
  let currentDate = todayISO();
  let currentFilter: Filter = 'all';

  const titleEl = h('div', { style: 'padding: 16px 20px 0; font-size: 20px; font-weight: 700' }, 'History');

  const dateNav = h('div', { class: 'date-nav' });
  const dateLabelEl = h('span', { class: 'date-nav-label' }, '');

  const prevBtn = h('button', { class: 'date-nav-btn', onClick: () => changeDate(-1) }, '‹');
  const nextBtn = h('button', { class: 'date-nav-btn', onClick: () => changeDate(1) }, '›');

  dateNav.appendChild(prevBtn);
  dateNav.appendChild(dateLabelEl);
  dateNav.appendChild(nextBtn);

  const filterBar = h('div', { class: 'filter-bar' });
  const filters: Array<[Filter, string, string]> = [
    ['all', 'All', ''],
    ['sleep', 'Sleep', 'sleep'],
    ['feeding', 'Feeding', 'feeding'],
    ['diaper', 'Diaper', 'diaper'],
  ];

  const filterBtns: HTMLButtonElement[] = [];
  for (const [id, label, cls] of filters) {
    const btn = h('button', {
      class: `pill${id === currentFilter ? ' active' : ''}${cls ? ` pill-${cls}` : ''}`,
      onClick: () => {
        currentFilter = id;
        filterBtns.forEach(b => b.classList.remove('active'));
        btn.classList.add('active');
        renderTimeline();
      },
    }, label) as HTMLButtonElement;
    filterBtns.push(btn);
    filterBar.appendChild(btn);
  }

  const timelineWrap = h('div', {});

  const changeDate = (delta: number) => {
    const d = new Date(currentDate);
    d.setDate(d.getDate() + delta);
    currentDate = d.toISOString().slice(0, 10);
    updateDateLabel();
    renderTimeline();
    // Disable next if at today
    nextBtn.disabled = currentDate >= todayISO();
  };

  const updateDateLabel = () => {
    const d = new Date(currentDate + 'T12:00:00');
    if (currentDate === todayISO()) {
      dateLabelEl.textContent = 'Today';
    } else {
      const yesterday = new Date();
      yesterday.setDate(yesterday.getDate() - 1);
      if (currentDate === yesterday.toISOString().slice(0, 10)) {
        dateLabelEl.textContent = 'Yesterday';
      } else {
        dateLabelEl.textContent = d.toLocaleDateString([], { month: 'short', day: 'numeric' });
      }
    }
    nextBtn.disabled = currentDate >= todayISO();
  };

  const confirmDelete = (msg: string, action: () => Promise<void>, successMsg: string) => async () => {
    if (!confirm(msg)) return;
    await action();
    showToast(successMsg);
    renderTimeline();
  };

  const renderTimeline = async () => {
    timelineWrap.innerHTML = '<div style="padding: 24px; text-align: center; color: var(--color-text-muted)">Loading...</div>';

    try {
      const promises: Promise<any[]>[] = [];
      if (currentFilter === 'all' || currentFilter === 'sleep') promises.push(api.getSleep(currentDate));
      if (currentFilter === 'all' || currentFilter === 'feeding') promises.push(api.getFeeding(currentDate));
      if (currentFilter === 'all' || currentFilter === 'diaper') promises.push(api.getDiaper(currentDate));

      const results = await Promise.all(promises);

      type TimelineItem = { time: string; el: HTMLElement };
      const items: TimelineItem[] = [];
      let ri = 0;

      if (currentFilter === 'all' || currentFilter === 'sleep') {
        const sleepLogs: SleepLog[] = results[ri++];
        for (const s of sleepLogs) {
          const detail = s.end_time
            ? `${formatTime(s.start_time)} → ${formatTime(s.end_time)}${s.duration_minutes ? ` · ${formatDuration(s.duration_minutes)}` : ''}`
            : `${formatTime(s.start_time)} — in progress`;
          const onEdit = () => {
            document.getElementById('app')!.appendChild(renderSleepEditModal(s, renderTimeline));
          };
          const onDelete = confirmDelete('Delete this sleep log?', () => api.deleteSleep(s.id), 'Sleep deleted');
          items.push({
            time: s.start_time,
            el: timelineItem('sleep', '😴', 'Sleep', detail, s.start_time, onEdit, onDelete),
          });
        }
      }

      if (currentFilter === 'all' || currentFilter === 'feeding') {
        const feedingLogs: FeedingLog[] = results[ri++];
        const typeLabel: Record<string, string> = { breast_left: '◀ Left breast', breast_right: '▶ Right breast', bottle: '🍼 Bottle' };
        for (const f of feedingLogs) {
          const detail = f.quantity_ml
            ? `${typeLabel[f.feed_type]} · ${f.quantity_ml}ml`
            : `${typeLabel[f.feed_type]}${f.duration_minutes ? ` · ${formatDuration(f.duration_minutes)}` : ''}`;
          const onEdit = () => {
            const modal = f.feed_type === 'bottle'
              ? renderBottleEditModal(f, renderTimeline)
              : renderBreastFeedEditModal(f, renderTimeline);
            document.getElementById('app')!.appendChild(modal);
          };
          const onDelete = confirmDelete('Delete this feeding log?', () => api.deleteFeeding(f.id), 'Feeding deleted');
          items.push({
            time: f.start_time,
            el: timelineItem('feeding', '🍼', 'Feeding', detail, f.start_time, onEdit, onDelete),
          });
        }
      }

      if (currentFilter === 'all' || currentFilter === 'diaper') {
        const diaperLogs: DiaperLog[] = results[ri++];
        const dLabel: Record<string, string> = { wet: 'Wet 💧', dirty: 'Dirty 💩', mixed: 'Mixed 🔄' };
        for (const d of diaperLogs) {
          const onEdit = () => {
            document.getElementById('app')!.appendChild(renderDiaperEditModal(d, renderTimeline));
          };
          const onDelete = confirmDelete('Delete this diaper log?', () => api.deleteDiaper(d.id), 'Diaper log deleted');
          items.push({
            time: d.changed_at,
            el: timelineItem('diaper', '🚼', 'Diaper', dLabel[d.diaper_type] ?? d.diaper_type, d.changed_at, onEdit, onDelete),
          });
        }
      }

      items.sort((a, b) => b.time.localeCompare(a.time));

      timelineWrap.innerHTML = '';
      if (items.length === 0) {
        timelineWrap.appendChild(h('div', { class: 'empty-state' },
          h('div', { class: 'empty-state-icon' }, '📋'),
          h('div', { class: 'empty-state-text' }, 'No events logged for this day.'),
        ));
      } else {
        const timeline = h('div', { class: 'timeline' });
        for (const item of items) timeline.appendChild(item.el);
        timelineWrap.appendChild(timeline);
      }
    } catch (e: any) {
      timelineWrap.innerHTML = '';
      timelineWrap.appendChild(h('div', { class: 'empty-state' },
        h('div', { class: 'empty-state-text' }, 'Failed to load history.'),
      ));
    }
  };

  updateDateLabel();
  renderTimeline();

  screen.appendChild(titleEl);
  screen.appendChild(dateNav);
  screen.appendChild(filterBar);
  screen.appendChild(timelineWrap);
  screen.appendChild(renderNav());
  screen.appendChild(renderQuickAdd(renderTimeline));

  return screen;
}

function timelineItem(
  category: string,
  emoji: string,
  title: string,
  detail: string,
  time: string,
  onEdit: (() => void) | null = null,
  onDelete: (() => void) | null = null,
  showTime: boolean = true,
): HTMLElement {
  const actions = h('div', { class: 'timeline-actions' });
  if (onEdit) {
    actions.appendChild(h('button', {
      class: 'timeline-action-btn',
      onClick: (e: Event) => { e.stopPropagation(); onEdit(); },
    }, '✏️'));
  }
  if (onDelete) {
    actions.appendChild(h('button', {
      class: 'timeline-action-btn timeline-action-delete',
      onClick: (e: Event) => { e.stopPropagation(); onDelete(); },
    }, '🗑️'));
  }

  const right = h('div', { class: 'timeline-right' });
  if (showTime) right.appendChild(h('span', { class: 'timeline-time' }, formatTime(time)));
  right.appendChild(actions);

  return h('div', { class: 'timeline-item' },
    h('div', { class: `timeline-dot timeline-dot-${category}` }, emoji),
    h('div', { class: 'timeline-content' },
      h('div', { class: 'timeline-text' },
        h('span', { class: 'timeline-title' }, title),
        h('span', { class: 'timeline-sep' }, ' · '),
        h('span', { class: 'timeline-detail' }, detail),
      ),
      right,
    ),
  );
}
