/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { highlightElement, highlightText, test } from '../../fixtures/demo';
import { ensureMount, rootApi } from '../../fixtures/root-api';

const MOUNT = 'generic-secret-migration';
const CREATE_PATH = 'created';
const LIST_DELETE_PATH = 'list-delete';
const SHOW_DELETE_PATH = 'show-delete';
const LIST_PATH = `/v1/${MOUNT}/`;
const PAGINATION_PATHS = Array.from(
  { length: 21 },
  (_, index) => `zz-pagination-${String(index).padStart(2, '0')}`
);

const cleanSecret = async (api: Awaited<ReturnType<typeof rootApi>>, path: string) => {
  const response = await api.delete(`/v1/${MOUNT}/${path}`);
  expect([204, 404]).toContain(response.status());
};

test.beforeAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  try {
    await ensureMount(api, MOUNT, 'generic');
  } finally {
    await api.dispose();
  }
});

test.beforeEach(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  try {
    for (const path of [
      CREATE_PATH,
      LIST_DELETE_PATH,
      SHOW_DELETE_PATH,
      'existing',
      'nested/item',
      ...PAGINATION_PATHS,
    ]) {
      await cleanSecret(api, path);
    }

    await expect(await api.post(`${LIST_PATH}existing`, { data: { seeded: 'existing-value' } })).toBeOK();
    await expect(await api.post(`${LIST_PATH}nested/item`, { data: { seeded: 'nested-value' } })).toBeOK();
    await expect(
      await api.post(`${LIST_PATH}${LIST_DELETE_PATH}`, { data: { seeded: 'list-delete-value' } })
    ).toBeOK();
    await expect(
      await api.post(`${LIST_PATH}${SHOW_DELETE_PATH}`, { data: { seeded: 'show-delete-value' } })
    ).toBeOK();
    await Promise.all(
      PAGINATION_PATHS.map(async (path) => {
        await expect(await api.post(`${LIST_PATH}${path}`, { data: { seeded: path } })).toBeOK();
      })
    );
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

// KV v1 users need list, edit, wrapping, and persistence behavior across visits.
test('superuser can list, create, update, and wrap a generic secret', async ({
  page,
  playwright,
}, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  try {
    await test.step('list secrets and nested folders from the backend', async () => {
      const listed = page.waitForResponse((response) => {
        const url = new URL(response.url());
        return (
          response.request().method() === 'GET' &&
          url.pathname.includes(`/v1/${MOUNT}/`) &&
          url.searchParams.get('list') === 'true'
        );
      });
      await page.goto(`secrets-engines/${MOUNT}/list/`);
      expect((await listed).ok()).toBe(true);
      await expect(page.getByRole('link', { name: 'existing', exact: true })).toBeVisible();
      await expect(page.getByRole('link', { name: 'nested/', exact: true })).toBeVisible();
      await highlightElement(page, page.getByRole('link', { name: 'existing', exact: true }));

      await page.goto(`secrets-engines/${MOUNT}/list/?pageFilter=existing`);
      await expect(page.getByRole('link', { name: 'existing', exact: true })).toBeVisible();
      await expect(page.getByRole('link', { name: 'nested/', exact: true })).toHaveCount(0);

      await page.goto(`secrets-engines/${MOUNT}/list/?page=2`);
      await expect(page.getByRole('link', { name: 'zz-pagination-20', exact: true })).toBeVisible();
      await expect(page.getByRole('button', { name: 'Previous page' })).toBeVisible();
    });

    await test.step('create a secret and verify its persisted request body', async () => {
      await page.goto(`secrets-engines/${MOUNT}/create/`);
      await page.getByRole('textbox', { name: 'Path for this secret' }).fill(CREATE_PATH);
      await page.getByRole('textbox', { name: 'Key' }).fill('created-key');
      await page.getByRole('textbox', { name: 'created-key' }).fill('created-value');

      const saved = page.waitForResponse((response) => {
        const path = decodeURIComponent(new URL(response.url()).pathname);
        return response.request().method() === 'POST' && path.endsWith(`/v1/${MOUNT}/${CREATE_PATH}`);
      });
      await page.getByRole('button', { name: 'Save', exact: true }).click();
      const response = await saved;
      expect(response.ok()).toBe(true);
      expect(response.request().postDataJSON()).toEqual({ 'created-key': 'created-value' });
      await expect(page.getByText('created-key', { exact: true })).toBeVisible();

      const readBack = await api.get(`${LIST_PATH}${CREATE_PATH}`);
      await expect(readBack).toBeOK();
      expect((await readBack.json()).data).toEqual({ 'created-key': 'created-value' });
    });

    await test.step('update the secret and verify it after reload', async () => {
      await page.getByRole('link', { name: 'Edit secret', exact: true }).click();
      await expect(page.getByRole('form', { name: 'secret edit form' })).toBeVisible();
      await page.getByRole('switch', { name: 'JSON' }).check();
      const editor = page.getByRole('textbox', { name: 'Secret Data' });
      await expect(editor).toBeVisible();
      await editor.fill('{"updated-key":"updated-value"}');

      const saved = page.waitForResponse((response) => {
        const path = decodeURIComponent(new URL(response.url()).pathname);
        return response.request().method() === 'POST' && path.endsWith(`/v1/${MOUNT}/${CREATE_PATH}`);
      });
      await page.getByRole('button', { name: 'Save', exact: true }).click();
      const response = await saved;
      expect(response.ok()).toBe(true);
      expect(response.request().postDataJSON()).toEqual({ 'updated-key': 'updated-value' });

      const readBack = await api.get(`${LIST_PATH}${CREATE_PATH}`);
      await expect(readBack).toBeOK();
      expect((await readBack.json()).data).toEqual({ 'updated-key': 'updated-value' });
      await page.goto(`secrets-engines/${MOUNT}/show/${CREATE_PATH}`);
      await page.reload();
      await expect(page).toHaveURL(new RegExp(`/secrets-engines/${MOUNT}/show/${CREATE_PATH}(?:\\?|$)`));
      await page.getByText('UI', { exact: true }).click();
      await expect(page.getByText('updated-key', { exact: true })).toBeVisible();
      await page.getByRole('button', { name: 'show value' }).click();
      await expect(page.locator('[data-test-masked-input] pre')).toHaveText('updated-value');
      await highlightText(page, page.locator('#app-main-content'), 'updated-key');
    });

    await test.step('wrap the secret and verify the wrapped payload can be unwrapped', async () => {
      await page.locator('[data-test-copy-menu-trigger]').click();
      const wrapped = page.waitForResponse((response) => {
        const request = response.request();
        return (
          request.method() === 'GET' &&
          new URL(response.url()).pathname.includes(`/v1/${MOUNT}/${CREATE_PATH}`) &&
          request.headers()['x-vault-wrap-ttl'] === '1800'
        );
      });
      await page.locator('[data-test-button="wrap"]').click();
      const response = await wrapped;
      expect(response.ok()).toBe(true);
      const token = (await response.json()).wrap_info.token as string;
      expect(token).not.toBe('');
      await expect(page.getByText('Secret successfully wrapped!', { exact: true })).toBeVisible();

      // The menu closes after wrapping, so reopen it to reveal the returned token.
      await page.locator('[data-test-copy-menu-trigger]').click();
      // Both the secret and wrapping token are masked; the token's reveal button is last.
      await page.getByRole('button', { name: 'show value', exact: true }).last().click();
      await expect(page.getByText(token, { exact: true })).toBeVisible();

      const unwrapped = await api.post('/v1/sys/wrapping/unwrap', { data: { token } });
      await expect(unwrapped).toBeOK();
      expect((await unwrapped.json()).data).toEqual({ 'updated-key': 'updated-value' });
    });
  } finally {
    await api.dispose();
  }
});

// Updating a secret should return users to its persisted Show view.
test('superuser returns to the show view after updating a generic secret', async ({
  page,
  playwright,
}, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  try {
    await test.step('save an update and verify the show view remains after reload', async () => {
      await page.goto(`secrets-engines/${MOUNT}/show/existing`);
      await page.getByRole('link', { name: 'Edit secret', exact: true }).click();
      await expect(page.getByRole('form', { name: 'secret edit form' })).toBeVisible();
      await page.getByRole('switch', { name: 'JSON' }).check();
      await page.getByRole('textbox', { name: 'Secret Data' }).fill('{"redirect-key":"redirect-value"}');

      const saved = page.waitForResponse((response) => {
        const path = decodeURIComponent(new URL(response.url()).pathname);
        return response.request().method() === 'POST' && path.endsWith(`/v1/${MOUNT}/existing`);
      });
      await page.getByRole('button', { name: 'Save', exact: true }).click();
      const response = await saved;
      expect(response.ok()).toBe(true);
      expect(response.request().postDataJSON()).toEqual({ 'redirect-key': 'redirect-value' });
      await expect(page).toHaveURL(new RegExp(`/secrets-engines/${MOUNT}/show/existing(?:\\?|$)`));

      const readBack = await api.get(`${LIST_PATH}existing`);
      await expect(readBack).toBeOK();
      expect((await readBack.json()).data).toEqual({ 'redirect-key': 'redirect-value' });

      await page.reload();
      await expect(page).toHaveURL(new RegExp(`/secrets-engines/${MOUNT}/show/existing(?:\\?|$)`));
      await page.getByText('UI', { exact: true }).click();
      await expect(page.getByText('redirect-key', { exact: true })).toBeVisible();
      const secretValue = page.locator('[data-test-masked-input]').last();
      await secretValue.getByRole('button', { name: 'show value' }).click();
      await expect(secretValue.locator('pre')).toHaveText('redirect-value');
    });
  } finally {
    await api.dispose();
  }
});

// Both delete entry points must remove persisted KV v1 data, not just update the current page.
test('superuser can delete generic secrets from the list and show views', async ({
  page,
  playwright,
}, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  try {
    await test.step('delete a secret from its list action', async () => {
      await page.goto(`secrets-engines/${MOUNT}/list/`);
      await page.getByRole('button', { name: `Manage database ${LIST_DELETE_PATH}` }).click();
      await page.getByRole('button', { name: 'Delete', exact: true }).click();
      const deleted = page.waitForResponse((response) => {
        const path = decodeURIComponent(new URL(response.url()).pathname);
        return response.request().method() === 'DELETE' && path.endsWith(`/v1/${MOUNT}/${LIST_DELETE_PATH}`);
      });
      await page.getByRole('dialog').getByRole('button', { name: 'Confirm', exact: true }).click();
      expect((await deleted).ok()).toBe(true);
      expect((await api.get(`${LIST_PATH}${LIST_DELETE_PATH}`)).status()).toBe(404);
    });

    await test.step('delete a secret from its show view', async () => {
      await page.goto(`secrets-engines/${MOUNT}/show/${SHOW_DELETE_PATH}`);
      await page.getByRole('button', { name: 'Delete', exact: true }).click();
      const deleted = page.waitForResponse((response) => {
        const path = decodeURIComponent(new URL(response.url()).pathname);
        return response.request().method() === 'DELETE' && path.endsWith(`/v1/${MOUNT}/${SHOW_DELETE_PATH}`);
      });
      await page.getByRole('dialog').getByRole('button', { name: 'Confirm', exact: true }).click();
      expect((await deleted).ok()).toBe(true);
      expect((await api.get(`${LIST_PATH}${SHOW_DELETE_PATH}`)).status()).toBe(404);
    });
  } finally {
    await api.dispose();
  }
});
