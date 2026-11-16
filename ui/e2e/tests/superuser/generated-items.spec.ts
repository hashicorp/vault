/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { randomUUID } from 'node:crypto';
import { expect } from '@playwright/test';
import { test } from '../../fixtures/demo';
import { rootApi } from '../../fixtures/root-api';

import type { APIRequestContext, Page, Request } from '@playwright/test';

// Generated item pages (auth method users, groups, roles, certificates) are built at runtime from each
// mount's OpenAPI help, so these scenarios pin the request contract and rendered values for two methods.

let api: APIRequestContext;
let mount: string;

const writeRequest = (page: Page, path: string, method: 'POST' | 'PUT' | 'DELETE') =>
  page.waitForRequest((request: Request) => request.url().endsWith(path) && request.method() === method);

test.beforeEach(async ({ playwright }, testInfo) => {
  api = await rootApi(playwright, testInfo);
  mount = `generated-${randomUUID().slice(0, 8)}`;
});

test.afterEach(async () => {
  try {
    await api.delete(`/v1/sys/auth/${mount}`);
  } finally {
    await api.dispose();
  }
});

test('userpass user create, show and edit send the schema payload', async ({ page }) => {
  const userPath = `/v1/auth/${mount}/users/gi-user`;
  await expect(await api.post(`/v1/sys/auth/${mount}`, { data: { type: 'userpass' } })).toBeOK();

  await test.step('create sends only edited values, without the path identifier', async () => {
    await page.goto(`access/${mount}/item/user/create`);
    await page.getByRole('textbox', { name: 'Username', exact: true }).fill('gi-user');
    await page.getByRole('textbox', { name: 'password', exact: true }).fill('disposable-password');
    await page.getByRole('button', { name: 'Tokens', exact: true }).click();
    await page.getByRole('textbox', { name: 'Maximum Uses of Generated Tokens', exact: true }).fill('5');
    await page
      .getByRole('textbox', { name: "Generated Token's Policies list item 0", exact: true })
      .fill('default');
    const created = writeRequest(page, userPath, 'POST');
    await page.getByRole('button', { name: 'Save', exact: true }).click();
    const request = await created;
    expect(request.postDataJSON()).toEqual({
      password: 'disposable-password',
      token_num_uses: 5,
      token_policies: ['default'],
    });
    expect((await request.response())?.ok()).toBe(true);
    await expect(page).toHaveURL(new RegExp(`/access/${mount}/item/user$`));
    await expect(page.getByText('Successfully saved user gi-user.', { exact: true })).toBeVisible();
    const { data } = await (await api.get(userPath)).json();
    expect(data.token_num_uses).toBe(5);
    expect(data.token_policies).toEqual(['default']);
  });

  await test.step('show renders schema values and hides the path identifier', async () => {
    await page.getByRole('link', { name: 'gi-user', exact: true }).click();
    await expect(page.getByRole('heading', { name: 'gi-user', exact: true })).toBeVisible();
    await expect(page.locator('[data-test-row-value="Maximum Uses of Generated Tokens"]')).toHaveText('5');
    await expect(page.locator('[data-test-row-value="Password"]')).toHaveText(
      'This value is sensitive and cannot be shown.'
    );
    await expect(page.locator('[data-test-row-label="Username"]')).toHaveCount(0);
  });

  await test.step('edit keeps the identifier read-only and sends every non-blank value', async () => {
    await page.getByRole('link', { name: 'Edit user', exact: true }).click();
    const username = page.getByRole('textbox', { name: 'Username', exact: true });
    await expect(username).toHaveValue('gi-user');
    await expect(username).not.toBeEditable();
    await page.getByRole('button', { name: 'Tokens', exact: true }).click();
    const numUses = page.getByRole('textbox', { name: 'Maximum Uses of Generated Tokens', exact: true });
    await expect(numUses).toHaveValue('5');
    await numUses.fill('7');
    const updated = writeRequest(page, userPath, 'PUT');
    await page.getByRole('button', { name: 'Save', exact: true }).click();
    const request = await updated;
    expect(request.postDataJSON()).toEqual({
      alias_metadata: {},
      token_explicit_max_ttl: 0,
      token_max_ttl: 0,
      token_no_default_policy: false,
      token_num_uses: 7,
      token_period: 0,
      token_policies: ['default'],
      token_ttl: 0,
      token_type: 'default',
    });
    expect((await request.response())?.ok()).toBe(true);
    await expect(page).toHaveURL(new RegExp(`/access/${mount}/item/user$`));
    expect((await (await api.get(userPath)).json()).data.token_num_uses).toBe(7);
    await page.goto(`access/${mount}/item/user/show/gi-user`);
    await expect(page.locator('[data-test-row-value="Maximum Uses of Generated Tokens"]')).toHaveText('7');
  });
});

test('ldap groups render method tabs, create and delete from the list', async ({ page }) => {
  const groupPath = `/v1/auth/${mount}/groups/devs`;
  await expect(await api.post(`/v1/sys/auth/${mount}`, { data: { type: 'ldap' } })).toBeOK();

  await test.step('an empty item type shows its tabs and empty state', async () => {
    await page.goto(`access/${mount}/item/group`);
    const tabs = page.getByRole('navigation', { name: `navigation to manage ${mount}` });
    await expect(tabs.getByRole('link', { name: 'Groups', exact: true })).toBeVisible();
    await expect(tabs.getByRole('link', { name: 'Users', exact: true })).toBeVisible();
    await expect(page.getByRole('heading', { name: 'No groups yet', exact: true })).toBeVisible();
  });

  await test.step('create posts to the group path and lists the new group', async () => {
    await page.getByRole('link', { name: 'Create group', exact: true }).first().click();
    await page.getByRole('textbox', { name: 'Name', exact: true }).fill('devs');
    await page.getByRole('textbox', { name: 'Policies list item 0', exact: true }).fill('default');
    const created = writeRequest(page, groupPath, 'POST');
    await page.getByRole('button', { name: 'Save', exact: true }).click();
    const request = await created;
    expect(request.postDataJSON()).toEqual({ policies: ['default'] });
    expect((await request.response())?.ok()).toBe(true);
    await expect(page.getByRole('link', { name: 'devs', exact: true })).toBeVisible();
    expect((await (await api.get(groupPath)).json()).data.policies).toEqual(['default']);
  });

  await test.step('delete from the list removes the group from the backend', async () => {
    const row = page.locator('[data-test-list-item-link]').filter({ hasText: 'devs' });
    await row.getByRole('button', { name: 'Overflow options', exact: true }).click();
    await page.getByRole('button', { name: 'Delete group', exact: true }).click();
    const deleted = writeRequest(page, groupPath, 'DELETE');
    const dialog = page.getByRole('dialog', { name: 'Delete group?', exact: true });
    await expect(dialog).toHaveCount(1);
    await dialog.getByRole('button', { name: 'Confirm', exact: true }).click();
    expect((await (await deleted).response())?.ok()).toBe(true);
    await expect(page.getByText('Successfully deleted group devs.', { exact: true })).toBeVisible();
    expect((await api.get(groupPath)).status()).toBe(404);
    await page.reload();
    await expect(page.getByRole('heading', { name: 'No groups yet', exact: true })).toBeVisible();
  });

  await test.step('the second item type tab lists users', async () => {
    await page
      .getByRole('navigation', { name: `navigation to manage ${mount}` })
      .getByRole('link', { name: 'Users', exact: true })
      .click();
    await expect(page).toHaveURL(new RegExp(`/access/${mount}/item/user$`));
    await expect(page.getByRole('heading', { name: 'No users yet', exact: true })).toBeVisible();
  });
});

test('userpass user list paginates beyond one page', async ({ page }) => {
  await expect(await api.post(`/v1/sys/auth/${mount}`, { data: { type: 'userpass' } })).toBeOK();
  const names = Array.from({ length: 101 }, (_, i) => `user-${String(i).padStart(3, '0')}`);
  for (const name of names) {
    await expect(await api.post(`/v1/auth/${mount}/users/${name}`, { data: { password: 'p' } })).toBeOK();
  }

  await page.goto(`access/${mount}/item/user`);
  await expect(page.getByText('1–100 of 101', { exact: true })).toBeVisible();
  await expect(page.locator('[data-test-list-item-link]')).toHaveCount(100);
  await expect(page.getByRole('link', { name: 'user-000', exact: true })).toBeVisible();

  await page.getByRole('link', { name: 'page 2', exact: true }).click();
  await expect(page.getByText('101–101 of 101', { exact: true })).toBeVisible();
  await expect(page.locator('[data-test-list-item-link]')).toHaveCount(1);
  await expect(page.getByRole('link', { name: 'user-100', exact: true })).toBeVisible();
});
