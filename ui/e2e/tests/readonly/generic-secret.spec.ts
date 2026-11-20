/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { test } from '../../fixtures/demo';
import { ensureMount, rootApi } from '../../fixtures/root-api';

const MOUNT = 'generic-secret-readonly';
const SECRET = 'readable';

test.beforeAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  try {
    await ensureMount(api, MOUNT, 'generic');
    await expect(await api.post(`/v1/${MOUNT}/${SECRET}`, { data: { seeded: 'original-value' } })).toBeOK();
  } finally {
    await api.dispose();
  }
});

test.afterAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  try {
    const response = await api.delete(`/v1/sys/mounts/${MOUNT}`);
    expect([204, 404]).toContain(response.status());
  } finally {
    await api.dispose();
  }
});

// Read-only users must retain list/read access without receiving secret mutation capabilities.
test('read-only user can read generic secrets but cannot create, update, or delete them', async ({
  page,
  playwright,
}, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  try {
    await test.step('list and view a readable secret without edit or delete actions', async () => {
      await page.goto(`secrets-engines/${MOUNT}/list/`);
      await expect(page.getByRole('link', { name: SECRET, exact: true })).toBeVisible();
      await page.getByRole('button', { name: `Manage database ${SECRET}` }).click();
      await expect(page.getByRole('link', { name: 'Details', exact: true })).toBeVisible();
      await expect(page.getByRole('link', { name: 'Edit', exact: true })).toHaveCount(0);
      await expect(page.getByRole('button', { name: 'Delete', exact: true })).toHaveCount(0);

      await page.goto(`secrets-engines/${MOUNT}/show/${SECRET}`);
      await expect(page.getByText('seeded', { exact: true })).toBeVisible();
      const secretValue = page.locator('[data-test-masked-input]').last();
      await secretValue.getByRole('button', { name: 'show value' }).click();
      await expect(secretValue.locator('pre')).toHaveText('original-value');
      await expect(page.getByRole('link', { name: 'Edit secret', exact: true })).toHaveCount(0);
      await expect(page.getByRole('button', { name: 'Delete', exact: true })).toHaveCount(0);
    });

    await test.step('deny a create through the UI and preserve existing data', async () => {
      await page.goto(`secrets-engines/${MOUNT}/create/`);
      await page.getByRole('textbox', { name: 'Path for this secret' }).fill('denied-create');
      await page.getByRole('textbox', { name: 'Key' }).fill('key');
      await page.getByRole('textbox', { name: 'key', exact: true }).fill('value');
      const denied = page.waitForResponse((response) => {
        const path = decodeURIComponent(new URL(response.url()).pathname);
        return response.request().method() === 'POST' && path.endsWith(`/v1/${MOUNT}/denied-create`);
      });
      await page.getByRole('button', { name: 'Save', exact: true }).click();
      expect((await denied).status()).toBe(403);
      await expect(page.locator('[data-test-message-error]')).toContainText(/permission denied/i);
      expect((await api.get(`/v1/${MOUNT}/denied-create`)).status()).toBe(404);
    });

    await test.step('deny direct update and delete requests without changing the secret', async () => {
      const token = await page.evaluate(() => {
        const storageKey = Object.keys(localStorage).find((key) => key.startsWith('vault-token'));
        if (!storageKey) throw new Error('No stored Vault token was found');
        const stored = JSON.parse(localStorage.getItem(storageKey) ?? '{}') as { token?: string };
        return stored.token ?? '';
      });
      expect(token).not.toBe('');
      const headers = { 'X-Vault-Token': token };

      const deniedWrite = await page.request.post(`/v1/${MOUNT}/${SECRET}`, {
        headers,
        data: { seeded: 'unauthorized-value' },
      });
      expect(deniedWrite.status()).toBe(403);
      const deniedDelete = await page.request.delete(`/v1/${MOUNT}/${SECRET}`, { headers });
      expect(deniedDelete.status()).toBe(403);

      const unchanged = await api.get(`/v1/${MOUNT}/${SECRET}`);
      await expect(unchanged).toBeOK();
      expect((await unchanged.json()).data).toEqual({ seeded: 'original-value' });
      await page.goto(`secrets-engines/${MOUNT}/show/${SECRET}`);
      await page.reload();
      await expect(page.getByText('seeded', { exact: true })).toBeVisible();
      const unchangedValue = page.locator('[data-test-masked-input]').last();
      await unchangedValue.getByRole('button', { name: 'show value' }).click();
      await expect(unchangedValue.locator('pre')).toHaveText('original-value');
    });
  } finally {
    await api.dispose();
  }
});
