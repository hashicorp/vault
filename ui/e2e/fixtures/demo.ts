/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * Drop-in replacement for the Playwright `test` import.
 *
 * Identical to `@playwright/test` during normal runs. When PW_VIDEO is set it also paints a
 * synthetic cursor and click ripple so a recorded run is followable. Import this instead of
 * `@playwright/test` in any spec you may want to demo.
 */
import { test as base, expect } from '@playwright/test';
import { cursorInitScript, cursorTimings, isVideoEnabled } from '../video-config';

import type { Locator, Page } from '@playwright/test';

/**
 * Pauses a recording so a result stays on screen long enough to read. No-op when not recording.
 */
export const holdForVideo = async (page: Page, ms = 1_500) => {
  if (isVideoEnabled) await page.waitForTimeout(ms);
};

/**
 * Drags the recorded cursor across text inside `scope`, from the start of `fromText` to the end
 * of `toText`, so viewers can find it, then holds. No-op when not recording.
 */
export const highlightText = async (page: Page, scope: Locator, fromText: string, toText = fromText) => {
  if (!isVideoEnabled) return;
  const bounds = await scope.evaluate(
    (root, { fromText, toText }) => {
      const walker = document.createTreeWalker(root, NodeFilter.SHOW_TEXT);
      let start: { node: Node; offset: number } | undefined;
      let end: { node: Node; offset: number } | undefined;
      for (let node = walker.nextNode(); node && !end; node = walker.nextNode()) {
        const text = node.textContent ?? '';
        if (!start && text.includes(fromText)) start = { node, offset: text.indexOf(fromText) };
        if (start && text.includes(toText)) end = { node, offset: text.indexOf(toText) + toText.length };
      }
      if (!start || !end) return null;
      root.scrollIntoView({ block: 'center' });
      const range = document.createRange();
      range.setStart(start.node, start.offset);
      range.setEnd(end.node, end.offset);
      const rects = [...range.getClientRects()].filter((rect) => rect.width > 0);
      const first = rects[0];
      const last = rects[rects.length - 1];
      return {
        from: { x: first.left + 1, y: first.top + first.height / 2 },
        to: { x: last.right - 1, y: last.top + last.height / 2 },
      };
    },
    { fromText, toText }
  );
  if (!bounds) throw new Error(`Text "${fromText}" was not found to highlight`);
  await page.mouse.move(bounds.from.x, bounds.from.y, { steps: 15 });
  await page.waitForTimeout(400);
  await page.mouse.down();
  await page.mouse.move(bounds.to.x, bounds.to.y, { steps: 40 });
  await page.mouse.up();
  await page.waitForTimeout(1_800);
};

let clickPatched = false;

/**
 * Makes the pointer arrive before the click instead of chasing it.
 *
 * `click()` dispatches mousemove and mousedown 1ms apart, so an animated cursor is always
 * mid-flight when the button is already pressed. Hovering first gives the glide somewhere to
 * happen, which also keeps the ripple on the view the click belongs to.
 *
 * Patches the shared Locator prototype so it covers `getByRole` and friends.
 */
const leadPointerBeforeClick = (page: Page, glideMs: number) => {
  if (clickPatched || glideMs <= 0) return;
  clickPatched = true;

  const prototype = Object.getPrototypeOf(page.locator('body')) as Locator;
  const originalClick = prototype.click;

  prototype.click = async function patchedClick(this: Locator, ...args: Parameters<Locator['click']>) {
    try {
      await this.hover({ timeout: 2_000 });
      await this.page().waitForTimeout(glideMs);
    } catch {
      // Not every target can be hovered. The flourish is optional, so fall through and
      // click exactly as Playwright would have.
    }
    return originalClick.apply(this, args);
  };
};

export const test = base.extend({
  page: async ({ page }, use) => {
    if (isVideoEnabled) {
      // must run before the first navigation so it survives every document
      await page.addInitScript(cursorInitScript, cursorTimings);
      leadPointerBeforeClick(page, cursorTimings.glideMs);
    }
    await use(page);
  },
});

export { expect };
