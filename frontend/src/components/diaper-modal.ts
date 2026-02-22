import { h } from '../utils/dom';
import { api } from '../api';
import { showToast } from './toast';
import { nowISO } from '../utils/date';

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

  const logDiaper = async (type: DiaperType) => {
    try {
      await api.createDiaper({ diaper_type: type, changed_at: nowISO() });
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

  modal.appendChild(grid);

  setTimeout(() => overlay.classList.add('open'), 10);
  overlay.addEventListener('click', (e) => { if (e.target === overlay) close(); });

  return overlay;
}
