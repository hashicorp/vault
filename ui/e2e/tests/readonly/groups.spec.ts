/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { test } from '../../fixtures/demo';
import { rootApi } from '../../fixtures/root-api';

const GROUP = 'readonly-group';
const EXTERNAL_GROUP = 'readonly-external-group';
const ALIAS = 'readonly-alias';

// Seeding is idempotent because Playwright re-runs beforeAll in a fresh worker after a failure.
test.beforeAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  await expect(await api.post(`/v1/identity/group/name/${GROUP}`)).toBeOK();
  await expect(
    await api.post(`/v1/identity/group/name/${EXTERNAL_GROUP}`, { data: { type: 'external' } })
  ).toBeOK();
  const { data: group } = await (await api.get(`/v1/identity/group/name/${EXTERNAL_GROUP}`)).json();
  if (!group.alias?.id) {
    const auth = await (await api.get('/v1/sys/auth')).json();
    const accessor = auth.data['token/'].accessor;
    await expect(
      await api.post('/v1/identity/group-alias', {
        data: { name: ALIAS, mount_accessor: accessor, canonical_id: group.id },
      })
    ).toBeOK();
  }
  await api.dispose();
});

// Read-only tokens must not see create or edit actions for groups or group aliases.
test('read-only user cannot create or edit groups', async ({ page }) => {
  await page.goto('access/identity/groups');
  await expect(page.getByRole('heading', { name: 'Groups', level: 1 })).toBeVisible();
  await expect(page.getByRole('link', { name: GROUP })).toBeVisible();

  await test.step('header create action is hidden', async () => {
    await expect(page.getByRole('link', { name: 'Create group' })).toHaveCount(0);
  });

  await test.step('group details hide the edit action', async () => {
    await page.getByRole('link', { name: GROUP }).click();
    await expect(page.getByRole('heading', { name: GROUP })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Details' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Edit group' })).toHaveCount(0);
  });
});

test('read-only user cannot create groups or edit group aliases from the aliases tab', async ({ page }) => {
  await page.goto('access/identity/groups/aliases');
  await expect(page.getByRole('link', { name: ALIAS })).toBeVisible();

  await test.step('header create action is hidden', async () => {
    await expect(page.getByRole('link', { name: 'Create group' })).toHaveCount(0);
  });

  await test.step('alias details hide the edit action', async () => {
    await page.getByRole('link', { name: ALIAS }).click();
    await expect(page.getByRole('heading', { name: ALIAS })).toBeVisible();
    await expect(page.getByRole('link', { name: /^Edit/ })).toHaveCount(0);
  });
});
