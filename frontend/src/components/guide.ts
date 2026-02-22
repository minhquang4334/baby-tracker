import { h } from '../utils/dom';
import { renderNav } from './nav';

export function renderGuideScreen(): HTMLElement {
  const screen = h('div', { class: 'screen guide-screen' });
  screen.appendChild(h('div', { class: 'guide-title' }, '📖 Hướng Dẫn Phát Triển Trẻ'));
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
  wrap.appendChild(section('🍼', 'Lịch Bú, Ngủ & Tã — Theo Tháng',
    tbl(
      ['Tháng', 'Số lần bú', 'Lượng/lần', 'Ngủ/ngày', 'Tã ướt', 'Tã bẩn'],
      [
        ['0', '8–12 lần', '30–60 ml', '16–18 h', '6–8', '3–4'],
        ['1', '8–10 lần', '60–90 ml', '15–17 h', '6–8', '3–4'],
        ['2', '7–9 lần', '90–120 ml', '14–16 h', '5–6', '2–3'],
        ['3', '6–8 lần', '120–150 ml', '14–16 h', '5–6', '2–3'],
        ['4', '6–7 lần', '120–180 ml', '14–15 h', '4–6', '1–3'],
        ['5', '5–6 lần', '150–180 ml', '14–15 h', '4–6', '1–3'],
        ['6', '4–5 + dặm', '150–210 ml', '13–15 h', '4–6', '1–2'],
        ['7', '4–5 + dặm', '180–210 ml', '13–14 h', '4–6', '1–2'],
        ['8', '4–5 + dặm', '180–210 ml', '13–14 h', '4–6', '1–2'],
        ['9', '3–4 + dặm', '180–240 ml', '12–14 h', '4–5', '1–2'],
        ['10', '3–4 + dặm', '180–240 ml', '12–14 h', '4–5', '1–2'],
        ['11', '3–4 + dặm', '180–240 ml', '12–14 h', '4–5', '1–2'],
        ['12', '3–4 + dặm', '180–240 ml', '12–14 h', '4–5', '1–2'],
      ],
    ),
    note('💡 Ăn dặm bắt đầu từ tháng 6 theo khuyến cáo Bộ Y tế & WHO. Bú mẹ hoàn toàn: theo nhu cầu (on-demand).'),
  ));

  // 2. WHO weight / height
  const whoHeaders = ['Tháng', 'TB (kg)', '-2SD', '+2SD', 'TB (cm)', '-2SD', '+2SD'];
  wrap.appendChild(section('📏', 'Cân Nặng & Chiều Cao Chuẩn WHO',
    h('p', { class: 'ref-sub-title' }, '🩷 Bé Gái'),
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
    h('p', { class: 'ref-sub-title', style: 'margin-top:12px' }, '💙 Bé Trai'),
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
    note('TB = trung bình · -2SD = ngưỡng thấp (cần gặp bác sĩ) · +2SD = ngưỡng cao'),
  ));

  // 3. Milestones
  wrap.appendChild(section('📌', 'Milestones Phát Triển Quan Trọng',
    tbl(
      ['Tháng', 'Vận động', 'Ngôn ngữ', 'Xã hội'],
      [
        ['1–2',   'Ngóc đầu khi nằm sấp',        'Phát âm "ê, a"',         'Nhìn theo mặt người'],
        ['3–4',   'Giữ đầu vững, lật người',      'Cười to',                'Nhận ra bố mẹ'],
        ['5–6',   'Ngồi có đỡ, cầm đồ vật',       'Bập bẹ "ba ba, ma ma"',  'Bắt đầu sợ người lạ'],
        ['7–9',   'Bò, kéo đứng',                 'Hiểu "không"',           'Chơi peek-a-boo'],
        ['10–12', 'Đứng vịn, bước đi đầu tiên',   '1–2 từ có nghĩa',        'Vỗ tay, vẫy tay'],
      ],
    ),
  ));

  // 4. Vaccination
  wrap.appendChild(section('💉', 'Lịch Tiêm Chủng Mở Rộng (0–12 tháng)',
    tbl(
      ['Tháng tuổi', 'Vắc-xin'],
      [
        ['Sơ sinh (24h)', 'Viêm gan B mũi 1, BCG (lao)'],
        ['2 tháng',       '5 trong 1 (DPT-VGB-Hib) mũi 1, OPV1'],
        ['3 tháng',       '5 trong 1 mũi 2, OPV2'],
        ['4 tháng',       '5 trong 1 mũi 3, OPV3'],
        ['9 tháng',       'Sởi mũi 1'],
        ['12 tháng',      'Sởi – Rubella (MR)'],
      ],
    ),
    note('Tham khảo thêm tại VNVC hoặc trạm y tế phường/xã gần nhất.'),
  ));

  // 5. Diaper sizes
  wrap.appendChild(section('🛒', 'Cỡ Tã Bỉm Theo Tháng',
    tbl(
      ['Size', 'Cân nặng', 'Tháng tham khảo'],
      [
        ['NB (Newborn)', '< 5 kg',   '0–1 tháng'],
        ['S',            '3–8 kg',   '1–3 tháng'],
        ['M',            '6–11 kg',  '3–8 tháng'],
        ['L',            '9–14 kg',  '8–12 tháng+'],
      ],
    ),
  ));

  // 6. Warning signs
  const warnList = h('ul', { class: 'ref-warn-list' });
  [
    'Cân nặng dưới -2SD liên tục 2 tháng',
    'Bé không tăng cân sau 2 tuần đầu',
    'Số tã ướt < 6/ngày trong tháng đầu',
    'Bé không phản ứng với âm thanh sau 3 tháng',
    'Sốt > 38°C ở trẻ dưới 3 tháng — đến viện ngay',
    'Không biết lật sau 5 tháng, không ngồi được sau 9 tháng',
  ].forEach(w => warnList.appendChild(h('li', {}, w)));
  wrap.appendChild(section('⚠️', 'Khi Nào Cần Gặp Bác Sĩ?', warnList));

  // 7. Sources
  const srcList = h('ul', { class: 'ref-src-list' });
  ([
    ['Viện Dinh Dưỡng Quốc Gia VN', 'http://viendinhduong.vn'],
    ['WHO Child Growth Standards', 'https://www.who.int/tools/child-growth-standards'],
    ['Bộ Y tế VN — Hướng dẫn nuôi con bằng sữa mẹ', 'https://moh.gov.vn'],
  ] as [string, string][]).forEach(([label, url]) =>
    srcList.appendChild(h('li', {}, h('a', { href: url, target: '_blank', rel: 'noopener' }, label)))
  );
  wrap.appendChild(section('📚', 'Nguồn Tham Khảo', srcList));

  wrap.appendChild(h('p', { class: 'ref-disclaimer' },
    'Tài liệu mang tính tham khảo. Luôn tham vấn bác sĩ nhi khoa cho tình trạng cụ thể của bé.',
  ));

  return wrap;
}
