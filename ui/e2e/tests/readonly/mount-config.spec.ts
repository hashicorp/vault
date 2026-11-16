/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect, type Page } from '@playwright/test';
import { highlightElement, holdForVideo, test } from '../../fixtures/demo';
import { ensureMount, rootApi } from '../../fixtures/root-api';

const MOUNT = 'readonly-kv/mount-config';

interface TuneConfig {
  default_lease_ttl: number;
  max_lease_ttl: number;
}

const readToken = (page: Page) =>
  page.evaluate(() => {
    const key = Object.keys(localStorage).find((item) => item.startsWith('vault-token'));
    if (!key) throw new Error('No stored Vault token');
    const value = localStorage.getItem(key);
    if (!value) throw new Error('No stored Vault token value');
    const token = (JSON.parse(value) as { token?: string }).token;
    if (!token) throw new Error('Stored Vault token is empty');
    return token;
  });

const readTuneConfig = async (page: Page): Promise<TuneConfig> => {
  const response = await page.request.get(`/v1/sys/mounts/${MOUNT}/tune`, {
    headers: { 'X-Vault-Token': await readToken(page) },
  });
  await expect(response).toBeOK();
  const { data }: { data: TuneConfig } = await response.json();
  return data;
};

const removeMountIfPresent = async (api: Awaited<ReturnType<typeof rootApi>>) => {
  const response = await api.get('/v1/sys/mounts');
  await expect(response).toBeOK();
  const { data } = await response.json();
  if (Object.prototype.hasOwnProperty.call(data, `${MOUNT}/`)) {
    await expect(await api.delete(`/v1/sys/mounts/${MOUNT}`)).toBeOK();
  }
};

test.beforeAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  try {
    await ensureMount(api, MOUNT, 'kv-v2');
    await expect(
      await api.post(`/v1/sys/mounts/${MOUNT}/tune`, {
        data: { default_lease_ttl: '2h', max_lease_ttl: '4h' },
      })
    ).toBeOK();
  } finally {
    await api.dispose();
  }
});

test.afterAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  try {
    await removeMountIfPresent(api);
  } finally {
    await api.dispose();
  }
});

// A read-only token may inspect mount settings but must not persist configuration changes.
test('read-only user cannot tune mount configuration', async ({ page }) => {
  let configurationUrl = '';
  await test.step('open the seeded mount configuration', async () => {
    await page.goto('secrets-engines');
    const mountRow = page.getByRole('row').filter({ hasText: `${MOUNT}/` });
    await mountRow.getByRole('button', { name: 'Options' }).click();
    await page.getByRole('link', { name: 'View configuration' }).click();

    await expect(page.getByRole('form', { name: 'general settings form' })).toBeVisible();
    await expect(page.getByLabel('Default time-to-live (TTL)time')).toHaveValue('2');
    await expect(page.getByLabel('Default time-to-live (TTL)unit')).toHaveValue('h');
    await expect(page.getByLabel('Maximum time-to-live (TTL)time')).toHaveValue('4');
    await expect(page.getByLabel('Maximum time-to-live (TTL)unit')).toHaveValue('h');
    configurationUrl = page.url();
  });

  await test.step('denied tune leaves backend configuration unchanged', async () => {
    await page.goto(configurationUrl);
    await page.getByLabel('Default time-to-live (TTL)time').fill('3');
    const tuneResponsePromise = page.waitForResponse(
      (response) =>
        decodeURIComponent(new URL(response.url()).pathname) === `/v1/sys/mounts/${MOUNT}/tune` &&
        response.request().method() === 'POST'
    );
    await page.getByRole('button', { name: 'Save changes' }).click();
    const response = await tuneResponsePromise;
    expect(response.status()).toBe(403);
    await expect(page.getByText(/permission denied/i)).toBeVisible();
    await highlightElement(page, page.getByText(/permission denied/i));
    await holdForVideo(page);
    await test.info().attach('denied-mount-tune', {
      body: JSON.stringify({
        method: response.request().method(),
        status: response.status(),
        payload: response.request().postDataJSON(),
      }),
      contentType: 'application/json',
    });

    const persisted = await readTuneConfig(page);
    expect(persisted).toMatchObject({ default_lease_ttl: 7200, max_lease_ttl: 14400 });

    await page.goto('secrets-engines');
    const mountRow = page.getByRole('row').filter({ hasText: `${MOUNT}/` });
    await mountRow.getByRole('button', { name: 'Options' }).click();
    await page.getByRole('link', { name: 'View configuration' }).click();
    await expect(page.getByLabel('Default time-to-live (TTL)time')).toHaveValue('2');
    await expect(page.getByLabel('Maximum time-to-live (TTL)time')).toHaveValue('4');
  });
});
