type RouteHandler = () => HTMLElement;

const routes = new Map<string, RouteHandler>();
let currentRoute = '';

export function addRoute(path: string, handler: RouteHandler): void {
  routes.set(path, handler);
}

export function navigate(path: string): void {
  window.location.hash = '#' + path;
}

export function initRouter(container: HTMLElement): void {
  const render = () => {
    const hash = window.location.hash.slice(1) || '/onboarding';
    const path = hash.split('?')[0];
    if (path === currentRoute) return;
    currentRoute = path;

    const handler = routes.get(path) ?? routes.get('/dashboard');
    if (!handler) return;

    const el = handler();
    container.innerHTML = '';
    container.appendChild(el);
  };

  window.addEventListener('hashchange', render);
  render();
}

export function currentPath(): string {
  return window.location.hash.slice(1).split('?')[0] || '/onboarding';
}
