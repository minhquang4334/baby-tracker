import { h } from '../utils/dom';
import { renderSleepModal } from './sleep-modal';
import { renderFeedingModal } from './feeding-modal';
import { renderDiaperModal } from './diaper-modal';
import { renderGrowthModal } from './growth-modal';

interface FabItem {
  id: string;
  emoji: string;
  label: string;
  color: string;
  bg: string;
  action: (close: () => void, refresh: () => void) => void;
}

export function renderQuickAdd(onRefresh: () => void): HTMLElement {
  let isOpen = false;

  const overlay = h('div', { class: 'fab-overlay' });
  const fabBtn = h('button', { class: 'fab', 'aria-label': 'Quick add' }, '+');
  const menu = h('div', { class: 'fab-menu' });

  const toggle = () => {
    isOpen = !isOpen;
    fabBtn.classList.toggle('open', isOpen);
    menu.classList.toggle('open', isOpen);
    overlay.classList.toggle('open', isOpen);
  };

  const close = () => {
    isOpen = false;
    fabBtn.classList.remove('open');
    menu.classList.remove('open');
    overlay.classList.remove('open');
  };

  fabBtn.addEventListener('click', toggle);
  overlay.addEventListener('click', close);

  const items: FabItem[] = [
    {
      id: 'sleep', emoji: '😴', label: 'Sleep', color: '#8B5CF6', bg: '#F3EFFE',
      action: (closeFn, refresh) => {
        closeFn();
        const modal = renderSleepModal(refresh);
        document.getElementById('app')!.appendChild(modal);
      },
    },
    {
      id: 'feed', emoji: '🍼', label: 'Feed', color: '#EC4899', bg: '#FDE8F3',
      action: (closeFn, refresh) => {
        closeFn();
        const modal = renderFeedingModal(refresh);
        document.getElementById('app')!.appendChild(modal);
      },
    },
    {
      id: 'diaper', emoji: '🚼', label: 'Diaper', color: '#10B981', bg: '#D1FAF0',
      action: (closeFn, refresh) => {
        closeFn();
        const modal = renderDiaperModal(refresh);
        document.getElementById('app')!.appendChild(modal);
      },
    },
    {
      id: 'growth', emoji: '📏', label: 'Growth', color: '#F59E0B', bg: '#FEF3DC',
      action: (closeFn, refresh) => {
        closeFn();
        const modal = renderGrowthModal(refresh);
        document.getElementById('app')!.appendChild(modal);
      },
    },
  ];

  for (const item of items) {
    const row = h('div', { class: 'fab-item' },
      h('span', { class: 'fab-item-label' }, item.label),
      h('button', {
        class: 'fab-item-btn',
        style: `background:${item.bg}; color:${item.color}`,
        onClick: () => item.action(close, onRefresh),
      }, item.emoji),
    );
    menu.appendChild(row);
  }

  const container = h('div', { class: 'fab-container' });
  container.appendChild(overlay);
  container.appendChild(menu);
  container.appendChild(fabBtn);
  return container;
}
