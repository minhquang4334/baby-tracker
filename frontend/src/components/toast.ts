let toastEl: HTMLElement | null = null;
let hideTimer: ReturnType<typeof setTimeout> | null = null;

export function showToast(msg: string, type: 'success' | 'error' = 'success'): void {
  if (!toastEl) {
    toastEl = document.createElement('div');
    toastEl.className = 'toast';
    document.body.appendChild(toastEl);
  }

  if (hideTimer) clearTimeout(hideTimer);

  toastEl.textContent = msg;
  toastEl.className = `toast${type === 'error' ? ' error' : ''}`;

  // Force reflow
  void toastEl.offsetHeight;
  toastEl.classList.add('show');

  hideTimer = setTimeout(() => {
    toastEl?.classList.remove('show');
  }, 2500);
}
