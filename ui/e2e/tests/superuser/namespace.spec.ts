/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';

test('namespace workflow', async ({ page }) => {
  await page.goto('dashboard');
  // nav to namespaces and create a new namespace
  await page.getByRole('link', { name: 'Access control' }).click();
  await page.getByRole('link', { name: 'Namespaces' }).click();
  // skip guided tour if it appears
  await page.getByRole('button', { name: 'Skip' }).click();

  await page.getByRole('link', { name: 'Create namespace' }).click();
  await page.getByRole('textbox', { name: 'Path' }).fill('testNamespace');
  await page.getByRole('button', { name: 'Save' }).click();

  // click on the namespace picker in the top navbar and switch to the new namespace
  await page.getByRole('button', { name: 'root' }).click();
  await page.getByRole('option', { name: 'testNamespace' }).click();

  // verify that we are switched into the new namespace by checking for the namespace name in the header
  await expect(page.locator('#app-main-content').getByText('testNamespace')).toBeVisible();
});
