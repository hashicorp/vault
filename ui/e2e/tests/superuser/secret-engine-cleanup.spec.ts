/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect, highlightElement, holdForVideo, test } from '../../fixtures/demo';
import { ensureMount, rootApi } from '../../fixtures/root-api';

const MOUNT_PATH = 'edm-item21-kv';

test.beforeAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await ensureMount(api, MOUNT_PATH, 'kv-v2');
  await api.dispose();
});

test.afterAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  try {
    const response = await api.get(`/v1/sys/mounts/${MOUNT_PATH}`);
    if (response.ok()) await expect(await api.delete(`/v1/sys/mounts/${MOUNT_PATH}`)).toBeOK();
  } finally {
    await api.dispose();
  }
});

// Logging out after loading engine Resources must not unload a removed Ember Data model.
test('superuser can log out after viewing secrets engines', async ({ page }) => {
  await test.step('load and reload the mounted engine list', async () => {
    await page.goto('secrets-engines');
    await expect(page.getByRole('heading', { name: 'Secrets engines', level: 1 })).toBeVisible();

    const mountLink = page.getByRole('link', { name: `${MOUNT_PATH}/`, exact: true });
    await expect(mountLink).toBeVisible();
    await page.reload();
    await expect(mountLink).toBeVisible();
    await highlightElement(page, mountLink);
    await holdForVideo(page);
  });

  await test.step('return to authentication after logout', async () => {
    await page.getByRole('button', { name: 'User menu' }).click();
    await page.getByRole('link', { name: 'Log out' }).click();
    const tokenField = page.getByRole('textbox', { name: 'Token' });
    await expect(tokenField).toBeVisible();
    await highlightElement(page, tokenField);
    await holdForVideo(page);
  });
});
