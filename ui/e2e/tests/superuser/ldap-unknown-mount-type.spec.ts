/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { test } from '../../fixtures/demo';
import { ensureMount, rootApi } from '../../fixtures/root-api';
import { isVideoEnabled } from '../../video-config';

import type { Locator, Page } from '@playwright/test';

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

// Drags the mouse from the start of `fromText` to the end of `toText` inside `scope`, so the recording
// selects exactly the text being demoed. Text nodes are searched in document order.
const dragSelect = async (page: Page, scope: Locator, fromText: string, toText = fromText) => {
  const bounds = await scope.evaluate(
    (root, { fromText, toText }) => {
      const walker = document.createTreeWalker(root, NodeFilter.SHOW_TEXT);
      let start: { node: Node; offset: number } | undefined;
      let end: { node: Node; offset: number } | undefined;
      for (let node = walker.nextNode(); node && !end; node = walker.nextNode()) {
        const text = node.textContent ?? '';
        if (!start && text.includes(fromText)) start = { node, offset: text.indexOf(fromText) };
        if (start && text.includes(toText)) end = { node, offset: text.indexOf(toText) + toText.length };
      }
      if (!start || !end) return null;
      root.scrollIntoView({ block: 'center' });
      const range = document.createRange();
      range.setStart(start.node, start.offset);
      range.setEnd(end.node, end.offset);
      const rects = [...range.getClientRects()].filter((rect) => rect.width > 0);
      const first = rects[0];
      const last = rects[rects.length - 1];
      return {
        from: { x: first.left + 1, y: first.top + first.height / 2 },
        to: { x: last.right - 1, y: last.top + last.height / 2 },
      };
    },
    { fromText, toText }
  );
  if (!bounds) throw new Error(`Text "${fromText}" was not found to highlight`);
  await page.mouse.move(bounds.from.x, bounds.from.y, { steps: 15 });
  await page.waitForTimeout(400);
  await page.mouse.down();
  await page.mouse.move(bounds.to.x, bounds.to.y, { steps: 40 });
  await page.mouse.up();
  await page.waitForTimeout(1_800);
};

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
    if (isVideoEnabled) await page.waitForTimeout(1_500);
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
    if (isVideoEnabled) {
      const notice = page.locator('.hds-alert').filter({ hasText: NOTICE_TITLE });
      await dragSelect(page, notice, NOTICE_TITLE, 'when creating static roles.');
    }
  });
});
