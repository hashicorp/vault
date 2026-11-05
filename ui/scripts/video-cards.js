/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-env node */

/**
 * Renders the full-frame cards that separate a demo video's scenarios: an intro, one card per scenario, and a
 * closing summary. Shared by scripts/record-e2e-video.js and any tool that adds cards to a recorded demo, so
 * every card looks the same.
 *
 * ffmpeg builds without libass/drawtext can't draw text, so each card is a PNG screenshot of a small HTML page.
 * Text is always set through textContent, never as HTML.
 */

const { chromium } = require('@playwright/test');

// Card background, also used for fades between cards and clips.
const CARD_COLOR = '0x15181E';

/**
 * @typedef {object} Card
 * @property {string} file - PNG to write
 * @property {string} title
 * @property {string} [eyebrow] - small upper-case line above the title
 * @property {string} [subtitle]
 * @property {string[]} [items] - a list under the title
 * @property {'number' | 'check'} [marker] - list marker, numbered by default
 */

/**
 * Writes each card as a PNG of `width`×`height`, with its content centred in the top `contentHeight` pixels
 * (narrated cards leave the caption strip below free).
 *
 * @param {Card[]} cards
 * @param {{ width: number, height: number, contentHeight?: number, fontPx?: number }} size
 */
const renderCards = async (
  cards,
  { width, height, contentHeight = height, fontPx = Math.round(contentHeight * 0.036) }
) => {
  if (!cards.length) return;
  const browser = await chromium.launch();
  try {
    const page = await browser.newPage({ viewport: { width, height } });
    for (const card of cards) {
      await page.setContent(
        `<body style="margin:0;background:#${CARD_COLOR.slice(2)};width:${width}px;height:${height}px;` +
          'font-family:Arial,Helvetica,sans-serif;color:#fff">' +
          `<div style="height:${contentHeight}px;display:flex;flex-direction:column;align-items:center;justify-content:center">` +
          `<div id="eyebrow" style="color:#8C96A6;font:600 ${Math.round(
            fontPx * 0.8
          )}px Arial;letter-spacing:0.12em;` +
          'text-transform:uppercase;margin-bottom:22px"></div>' +
          `<h1 id="title" style="margin:0 120px;font:700 ${Math.round(
            fontPx * 2.1
          )}px/1.2 Arial;text-align:center"></h1>` +
          '<div style="width:72px;height:4px;background:#1C7DFF;border-radius:2px;margin:28px 0"></div>' +
          `<p id="subtitle" style="margin:0 160px;color:#D5D9E0;font:400 ${Math.round(
            fontPx * 1.05
          )}px/1.4 Arial;` +
          'text-align:center"></p>' +
          `<ol id="items" style="list-style:none;margin:0 160px;padding:0;color:#D5D9E0;` +
          `font:400 ${Math.round(fontPx * 1.05)}px/1.45 Arial"></ol></div></body>`
      );
      await page.evaluate((c) => {
        const set = (id, text) => {
          const el = document.getElementById(id);
          el.textContent = text || '';
          if (!text) el.style.display = 'none';
        };
        set('eyebrow', c.eyebrow);
        set('title', c.title);
        set('subtitle', c.subtitle);
        const list = document.getElementById('items');
        if (!c.items?.length) list.style.display = 'none';
        (c.items || []).forEach((item, i) => {
          const li = document.createElement('li');
          Object.assign(li.style, { display: 'flex', gap: '16px', margin: '10px 0' });
          const marker = document.createElement('span');
          Object.assign(marker.style, { color: '#1C7DFF', fontWeight: '700', minWidth: '1.2em' });
          marker.textContent = c.marker === 'check' ? '✓' : `${i + 1}`;
          const text = document.createElement('span');
          text.textContent = item;
          li.append(marker, text);
          list.append(li);
        });
      }, card);
      await page.screenshot({ path: card.file });
    }
  } finally {
    await browser.close();
  }
};

// Test titles are written for reports ("superuser can…"); a card shows them as a sentence-case heading.
const cardTitle = (testTitle) => {
  const text = String(testTitle || '').trim();
  return text ? text[0].toUpperCase() + text.slice(1) : '';
};

module.exports = { CARD_COLOR, cardTitle, renderCards };
