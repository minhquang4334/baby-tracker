import { h } from '../utils/dom';
import { api } from '../api';
import { state } from '../state';
import { showToast } from './toast';
import { formatTime, formatElapsed, elapsedSeconds, nowISO, nowForInput, localInputToISO } from '../utils/date';

export function renderSleepModal(onSave: () => void): HTMLElement {
  const activeSleep = state.activeSleep.get();
  let timerInterval: ReturnType<typeof setInterval> | null = null;

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
    h('h2', { class: 'modal-title' }, activeSleep ? 'Sleep in Progress' : 'Start Sleep'),
    h('button', { class: 'modal-close', onClick: close }, '×'),
  );
  modal.appendChild(header);

  const body = h('div', { class: 'modal-body' });
  modal.appendChild(body);

  if (activeSleep) {
    // Show active timer with stop button
    const timerEl = h('div', { class: 'timer-elapsed' }, formatElapsed(elapsedSeconds(activeSleep.start_time)));
    const startedEl = h('p', { class: 'timer-started' }, `Started at ${formatTime(activeSleep.start_time)}`);

    body.appendChild(h('div', { class: 'timer-display' }, timerEl, startedEl));

    timerInterval = setInterval(() => {
      timerEl.textContent = formatElapsed(elapsedSeconds(activeSleep.start_time));
    }, 1000);

    const stopBtn = h('button', {
      class: 'btn btn-primary btn-full',
      onClick: async () => {
        stopBtn.disabled = true;
        stopBtn.textContent = 'Stopping...';
        try {
          const updated = await api.updateSleep(activeSleep.id, { end_time: nowISO() });
          state.activeSleep.set(null);
          showToast(`Sleep logged: ${updated.duration_minutes}m`);
          onSave();
          close();
        } catch (e: any) {
          showToast(e.message, 'error');
          stopBtn.disabled = false;
          stopBtn.textContent = 'Stop Sleep';
        }
      },
    }, 'Stop Sleep');

    const deleteBtn = h('button', {
      class: 'btn btn-ghost btn-full',
      style: 'margin-top: 8px',
      onClick: async () => {
        if (!confirm('Delete this sleep session?')) return;
        await api.deleteSleep(activeSleep.id);
        state.activeSleep.set(null);
        onSave();
        close();
      },
    }, 'Cancel Session');

    modal.appendChild(h('div', { class: 'modal-actions', style: 'flex-direction: column' }, stopBtn, deleteBtn));
  } else {
    // Manual entry or quick start
    const startInput = h('input', {
      class: 'form-input',
      type: 'datetime-local',
      id: 'sleep-start',
    }) as HTMLInputElement;
    startInput.value = nowForInput();

    const endInput = h('input', {
      class: 'form-input',
      type: 'datetime-local',
      id: 'sleep-end',
    }) as HTMLInputElement;

    body.appendChild(h('div', { class: 'form-group' },
      h('label', { class: 'form-label', for: 'sleep-start' }, 'Start time'),
      startInput,
    ));
    body.appendChild(h('div', { class: 'form-group' },
      h('label', { class: 'form-label', for: 'sleep-end' }, 'End time (leave blank to start timer)'),
      endInput,
    ));

    const startBtn = h('button', {
      class: 'btn btn-primary btn-full',
      onClick: async () => {
        startBtn.disabled = true;
        startBtn.innerHTML = '<div class="spinner"></div>';
        try {
          const startISO = localInputToISO(startInput.value);
          const endISO = endInput.value ? localInputToISO(endInput.value) : '';

          if (endISO) {
            // Log a completed sleep: create then immediately stop it
            const res = await api.createSleep({ start_time: startISO });
            await api.updateSleep(res.id, { end_time: endISO });
            state.activeFeeding.set(null);
            if (res.stopped_feeding) {
              const side = res.stopped_feeding.feed_type === 'breast_left' ? 'Left' : 'Right';
              showToast(`${side} feed stopped — sleep logged`);
            } else {
              showToast('Sleep logged');
            }
          } else {
            // Start a live timer
            const res = await api.createSleep({ start_time: startISO });
            state.activeSleep.set(res);
            state.activeFeeding.set(null);
            if (res.stopped_feeding) {
              const side = res.stopped_feeding.feed_type === 'breast_left' ? 'Left' : 'Right';
              showToast(`${side} feed stopped (${res.stopped_feeding.duration_minutes}m) — sleep started`);
            } else {
              showToast('Sleep started');
            }
          }
          onSave();
          close();
        } catch (e: any) {
          showToast(e.message, 'error');
          startBtn.disabled = false;
          startBtn.textContent = endInput.value ? 'Log Sleep' : 'Start Sleep';
        }
      },
    }, 'Start Sleep');

    // Update button label dynamically based on whether end time is filled
    endInput.addEventListener('input', () => {
      startBtn.textContent = endInput.value ? 'Log Sleep' : 'Start Sleep';
    });

    modal.appendChild(h('div', { class: 'modal-actions' }, startBtn));
  }

  // Animate in
  setTimeout(() => overlay.classList.add('open'), 10);

  overlay.addEventListener('click', (e) => {
    if (e.target === overlay) close();
  });

  return overlay;
}
