/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { highlightText, holdForVideo, test } from '../../fixtures/demo';
import { ensureMount, rootApi } from '../../fixtures/root-api';

import type { Locator, Page } from '@playwright/test';

const SSH = 'ttl-max-ssh';
const TRANSIT = 'ttl-max-transit';
const K8S = 'ttl-max-k8s';
const TOO_LARGE = (label: string) => `${label} is too large. Maximum allowed value is 106,751 days.`;

// Setup and cleanup live in hooks so the mounts are removed even when a test fails or retries.
test.beforeAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await ensureMount(api, SSH, 'ssh');
  await ensureMount(api, TRANSIT, 'transit');
  await ensureMount(api, K8S, 'kubernetes');
  // No Kubernetes cluster is contacted; the config and role only let the credentials form load.
  await expect(
    await api.post(`/v1/${K8S}/config`, {
      data: { kubernetes_host: 'https://127.0.0.1:1', disable_local_ca_jwt: true },
    })
  ).toBeOK();
  await expect(
    await api.post(`/v1/${K8S}/roles/ttl-max-role`, {
      data: { allowed_kubernetes_namespaces: '*', service_account_name: 'sa' },
    })
  ).toBeOK();
  await expect(await api.post(`/v1/${TRANSIT}/keys/rotating`, { data: {} })).toBeOK();
  await expect(await api.post(`/v1/${TRANSIT}/keys/exportable`, { data: { exportable: true } })).toBeOK();
  await api.dispose();
});

test.afterAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  for (const mount of [SSH, TRANSIT, K8S]) await api.delete(`/v1/sys/mounts/${mount}`);
  await api.dispose();
});

// Records every request the page sends to a mount, so a test can prove nothing was saved.
const recordWrites = (page: Page, mount: string) => {
  const writes: string[] = [];
  page.on('request', (request) => {
    if (request.method() !== 'GET' && request.url().includes(`/v1/${mount}/`)) {
      writes.push(`${request.method()} ${new URL(request.url()).pathname}`);
    }
  });
  return writes;
};

// Clicking first moves the recorded cursor to the input; fill alone types without it.
const fillTtl = async (picker: Locator, value: string, unit?: string) => {
  const units = picker.locator('[data-test-ttl-value]');
  await units.click();
  await units.fill(value);
  if (unit) await picker.getByRole('combobox').selectOption(unit);
};

// An over-max TTL should be explained where it was typed and must never reach the API.
test('an SSH role blocks a TTL longer than Vault supports', async ({ page }) => {
  const writes = recordWrites(page, SSH);
  await page.goto(`secrets-engines/${SSH}/create`);
  await page.getByLabel('Role name').fill('ttl-max-role');
  await page.getByRole('button', { name: 'Options' }).click();
  const ttl = page.locator('[data-test-field="ttl"]');

  await test.step('an over-max TTL is flagged and the role is not saved', async () => {
    await ttl.locator('[data-test-toggle-label]').click();
    await fillTtl(ttl, '9999999999', 'd');
    await page.getByRole('button', { name: 'Create role' }).click();

    await expect(ttl.getByText(TOO_LARGE('TTL'))).toBeVisible();
    await highlightText(page, ttl, TOO_LARGE('TTL'));
    expect(writes).toEqual([]);
  });

  await test.step('a TTL left over from the CA type does not block an OTP role', async () => {
    await page.locator('[data-test-input="key_type"]').selectOption('otp');
    await page.locator('[data-test-input="default_user"]').fill('admin');
    await page.getByRole('button', { name: 'Create role' }).click();

    await expect(page).toHaveURL(new RegExp(`${SSH}/show/ttl-max-role`));
    expect(writes).toEqual([`POST /v1/${SSH}/roles/ttl-max-role`]);
    await holdForVideo(page);
  });
});

test('a transit key blocks an auto-rotation period longer than Vault supports', async ({ page }) => {
  const writes = recordWrites(page, TRANSIT);
  const period = () => page.getByRole('group', { name: 'Auto-rotation period' });

  await test.step('creating a key', async () => {
    await page.goto(`secrets-engines/${TRANSIT}/create`);
    await page.getByLabel('Name').fill('new-key');
    await page.getByText('Key will never be automatically rotated').click();
    await fillTtl(period(), '9999999999', 'd');
    await page.getByRole('button', { name: 'Create key' }).click();

    await expect(page.getByText(TOO_LARGE('Auto-rotation period'))).toBeVisible();
    await highlightText(page, page.locator('form'), TOO_LARGE('Auto-rotation period'));
    expect(writes).toEqual([]);
  });

  await test.step('editing a key', async () => {
    await page.goto(`secrets-engines/${TRANSIT}/edit/rotating`);
    await page.getByText('Key will never be automatically rotated').click();
    await fillTtl(period(), '9999999999', 'd');
    await page.getByRole('button', { name: 'Update transit key' }).click();

    await expect(page.getByText(TOO_LARGE('Auto-rotation period'))).toBeVisible();
    await highlightText(page, page.locator('form'), TOO_LARGE('Auto-rotation period'));
    expect(writes).toEqual([]);

    await fillTtl(period(), '30', 'd');
    await page.getByRole('button', { name: 'Update transit key' }).click();
    await expect.poll(() => writes).toEqual([`POST /v1/${TRANSIT}/keys/rotating/config`]);
    await holdForVideo(page);
  });
});

test('a transit export blocks a wrap TTL longer than Vault supports', async ({ page }) => {
  const exports: string[] = [];
  page.on('request', (request) => {
    if (request.url().includes(`/v1/${TRANSIT}/export/`)) exports.push(request.url());
  });
  await page.goto(`secrets-engines/${TRANSIT}/show/exportable`);
  await page.getByRole('link', { name: 'Key Actions' }).click();
  await page.getByText('Export Key', { exact: true }).click();

  await fillTtl(page.getByRole('group', { name: 'Wrap response' }), '9999999999', 'd');
  await page.getByRole('button', { name: 'Export key' }).click();

  await expect(page.getByText(TOO_LARGE('Wrap response'))).toBeVisible();
  await highlightText(page, page.locator('form'), TOO_LARGE('Wrap response'));
  expect(exports).toEqual([]);

  await fillTtl(page.getByRole('group', { name: 'Wrap response' }), '30', 'm');
  await page.getByRole('button', { name: 'Export key' }).click();
  await expect.poll(() => exports.length).toBe(1);
  await holdForVideo(page);
});

test('generating Kubernetes credentials blocks a TTL longer than Vault supports', async ({ page }) => {
  const writes = recordWrites(page, K8S);
  await page.goto(`secrets-engines/${K8S}/kubernetes/roles/ttl-max-role/credentials`);
  await page.getByLabel('Kubernetes namespace').fill('default');
  const ttl = page.getByRole('group', { name: 'Time-to-Live (TTL)' });

  await test.step('an over-max TTL is flagged and nothing is requested', async () => {
    await ttl.getByText('Time-to-Live (TTL)').click();
    await fillTtl(ttl, '9999999999', 'd');
    await page.getByRole('button', { name: 'Generate credentials' }).click();

    await expect(page.getByText(TOO_LARGE('Time-to-Live (TTL)'))).toBeVisible();
    await highlightText(page, page.locator('form'), TOO_LARGE('Time-to-Live (TTL)'));
    expect(writes).toEqual([]);
  });

  await test.step('turning the TTL off sends no TTL at all', async () => {
    await ttl.getByText('Time-to-Live (TTL)').click();
    await page.getByRole('button', { name: 'Generate credentials' }).click();
    await expect.poll(() => writes).toEqual([`POST /v1/${K8S}/creds/ttl-max-role`]);
    await holdForVideo(page);
  });
});
