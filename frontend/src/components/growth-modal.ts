import { h } from '../utils/dom';
import { api } from '../api';
import { showToast } from './toast';
import { todayISO } from '../utils/date';

export function renderGrowthModal(onSave: () => void): HTMLElement {
  const overlay = h('div', { class: 'modal-overlay' });
  const modal = h('div', { class: 'modal' });

  const close = () => {
    overlay.classList.remove('open');
    setTimeout(() => overlay.remove(), 300);
  };

  overlay.appendChild(modal);
  modal.appendChild(h('div', { class: 'modal-handle' }));
  modal.appendChild(h('div', { class: 'modal-header' },
    h('h2', { class: 'modal-title' }, 'Growth Measurement'),
    h('button', { class: 'modal-close', onClick: close }, '×'),
  ));

  const dateInput = h('input', {
    class: 'form-input',
    type: 'date',
    value: todayISO(),
    max: todayISO(),
  }) as HTMLInputElement;

  const weightInput = h('input', {
    class: 'form-input',
    type: 'number',
    placeholder: 'e.g. 4200',
    min: '0',
    step: '10',
  }) as HTMLInputElement;

  const lengthInput = h('input', {
    class: 'form-input',
    type: 'number',
    placeholder: 'e.g. 520',
    min: '0',
    step: '1',
  }) as HTMLInputElement;

  const headInput = h('input', {
    class: 'form-input',
    type: 'number',
    placeholder: 'e.g. 340',
    min: '0',
    step: '1',
  }) as HTMLInputElement;

  const saveBtn = h('button', {
    class: 'btn btn-primary btn-full',
    onClick: async () => {
      const date = dateInput.value;
      if (!date) { showToast('Date is required', 'error'); return; }

      const toInt = (s: string) => s ? parseInt(s, 10) : undefined;
      const weight = toInt(weightInput.value);
      const length = toInt(lengthInput.value);
      const head = toInt(headInput.value);

      if (!weight && !length && !head) {
        showToast('Enter at least one measurement', 'error');
        return;
      }

      saveBtn.disabled = true;
      saveBtn.innerHTML = '<div class="spinner"></div>';

      try {
        await api.createGrowth({
          measured_on: date,
          weight_grams: weight ?? null,
          length_mm: length ?? null,
          head_circumference_mm: head ?? null,
        });
        showToast('Growth recorded');
        onSave();
        close();
      } catch (e: any) {
        showToast(e.message, 'error');
        saveBtn.disabled = false;
        saveBtn.textContent = 'Save';
      }
    },
  }, 'Save');

  modal.appendChild(h('div', { class: 'modal-body' },
    h('div', { class: 'form-group' },
      h('label', { class: 'form-label' }, 'Date'),
      dateInput,
    ),
    h('div', { class: 'form-group' },
      h('label', { class: 'form-label' }, 'Weight (grams)'),
      weightInput,
    ),
    h('div', { class: 'form-group' },
      h('label', { class: 'form-label' }, 'Length (mm)'),
      lengthInput,
    ),
    h('div', { class: 'form-group' },
      h('label', { class: 'form-label' }, 'Head circumference (mm)'),
      headInput,
    ),
  ));

  modal.appendChild(h('div', { class: 'modal-actions' }, saveBtn));

  setTimeout(() => overlay.classList.add('open'), 10);
  overlay.addEventListener('click', (e) => { if (e.target === overlay) close(); });

  return overlay;
}
