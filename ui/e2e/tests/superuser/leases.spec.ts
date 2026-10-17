/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';

test('leases list navigation and prefix revoke buttons', async ({ page }) => {
  await page.goto('dashboard');
  await page.getByRole('link', { name: 'Access control' }).click();
  await page.getByRole('link', { name: 'Leases' }).click();

  // The list renders the Leases heading and at least the auth/ prefix link from setup tokens
  await expect(page.getByRole('heading', { name: 'Leases' })).toContainText('Leases');
  const authLink = page.getByRole('link', { name: 'auth/' });
  await expect(authLink).toBeVisible();

  // Navigate into auth/token/create/ — the setup project creates tokens, so this prefix exists
  await authLink.click();
  await page.getByRole('link', { name: 'token/' }).click();
  await page.getByRole('link', { name: 'create/' }).click();

  // At the create/ prefix level: both revoke-prefix action buttons are visible (capability-guarded)
  await expect(page.getByRole('button', { name: 'Force revoke prefix' })).toBeVisible();
  await expect(page.getByRole('button', { name: 'Revoke prefix', exact: true })).toBeVisible();

  // The list page itself rendered correctly at this prefix (heading present, no error state)
  await expect(page.getByRole('heading', { name: 'Leases' })).toBeVisible();

  // Back to Leases button navigates back to list root and is no longer visible at root
  await page.getByRole('link', { name: 'Back to Leases' }).click();
  await expect(page.getByRole('link', { name: 'Back to Leases' })).not.toBeVisible();
});
