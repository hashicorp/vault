/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { test } from '../../fixtures/raft-cluster';

// Login and unseal enter secrets; do not retain them in trace/video/screenshot artifacts.
test.use({
  storageState: { cookies: [], origins: [] },
  trace: 'off',
  video: 'off',
  screenshot: 'off',
});

// The cluster state is polled every 10 seconds, so a seal that happens outside the UI must be noticed
// without any user action. This runs on a disposable cluster because it seals the node it talks to.
test('moves to the unseal page when the cluster is sealed elsewhere, then unseals through the UI', async ({
  page,
  raftCluster,
}) => {
  test.setTimeout(120_000);
  const { leader, token, unsealKey } = raftCluster;
  const readSealed = async () => {
    const response = await page.request.get(`${leader.apiURL}/v1/sys/seal-status`);
    await expect(response).toBeOK();
    return (await response.json()).sealed;
  };

  await test.step('signs in to an unsealed cluster', async () => {
    await page.goto(`${leader.apiURL}/ui/vault/auth`);
    await page.getByRole('textbox', { name: 'Token' }).fill(token);
    await page.getByRole('button', { name: 'Sign in' }).click();
    await expect(page).toHaveURL(/\/ui\/vault\/dashboard/);
    expect(await readSealed()).toBe(false);
  });

  await test.step('notices a seal made through the API', async () => {
    const sealed = await page.request.put(`${leader.apiURL}/v1/sys/seal`, {
      headers: { 'X-Vault-Token': token },
    });
    await expect(sealed).toBeOK();
    expect(await readSealed()).toBe(true);
    await expect(page).toHaveURL(/\/ui\/vault\/unseal/, { timeout: 45_000 });
    await expect(page.getByRole('heading', { name: 'Unseal Vault' })).toBeVisible();
    await expect(page.getByText('Vault is sealed')).toBeVisible();
  });

  await test.step('unseals through the UI and returns to the dashboard, still signed in', async () => {
    await page.getByRole('textbox', { name: 'Unseal Key Portion' }).fill(unsealKey);
    await page.getByRole('button', { name: 'Unseal', exact: true }).click();
    await expect.poll(readSealed, { timeout: 30_000 }).toBe(false);
    await expect(page).toHaveURL(/\/ui\/vault\/dashboard/, { timeout: 45_000 });
    await expect(page.getByRole('heading', { name: /^Vault v/ })).toBeVisible();
  });
});
