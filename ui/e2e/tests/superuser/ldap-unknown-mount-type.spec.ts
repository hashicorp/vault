/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { highlightText, holdForVideo, test } from '../../fixtures/demo';
import { ensureMount, rootApi } from '../../fixtures/root-api';

import type { Page } from '@playwright/test';

const LDAP_PATH = 'ldap-unknown-type';
const POLICY = 'ldap-role-writer';
// The alert says the same thing in lower case, so match the helper text case-sensitively.
const NOTICE_TITLE = 'Unable to determine LDAP configuration mode';
const NOTICE_DESCRIPTION =
  "You don't have permission to read the LDAP secrets engine configuration. Self-managed engines don't support dynamic roles and require password and DN when creating static roles.";

// Can create and manage roles but cannot read the mount config, so it cannot tell whether the
// mount is self-managed.
const WRITER_POLICY = `
path "${LDAP_PATH}/static-role/*" { capabilities = ["create", "read", "update", "list"] }
path "${LDAP_PATH}/static-role" { capabilities = ["list"] }
path "${LDAP_PATH}/role/*" { capabilities = ["create", "read", "update", "list"] }
path "${LDAP_PATH}/role" { capabilities = ["list"] }
`;

let writerToken = '';

test.beforeAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await ensureMount(api, LDAP_PATH, 'ldap');
  // A self-managed mount never binds as an admin, so no LDAP server is contacted here.
  await expect(
    await api.post(`/v1/${LDAP_PATH}/config`, { data: { url: 'ldap://127.0.0.1:389', self_managed: true } })
  ).toBeOK();
  await expect(await api.put(`/v1/sys/policies/acl/${POLICY}`, { data: { policy: WRITER_POLICY } })).toBeOK();
  const tokenResponse = await api.post('/v1/auth/token/create', { data: { policies: [POLICY], ttl: '1h' } });
  await expect(tokenResponse).toBeOK();
  writerToken = (await tokenResponse.json()).auth.client_token;
  await api.dispose();
});

test.afterAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await api.delete(`/v1/sys/mounts/${LDAP_PATH}`);
  await api.delete(`/v1/sys/policies/acl/${POLICY}`);
  await api.dispose();
});

// Swaps the persona's stored token rather than signing in, so the token is never typed into a
// field that a recording could capture.
const useToken = async (page: Page, token: string, policies: string[]) => {
  await page.evaluate(
    ({ token, policies }) => {
      const key = Object.keys(localStorage).find((k) => k.startsWith('vault-token'));
      if (!key) throw new Error('No stored Vault token to replace');
      const stored = JSON.parse(localStorage.getItem(key) as string);
      localStorage.setItem(key, JSON.stringify({ ...stored, token, policies }));
    },
    { token, policies }
  );
};

test('LDAP create role flags an unknown mount type when the config cannot be read', async ({ page }) => {
  await test.step('a user who can read the config gets the self-managed form', async () => {
    await page.goto(`secrets-engines/${LDAP_PATH}/ldap/roles/create`);
    await expect(page.getByRole('heading', { name: 'Create static role', level: 1 })).toBeVisible();
    await expect(page.getByText(NOTICE_TITLE)).toHaveCount(0);
    await expect(page.getByRole('radio', { name: /Static role/ })).toHaveCount(0);
    await holdForVideo(page);
  });

  await useToken(page, writerToken, [POLICY]);

  await test.step('a user who can create roles is still offered Create role', async () => {
    await page.goto(`secrets-engines/${LDAP_PATH}/ldap/roles`);
    await expect(page.getByRole('link', { name: 'Create role' })).toBeVisible();
    await page.getByRole('link', { name: 'Create role' }).click();
  });

  await test.step('a user who cannot read the config is told the mount type is unknown', async () => {
    await expect(page.getByRole('heading', { name: 'Create Role', level: 1 })).toBeVisible();
    await expect(page.getByText(NOTICE_TITLE)).toBeVisible();
    await expect(page.getByText(NOTICE_DESCRIPTION)).toBeVisible();
    // Both role types stay available because the mount could be root-managed.
    await expect(page.getByRole('radio', { name: /Static role/ })).toBeVisible();
    await expect(page.getByRole('radio', { name: /Dynamic role/ })).toBeVisible();
    const notice = page.locator('.hds-alert').filter({ hasText: NOTICE_TITLE });
    await highlightText(page, notice, NOTICE_TITLE, 'when creating static roles.');
  });
});
