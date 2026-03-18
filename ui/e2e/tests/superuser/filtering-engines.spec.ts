/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';
import { BasePage } from '../../pages/base';

test('filtering secrets engines workflow', async ({ page }) => {
  const basePage = new BasePage(page);
  // nav to secrets engines page and enable a couple of engines to test filtering
  await page.goto('dashboard');
  await page.getByRole('link', { name: 'Secrets', exact: true }).click();

  // skip intro page if it appears
  if (await page.getByRole('button', { name: 'Skip' }).isVisible()) {
    await page.getByRole('button', { name: 'Skip' }).click();
  }

  // enable transit
  await page.getByRole('link', { name: 'Enable new engine' }).click();
  await page.getByLabel('Transit - enabled engine type').click();
  await page.getByRole('button', { name: 'Enable engine' }).click();
  await basePage.dismissFlashMessages();
  await page.getByLabel('breadcrumbs').getByRole('link', { name: 'Secrets engines' }).click();

  // enable kv
  await page.getByRole('link', { name: 'Enable new engine' }).click();
  await page.getByRole('heading', { name: 'KV' }).click();
  await page.getByRole('button', { name: 'Enable engine' }).click();
  await basePage.dismissFlashMessages();
  await page.getByLabel('breadcrumbs').getByRole('link', { name: 'Secrets engines' }).click();

  // search for transit engine
  await page.getByRole('searchbox', { name: 'Search' }).fill('transit');

  // assert theres only 1 result and its the transit engine
  await expect(page.getByText('1–1 of 1 page 1 Items per')).toBeVisible();
  await expect(page.getByRole('link', { name: 'transit/' })).toBeVisible();

  // assert user can click into transit engine from search results and nav back
  await page.getByRole('link', { name: 'transit/' }).click();
  await expect(page.getByRole('heading', { name: 'transit', exact: true })).toBeVisible();
  await page.getByLabel('breadcrumbs').getByRole('link', { name: 'Secrets engines' }).click();

  // filter for cubbyhole engine type
  await page.getByRole('button', { name: 'Engine type' }).click();
  await page.getByRole('checkbox', { name: 'cubbyhole', exact: true }).check();
  await page.getByRole('button', { name: 'Engine type' }).click();
  await expect(page.getByRole('link', { name: 'cubbyhole/' })).toBeVisible();

  // clear filter
  await page.getByRole('button', { name: 'Clear all' }).click();

  // filter by engine type and version
  await page.getByRole('button', { name: 'Engine type' }).click();
  await page.getByRole('checkbox', { name: 'cubbyhole', exact: true }).check();
  await page.getByRole('button', { name: 'Version', exact: true }).click();
  // clicking the label here so it closes the dropdown, the checkbox is inside the label so it will still be checked
  // Note: /\+builtin\.vault/ is regex to partial match the builtin label since the version will change
  await page
    .locator('label')
    .filter({ hasText: /\+builtin\.vault/ })
    .click();

  // assert filter chips are visible and applied correctly
  await expect(
    page
      .locator('span')
      .filter({ hasText: /\+builtin\.vault/ })
      .nth(1)
  ).toBeVisible();
  await expect(page.locator('span').filter({ hasText: 'cubbyhole' }).nth(1)).toBeVisible();
  await expect(page.getByRole('link', { name: 'cubbyhole/' })).toBeVisible();
  await expect(page.getByRole('gridcell', { name: /\+builtin\.vault/ })).toBeVisible();

  // clear filters and assert all engines are visible again
  await page.getByRole('button', { name: 'Clear all' }).click();
  await expect(page.getByRole('link', { name: 'kv/' })).toBeVisible();
  await expect(page.getByRole('link', { name: 'transit/' })).toBeVisible();

  await basePage.disableEngine('transit');
  await basePage.disableEngine('kv');
  await basePage.dismissFlashMessages();
});
