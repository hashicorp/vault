/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { test } from '../../fixtures/demo';
import { ensureMount, rootApi } from '../../fixtures/root-api';

const KV_PATH = 'actions-kv';

test.beforeAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await ensureMount(api, KV_PATH, 'kv-v2');
  await api.dispose();
});

test.afterAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await api.delete(`/v1/sys/mounts/${KV_PATH}`);
  await api.dispose();
});

// Capability gating must not hide actions from a user who is allowed to use them.
test('privileged user keeps create and delete actions on list views', async ({ page }) => {
  await test.step('secrets engines offer enable and delete', async () => {
    await page.goto('secrets-engines');
    await expect(page.getByRole('link', { name: 'Enable new engine' })).toBeVisible();
    await page
      .getByRole('row')
      .filter({ hasText: `${KV_PATH}/` })
      .getByRole('button', { name: 'Options' })
      .click();
    await expect(page.getByRole('button', { name: 'Delete engine path' })).toBeVisible();
    await page.keyboard.press('Escape');
  });
});

// The cubbyhole/ engine is a built-in per-token mount that Vault refuses to disable.
test('cubbyhole engine never offers delete, even to a privileged user', async ({ page }) => {
  await page.goto('secrets-engines');
  await page
    .getByRole('row')
    .filter({ hasText: 'cubbyhole/' })
    .getByRole('button', { name: 'Options' })
    .click();
  await expect(page.getByRole('link', { name: 'View configuration' })).toBeVisible();
  await expect(page.getByRole('button', { name: 'Delete engine path' })).toHaveCount(0);
});
