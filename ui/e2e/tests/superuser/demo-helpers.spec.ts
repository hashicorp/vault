/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { HIGHLIGHT_ID, highlightElement, test } from '../../fixtures/demo';
import { isVideoEnabled } from '../../video-config';

// The demo helpers must never change normal test runs; they only draw while a video is being recorded.
// `#${HIGHLIGHT_ID}` targets the helper's own overlay element, not app markup.
test.describe('demo recording helpers', () => {
  test('highlightElement draws nothing and returns at once when not recording', async ({ page }) => {
    test.skip(isVideoEnabled, 'only meaningful outside recordings');
    await page.setContent('<button type="button">Save</button>');
    const started = Date.now();
    await highlightElement(page, page.getByRole('button', { name: 'Save' }));
    expect(Date.now() - started).toBeLessThan(500);
    await expect(page.locator(`#${HIGHLIGHT_ID}`)).toHaveCount(0);
  });

  test('highlightElement outlines the component while recording, then removes the outline', async ({
    page,
  }) => {
    test.skip(!isVideoEnabled, 'only runs while recording (PW_VIDEO set)');
    await page.setContent('<button type="button">Save</button>');
    const done = highlightElement(page, page.getByRole('button', { name: 'Save' }), 1_000);
    await expect(page.locator(`#${HIGHLIGHT_ID}`)).toHaveCount(1);
    await done;
    await expect(page.locator(`#${HIGHLIGHT_ID}`)).toHaveCount(0);
  });
});
