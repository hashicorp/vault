/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { test } from '../../fixtures/demo';
import { ensureMount, ensureNamespace, rootApi } from '../../fixtures/root-api';

const NAMESPACE = 'actions-ns';
const KV_PATH = 'actions-kv';
const GROUP = 'actions-group';

test.beforeAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await ensureNamespace(api, NAMESPACE);
  await ensureMount(api, KV_PATH, 'kv-v2');
  // External, so the details view offers "Add alias" when the token may add one.
  await expect(await api.post(`/v1/identity/group/name/${GROUP}`, { data: { type: 'external' } })).toBeOK();
  await api.dispose();
});

test.afterAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await api.delete(`/v1/sys/namespaces/${NAMESPACE}`);
  await api.delete(`/v1/sys/mounts/${KV_PATH}`);
  await api.delete(`/v1/identity/group/name/${GROUP}`);
  await api.dispose();
});

// Capability gating must not hide actions from a user who is allowed to use them.
test('privileged user keeps create and delete actions on list views', async ({ page }) => {
  await test.step('namespaces offer create, switch, and delete', async () => {
    await page.goto('access/namespaces');
    await expect(page.getByRole('link', { name: 'Create namespace' })).toBeVisible();
    await page.getByPlaceholder('Filter by namespace path').fill(NAMESPACE);
    await page.getByRole('button', { name: 'More options' }).click();
    await expect(page.getByRole('button', { name: 'Switch to namespace' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Delete' })).toBeVisible();
    await page.keyboard.press('Escape');
  });

  await test.step('groups offer create, and group details offer edit and add alias', async () => {
    await page.goto('access/identity/groups');
    await expect(page.getByRole('link', { name: 'Create group' })).toBeVisible();
    await page.getByRole('link', { name: GROUP }).click();
    await expect(page.getByRole('link', { name: 'Edit group' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Add alias' })).toBeVisible();
  });

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
