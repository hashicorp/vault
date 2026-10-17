/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { test } from '../../fixtures/demo';
import { ensureMount, rootApi } from '../../fixtures/root-api';
import { isVideoEnabled } from '../../video-config';

import type { APIRequestContext } from '@playwright/test';

const LDAP_PATH = 'ldap-role-actions';

// Writing the same static role twice fails because skip_import_rotation is create-only, and setup
// re-runs in a fresh worker after a failure.
const ensureRole = async (api: APIRequestContext, path: string, data: Record<string, unknown>) => {
  if ((await api.get(`/v1/${LDAP_PATH}/${path}`)).ok()) return;
  await expect(await api.post(`/v1/${LDAP_PATH}/${path}`, { data })).toBeOK();
};

test.beforeAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await ensureMount(api, LDAP_PATH, 'ldap');
  // No LDAP server is needed: roles below are created without contacting it.
  await expect(
    await api.post(`/v1/${LDAP_PATH}/config`, {
      data: { url: 'ldap://127.0.0.1:389', binddn: 'cn=admin,dc=example,dc=com', bindpass: 'admin' },
    })
  ).toBeOK();
  await ensureRole(api, 'static-role/actions-static', {
    username: 'actions-static',
    dn: 'cn=actions-static,dc=example,dc=com',
    rotation_period: '24h',
    skip_import_rotation: true,
  });
  await ensureRole(api, 'role/actions-dynamic', {
    creation_ldif: 'dn: cn={{.Username}},dc=example,dc=com\nobjectClass: person\ncn: {{.Username}}\nsn: x\n',
    deletion_ldif: 'dn: cn={{.Username}},dc=example,dc=com\nchangetype: delete\n',
  });
  await api.dispose();
});

// Hiding an empty Manage menu must not hide it from users who are allowed to act on the role.
test('privileged user is offered the Manage menu on LDAP role details', async ({ page }) => {
  await page.goto(`secrets-engines/${LDAP_PATH}/ldap/roles`);

  for (const [type, name, actions] of [
    ['static', 'actions-static', ['Edit role', 'Rotate credentials', 'Delete role']],
    ['dynamic', 'actions-dynamic', ['Edit role', 'Delete role']],
  ] as const) {
    await test.step(`${type} role details offer ${actions.join(', ')}`, async () => {
      // The row's inner link is named after the role and its type badge.
      await page.getByRole('link', { name: `${name} ${type}`, exact: true }).click();
      await expect(page.getByRole('heading', { name, level: 1 })).toBeVisible();

      await page.getByRole('button', { name: 'Manage' }).click();
      if (isVideoEnabled) await page.waitForTimeout(1_500);
      // HDS renders the items as links and buttons inside the opened list.
      for (const action of actions) {
        await expect(page.getByRole('listitem').filter({ hasText: action })).toBeVisible();
      }
      await page.keyboard.press('Escape');

      await page.getByRole('link', { name: 'Roles' }).click();
    });
  }
});
