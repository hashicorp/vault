/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { test } from '../../fixtures/demo';
import { isVideoEnabled } from '../../video-config';
import { ensureMount, rootApi } from '../../fixtures/root-api';

const MOUNT = 'pki-routes';
const ROLE = 'routes-role';

// Seeding is idempotent because Playwright re-runs beforeAll in a fresh worker after a failure.
test.beforeAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await ensureMount(api, MOUNT, 'pki');
  const issuers = await api.get(`/v1/${MOUNT}/issuers?list=true`);
  if (!issuers.ok()) {
    await expect(
      await api.post(`/v1/${MOUNT}/root/generate/internal`, {
        data: { common_name: 'routes.example.com', issuer_name: 'routes-issuer', key_name: 'routes-key' },
      })
    ).toBeOK();
  }
  await expect(
    await api.post(`/v1/${MOUNT}/roles/${ROLE}`, { data: { issuer_ref: 'default', allow_any_name: true } })
  ).toBeOK();
  await api.dispose();
});

// Pins every PKI URL and its landing page so inlining the engine cannot change a route.
test('can open every PKI route by URL and by navigation', async ({ page, playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  const [issuerId] = (await (await api.get(`/v1/${MOUNT}/issuers?list=true`)).json()).data.keys;
  const [keyId] = (await (await api.get(`/v1/${MOUNT}/keys?list=true`)).json()).data.keys;
  await api.dispose();
  const base = `secrets-engines/${MOUNT}/pki`;
  const heading = page.getByRole('heading', { level: 1 });

  const routes: Array<[string, string]> = [
    ['overview', '/overview'],
    ['roles', '/roles'],
    ['roles/create', '/roles/create'],
    [`roles/${ROLE}/details`, `/roles/${ROLE}/details`],
    [`roles/${ROLE}/edit`, `/roles/${ROLE}/edit`],
    [`roles/${ROLE}/generate`, `/roles/${ROLE}/generate`],
    [`roles/${ROLE}/sign`, `/roles/${ROLE}/sign`],
    ['issuers', '/issuers'],
    ['issuers/import', '/issuers/import'],
    ['issuers/generate-root', '/issuers/generate-root'],
    ['issuers/generate-intermediate', '/issuers/generate-intermediate'],
    [`issuers/${issuerId}/details`, `/issuers/${issuerId}/details`],
    [`issuers/${issuerId}/edit`, `/issuers/${issuerId}/edit`],
    [`issuers/${issuerId}/sign`, `/issuers/${issuerId}/sign`],
    ['keys', '/keys'],
    ['keys/create', '/keys/create'],
    ['keys/import', '/keys/import'],
    [`keys/${keyId}/details`, `/keys/${keyId}/details`],
    [`keys/${keyId}/edit`, `/keys/${keyId}/edit`],
    ['certificates', '/certificates'],
    ['tidy', '/tidy'],
    ['tidy/auto/configure', '/tidy/auto/configure'],
    ['tidy/manual', '/tidy/manual'],
  ];

  for (const [visit, landing] of routes) {
    await test.step(`${visit} lands on ${landing}`, async () => {
      await page.goto(`${base}/${visit}`);
      await expect(page).toHaveURL(new RegExp(`/ui/vault/${base}${landing}$`));
      await expect(heading).toBeVisible();
    });
  }

  await test.step('the engine root redirects to overview', async () => {
    await page.goto(base);
    await expect(page).toHaveURL(new RegExp(`/ui/vault/${base}/overview$`));
    // Leave the landing page on screen long enough to read in a recording.
    if (isVideoEnabled) await page.waitForTimeout(1500);
  });

  await test.step('tabs navigate between roles, issuers, keys and certificates', async () => {
    await page.goto(`${base}/overview`);
    for (const [tab, landing] of [
      ['Roles', '/roles'],
      ['Issuers', '/issuers'],
      ['Keys', '/keys'],
      ['Certificates', '/certificates'],
    ]) {
      await page.getByRole('link', { name: tab, exact: true }).click();
      await expect(page).toHaveURL(new RegExp(`/ui/vault/${base}${landing}$`));
    }
  });
});
