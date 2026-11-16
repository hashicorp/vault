/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { test } from '../../fixtures/demo';
import { ensureMount, rootApi } from '../../fixtures/root-api';

const KV_PATH = 'readonly-kv';

test.beforeAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await ensureMount(api, KV_PATH, 'kv-v2');
  await api.dispose();
});

// Read-only tokens must not be offered enable/delete actions, and cubbyhole/ is never deletable.
test('read-only user cannot enable or delete secrets engines', async ({ page }) => {
  await page.goto('secrets-engines');
  await expect(page.getByRole('heading', { name: 'Secrets engines', level: 1 })).toBeVisible();

  await test.step('header enable action is hidden', async () => {
    await expect(page.getByRole('link', { name: 'Enable new engine' })).toHaveCount(0);
  });

  for (const path of ['cubbyhole/', `${KV_PATH}/`]) {
    await test.step(`${path} row menu only offers view configuration`, async () => {
      await page.getByRole('row').filter({ hasText: path }).getByRole('button', { name: 'Options' }).click();
      await expect(page.getByRole('link', { name: 'View configuration' })).toBeVisible();
      await expect(page.getByRole('button', { name: 'Delete engine path' })).toHaveCount(0);
      await page.keyboard.press('Escape');
    });
  }
});

// Read-only users still need the mount data returned by the secrets-engine list route.
test('read-only user can read a mounted secrets engine', async ({ page }) => {
  await page.goto('secrets-engines');

  const mountLink = page.getByRole('link', { name: `${KV_PATH}/`, exact: true });
  await expect(mountLink).toBeVisible();
  await page.reload();
  await expect(mountLink).toBeVisible();
});
