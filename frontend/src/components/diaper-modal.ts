import { h } from '../utils/dom';
import { api } from '../api';
import { showToast } from './toast';
import { nowISO, nowForInput, isoToLocalInput, localInputToISO } from '../utils/date';
import type { DiaperLog } from '../types/models';

type DiaperType = 'wet' | 'dirty' | 'mixed';

export function renderDiaperModal(onSave: () => void): HTMLElement {
  const overlay = h('div', { class: 'modal-overlay' });
  const modal = h('div', { class: 'modal' });

  const close = () => {
    overlay.classList.remove('open');
    setTimeout(() => overlay.remove(), 300);
  };

  overlay.appendChild(modal);
  modal.appendChild(h('div', { class: 'modal-handle' }));
  modal.appendChild(h('div', { class: 'modal-header' },
    h('h2', { class: 'modal-title' }, 'Diaper Change'),
    h('button', { class: 'modal-close', onClick: close }, '×'),
  ));

  const timeInput = h('input', {
    class: 'form-input', type: 'datetime-local', value: nowForInput(),
  }) as HTMLInputElement;

  const logDiaper = async (type: DiaperType) => {
    try {
      const changedAt = timeInput.value ? localInputToISO(timeInput.value) : nowISO();
      await api.createDiaper({ diaper_type: type, changed_at: changedAt });
      const labels: Record<DiaperType, string> = { wet: 'Wet', dirty: 'Dirty', mixed: 'Mixed' };
      showToast(`${labels[type]} diaper logged`);
      onSave();
      close();
    } catch (e: any) {
      showToast(e.message, 'error');
    }
  };

  const options: Array<[DiaperType, string, string]> = [
    ['wet', '💧', 'Wet'],
    ['dirty', '💩', 'Dirty'],
    ['mixed', '🔄', 'Mixed'],
  ];

  const grid = h('div', { class: 'diaper-options' });
  for (const [type, emoji, label] of options) {
    grid.appendChild(h('button', {
      class: 'diaper-btn',
      onClick: () => logDiaper(type),
    },
      h('span', { class: 'diaper-btn-emoji' }, emoji),
      label,
    ));
  }

  modal.appendChild(h('div', { class: 'modal-body' },
    h('div', { class: 'form-group' }, h('label', { class: 'form-label' }, 'Time'), timeInput),
  ));
  modal.appendChild(grid);

  setTimeout(() => overlay.classList.add('open'), 10);
  overlay.addEventListener('click', (e) => { if (e.target === overlay) close(); });

  return overlay;
}

export function renderDiaperEditModal(diaper: DiaperLog, onSave: () => void): HTMLElement {
  const overlay = h('div', { class: 'modal-overlay' });
  const modal = h('div', { class: 'modal' });
  overlay.appendChild(modal);

  const close = () => {
    overlay.classList.remove('open');
    setTimeout(() => overlay.remove(), 300);
  };
  overlay.addEventListener('click', (e) => { if (e.target === overlay) close(); });

  let selectedType = diaper.diaper_type as DiaperType;
  const btns: HTMLButtonElement[] = [];

  const grid = h('div', { class: 'diaper-options' });
  const options: Array<[DiaperType, string, string]> = [
    ['wet', '💧', 'Wet'],
    ['dirty', '💩', 'Dirty'],
    ['mixed', '🔄', 'Mixed'],
  ];
  for (const [type, emoji, label] of options) {
    const btn = h('button', {
      class: `diaper-btn${type === selectedType ? ' active' : ''}`,
      style: type === selectedType ? 'outline: 2px solid var(--color-diaper)' : '',
      onClick: () => {
        selectedType = type;
        btns.forEach((b, i) => {
          b.style.outline = options[i][0] === selectedType ? '2px solid var(--color-diaper)' : '';
        });
      },
    },
      h('span', { class: 'diaper-btn-emoji' }, emoji),
      label,
    ) as HTMLButtonElement;
    btns.push(btn);
    grid.appendChild(btn);
  }

  const timeInput = h('input', {
    class: 'form-input', type: 'datetime-local', value: isoToLocalInput(diaper.changed_at),
  }) as HTMLInputElement;

  const saveBtn = h('button', {
    class: 'btn btn-primary', style: 'flex: 1',
    onClick: async () => {
      saveBtn.disabled = true;
      try {
        await api.updateDiaper(diaper.id, {
          diaper_type: selectedType,
          changed_at: localInputToISO(timeInput.value),
        });
        showToast('Diaper updated');
        close();
        onSave();
      } catch (e: any) {
        showToast(e.message ?? 'Failed to update', 'error');
        saveBtn.disabled = false;
      }
    },
  }, 'Save');

  modal.appendChild(h('div', { class: 'modal-handle' }));
  modal.appendChild(h('div', { class: 'modal-header' },
    h('div', { class: 'modal-title' }, '✏️ Edit Diaper'),
    h('button', { class: 'modal-close', onClick: close }, '×'),
  ));
  modal.appendChild(grid);
  modal.appendChild(h('div', { class: 'modal-body' },
    h('div', { class: 'form-group' }, h('label', { class: 'form-label' }, 'Time'), timeInput),
  ));
  modal.appendChild(h('div', { class: 'modal-actions' }, saveBtn));

  setTimeout(() => overlay.classList.add('open'), 10);
  return overlay;
}
