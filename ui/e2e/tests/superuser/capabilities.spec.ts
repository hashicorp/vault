/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';

test('capabilities-backed routes render actions', async ({ page }) => {
  await page.goto('dashboard');

  await test.step('leases routes expose revoke actions', async () => {
    await page.getByRole('link', { name: 'Access control' }).click();
    await page.getByRole('link', { name: 'Leases' }).click();
    await expect(page.getByRole('heading', { name: 'Leases' })).toContainText('Leases');
    await page.getByRole('link', { name: 'auth/' }).click();
    await page.getByRole('link', { name: 'token/' }).click();
    await page.getByRole('link', { name: 'create/' }).click();
    await expect(page.getByRole('button', { name: 'Force revoke prefix' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Revoke prefix', exact: true })).toBeVisible();
  });

  await test.step('seal route exposes seal action', async () => {
    await page.goto('dashboard');
    await page.getByRole('link', { name: 'Resilience and recovery' }).click();
    await page.getByRole('link', { name: 'Seal Vault' }).click();
    await expect(page.getByRole('button', { name: 'Seal' })).toBeVisible();
  });

  await test.step('acl policy route exposes create action', async () => {
    await page.goto('dashboard');
    await page.getByRole('link', { name: 'Access control' }).click();
    await expect(page.getByRole('link', { name: 'Create ACL policy' })).toBeVisible();
  });
});
