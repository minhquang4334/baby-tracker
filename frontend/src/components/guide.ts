import { h } from '../utils/dom';
import { renderNav } from './nav';

export function renderGuideScreen(): HTMLElement {
  const screen = h('div', { class: 'screen guide-screen' });
  screen.appendChild(h('div', { class: 'guide-title' }, '📖 Baby Development Guide'));
  screen.appendChild(renderReferenceGuide());
  screen.appendChild(renderNav());
  return screen;
}

function renderReferenceGuide(): HTMLElement {
  const wrap = h('div', { class: 'ref-guide' });

  const section = (emoji: string, title: string, ...children: HTMLElement[]) =>
    h('div', { class: 'ref-section' },
      h('div', { class: 'ref-section-title' }, `${emoji} ${title}`),
      ...children,
    );

  const note = (text: string) => h('p', { class: 'ref-note' }, text);

  const tbl = (headers: string[], rows: string[][], highlight?: number) => {
    const t = h('div', { class: 'ref-table-wrap' });
    const head = h('div', { class: 'ref-table-row ref-table-head' });
    headers.forEach(hd => head.appendChild(h('span', {}, hd)));
    t.appendChild(head);
    rows.forEach((row, ri) => {
      const tr = h('div', { class: `ref-table-row${ri % 2 === 1 ? ' ref-table-row-alt' : ''}` });
      row.forEach((cell, ci) =>
        tr.appendChild(h('span', { class: ci === highlight ? 'ref-cell-hi' : '' }, cell))
      );
      t.appendChild(tr);
    });
    return t;
  };

  // 1. Feeding / sleep / diaper schedule
  wrap.appendChild(section('🍼', 'Feeding, Sleep & Diapers by Month',
    tbl(
      ['Month', 'Feeds/day', 'Per feed', 'Sleep/day', 'Wet diapers', 'Dirty diapers'],
      [
        ['0',  '8–12',        '30–60 ml',  '16–18 h', '6–8', '3–4'],
        ['1',  '8–10',        '60–90 ml',  '15–17 h', '6–8', '3–4'],
        ['2',  '7–9',         '90–120 ml', '14–16 h', '5–6', '2–3'],
        ['3',  '6–8',         '120–150 ml','14–16 h', '5–6', '2–3'],
        ['4',  '6–7',         '120–180 ml','14–15 h', '4–6', '1–3'],
        ['5',  '5–6',         '150–180 ml','14–15 h', '4–6', '1–3'],
        ['6',  '4–5 + solids','150–210 ml','13–15 h', '4–6', '1–2'],
        ['7',  '4–5 + solids','180–210 ml','13–14 h', '4–6', '1–2'],
        ['8',  '4–5 + solids','180–210 ml','13–14 h', '4–6', '1–2'],
        ['9',  '3–4 + solids','180–240 ml','12–14 h', '4–5', '1–2'],
        ['10', '3–4 + solids','180–240 ml','12–14 h', '4–5', '1–2'],
        ['11', '3–4 + solids','180–240 ml','12–14 h', '4–5', '1–2'],
        ['12', '3–4 + solids','180–240 ml','12–14 h', '4–5', '1–2'],
      ],
    ),
    note('💡 Solid foods are recommended from month 6 (WHO & AAP guidelines). Breastfed babies feed on demand.'),
  ));

  // 2. WHO weight / height
  const whoHeaders = ['Month', 'Avg (kg)', '-2SD', '+2SD', 'Avg (cm)', '-2SD', '+2SD'];
  wrap.appendChild(section('📏', 'WHO Weight & Length Standards',
    h('p', { class: 'ref-sub-title' }, '🩷 Girls'),
    tbl(whoHeaders, [
      ['0',  '3.2', '2.4', '4.2',  '49.1', '45.6', '52.7'],
      ['1',  '4.2', '3.2', '5.5',  '53.7', '50.0', '57.4'],
      ['2',  '5.1', '3.9', '6.6',  '57.1', '53.2', '61.1'],
      ['3',  '5.8', '4.5', '7.5',  '59.8', '55.8', '63.8'],
      ['4',  '6.4', '5.0', '8.1',  '62.1', '58.0', '66.2'],
      ['5',  '6.9', '5.4', '8.8',  '64.0', '59.9', '68.2'],
      ['6',  '7.3', '5.7', '9.3',  '65.7', '61.5', '70.0'],
      ['7',  '7.6', '6.0', '9.8',  '67.3', '63.0', '71.6'],
      ['8',  '7.9', '6.3', '10.2', '68.7', '64.4', '73.2'],
      ['9',  '8.2', '6.5', '10.5', '70.1', '65.6', '74.7'],
      ['10', '8.5', '6.7', '10.9', '71.5', '66.8', '76.2'],
      ['11', '8.7', '6.9', '11.2', '72.8', '68.0', '77.6'],
      ['12', '8.9', '7.1', '11.5', '74.0', '69.2', '78.9'],
    ], 1),
    h('p', { class: 'ref-sub-title', style: 'margin-top:12px' }, '💙 Boys'),
    tbl(whoHeaders, [
      ['0',  '3.3', '2.5', '4.4',  '49.9', '46.3', '53.4'],
      ['1',  '4.5', '3.4', '5.8',  '54.7', '50.8', '58.6'],
      ['2',  '5.6', '4.3', '7.1',  '58.4', '54.4', '62.4'],
      ['3',  '6.4', '5.0', '8.0',  '61.4', '57.3', '65.5'],
      ['4',  '7.0', '5.5', '8.7',  '63.9', '59.7', '68.0'],
      ['5',  '7.5', '5.9', '9.3',  '65.9', '61.7', '70.1'],
      ['6',  '7.9', '6.2', '9.8',  '67.6', '63.3', '71.9'],
      ['7',  '8.3', '6.5', '10.3', '69.2', '64.8', '73.5'],
      ['8',  '8.6', '6.8', '10.7', '70.6', '66.2', '75.0'],
      ['9',  '8.9', '7.1', '11.0', '72.0', '67.5', '76.5'],
      ['10', '9.2', '7.3', '11.4', '73.3', '68.7', '77.9'],
      ['11', '9.4', '7.5', '11.7', '74.5', '69.9', '79.2'],
      ['12', '9.6', '7.7', '12.0', '75.7', '71.0', '80.5'],
    ], 1),
    note('Avg = average · -2SD = low threshold (consult a doctor) · +2SD = high threshold'),
  ));

  // 3. Milestones
  wrap.appendChild(section('📌', 'Key Development Milestones',
    tbl(
      ['Month', 'Motor', 'Language', 'Social'],
      [
        ['1–2',   'Lifts head during tummy time',  'Cooing sounds',           'Follows faces with eyes'],
        ['3–4',   'Holds head steady, rolls over', 'Laughs out loud',         'Recognises parents'],
        ['5–6',   'Sits with support, grasps toys','Babbles "ba-ba, ma-ma"',  'Stranger anxiety begins'],
        ['7–9',   'Crawls, pulls to stand',        'Understands "no"',        'Plays peek-a-boo'],
        ['10–12', 'Stands holding on, first steps','1–2 meaningful words',    'Claps hands, waves bye'],
      ],
    ),
  ));

  // 4. Vaccination
  wrap.appendChild(section('💉', 'Vaccination Schedule (0–12 months)',
    tbl(
      ['Age', 'Vaccines'],
      [
        ['Birth (24h)',  'Hepatitis B dose 1, BCG (tuberculosis)'],
        ['2 months',    'DTaP-HepB-Hib dose 1, OPV1, PCV1'],
        ['3 months',    'DTaP-HepB-Hib dose 2, OPV2, PCV2'],
        ['4 months',    'DTaP-HepB-Hib dose 3, OPV3, PCV3'],
        ['9 months',    'Measles dose 1'],
        ['12 months',   'MMR (Measles–Mumps–Rubella)'],
      ],
    ),
    note('Schedules may vary by country. Always follow your paediatrician\'s recommendations.'),
  ));

  // 5. Diaper sizes
  wrap.appendChild(section('🛒', 'Diaper Size Guide',
    tbl(
      ['Size', 'Weight', 'Typical age'],
      [
        ['NB (Newborn)', '< 5 kg',   '0–1 month'],
        ['S',            '3–8 kg',   '1–3 months'],
        ['M',            '6–11 kg',  '3–8 months'],
        ['L',            '9–14 kg',  '8–12 months+'],
      ],
    ),
  ));

  // 6. Warning signs
  const warnList = h('ul', { class: 'ref-warn-list' });
  [
    'Weight below -2SD for two consecutive months',
    'No weight gain after the first 2 weeks',
    'Fewer than 6 wet diapers per day in the first month',
    'No reaction to sounds by 3 months',
    'Fever > 38 °C in a baby under 3 months — seek care immediately',
    'Not rolling by 5 months, not sitting independently by 9 months',
  ].forEach(w => warnList.appendChild(h('li', {}, w)));
  wrap.appendChild(section('⚠️', 'When to See a Doctor', warnList));

  // 7. Sources
  const srcList = h('ul', { class: 'ref-src-list' });
  ([
    ['WHO Child Growth Standards', 'https://www.who.int/tools/child-growth-standards'],
    ['CDC Developmental Milestones', 'https://www.cdc.gov/ncbddd/actearly/milestones/index.html'],
    ['AAP — Infant Feeding Guidelines', 'https://www.healthychildren.org'],
  ] as [string, string][]).forEach(([label, url]) =>
    srcList.appendChild(h('li', {}, h('a', { href: url, target: '_blank', rel: 'noopener' }, label)))
  );
  wrap.appendChild(section('📚', 'Sources', srcList));

  wrap.appendChild(h('p', { class: 'ref-disclaimer' },
    'For reference only. Always consult your paediatrician for advice specific to your baby.',
  ));

  return wrap;
}
