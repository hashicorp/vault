/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { test } from '../../fixtures/demo';
import { ensureNamespace, rootApi } from '../../fixtures/root-api';

const NAMESPACE = 'readonly-ns';

test.beforeAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await ensureNamespace(api, NAMESPACE);
  await api.dispose();
});

// A read-only token must not be offered create/delete actions it cannot perform,
// while still being able to switch into a namespace it can read.
test('read-only user can switch to a namespace but cannot create or delete one', async ({ page }) => {
  await page.goto('access/namespaces');
  await expect(page.getByRole('heading', { name: 'Namespaces', level: 1 })).toBeVisible();
  await expect(page.getByRole('gridcell', { name: NAMESPACE })).toBeVisible();

  await test.step('header create action is hidden', async () => {
    await expect(page.getByRole('link', { name: 'Create namespace' })).toHaveCount(0);
  });

  await test.step('row menu offers switch but not delete', async () => {
    await page.getByRole('button', { name: 'More options' }).click();
    await expect(page.getByRole('button', { name: 'Switch to namespace' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Delete' })).toHaveCount(0);
  });

  await test.step('switching lands on the namespace dashboard', async () => {
    await page.getByRole('button', { name: 'Switch to namespace' }).click();
    await expect(page).toHaveURL(new RegExp(`/dashboard\\?namespace=${NAMESPACE}`));
  });
});
