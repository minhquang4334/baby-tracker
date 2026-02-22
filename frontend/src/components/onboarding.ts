import { h } from '../utils/dom';
import { api } from '../api';
import { state } from '../state';
import { navigate } from '../router';
import { showToast } from './toast';

export function renderOnboarding(): HTMLElement {
  let selectedGender = 'female';

  const screen = h('div', { class: 'screen onboarding' });

  const header = h('div', { class: 'onboarding-header' },
    h('div', { class: 'onboarding-emoji' }, '👶'),
    h('h1', { class: 'onboarding-title' }, 'Welcome!'),
    h('p', { class: 'onboarding-subtitle' }, 'Set up your baby\'s profile to get started tracking.'),
  );

  const nameInput = h('input', {
    class: 'form-input',
    type: 'text',
    placeholder: 'e.g. Emma',
    id: 'baby-name',
  }) as HTMLInputElement;

  const dobInput = h('input', {
    class: 'form-input',
    type: 'date',
    id: 'baby-dob',
    max: new Date().toISOString().slice(0, 10),
  }) as HTMLInputElement;

  const genderBtns: HTMLButtonElement[] = [];
  const genderGroup = h('div', { class: 'gender-pills' });
  for (const [val, label] of [['female', 'Girl'], ['male', 'Boy'], ['other', 'Other']] as const) {
    const btn = h('button', {
      class: `gender-pill${val === selectedGender ? ' selected' : ''}`,
      type: 'button',
      onClick: () => {
        selectedGender = val;
        genderBtns.forEach(b => b.classList.remove('selected'));
        btn.classList.add('selected');
      },
    }, label) as HTMLButtonElement;
    genderBtns.push(btn);
    genderGroup.appendChild(btn);
  }

  const submitBtn = h('button', {
    class: 'btn btn-primary btn-full',
    type: 'button',
    onClick: async () => {
      const name = nameInput.value.trim();
      const dob = dobInput.value;
      if (!name) { showToast('Please enter a name', 'error'); return; }
      if (!dob) { showToast('Please enter a date of birth', 'error'); return; }

      submitBtn.disabled = true;
      submitBtn.innerHTML = '<div class="spinner"></div>';

      try {
        const child = await api.createChild({ name, date_of_birth: dob, gender: selectedGender });
        state.child.set(child);
        navigate('/dashboard');
      } catch (e: any) {
        showToast(e.message ?? 'Something went wrong', 'error');
        submitBtn.disabled = false;
        submitBtn.textContent = 'Get Started';
      }
    },
  }, 'Get Started') as HTMLButtonElement;

  const form = h('div', { class: 'onboarding-form' },
    h('div', { class: 'form-group' },
      h('label', { class: 'form-label', for: 'baby-name' }, 'Baby\'s name'),
      nameInput,
    ),
    h('div', { class: 'form-group' },
      h('label', { class: 'form-label', for: 'baby-dob' }, 'Date of birth'),
      dobInput,
    ),
    h('div', { class: 'form-group' },
      h('label', { class: 'form-label' }, 'Gender'),
      genderGroup,
    ),
    h('div', { class: 'onboarding-footer' },
      submitBtn,
    ),
  );

  screen.appendChild(header);
  screen.appendChild(form);
  return screen;
}
