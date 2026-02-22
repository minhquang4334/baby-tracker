type TagName = keyof HTMLElementTagNameMap;
type Props = Record<string, string | number | boolean | EventListener | null | undefined>;

export function h<T extends TagName>(
  tag: T,
  props?: Props,
  ...children: (HTMLElement | string | null | undefined)[]
): HTMLElementTagNameMap[T] {
  const el = document.createElement(tag);
  if (props) {
    for (const [key, val] of Object.entries(props)) {
      if (val == null) continue;
      if (key.startsWith('on') && typeof val === 'function') {
        el.addEventListener(key.slice(2).toLowerCase(), val as EventListener);
      } else if (key === 'class' || key === 'className') {
        el.className = String(val);
      } else if (key === 'style' && typeof val === 'string') {
        el.setAttribute('style', val);
      } else if (typeof val === 'boolean') {
        if (val) el.setAttribute(key, '');
      } else {
        el.setAttribute(key, String(val));
      }
    }
  }
  for (const child of children) {
    if (child == null) continue;
    if (typeof child === 'string') {
      el.appendChild(document.createTextNode(child));
    } else {
      el.appendChild(child);
    }
  }
  return el;
}

export function $(selector: string, root: ParentNode = document): Element | null {
  return root.querySelector(selector);
}

export function $$(selector: string, root: ParentNode = document): Element[] {
  return Array.from(root.querySelectorAll(selector));
}

export function mount(el: HTMLElement, target: string | HTMLElement): void {
  const container = typeof target === 'string'
    ? document.querySelector(target) as HTMLElement
    : target;
  if (!container) throw new Error(`Mount target not found: ${target}`);
  container.innerHTML = '';
  container.appendChild(el);
}
