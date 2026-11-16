/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { test } from '../../fixtures/demo';
import { ensureMount, rootApi } from '../../fixtures/root-api';
import { isVideoEnabled } from '../../video-config';

const MOUNT = 'kv-routes';
const SECRET = 'app/first-secret';
const SECRET_URL = `secrets-engines/${MOUNT}/kv/${encodeURIComponent(SECRET)}`;

const hold = async (page: import('@playwright/test').Page, ms = 1500) => {
  if (isVideoEnabled) await page.waitForTimeout(ms);
};

test.beforeAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await ensureMount(api, MOUNT, 'kv-v2');
  // Only seed once so repeated runs see the same versions.
  const existing = await api.get(`/v1/${MOUNT}/metadata/${SECRET}`);
  if (!existing.ok()) {
    for (const version of [1, 2, 3]) {
      await expect(
        await api.post(`/v1/${MOUNT}/data/${SECRET}`, { data: { data: { [`key-${version}`]: 'value' } } })
      ).toBeOK();
    }
  }
  await api.dispose();
});

test.afterAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await api.delete(`/v1/sys/mounts/${MOUNT}`);
  await api.dispose();
});

// Deep links to every kv page must keep working when the engine is inlined into the main app.
test('kv urls resolve to the expected pages', async ({ page }) => {
  await test.step('engine root redirects to the list', async () => {
    await page.goto(`secrets-engines/${MOUNT}/kv`);
    await expect(page).toHaveURL(new RegExp(`/ui/vault/secrets-engines/${MOUNT}/kv/list$`));
    await expect(page.getByRole('heading', { name: MOUNT })).toBeVisible();
    await hold(page);
  });

  await test.step('directory list shows nested secrets', async () => {
    await page.goto(`secrets-engines/${MOUNT}/kv/list/app`);
    await expect(page).toHaveURL(new RegExp(`/kv/list/app$`));
    await expect(page.getByRole('link', { name: 'first-secret' })).toBeVisible();
    await hold(page);
  });

  await test.step('create secret form', async () => {
    await page.goto(`secrets-engines/${MOUNT}/kv/create`);
    await expect(page.getByRole('textbox', { name: 'Path for this secret' })).toBeVisible();
    await hold(page, 1000);
  });

  await test.step('secret overview, details, paths', async () => {
    await page.goto(SECRET_URL);
    await expect(page.getByText('Current version', { exact: true })).toBeVisible();
    await hold(page);

    await page.goto(`${SECRET_URL}/details?version=2`);
    await expect(page).toHaveURL(/\/details\?version=2$/);
    await expect(page.getByText('key-2')).toBeVisible();
    await hold(page);

    await page.goto(`${SECRET_URL}/paths`);
    await expect(page).toHaveURL(/\/paths$/);
    await expect(page.getByRole('heading', { name: /paths/i }).first()).toBeVisible();
    await hold(page, 1000);
  });

  await test.step('metadata, version history and diff', async () => {
    await page.goto(`${SECRET_URL}/metadata`);
    await expect(page).toHaveURL(/\/metadata$/);
    await hold(page, 1000);

    await page.goto(`${SECRET_URL}/metadata/versions`);
    await expect(page.getByText('Version 3')).toBeVisible();
    await expect(page.getByText('Version 1')).toBeVisible();
    await hold(page);

    await page.goto(`${SECRET_URL}/metadata/diff`);
    await expect(page).toHaveURL(/\/metadata\/diff$/);
    await hold(page, 1000);
  });
});

// Secrets are created through the form and land on the overview with the encoded path preserved.
test('can create a secret and a new version, then edit metadata', async ({ page }) => {
  const path = `app/created-${Date.now()}`;

  await test.step('create the secret', async () => {
    await page.goto(`secrets-engines/${MOUNT}/kv/create`);
    await page.getByRole('textbox', { name: 'Path for this secret' }).fill(path);
    await page.getByRole('textbox', { name: 'key' }).fill('foo');
    await page.getByRole('textbox', { name: 'foo' }).fill('bar');
    await page.getByRole('button', { name: 'Save' }).click();
    await expect(page).toHaveURL(new RegExp(`/kv/${encodeURIComponent(path)}`));
    await expect(page.getByText('Current version', { exact: true })).toBeVisible();
    await hold(page);
  });

  await test.step('create a new version', async () => {
    await page.goto(`secrets-engines/${MOUNT}/kv/${encodeURIComponent(path)}/details/edit`);
    await page.getByRole('textbox', { name: 'key' }).first().fill('foo-v2');
    await page.getByRole('textbox', { name: 'foo-v2' }).fill('bar-v2');
    await page.getByRole('button', { name: 'Save' }).click();
    await expect(page).toHaveURL(new RegExp(`/kv/${encodeURIComponent(path)}$`));
    await expect(page.getByText('The current version of this secret.')).toBeVisible();
    await hold(page);
  });

  await test.step('edit custom metadata', async () => {
    await page.goto(`secrets-engines/${MOUNT}/kv/${encodeURIComponent(path)}/metadata/edit`);
    await page.getByRole('textbox', { name: 'key' }).fill('meta');
    await page.getByRole('textbox', { name: 'value' }).fill('data');
    await page.getByRole('button', { name: 'Update' }).click();
    await expect(page).toHaveURL(new RegExp(`${encodeURIComponent(path)}/metadata$`));
    await expect(page.getByText('meta', { exact: true })).toBeVisible();
    await hold(page);
  });
});

// The configuration pages are the only entry to engine-level settings.
test('can edit the kv configuration and exit back to the list', async ({ page }) => {
  await page.goto(`secrets-engines/${MOUNT}/kv/configuration`);
  await expect(page.getByRole('heading', { name: `${MOUNT} configuration` })).toBeVisible();

  await page.getByRole('link', { name: 'Edit configuration' }).click();
  await expect(page).toHaveURL(/\/kv\/configure$/);
  await page.locator('#max_versions').fill('7');
  await page.getByRole('button', { name: 'Save' }).click();
  await expect(page).toHaveURL(/\/kv\/configuration$/);
  await expect(page.getByText('Maximum number of versions')).toBeVisible();
  await hold(page);

  await page.getByRole('link', { name: 'Exit configuration' }).click();
  await expect(page).toHaveURL(/\/kv\/list$/);
});
