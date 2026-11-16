/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { test } from '../../fixtures/demo';
import { ensureMount, rootApi } from '../../fixtures/root-api';
import { isVideoEnabled } from '../../video-config';

import type { APIRequestContext, Page } from '@playwright/test';

// Leaves the result on screen long enough to read in a recording; a no-op otherwise.
const hold = async (page: Page, ms = 1500) => {
  if (isVideoEnabled) await page.waitForTimeout(ms);
};

const LDAP_PATH = 'ldap-routing';
const UNCONFIGURED_PATH = 'ldap-routing-unconfigured';
const BASE = `secrets-engines/${LDAP_PATH}/ldap`;

const ROLE_LDIF = {
  creation_ldif: 'dn: cn={{.Username}},dc=example,dc=com\nobjectClass: person\ncn: {{.Username}}\nsn: x\n',
  deletion_ldif: 'dn: cn={{.Username}},dc=example,dc=com\nchangetype: delete\n',
};

// Writing the same static role twice fails because skip_import_rotation is create-only, and setup
// re-runs in a fresh worker after a failure.
const ensureWrite = async (api: APIRequestContext, path: string, data: Record<string, unknown>) => {
  if ((await api.get(`/v1/${LDAP_PATH}/${path}`)).ok()) return;
  await expect(await api.post(`/v1/${LDAP_PATH}/${path}`, { data })).toBeOK();
};

const staticRole = (name: string) => ({
  username: name,
  dn: `cn=${name},dc=example,dc=com`,
  rotation_period: '24h',
  skip_import_rotation: true,
});

test.beforeAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await ensureMount(api, LDAP_PATH, 'ldap');
  await ensureMount(api, UNCONFIGURED_PATH, 'ldap');
  // No LDAP server is needed: roles below are created without contacting it. Libraries cannot be
  // created without a reachable LDAP server, so library detail routes are covered by the QUnit acceptance tests.
  await expect(
    await api.post(`/v1/${LDAP_PATH}/config`, {
      data: { url: 'ldap://127.0.0.1:389', binddn: 'cn=admin,dc=example,dc=com', bindpass: 'admin' },
    })
  ).toBeOK();
  await ensureWrite(api, 'static-role/routing-static', staticRole('routing-static'));
  await ensureWrite(api, 'static-role/admin/child-static', staticRole('child-static'));
  await ensureWrite(api, 'role/routing-dynamic', ROLE_LDIF);
  await ensureWrite(api, 'role/admin/child-dynamic', ROLE_LDIF);
  await api.dispose();
});

test.afterAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await api.delete(`/v1/sys/mounts/${LDAP_PATH}`);
  await api.delete(`/v1/sys/mounts/${UNCONFIGURED_PATH}`);
  await api.dispose();
});

test('can follow the LDAP tabs and deep link to each top level page', async ({ page }) => {
  await test.step('overview links to the roles and libraries tabs', async () => {
    await page.goto(`${BASE}/overview`);
    await expect(page).toHaveURL(new RegExp(`/${BASE}/overview$`));
    await expect(page.getByRole('heading', { name: LDAP_PATH, level: 1 })).toBeVisible();
    await hold(page);

    await page.getByRole('link', { name: 'Roles', exact: true }).first().click();
    await expect(page).toHaveURL(new RegExp(`/${BASE}/roles$`));
    await page.getByRole('link', { name: 'Libraries', exact: true }).first().click();
    await expect(page).toHaveURL(new RegExp(`/${BASE}/libraries$`));
    await page.getByRole('link', { name: 'Overview', exact: true }).first().click();
    await expect(page).toHaveURL(new RegExp(`/${BASE}/overview$`));
  });

  await test.step('configuration and configure deep links load', async () => {
    await page.goto(`${BASE}/configuration`);
    await expect(page).toHaveURL(new RegExp(`/${BASE}/configuration$`));
    await expect(page.getByText('cn=admin,dc=example,dc=com')).toBeVisible();
    await hold(page);

    await page.goto(`${BASE}/configure`);
    await expect(page).toHaveURL(new RegExp(`/${BASE}/configure$`));
    await expect(page.getByRole('button', { name: 'Save' })).toBeVisible();
  });

  await test.step('create pages load', async () => {
    await page.goto(`${BASE}/roles/create`);
    await expect(page).toHaveURL(new RegExp(`/${BASE}/roles/create$`));
    await expect(page.getByRole('button', { name: 'Create role' })).toBeVisible();

    await page.goto(`${BASE}/libraries/create`);
    await expect(page).toHaveURL(new RegExp(`/${BASE}/libraries/create$`));
    await expect(page.getByRole('button', { name: 'Create library' })).toBeVisible();
  });
});

test('can deep link to static and dynamic role pages', async ({ page }) => {
  for (const [type, name] of [
    ['static', 'routing-static'],
    ['dynamic', 'routing-dynamic'],
  ]) {
    await test.step(`${type} role: details, edit and credentials`, async () => {
      await page.goto(`${BASE}/roles/${type}/${name}/details`);
      await expect(page).toHaveURL(new RegExp(`/${BASE}/roles/${type}/${name}/details$`));
      await expect(page.getByRole('heading', { name, level: 1 })).toBeVisible();
      await hold(page);

      await page.goto(`${BASE}/roles/${type}/${name}/edit`);
      await expect(page).toHaveURL(new RegExp(`/${BASE}/roles/${type}/${name}/edit$`));
      await expect(page.getByRole('button', { name: 'Save' })).toBeVisible();

      await page.goto(`${BASE}/roles/${type}/${name}/credentials`);
      await expect(page).toHaveURL(new RegExp(`/${BASE}/roles/${type}/${name}/credentials$`));
    });
  }

  await test.step('a role url without a tab lands on details', async () => {
    await page.goto(`${BASE}/roles/static/routing-static`);
    await expect(page).toHaveURL(new RegExp(`/${BASE}/roles/static/routing-static/details$`));
  });

  await test.step('role list rows open the role details', async () => {
    await page.goto(`${BASE}/roles`);
    await page.getByRole('link', { name: 'routing-dynamic dynamic', exact: true }).click();
    await expect(page).toHaveURL(new RegExp(`/${BASE}/roles/dynamic/routing-dynamic/details$`));
  });
});

test('can traverse hierarchical roles', async ({ page }) => {
  await test.step('the role list shows the folder and opens its subdirectory', async () => {
    await page.goto(`${BASE}/roles`);
    await expect(page.getByRole('link', { name: 'admin/', exact: false }).first()).toBeVisible();
    await page.goto(`${BASE}/roles/static/subdirectory/admin/`);
    await expect(page).toHaveURL(new RegExp(`/${BASE}/roles/static/subdirectory/admin/$`));
    await expect(page.getByRole('link', { name: 'child-static static', exact: true })).toBeVisible();
    await hold(page);
  });

  await test.step('a child role opens with an encoded name in the url', async () => {
    await page.getByRole('link', { name: 'child-static static', exact: true }).click();
    await expect(page).toHaveURL(new RegExp(`/${BASE}/roles/static/admin%2Fchild-static/details$`));
    await expect(page.getByRole('heading', { name: 'admin/child-static', level: 1 })).toBeVisible();
  });

  await test.step('dynamic roles have their own subdirectory url', async () => {
    await page.goto(`${BASE}/roles/dynamic/subdirectory/admin/`);
    await expect(page.getByRole('link', { name: 'child-dynamic dynamic', exact: true })).toBeVisible();
  });
});

test('can reach the library list, create and subdirectory pages', async ({ page }) => {
  await test.step('the library list is reachable from the tab and offers creation', async () => {
    await page.goto(`${BASE}/overview`);
    await page.getByRole('link', { name: 'Libraries', exact: true }).first().click();
    await expect(page).toHaveURL(new RegExp(`/${BASE}/libraries$`));
    await hold(page);
  });

  await test.step('the create page loads and the subdirectory url resolves', async () => {
    await page.goto(`${BASE}/libraries/create`);
    await expect(page).toHaveURL(new RegExp(`/${BASE}/libraries/create$`));
    await expect(page.getByRole('button', { name: 'Create library' })).toBeVisible();

    await page.goto(`${BASE}/libraries/subdirectory/admin/`);
    await expect(page).toHaveURL(new RegExp(`/${BASE}/libraries/subdirectory/admin/$`));
  });
});

test('is sent to configure when the mount has no configuration', async ({ page }) => {
  await page.goto(`secrets-engines/${UNCONFIGURED_PATH}/ldap/configuration`);
  await expect(page).toHaveURL(new RegExp(`/secrets-engines/${UNCONFIGURED_PATH}/ldap/configure$`));
  await expect(page.getByRole('button', { name: 'Save' })).toBeVisible();
  await hold(page);
});
