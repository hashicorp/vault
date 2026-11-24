/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { test } from '../../fixtures/demo';

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

// Invalid lease IDs must show Vault's specific error instead of a generic API failure.
test('invalid lease ID shows the specific error message', async ({ page }) => {
  await test.step('open a nonexistent lease', async () => {
    await page.goto('access/leases/show/not-found');
  });

  await test.step('show the invalid lease ID message', async () => {
    await expect(page.getByText('not-found is not a valid lease ID.', { exact: true })).toBeVisible();
  });
});
