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
