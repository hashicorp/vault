/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';

test('leases workflow', async ({ page }) => {
  await page.goto('dashboard');
  await page.getByRole('link', { name: 'Access control' }).click();
  await page.getByRole('link', { name: 'Leases' }).click();
  await expect(page.getByRole('heading', { name: 'Leases' })).toContainText('Leases');
  await page.getByRole('link', { name: 'auth/' }).click();
  await page.getByRole('link', { name: 'token/' }).click();
  await page.getByRole('link', { name: 'create/' }).click();
  await expect(page.getByRole('button', { name: 'Force revoke prefix' })).toBeVisible();
  await expect(page.getByRole('button', { name: 'Revoke prefix', exact: true })).toBeVisible();
  await page.getByRole('link', { name: 'Back to Leases' }).click();
  await expect(page.getByRole('link', { name: 'Back to Leases' })).not.toBeVisible();
});
