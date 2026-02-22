import { h } from '../utils/dom';
import { api } from '../api';
import { showToast } from './toast';
import { isoToLocalInput, localInputToISO } from '../utils/date';
import type { SleepLog } from '../types/models';

export function renderSleepEditModal(sleep: SleepLog, onSave: () => void): HTMLElement {
  const overlay = h('div', { class: 'modal-overlay' });
  const modal = h('div', { class: 'modal' });

  const close = () => {
    overlay.classList.remove('open');
    setTimeout(() => overlay.remove(), 300);
  };

  overlay.addEventListener('click', (e) => {
    if (e.target === overlay) close();
  });

  const startVal = isoToLocalInput(sleep.start_time);
  const endVal = sleep.end_time ? isoToLocalInput(sleep.end_time) : '';

  const startInput = h('input', {
    type: 'datetime-local',
    class: 'form-input',
    value: startVal,
  }) as HTMLInputElement;

  const endInput = h('input', {
    type: 'datetime-local',
    class: 'form-input',
    value: endVal,
  }) as HTMLInputElement;

  const saveBtn = h('button', {
    class: 'btn btn-primary',
    style: 'flex: 1',
    onClick: async () => {
      const startTime = localInputToISO(startInput.value);
      const endTime = endInput.value ? localInputToISO(endInput.value) : '';
      try {
        await api.updateSleep(sleep.id, { start_time: startTime, end_time: endTime || undefined });
        showToast('Sleep updated');
        close();
        onSave();
      } catch (e: any) {
        showToast(e.message ?? 'Failed to update');
      }
    },
  }, 'Save');

  modal.appendChild(h('div', { class: 'modal-handle' }));
  modal.appendChild(
    h('div', { class: 'modal-header' },
      h('div', { class: 'modal-title' }, '✏️ Edit Sleep'),
      h('button', { class: 'modal-close', onClick: close }, '×'),
    )
  );
  modal.appendChild(
    h('div', { class: 'modal-body' },
      h('div', { class: 'form-group' },
        h('label', { class: 'form-label' }, 'Start time'),
        startInput,
      ),
      h('div', { class: 'form-group' },
        h('label', { class: 'form-label' }, 'End time'),
        endInput,
      ),
    )
  );
  modal.appendChild(
    h('div', { class: 'modal-actions' }, saveBtn)
  );

  overlay.appendChild(modal);
  requestAnimationFrame(() => overlay.classList.add('open'));

  return overlay;
}
