/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';

test('license page renders license details', async ({ page }) => {
  await page.goto('dashboard');
  await page.getByRole('link', { name: 'Reporting' }).click();
  await page.getByRole('link', { name: 'License' }).click();

  await expect(page.getByRole('heading', { name: 'License' })).toBeVisible();

  await test.step('shows license detail rows', async () => {
    // data-test-row-label is set by InfoTableRow on the label span — reliable across all builds
    await expect(page.locator('[data-test-row-label="License ID"]')).toBeVisible();
    await expect(page.locator('[data-test-row-label="Valid from"]')).toBeVisible();
    await expect(page.locator('[data-test-row-label="License state"]')).toBeVisible();
  });

  await test.step('shows feature rows', async () => {
    // HSM is always the first feature in the allFeatures() list
    await expect(page.locator('[data-test-row-label="HSM"]')).toBeVisible();
  });

  await test.step('shows at least one active feature status', async () => {
    // data-test-row-value is set by InfoTableRow on the value span
    await expect(page.locator('[data-test-row-value="HSM"]').filter({ hasText: 'Active' })).toBeVisible();
  });
});
