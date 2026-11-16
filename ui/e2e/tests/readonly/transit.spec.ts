/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { test } from '../../fixtures/demo';
import { rootApi } from '../../fixtures/root-api';

const MOUNT = 'readonly-transit';
const NAME = 'readonly-key';
const KEY_PATH = `/v1/${MOUNT}/keys/${NAME}`;

test('read-only Transit user can inspect a key but cannot change or use it', async ({
  page,
  playwright,
}, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  try {
    await expect(await api.post(`/v1/sys/mounts/${MOUNT}`, { data: { type: 'transit' } })).toBeOK();
    await expect(await api.post(KEY_PATH, { data: { type: 'aes256-gcm96' } })).toBeOK();
    await expect(await api.post(`${KEY_PATH}/config`, { data: { deletion_allowed: true } })).toBeOK();
    const original = await api.get(KEY_PATH);
    await expect(original).toBeOK();
    const originalData = (await original.json()).data;

    await page.goto(`secrets-engines/${MOUNT}/show/${NAME}?tab=details`);
    await expect(page.locator('[data-test-row-value="Type"]').last()).toHaveText('aes256-gcm96');
    await expect(page.getByRole('link', { name: 'Edit key', exact: true })).toHaveCount(0);
    await page.getByRole('link', { name: 'Versions', exact: true }).click();
    await expect(page.getByText('Version 1', { exact: true })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Rotate key', exact: true })).toHaveCount(0);

    const token = await page.evaluate(() => {
      const key = Object.keys(localStorage).find((key) => key.startsWith('vault-token'));
      return key ? (JSON.parse(localStorage.getItem(key) as string).token as string) : '';
    });
    expect(token).not.toBe('');
    const headers = { 'X-Vault-Token': token };
    expect(
      (await page.request.post(`${KEY_PATH}/config`, { headers, data: { deletion_allowed: false } })).status()
    ).toBe(403);
    expect((await page.request.post(`${KEY_PATH}/rotate`, { headers, data: {} })).status()).toBe(403);
    expect((await page.request.delete(KEY_PATH, { headers })).status()).toBe(403);

    await page.getByRole('link', { name: 'Key Actions', exact: true }).click();
    await page.getByRole('link', { name: /^Encrypt / }).click();
    await page.getByRole('textbox', { name: 'Plaintext' }).fill('permission denied');
    await page.getByRole('button', { name: 'Encrypt', exact: true }).click();
    await expect(page.locator('[data-test-message-error]')).toContainText('permission denied');
    const unchanged = await api.get(KEY_PATH);
    await expect(unchanged).toBeOK();
    expect((await unchanged.json()).data).toEqual(originalData);
    await page.reload();
    await expect(page.getByRole('button', { name: 'Encrypt', exact: true })).toBeVisible();
  } finally {
    await expect(await api.delete(`/v1/sys/mounts/${MOUNT}`)).toBeOK();
    await api.dispose();
  }
});
