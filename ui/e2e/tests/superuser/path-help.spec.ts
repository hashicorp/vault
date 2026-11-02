/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { randomUUID } from 'node:crypto';
import { expect } from '@playwright/test';
import { test } from '../../fixtures/demo';
import { rootApi } from '../../fixtures/root-api';

import type { APIRequestContext, Page } from '@playwright/test';

let api: APIRequestContext;
let mount: string;

const readToken = (page: Page): Promise<string> =>
  page.evaluate(() => {
    const key = Object.keys(localStorage).find((key) => key.startsWith('vault-token'));
    if (!key) throw new Error('No stored Vault token');
    return JSON.parse(localStorage.getItem(key) as string).token as string;
  });

test.beforeEach(async ({ playwright }, testInfo) => {
  api = await rootApi(playwright, testInfo);
  mount = `path-help-${randomUUID()}`;
  await expect(await api.post(`/v1/sys/auth/${mount}`, { data: { type: 'userpass' } })).toBeOK();
});

test.afterEach(async () => {
  try {
    await expect(await api.delete(`/v1/sys/auth/${mount}`)).toBeOK();
  } finally {
    await api.dispose();
  }
});

test('generated-user fields preserve create, edit and delete outcomes', async ({ page }) => {
  const username = 'metadata-user';
  const userPath = `/v1/auth/${mount}/users/${username}`;
  const listRoute = `access/${mount}/item/user`;
  let token: string;

  await test.step('schema fields submit a user and persist token policy values', async () => {
    await page.goto(`${listRoute}/create`);
    token = await readToken(page);
    await expect(page.getByRole('textbox', { name: 'Username', exact: true })).toBeVisible();
    await page.getByRole('textbox', { name: 'Username', exact: true }).fill(username);
    await page.getByRole('textbox', { name: 'password', exact: true }).fill('disposable-password');
    await page.getByRole('button', { name: 'Tokens', exact: true }).click();
    await page.locator('[data-test-input="tokenPolicies"] textarea').first().fill('default');
    const saved = page.waitForResponse(
      (response) => response.url().endsWith(userPath) && response.request().method() === 'POST'
    );
    await page.getByRole('button', { name: 'Save', exact: true }).click();
    expect((await saved).ok()).toBe(true);
    await expect(page.getByRole('link', { name: username, exact: true })).toBeVisible();
    const response = await page.request.get(userPath, { headers: { 'X-Vault-Token': token } });
    await expect(response).toBeOK();
    expect((await response.json()).data.token_policies).toEqual(['default']);
    await page.reload();
    await expect(page.getByRole('link', { name: username, exact: true })).toBeVisible();
  });

  await test.step('read and edit preserve schema field values after reload', async () => {
    await page.getByRole('link', { name: username, exact: true }).click();
    await expect(page.getByRole('heading', { name: username, exact: true })).toBeVisible();
    await page.getByRole('link', { name: 'Edit user', exact: true }).click();
    await page.getByRole('button', { name: 'Tokens', exact: true }).click();
    await expect(page.locator('[data-test-input="tokenPolicies"] textarea').first()).toHaveValue('default');
    await page.locator('[data-test-input="tokenPolicies"] textarea').nth(1).fill('metadata-policy');
    const saved = page.waitForResponse(
      (response) => response.url().endsWith(userPath) && response.request().method() === 'PUT'
    );
    await page.getByRole('button', { name: 'Save', exact: true }).click();
    expect((await saved).ok()).toBe(true);
    await expect(page.getByRole('link', { name: username, exact: true })).toBeVisible();
    const response = await page.request.get(userPath, { headers: { 'X-Vault-Token': token } });
    await expect(response).toBeOK();
    expect((await response.json()).data.token_policies).toEqual(['default', 'metadata-policy']);
    await page.goto(`${listRoute}/edit/${username}`);
    await page.getByRole('button', { name: 'Tokens', exact: true }).click();
    await expect(page.locator('[data-test-input="tokenPolicies"] textarea').nth(1)).toHaveValue(
      'metadata-policy'
    );
    await page.reload();
    await page.getByRole('button', { name: 'Tokens', exact: true }).click();
    await expect(page.locator('[data-test-input="tokenPolicies"] textarea').first()).toHaveValue('default');
    await expect(page.locator('[data-test-input="tokenPolicies"] textarea').nth(1)).toHaveValue(
      'metadata-policy'
    );
  });

  await test.step('delete removes the user from backend state and the reloaded list', async () => {
    await page.goto(listRoute);
    await page.getByRole('link', { name: username, exact: true }).click();
    const deleted = page.waitForResponse(
      (response) => response.url().endsWith(userPath) && response.request().method() === 'DELETE'
    );
    await page.getByRole('button', { name: 'Delete user', exact: true }).click();
    await page.getByRole('button', { name: 'Confirm', exact: true }).click();
    expect((await deleted).ok()).toBe(true);
    await expect(page).toHaveURL(new RegExp(`/access/${mount}/item/user$`));
    await expect(page.getByRole('link', { name: username, exact: true })).toHaveCount(0);
    expect((await page.request.get(userPath, { headers: { 'X-Vault-Token': token } })).status()).toBe(404);
    await page.reload();
    await expect(page.getByRole('heading', { name: 'No users yet', exact: true })).toBeVisible();
  });
});

test('generated-user create surfaces backend validation and can be retried', async ({ page }) => {
  const userPath = `/v1/auth/${mount}/users/validation-user`;
  await page.goto(`access/${mount}/item/user/create`);
  const token = await readToken(page);
  await page.getByRole('textbox', { name: 'Username', exact: true }).fill('validation-user');

  await test.step('a missing password fails without creating a user', async () => {
    const rejected = page.waitForResponse(
      (response) => response.url().endsWith(userPath) && response.request().method() === 'POST'
    );
    await page.getByRole('button', { name: 'Save', exact: true }).click();
    expect((await rejected).status()).toBe(400);
    await expect(page.locator('[data-test-message-error]')).toContainText(/password/i);
    expect((await page.request.get(userPath, { headers: { 'X-Vault-Token': token } })).status()).toBe(404);
  });

  await test.step('a corrected submission succeeds and survives reload', async () => {
    await page.getByRole('textbox', { name: 'password', exact: true }).fill('disposable-password');
    const saved = page.waitForResponse(
      (response) => response.url().endsWith(userPath) && response.request().method() === 'POST'
    );
    await page.getByRole('button', { name: 'Save', exact: true }).click();
    expect((await saved).ok()).toBe(true);
    await expect(page.getByRole('link', { name: 'validation-user', exact: true })).toBeVisible();
    await expect(await page.request.get(userPath, { headers: { 'X-Vault-Token': token } })).toBeOK();
    await page.reload();
    await expect(page.getByRole('link', { name: 'validation-user', exact: true })).toBeVisible();
  });
});

test('generated-user read-only token cannot change persisted data', async ({ page }) => {
  const policyName = `${mount}-readonly`;
  const userPath = `/v1/auth/${mount}/users/readonly-user`;
  await expect(
    await api.post(userPath, { data: { password: 'disposable-password', token_policies: ['default'] } })
  ).toBeOK();
  await expect(
    await api.put(`/v1/sys/policies/acl/${policyName}`, {
      data: {
        policy: `
          path "sys/auth" { capabilities = ["read"] }
          path "sys/auth/${mount}" { capabilities = ["read"] }
          path "sys/auth/${mount}/*" { capabilities = ["read"] }
          path "auth/${mount}/*" { capabilities = ["read", "list"] }
        `,
      },
    })
  ).toBeOK();
  const tokenResponse = await api.post('/v1/auth/token/create', {
    data: { policies: [policyName], ttl: '15m' },
  });
  await expect(tokenResponse).toBeOK();
  const token = (await tokenResponse.json()).auth.client_token as string;

  try {
    page.on('dialog', (dialog) => dialog.accept());
    await page.goto('dashboard');
    // Replace the session privately so a demo cannot capture token entry.
    await page.evaluate(
      ({ token, policyName }) => {
        const key = Object.keys(localStorage).find((key) => key.startsWith('vault-token'));
        if (!key) throw new Error('No stored Vault token to replace');
        const stored = JSON.parse(localStorage.getItem(key) as string);
        localStorage.setItem(key, JSON.stringify({ ...stored, token, policies: ['default', policyName] }));
      },
      { token, policyName }
    );

    await test.step('read-only values render but a UI update is denied without changing state', async () => {
      await page.goto(`access/${mount}/item/user/edit/readonly-user`);
      await page.getByRole('button', { name: 'Tokens', exact: true }).click();
      await expect(page.locator('[data-test-input="tokenPolicies"] textarea').first()).toHaveValue('default');
      await page.locator('[data-test-input="tokenPolicies"] textarea').first().fill('unauthorized-policy');
      const rejected = page.waitForResponse(
        (response) => response.url().endsWith(userPath) && response.request().method() === 'PUT'
      );
      await page.getByRole('button', { name: 'Save', exact: true }).click();
      expect((await rejected).status()).toBe(403);
      await expect(page.locator('[data-test-message-error]')).toContainText(/permission denied/i);
      const response = await page.request.get(userPath, { headers: { 'X-Vault-Token': token } });
      await expect(response).toBeOK();
      expect((await response.json()).data.token_policies).toEqual(['default']);
      await page.reload();
      await page.getByRole('button', { name: 'Tokens', exact: true }).click();
      await expect(page.locator('[data-test-input="tokenPolicies"] textarea').first()).toHaveValue('default');
    });

    await test.step('denied UI create does not create backend data', async () => {
      await page.goto(`access/${mount}/item/user/create`);
      await page.getByRole('textbox', { name: 'Username', exact: true }).fill('denied-user');
      await page.getByRole('textbox', { name: 'password', exact: true }).fill('disposable-password');
      const rejected = page.waitForResponse(
        (response) =>
          response.url().endsWith(`/v1/auth/${mount}/users/denied-user`) &&
          response.request().method() === 'POST'
      );
      await page.getByRole('button', { name: 'Save', exact: true }).click();
      expect((await rejected).status()).toBe(403);
      await expect(page.locator('[data-test-message-error]')).toContainText(/permission denied/i);
      expect((await api.get(`/v1/auth/${mount}/users/denied-user`)).status()).toBe(404);
    });

    await test.step('denied UI delete preserves the user after reload', async () => {
      await page.goto(`access/${mount}/item/user`);
      const row = page.locator('[data-test-list-item-link]').filter({ hasText: 'readonly-user' });
      await row.getByRole('button', { name: 'Overflow options', exact: true }).click();
      await page.getByRole('button', { name: 'Delete user', exact: true }).click();
      const rejected = page.waitForResponse(
        (response) => response.url().endsWith(userPath) && response.request().method() === 'DELETE'
      );
      // ListItem yields twice; the menu's second dialog owns its callMethod task.
      await page
        .getByRole('dialog', { name: 'Delete user?', exact: true })
        .last()
        .getByRole('button', { name: 'Confirm', exact: true })
        .click();
      expect((await rejected).status()).toBe(403);
      await expect(page.getByText('There was an error deleting this user', { exact: false })).toBeVisible();
      expect((await (await api.get(userPath)).json()).data.token_policies).toEqual(['default']);
      await page.reload();
      await expect(page.getByRole('link', { name: 'readonly-user', exact: true })).toBeVisible();
    });
  } finally {
    await expect(await api.post('/v1/auth/token/revoke', { data: { token } })).toBeOK();
    await expect(await api.delete(`/v1/sys/policies/acl/${policyName}`)).toBeOK();
  }
});
