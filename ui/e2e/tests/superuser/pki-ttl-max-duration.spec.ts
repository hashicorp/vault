/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { highlightText, holdForVideo, test } from '../../fixtures/demo';
import { ensureMount, rootApi } from '../../fixtures/root-api';

import type { Locator, Page } from '@playwright/test';

const MOUNT = 'pki-ttl-max';
const ROLE = 'ttl-max-role';
const TOO_LARGE = (label: string) => `${label} is too large. Maximum allowed value is 106,751 days.`;

// Setup and cleanup live in hooks so the mount is removed even when the test fails or retries.
test.beforeAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await ensureMount(api, MOUNT, 'pki');
  await expect(
    await api.post(`/v1/${MOUNT}/root/generate/internal`, {
      data: { common_name: 'ttl-max.example.com', ttl: '8760h', issuer_name: 'ttl-max-root' },
    })
  ).toBeOK();
  await expect(await api.post(`/v1/${MOUNT}/roles/${ROLE}`, { data: { allow_any_name: true } })).toBeOK();
  await api.dispose();
});

test.afterAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await api.delete(`/v1/sys/mounts/${MOUNT}`);
  await api.dispose();
});

// Records every request the page sends to the mount, so a test can prove nothing was saved.
const recordWrites = (page: Page) => {
  const writes: string[] = [];
  page.on('request', (request) => {
    if (request.method() !== 'GET' && request.url().includes(`/v1/${MOUNT}/`)) {
      writes.push(`${request.method()} ${new URL(request.url()).pathname}`);
    }
  });
  return writes;
};

// Clicking first moves the recorded cursor to the input; fill alone types without it.
const fillTtl = async (picker: Locator, value: string, unit?: string) => {
  const units = picker.getByLabel('Number of units');
  await units.click();
  await units.fill(value);
  if (unit) await picker.getByRole('combobox').selectOption(unit);
};

// An over-max TTL should be explained where it was typed and must never reach the API.
test('the pki configuration page blocks TTLs longer than Vault supports', async ({ page }) => {
  const writes = recordWrites(page);

  await test.step('an over-max CRL expiry is flagged and nothing is saved', async () => {
    await page.goto(`secrets-engines/${MOUNT}/pki/configuration/edit`);
    await fillTtl(page.getByRole('group', { name: 'Expiry' }), '9999999999', 'd');
    await page.getByRole('button', { name: 'Save' }).click();

    await expect(page.getByText(TOO_LARGE('Expiry'))).toBeVisible();
    await expect(page.getByText('There is an error with this form.')).toBeVisible();
    await highlightText(page, page.locator('[data-test-input="expiry"]'), TOO_LARGE('Expiry'));
    expect(writes).toEqual([]);
  });

  await test.step('an over-max ACME max TTL is flagged and nothing is saved', async () => {
    await page.goto(`secrets-engines/${MOUNT}/pki/configuration/edit`);
    await fillTtl(page.getByRole('group', { name: 'Max TTL' }), '9999999999', 'd');
    await page.getByRole('button', { name: 'Save' }).click();

    await expect(page.getByText(TOO_LARGE('Max TTL'))).toBeVisible();
    await highlightText(page, page.locator('[data-test-field="max_ttl"]'), TOO_LARGE('Max TTL'));
    expect(writes).toEqual([]);
  });

  await test.step('a valid value saves', async () => {
    // The form also saves the cluster config, which Vault rejects while its paths are empty.
    await page.locator('[data-test-input="path"]').fill(`https://vault.example.com/v1/${MOUNT}`);
    await page.locator('[data-test-input="aia_path"]').fill(`https://vault.example.com/v1/${MOUNT}`);
    await fillTtl(page.getByRole('group', { name: 'Max TTL' }), '90');
    await page.getByRole('button', { name: 'Save' }).click();
    await expect(page).toHaveURL(/pki\/configuration$/);
    await holdForVideo(page);
  });
});

test('the pki tidy page blocks durations longer than Vault supports', async ({ page }) => {
  const writes = recordWrites(page);
  const form = page.locator('form');
  await page.goto(`secrets-engines/${MOUNT}/pki/tidy/manual`);

  await test.step('an over-max pause duration is flagged', async () => {
    await fillTtl(page.getByRole('group', { name: 'Pause duration' }), '9999999999', 'd');
    await page.getByRole('button', { name: 'Perform tidy' }).click();

    await expect(page.getByText(TOO_LARGE('Pause duration'))).toBeVisible();
    await highlightText(page, form, TOO_LARGE('Pause duration'));
    expect(writes).toEqual([]);
  });

  await test.step('an over-max ACME safety buffer is flagged', async () => {
    await fillTtl(page.getByRole('group', { name: 'Pause duration' }), '0');
    await page.getByText('Tidy ACME disabled').click();
    await fillTtl(page.getByRole('group', { name: 'Tidy ACME enabled' }), '9999999999', 'd');
    await page.getByRole('button', { name: 'Perform tidy' }).click();

    await expect(page.getByText(TOO_LARGE('Tidy ACME enabled'))).toBeVisible();
    await highlightText(page, form, TOO_LARGE('Tidy ACME enabled'));
    expect(writes).toEqual([]);
  });

  await test.step('a valid value starts the tidy', async () => {
    await fillTtl(page.getByRole('group', { name: 'Tidy ACME enabled' }), '30', 'd');
    await page.getByRole('button', { name: 'Perform tidy' }).click();
    await expect.poll(() => writes).toEqual([`POST /v1/${MOUNT}/tidy`]);
    await holdForVideo(page);
  });
});

test('issuing a certificate blocks a TTL longer than Vault supports', async ({ page }) => {
  const writes = recordWrites(page);
  await page.goto(`secrets-engines/${MOUNT}/pki/roles/${ROLE}/generate`);
  const notValidAfter = page
    .getByRole('group')
    .filter({ hasText: 'Set relative certificate expiry with TTL' });

  await test.step('an over-max TTL is flagged and no certificate is issued', async () => {
    await page.getByLabel('Common name').fill('app.example.com');
    await fillTtl(notValidAfter, '9999999999', 'd');
    await page.getByRole('button', { name: 'Generate' }).click();

    await expect(page.getByText(TOO_LARGE('TTL'))).toBeVisible();
    await highlightText(page, page.locator('form'), TOO_LARGE('TTL'));
    expect(writes).toEqual([]);
  });

  await test.step('a valid TTL issues the certificate', async () => {
    await fillTtl(notValidAfter, '30');
    await page.getByRole('button', { name: 'Generate' }).click();
    await expect.poll(() => writes).toEqual([`POST /v1/${MOUNT}/issue/${ROLE}`]);
    await holdForVideo(page);
  });
});
