import { h } from '../utils/dom';
import { navigate, currentPath } from '../router';

interface NavItem {
  path: string;
  icon: string;
  label: string;
}

const NAV_ITEMS: NavItem[] = [
  { path: '/dashboard', icon: '🏠', label: 'Home' },
  { path: '/history', icon: '📋', label: 'History' },
  { path: '/analytics', icon: '📊', label: 'Analytics' },
  { path: '/growth', icon: '📈', label: 'Growth' },
  { path: '/guide', icon: '📖', label: 'Guide' },
];

export function renderNav(): HTMLElement {
  const current = currentPath();
  const nav = h('nav', { class: 'nav' });

  for (const item of NAV_ITEMS) {
    const isActive = current === item.path || (current === '/onboarding' && item.path === '/dashboard');
    const btn = h('button', {
      class: `nav-item${isActive ? ' active' : ''}`,
      onClick: () => navigate(item.path),
    },
      h('span', { class: 'nav-icon' }, item.icon),
      item.label,
    );
    nav.appendChild(btn);
  }

  return nav;
}
