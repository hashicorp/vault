/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect, type Page } from '@playwright/test';
import { highlightElement, holdForVideo, test } from '../../fixtures/demo';
import { rootApi } from '../../fixtures/root-api';
import { BasePage } from '../../pages/base';
import { ConfigurationSettingsPage } from '../../pages/configuration-settings';

const MOUNT = 'mount-config-superuser';

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
    await removeMountIfPresent(api);
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

// Mount configuration must remain readable and durable through the Resource-backed UI.
test('superuser can mount and persist secrets engine configuration', async ({ page }) => {
  const basePage = new BasePage(page);
  const configurationSettingsPage = new ConfigurationSettingsPage(page);

  await test.step('mount KV with known TTL values', async () => {
    const createResponsePromise = page.waitForResponse(
      (response) =>
        response.url().endsWith(`/v1/sys/mounts/${MOUNT}`) && response.request().method() === 'POST'
    );
    await basePage.enableEngine('kv', MOUNT, { skipEnable: true });
    await page.getByRole('button', { name: 'View additional settings' }).click();
    await page.locator('label').filter({ hasText: 'Default Lease TTL Vault will' }).click();
    await page.getByLabel('TTL unit for Default Lease TTL').selectOption('h');
    await page
      .getByRole('group', { name: 'Default Lease TTL Lease will' })
      .getByLabel('Number of units')
      .fill('2');
    await page.locator('label').filter({ hasText: 'Max Lease TTL Vault will' }).click();
    await page.getByLabel('TTL unit for Max Lease TTL').selectOption('h');
    await page
      .getByRole('group', { name: 'Max Lease TTL Lease will' })
      .getByLabel('Number of units')
      .fill('4');
    await page.getByRole('button', { name: 'Enable engine' }).click();

    const response = await createResponsePromise;
    expect(response.ok()).toBe(true);
    expect(new URL(response.url()).pathname).toBe(`/v1/sys/mounts/${MOUNT}`);
    expect(response.request().postDataJSON()).toMatchObject({
      type: 'kv',
      config: { default_lease_ttl: '7200s', max_lease_ttl: '14400s' },
      options: { version: 2 },
    });

    const persisted = await readTuneConfig(page);
    expect(persisted).toMatchObject({ default_lease_ttl: 7200, max_lease_ttl: 14400 });
    await test.info().attach('mount-create-readback', {
      body: JSON.stringify({ status: response.status(), config: persisted }),
      contentType: 'application/json',
    });
  });

  await test.step('mount remains listed after reload', async () => {
    await page.goto('secrets-engines');
    const mountLink = page.getByRole('link', { name: `${MOUNT}/`, exact: true });
    await expect(mountLink).toBeVisible();
    await page.reload();
    await expect(page.getByRole('link', { name: `${MOUNT}/`, exact: true })).toBeVisible();
  });

  await test.step('read and update mount TTL configuration', async () => {
    await page.goto('secrets-engines');
    await page.getByRole('link', { name: `${MOUNT}/`, exact: true }).click();
    await configurationSettingsPage.navigateToConfiguration(MOUNT);
    await configurationSettingsPage.navigateToGeneralSettings('kv');

    await expect(page.getByLabel('Default time-to-live (TTL)time')).toHaveValue('2');
    await expect(page.getByLabel('Default time-to-live (TTL)unit')).toHaveValue('h');
    await expect(page.getByLabel('Maximum time-to-live (TTL)time')).toHaveValue('4');
    await expect(page.getByLabel('Maximum time-to-live (TTL)unit')).toHaveValue('h');

    await page.getByLabel('Default time-to-live (TTL)time').fill('5');
    await page.getByLabel('Maximum time-to-live (TTL)time').fill('6');
    const tuneResponsePromise = page.waitForResponse(
      (response) =>
        response.url().endsWith(`/v1/sys/mounts/${MOUNT}/tune`) && response.request().method() === 'POST'
    );
    await page.getByRole('button', { name: 'Save changes' }).click();
    const response = await tuneResponsePromise;
    expect(response.ok()).toBe(true);
    expect(response.request().postDataJSON()).toMatchObject({
      default_lease_ttl: '5h',
      max_lease_ttl: '6h',
    });
    await expect(page.getByRole('alert', { name: 'Configuration saved' })).toBeVisible();

    const persisted = await readTuneConfig(page);
    expect(persisted).toMatchObject({ default_lease_ttl: 18000, max_lease_ttl: 21600 });
    await test.info().attach('mount-tune-readback', {
      body: JSON.stringify({
        method: response.request().method(),
        status: response.status(),
        payload: response.request().postDataJSON(),
        persisted,
      }),
      contentType: 'application/json',
    });
  });

  await test.step('persisted TTL values survive reopening General Settings', async () => {
    await page.goto('secrets-engines');
    await page.getByRole('link', { name: `${MOUNT}/`, exact: true }).click();
    await configurationSettingsPage.navigateToConfiguration(MOUNT);
    await configurationSettingsPage.navigateToGeneralSettings('kv');
    await expect(page.getByLabel('Default time-to-live (TTL)time')).toHaveValue('5');
    await expect(page.getByLabel('Default time-to-live (TTL)unit')).toHaveValue('h');
    await expect(page.getByLabel('Maximum time-to-live (TTL)time')).toHaveValue('6');
    await expect(page.getByLabel('Maximum time-to-live (TTL)unit')).toHaveValue('h');
    await highlightElement(page, page.getByLabel('Default time-to-live (TTL)time'));
    await holdForVideo(page);
  });
});
