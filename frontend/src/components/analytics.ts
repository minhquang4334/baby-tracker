import { h } from '../utils/dom';
import { api } from '../api';
import { renderNav } from './nav';
import { formatDuration, todayISO } from '../utils/date';
import type { DayStats } from '../types/models';

type Period = 7 | 14 | 30;

export function renderAnalyticsScreen(): HTMLElement {
  const screen = h('div', { class: 'screen analytics' });
  let period: Period = 7;
  let data: DayStats[] = [];

  // Header
  const header = h('div', { class: 'analytics-header' },
    h('div', { class: 'analytics-title' }, 'Analytics'),
  );

  // Period tabs
  const periodTabs = h('div', { class: 'period-tabs' });
  const periods: Period[] = [7, 14, 30];
  const tabEls: HTMLButtonElement[] = [];
  for (const p of periods) {
    const btn = h('button', {
      class: `period-tab${p === period ? ' active' : ''}`,
      onClick: () => {
        period = p;
        tabEls.forEach(b => b.classList.remove('active'));
        btn.classList.add('active');
        load();
      },
    }, `${p}d`) as HTMLButtonElement;
    tabEls.push(btn);
    periodTabs.appendChild(btn);
  }

  // Summary stat cards
  const statsRow = h('div', { class: 'analytics-stats' });

  // Chart sections
  const sleepChartSection = h('div', { class: 'analytics-chart-section' });
  const feedChartSection = h('div', { class: 'analytics-chart-section' });
  const diaperChartSection = h('div', { class: 'analytics-chart-section' });
  const bottleChartSection = h('div', { class: 'analytics-chart-section' });

  const render = () => {
    // Compute totals
    const totalSleepMinutes = data.reduce((s, d) => s + d.sleep_minutes, 0);
    const totalFeedings = data.reduce((s, d) => s + d.feeding_count, 0);
    const totalDiapers = data.reduce((s, d) => s + d.diaper_count, 0);
    const totalBottleML = data.reduce((s, d) => s + d.bottle_ml_total, 0);
    const bottleDays = data.filter(d => d.bottle_feed_count > 0).length || 1;
    const days = data.length || 1;

    const avgSleep = Math.round(totalSleepMinutes / days);
    const avgFeeding = (totalFeedings / days).toFixed(1);
    const avgDiaper = (totalDiapers / days).toFixed(1);
    const avgBottleL = totalBottleML > 0 ? (totalBottleML / bottleDays / 1000).toFixed(2) + ' L' : '—';

    statsRow.innerHTML = '';
    statsRow.appendChild(statCard('😴', formatDuration(avgSleep), 'Avg sleep/day'));
    statsRow.appendChild(statCard('🍼', avgFeeding, 'Avg feeds/day'));
    statsRow.appendChild(statCard('🚼', avgDiaper, 'Avg diapers/day'));
    statsRow.appendChild(statCard('🍶', avgBottleL, 'Avg bottle/day'));

    // Sleep chart
    sleepChartSection.innerHTML = '';
    const sleepCanvas = h('canvas', { height: '160' }) as HTMLCanvasElement;
    sleepChartSection.appendChild(h('div', { class: 'analytics-chart-title' }, '😴 Sleep (hours/day)'));
    sleepChartSection.appendChild(h('div', { class: 'analytics-chart-card' }, sleepCanvas));
    drawBarChart(sleepCanvas, data.map(d => d.date), data.map(d => +(d.sleep_minutes / 60).toFixed(1)), '#7C6AF0', 'h');

    // Feeding chart
    feedChartSection.innerHTML = '';
    const feedCanvas = h('canvas', { height: '160' }) as HTMLCanvasElement;
    feedChartSection.appendChild(h('div', { class: 'analytics-chart-title' }, '🍼 Feedings per day'));
    feedChartSection.appendChild(h('div', { class: 'analytics-chart-card' }, feedCanvas));
    drawBarChart(feedCanvas, data.map(d => d.date), data.map(d => d.feeding_count), '#E8507A', '');

    // Diaper chart (stacked wet + dirty)
    diaperChartSection.innerHTML = '';
    const diaperCanvas = h('canvas', { height: '160' }) as HTMLCanvasElement;
    diaperChartSection.appendChild(h('div', { class: 'analytics-chart-title' }, '🚼 Diaper changes per day'));
    diaperChartSection.appendChild(
      h('div', { class: 'analytics-chart-legend' },
        h('span', { class: 'legend-dot', style: 'background:#60A5FA' }),
        h('span', { class: 'legend-label' }, 'Wet'),
        h('span', { class: 'legend-dot', style: 'background:#A78BFA' }),
        h('span', { class: 'legend-label' }, 'Dirty'),
      ),
    );
    diaperChartSection.appendChild(h('div', { class: 'analytics-chart-card' }, diaperCanvas));
    drawStackedBarChart(
      diaperCanvas,
      data.map(d => d.date),
      data.map(d => d.wet_count),   '#60A5FA',
      data.map(d => d.dirty_count), '#A78BFA',
    );

    // Bottle chart (ml per day)
    bottleChartSection.innerHTML = '';
    const bottleCanvas = h('canvas', { height: '160' }) as HTMLCanvasElement;
    bottleChartSection.appendChild(h('div', { class: 'analytics-chart-title' }, '🍶 Bottle milk per day (ml)'));
    bottleChartSection.appendChild(h('div', { class: 'analytics-chart-card' }, bottleCanvas));
    drawBarChart(bottleCanvas, data.map(d => d.date), data.map(d => d.bottle_ml_total), '#F59E0B', 'ml');
  };

  const load = async () => {
    statsRow.innerHTML = '<div style="padding: 8px 0; color: var(--color-text-muted); font-size: 13px">Loading…</div>';
    sleepChartSection.innerHTML = '';
    feedChartSection.innerHTML = '';
    diaperChartSection.innerHTML = '';
    bottleChartSection.innerHTML = '';

    const today = todayISO();
    const fromDate = new Date(today + 'T12:00:00+07:00');
    fromDate.setDate(fromDate.getDate() - (period - 1));
    const from = fromDate.toISOString().slice(0, 10);

    try {
      data = await api.getAnalytics(from, today);
      render();
    } catch {
      statsRow.innerHTML = '<div style="padding: 8px; color: var(--color-text-muted)">Failed to load analytics.</div>';
    }
  };

  screen.appendChild(header);
  screen.appendChild(periodTabs);
  screen.appendChild(statsRow);
  screen.appendChild(sleepChartSection);
  screen.appendChild(feedChartSection);
  screen.appendChild(diaperChartSection);
  screen.appendChild(bottleChartSection);
  screen.appendChild(renderNav());

  load();
  return screen;
}

function statCard(icon: string, value: string, label: string): HTMLElement {
  return h('div', { class: 'analytics-stat-card' },
    h('div', { class: 'analytics-stat-icon' }, icon),
    h('div', { class: 'analytics-stat-value' }, value),
    h('div', { class: 'analytics-stat-label' }, label),
  );
}

interface ChartCtx {
  ctx: CanvasRenderingContext2D;
  padTop: number; padBottom: number; padLeft: number; padRight: number;
  chartW: number; chartH: number;
}

function setupChart(canvas: HTMLCanvasElement, maxVal: number, n: number): (ChartCtx & { barW: number; gap: number }) | null {
  const dpr = window.devicePixelRatio || 1;
  const cssW = canvas.parentElement?.clientWidth ?? 300;
  const cssH = parseInt(canvas.getAttribute('height') ?? '160');
  canvas.style.width = cssW + 'px';
  canvas.style.height = cssH + 'px';
  canvas.width = cssW * dpr;
  canvas.height = cssH * dpr;

  const ctx = canvas.getContext('2d');
  if (!ctx) return null;
  ctx.scale(dpr, dpr);

  const padTop = 20, padBottom = 28, padLeft = 36, padRight = 8;
  const chartW = cssW - padLeft - padRight;
  const chartH = cssH - padTop - padBottom;
  const barW = Math.max(4, Math.floor((chartW / n) * 0.6));
  const gap = (chartW - barW * n) / (n + 1);

  // Grid lines
  ctx.strokeStyle = 'rgba(0,0,0,0.06)';
  ctx.lineWidth = 1;
  for (let i = 0; i <= 4; i++) {
    const y = padTop + (chartH / 4) * i;
    ctx.beginPath(); ctx.moveTo(padLeft, y); ctx.lineTo(padLeft + chartW, y); ctx.stroke();
  }

  // Y axis labels
  ctx.fillStyle = 'rgba(0,0,0,0.35)';
  ctx.font = '10px sans-serif';
  ctx.textAlign = 'right';
  ctx.textBaseline = 'middle';
  for (let i = 0; i <= 2; i++) {
    const val = maxVal * (1 - i / 2);
    ctx.fillText(val % 1 === 0 ? String(Math.round(val)) : val.toFixed(1), padLeft - 4, padTop + (chartH / 2) * i);
  }

  return { ctx, padTop, padBottom, padLeft, padRight, chartW, chartH, barW, gap };
}

function drawXLabel(ctx: CanvasRenderingContext2D, label: string, x: number, barW: number, padTop: number, chartH: number, i: number, n: number) {
  if (n <= 10 || i === 0 || i === n - 1 || i % Math.ceil(n / 7) === 0) {
    ctx.fillStyle = 'rgba(0,0,0,0.35)';
    ctx.font = '10px sans-serif';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'top';
    ctx.fillText(label.slice(5), x + barW / 2, padTop + chartH + 6);
  }
}

function drawBarChart(canvas: HTMLCanvasElement, labels: string[], values: number[], color: string, unit: string) {
  const c = setupChart(canvas, Math.max(...values, 1), values.length);
  if (!c) return;
  const { ctx, padTop, chartH, barW, gap } = c;
  const maxVal = Math.max(...values, 1);

  ctx.textAlign = 'center';
  ctx.textBaseline = 'top';
  for (let i = 0; i < values.length; i++) {
    const x = c.padLeft + gap + i * (barW + gap);
    const barH = (values[i] / maxVal) * chartH;

    ctx.fillStyle = color;
    ctx.globalAlpha = 0.85;
    roundRect(ctx, x, padTop + chartH - barH, barW, barH, 3);
    ctx.fill();
    ctx.globalAlpha = 1;

    drawXLabel(ctx, labels[i] ?? '', x, barW, padTop, chartH, i, values.length);

    if (barH > 18 && values[i] > 0) {
      ctx.fillStyle = '#fff';
      ctx.font = 'bold 9px sans-serif';
      ctx.textBaseline = 'middle';
      ctx.fillText(`${values[i]}${unit}`, x + barW / 2, padTop + chartH - barH + Math.min(barH / 2, 10));
      ctx.textBaseline = 'top';
    }
  }
}

function drawStackedBarChart(
  canvas: HTMLCanvasElement,
  labels: string[],
  bottomValues: number[], bottomColor: string,
  topValues: number[],    topColor: string,
) {
  const totals = bottomValues.map((b, i) => b + topValues[i]);
  const c = setupChart(canvas, Math.max(...totals, 1), labels.length);
  if (!c) return;
  const { ctx, padTop, chartH, barW, gap } = c;
  const maxVal = Math.max(...totals, 1);

  for (let i = 0; i < labels.length; i++) {
    const x = c.padLeft + gap + i * (barW + gap);
    const bottomH = (bottomValues[i] / maxVal) * chartH;
    const topH    = (topValues[i]    / maxVal) * chartH;

    if (bottomH > 0) {
      ctx.fillStyle = bottomColor;
      ctx.globalAlpha = 0.85;
      roundRect(ctx, x, padTop + chartH - bottomH, barW, bottomH, topH > 0 ? 0 : 3);
      ctx.fill();
    }
    if (topH > 0) {
      ctx.fillStyle = topColor;
      ctx.globalAlpha = 0.85;
      roundRect(ctx, x, padTop + chartH - bottomH - topH, barW, topH, 3);
      ctx.fill();
    }
    ctx.globalAlpha = 1;

    const totalH = bottomH + topH;
    if (totalH > 18 && totals[i] > 0) {
      ctx.fillStyle = '#fff';
      ctx.font = 'bold 9px sans-serif';
      ctx.textBaseline = 'middle';
      ctx.fillText(String(totals[i]), x + barW / 2, padTop + chartH - totalH / 2);
    }

    drawXLabel(ctx, labels[i] ?? '', x, barW, padTop, chartH, i, labels.length);
  }
}

function roundRect(ctx: CanvasRenderingContext2D, x: number, y: number, w: number, h: number, r: number) {
  if (h < r * 2) r = h / 2;
  ctx.beginPath();
  ctx.moveTo(x + r, y);
  ctx.lineTo(x + w - r, y);
  ctx.quadraticCurveTo(x + w, y, x + w, y + r);
  ctx.lineTo(x + w, y + h);
  ctx.lineTo(x, y + h);
  ctx.lineTo(x, y + r);
  ctx.quadraticCurveTo(x, y, x + r, y);
  ctx.closePath();
}
