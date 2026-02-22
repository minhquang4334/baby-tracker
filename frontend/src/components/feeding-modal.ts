import { h } from '../utils/dom';
import { api } from '../api';
import { state } from '../state';
import { showToast } from './toast';
import { formatTime, formatElapsed, elapsedSeconds, nowISO } from '../utils/date';

type FeedTab = 'breast' | 'bottle';

export function renderFeedingModal(onSave: () => void): HTMLElement {
  const activeFeeding = state.activeFeeding.get();
  let timerInterval: ReturnType<typeof setInterval> | null = null;
  let currentTab: FeedTab = 'breast';

  const overlay = h('div', { class: 'modal-overlay' });
  const modal = h('div', { class: 'modal' });

  const close = () => {
    if (timerInterval) clearInterval(timerInterval);
    overlay.classList.remove('open');
    setTimeout(() => overlay.remove(), 300);
  };

  overlay.appendChild(modal);
  modal.appendChild(h('div', { class: 'modal-handle' }));

  const header = h('div', { class: 'modal-header' },
    h('h2', { class: 'modal-title' }, activeFeeding ? 'Feed in Progress' : 'Log Feeding'),
    h('button', { class: 'modal-close', onClick: close }, '×'),
  );
  modal.appendChild(header);

  if (activeFeeding) {
    // Active timer view
    const timerEl = h('div', { class: 'timer-elapsed' }, formatElapsed(elapsedSeconds(activeFeeding.start_time)));
    const sideLabel = activeFeeding.feed_type === 'breast_left' ? 'Left Breast' : 'Right Breast';
    const startedEl = h('p', { class: 'timer-started' }, `${sideLabel} — started ${formatTime(activeFeeding.start_time)}`);

    timerInterval = setInterval(() => {
      timerEl.textContent = formatElapsed(elapsedSeconds(activeFeeding.start_time));
    }, 1000);

    modal.appendChild(h('div', { class: 'modal-body' },
      h('div', { class: 'timer-display' }, timerEl, startedEl),
    ));

    const stopBtn = h('button', {
      class: 'btn btn-primary btn-full',
      onClick: async () => {
        stopBtn.disabled = true;
        stopBtn.textContent = 'Stopping...';
        try {
          const updated = await api.updateFeeding(activeFeeding.id, { end_time: nowISO() });
          state.activeFeeding.set(null);
          showToast(`Feed logged: ${updated.duration_minutes}m`);
          onSave();
          close();
        } catch (e: any) {
          showToast(e.message, 'error');
          stopBtn.disabled = false;
          stopBtn.textContent = 'Stop Feeding';
        }
      },
    }, 'Stop Feeding');

    modal.appendChild(h('div', { class: 'modal-actions', style: 'flex-direction: column' }, stopBtn));
  } else {
    // New feeding
    const tabBreast = h('button', { class: 'feed-type-tab active', onClick: () => switchTab('breast') }, '🤱 Breast') as HTMLButtonElement;
    const tabBottle = h('button', { class: 'feed-type-tab', onClick: () => switchTab('bottle') }, '🍼 Bottle') as HTMLButtonElement;

    modal.appendChild(h('div', { class: 'feed-type-tabs' }, tabBreast, tabBottle));

    const breastSection = renderBreastSection(close, onSave);
    const bottleSection = renderBottleSection(close, onSave);
    bottleSection.style.display = 'none';

    modal.appendChild(breastSection);
    modal.appendChild(bottleSection);

    function switchTab(tab: FeedTab) {
      currentTab = tab;
      tabBreast.classList.toggle('active', tab === 'breast');
      tabBottle.classList.toggle('active', tab === 'bottle');
      breastSection.style.display = tab === 'breast' ? '' : 'none';
      bottleSection.style.display = tab === 'bottle' ? '' : 'none';
    }
  }

  setTimeout(() => overlay.classList.add('open'), 10);
  overlay.addEventListener('click', (e) => { if (e.target === overlay) close(); });

  return overlay;
}

function renderBreastSection(close: () => void, onSave: () => void): HTMLElement {
  const wrap = h('div', {});

  const makeBreastBtn = (side: 'breast_left' | 'breast_right', emoji: string, label: string) =>
    h('button', {
      class: 'breast-btn',
      onClick: async () => {
        try {
          const res = await api.createFeeding({ feed_type: side, start_time: nowISO() });
          state.activeFeeding.set(res);
          state.activeSleep.set(null);
          if (res.stopped_sleep) {
            showToast(`Sleep stopped (${res.stopped_sleep.duration_minutes}m) — ${label} feed started`);
          } else {
            showToast(`${label} feed started`);
          }
          onSave();
          close();
        } catch (e: any) {
          showToast(e.message, 'error');
        }
      },
    },
      h('span', { class: 'breast-btn-emoji' }, emoji),
      label,
    );

  wrap.appendChild(h('div', { class: 'breast-options' },
    makeBreastBtn('breast_left', '◀️', 'Left'),
    makeBreastBtn('breast_right', '▶️', 'Right'),
  ));

  return wrap;
}

function renderBottleSection(close: () => void, onSave: () => void): HTMLElement {
  const lastML = state.lastBottleML.get();

  const mlInput = h('input', {
    class: 'form-input',
    type: 'number',
    value: String(lastML),
    min: '0',
    max: '500',
    step: '5',
    placeholder: 'ml',
  }) as HTMLInputElement;

  const saveBtn = h('button', {
    class: 'btn btn-primary btn-full',
    onClick: async () => {
      const ml = parseInt(mlInput.value, 10);
      if (isNaN(ml) || ml < 0) { showToast('Enter a valid quantity', 'error'); return; }
      saveBtn.disabled = true;
      saveBtn.innerHTML = '<div class="spinner"></div>';
      try {
        await api.createFeeding({ feed_type: 'bottle', start_time: nowISO(), quantity_ml: ml });
        state.lastBottleML.set(ml);
        showToast(`Bottle logged: ${ml}ml`);
        onSave();
        close();
      } catch (e: any) {
        showToast(e.message, 'error');
        saveBtn.disabled = false;
        saveBtn.textContent = 'Save';
      }
    },
  }, 'Save');

  return h('div', {},
    h('div', { class: 'modal-body' },
      h('div', { class: 'form-group' },
        h('label', { class: 'form-label' }, 'Amount (ml)'),
        mlInput,
      ),
    ),
    h('div', { class: 'modal-actions' }, saveBtn),
  );
}
