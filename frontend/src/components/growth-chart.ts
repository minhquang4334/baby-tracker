import { h } from '../utils/dom';
import { api } from '../api';
import { renderNav } from './nav';
import { renderQuickAdd } from './quick-add';
import { showToast } from './toast';
import type { GrowthLog } from '../types/models';

type ChartMetric = 'weight' | 'length';

export function renderGrowthScreen(): HTMLElement {
  const screen = h('div', { class: 'screen growth-screen' });
  let currentMetric: ChartMetric = 'weight';
  let logs: GrowthLog[] = [];

  const titleEl = h('div', { style: 'padding: 16px 20px 0; font-size: 20px; font-weight: 700' }, 'Growth');

  const toggleWt = h('button', {
    class: 'pill active pill-growth',
    onClick: () => setMetric('weight'),
  }, 'Weight') as HTMLButtonElement;

  const toggleLen = h('button', {
    class: 'pill',
    onClick: () => setMetric('length'),
  }, 'Length') as HTMLButtonElement;

  const canvasWrap = h('div', { class: 'chart-wrap' });
  const canvas = h('canvas', { class: 'chart-canvas', height: '200' }) as HTMLCanvasElement;
  canvasWrap.appendChild(canvas);

  const tableWrap = h('div', {});

  const setMetric = (m: ChartMetric) => {
    currentMetric = m;
    toggleWt.classList.toggle('active', m === 'weight');
    toggleLen.classList.toggle('active', m === 'length');
    render();
  };

  const render = () => {
    drawChart(canvas, logs, currentMetric);
    renderTable(tableWrap, logs);
  };

  const load = async () => {
    try {
      logs = await api.getGrowth();
      render();
    } catch (e: any) {
      showToast('Failed to load growth data', 'error');
    }
  };

  const refresh = () => load();

  screen.appendChild(titleEl);
  screen.appendChild(h('div', { class: 'chart-toggle' }, toggleWt, toggleLen));
  screen.appendChild(canvasWrap);
  screen.appendChild(tableWrap);
  screen.appendChild(renderNav());
  screen.appendChild(renderQuickAdd(refresh));

  load();
  return screen;
}

function drawChart(canvas: HTMLCanvasElement, logs: GrowthLog[], metric: ChartMetric): void {
  const parent = canvas.parentElement;
  if (!parent) return;

  const W = parent.clientWidth - 16;
  const H = 200;
  canvas.width = W;
  canvas.height = H;

  const ctx = canvas.getContext('2d');
  if (!ctx) return;

  ctx.clearRect(0, 0, W, H);

  const data = logs
    .filter(l => metric === 'weight' ? l.weight_grams != null : l.length_mm != null)
    .map(l => ({
      date: l.measured_on,
      val: metric === 'weight' ? l.weight_grams! : l.length_mm!,
    }))
    .sort((a, b) => a.date.localeCompare(b.date));

  if (data.length === 0) {
    ctx.fillStyle = '#9E8E89';
    ctx.font = '14px Inter, sans-serif';
    ctx.textAlign = 'center';
    ctx.fillText('No data yet — add measurements using the + button', W / 2, H / 2);
    return;
  }

  const pad = { top: 20, right: 16, bottom: 40, left: 50 };
  const chartW = W - pad.left - pad.right;
  const chartH = H - pad.top - pad.bottom;

  const vals = data.map(d => d.val);
  const minV = Math.min(...vals);
  const maxV = Math.max(...vals);
  const rangeV = maxV - minV || 1;

  const xScale = (i: number) => pad.left + (i / Math.max(data.length - 1, 1)) * chartW;
  const yScale = (v: number) => pad.top + chartH - ((v - minV) / rangeV) * chartH;

  // Grid lines
  ctx.strokeStyle = '#F0EBE8';
  ctx.lineWidth = 1;
  for (let i = 0; i <= 4; i++) {
    const y = pad.top + (i / 4) * chartH;
    ctx.beginPath();
    ctx.moveTo(pad.left, y);
    ctx.lineTo(pad.left + chartW, y);
    ctx.stroke();

    const labelV = maxV - (i / 4) * rangeV;
    ctx.fillStyle = '#9E8E89';
    ctx.font = '11px Inter, sans-serif';
    ctx.textAlign = 'right';
    ctx.fillText(metric === 'weight' ? `${(labelV / 1000).toFixed(1)}kg` : `${(labelV / 10).toFixed(1)}cm`, pad.left - 4, y + 4);
  }

  // Line
  ctx.strokeStyle = '#F59E0B';
  ctx.lineWidth = 2.5;
  ctx.lineJoin = 'round';
  ctx.beginPath();
  data.forEach((d, i) => {
    const x = xScale(i);
    const y = yScale(d.val);
    if (i === 0) ctx.moveTo(x, y); else ctx.lineTo(x, y);
  });
  ctx.stroke();

  // Fill under line
  ctx.fillStyle = 'rgba(245, 158, 11, 0.1)';
  ctx.beginPath();
  data.forEach((d, i) => {
    const x = xScale(i);
    const y = yScale(d.val);
    if (i === 0) ctx.moveTo(x, y); else ctx.lineTo(x, y);
  });
  ctx.lineTo(xScale(data.length - 1), pad.top + chartH);
  ctx.lineTo(xScale(0), pad.top + chartH);
  ctx.closePath();
  ctx.fill();

  // Dots
  ctx.fillStyle = '#F59E0B';
  data.forEach((d, i) => {
    ctx.beginPath();
    ctx.arc(xScale(i), yScale(d.val), 4, 0, Math.PI * 2);
    ctx.fill();
  });

  // X labels (every nth)
  const step = Math.ceil(data.length / 5);
  ctx.fillStyle = '#9E8E89';
  ctx.font = '11px Inter, sans-serif';
  ctx.textAlign = 'center';
  data.forEach((d, i) => {
    if (i % step !== 0 && i !== data.length - 1) return;
    const label = new Date(d.date + 'T12:00:00').toLocaleDateString([], { month: 'short', day: 'numeric' });
    ctx.fillText(label, xScale(i), H - 8);
  });
}

function renderTable(wrap: HTMLElement, logs: GrowthLog[]): void {
  wrap.innerHTML = '';
  if (logs.length === 0) return;

  const sorted = [...logs].sort((a, b) => b.measured_on.localeCompare(a.measured_on));

  const table = h('div', { class: 'growth-table' },
    h('div', { class: 'growth-table-header' },
      h('span', {}, 'Date'),
      h('span', {}, 'Weight'),
      h('span', {}, 'Length'),
      h('span', {}, 'Head'),
    ),
  );

  for (const l of sorted) {
    const date = new Date(l.measured_on + 'T12:00:00').toLocaleDateString([], { month: 'short', day: 'numeric' });
    const weight = l.weight_grams != null ? `${(l.weight_grams / 1000).toFixed(2)}kg` : '—';
    const length = l.length_mm != null ? `${(l.length_mm / 10).toFixed(1)}cm` : '—';
    const head = l.head_circumference_mm != null ? `${(l.head_circumference_mm / 10).toFixed(1)}cm` : '—';

    table.appendChild(h('div', { class: 'growth-table-row' },
      h('span', {}, date),
      h('span', {}, weight),
      h('span', {}, length),
      h('span', {}, head),
    ));
  }

  wrap.appendChild(table);
}
