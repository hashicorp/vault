/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { test } from '../../fixtures/demo';
import { ensureMount, rootApi } from '../../fixtures/root-api';
import { isVideoEnabled } from '../../video-config';

import type { APIRequestContext } from '@playwright/test';

const LDAP_PATH = 'readonly-ldap';

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
  await ensureRole(api, 'static-role/readonly-static', {
    username: 'readonly-static',
    dn: 'cn=readonly-static,dc=example,dc=com',
    rotation_period: '24h',
    skip_import_rotation: true,
  });
  await ensureRole(api, 'role/readonly-dynamic', {
    creation_ldif: 'dn: cn={{.Username}},dc=example,dc=com\nobjectClass: person\ncn: {{.Username}}\nsn: x\n',
    deletion_ldif: 'dn: cn={{.Username}},dc=example,dc=com\nchangetype: delete\n',
  });
  await api.dispose();
});

// With no edit, delete or rotate access the menu would be empty, so it is not offered at all.
test('read-only user is not offered the Manage menu on LDAP role details', async ({ page }) => {
  await page.goto(`secrets-engines/${LDAP_PATH}/ldap/roles`);

  for (const [type, name] of [
    ['static', 'readonly-static'],
    ['dynamic', 'readonly-dynamic'],
  ]) {
    await test.step(`${type} role details hide Manage`, async () => {
      // The row's inner link is named after the role and its type badge.
      await page.getByRole('link', { name: `${name} ${type}`, exact: true }).click();
      await expect(page.getByRole('heading', { name, level: 1 })).toBeVisible();

      // Point at the empty space left of Get credentials, where Manage would otherwise sit.
      const getCredentials = await page.getByRole('link', { name: 'Get credentials' }).boundingBox();
      if (getCredentials) {
        await page.mouse.move(getCredentials.x - 70, getCredentials.y + getCredentials.height / 2);
      }
      if (isVideoEnabled) await page.waitForTimeout(1_500);
      await expect(page.getByRole('button', { name: 'Manage' })).toHaveCount(0);

      await page.getByRole('link', { name: 'Roles' }).click();
    });
  }
});

// Creating a role needs create access, so a read-only user is not sent to a form it cannot submit.
test('read-only user is not offered Create role on LDAP pages', async ({ page }) => {
  await test.step('roles list has no Create role', async () => {
    await page.goto(`secrets-engines/${LDAP_PATH}/ldap/roles`);
    await expect(page.getByRole('link', { name: 'readonly-static static', exact: true })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Create role' })).toHaveCount(0);
  });

  await test.step('overview offers Create new for libraries only', async () => {
    await page.getByRole('link', { name: 'Overview' }).click();
    // Production builds strip data-test attributes, so the cards are told apart by their links.
    const createNew = page.getByRole('link', { name: 'Create new' });
    await expect(createNew).toHaveCount(1);
    await expect(createNew).toHaveAttribute('href', /\/libraries\/create$/);
    if (isVideoEnabled) {
      // The Roles card mirrors the Libraries card, so its Create new link would sit at the same
      // offset from the card's left edge. Circle that empty spot.
      const link = await createNew.boundingBox();
      const librariesText = await page
        .getByText('The total number of libraries that have been created')
        .boundingBox();
      const rolesText = await page.getByText('The total number of roles that have been set up').boundingBox();
      if (link && librariesText && rolesText) {
        const centerX = link.x + link.width / 2 - (librariesText.x - rolesText.x);
        const centerY = link.y + link.height / 2;
        const radius = 14;
        await page.mouse.move(centerX + radius, centerY, { steps: 20 });
        // Recording slows every Playwright call, so one loop of 12 points already reads as a smooth circle.
        for (let step = 1; step <= 12; step++) {
          const angle = (step / 12) * Math.PI * 2;
          await page.mouse.move(centerX + radius * Math.cos(angle), centerY + radius * Math.sin(angle));
        }
        await page.waitForTimeout(1_000);
      }
    }
  });
});
