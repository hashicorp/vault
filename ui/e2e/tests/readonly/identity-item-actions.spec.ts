/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { test } from '../../fixtures/demo';
import { rootApi } from '../../fixtures/root-api';

const ENTITY = 'readonly-actions-entity';
const MEMBER = 'readonly-actions-member';
const GROUP = 'readonly-actions-group';

// Seeding is idempotent because Playwright re-runs beforeAll in a fresh worker after a failure.
test.beforeAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  const upsert = async (type: string, name: string, data = {}) => {
    await expect(await api.post(`/v1/identity/${type}/name/${name}`, { data })).toBeOK();
    return (await (await api.get(`/v1/identity/${type}/name/${name}`)).json()).data;
  };
  const member = await upsert('entity', MEMBER);
  await upsert('entity', ENTITY, {
    disabled: true,
    policies: ['default'],
    metadata: { team: 'ui' },
  });
  await upsert('group', GROUP, {
    policies: ['default'],
    metadata: { team: 'ui' },
    member_entity_ids: [member.id],
  });
  await api.dispose();
});

// Read-only tokens may view policies, metadata and members, but never see the actions that change them.
test('read-only user cannot change entity policies, metadata or enabled state', async ({
  page,
  playwright,
}, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  const { data: entity } = await (await api.get(`/v1/identity/entity/name/${ENTITY}`)).json();
  await api.dispose();

  await test.step('details tab hides the Enable action', async () => {
    await page.goto(`access/identity/entities/${entity.id}/details`);
    await expect(page.getByText('is disabled. All associated tokens cannot be used')).toBeVisible();
    await expect(page.getByRole('button', { name: 'Enable' })).toHaveCount(0);
  });

  await test.step('policies tab hides the policy management menu', async () => {
    await page.goto(`access/identity/entities/${entity.id}/policies`);
    await expect(page.getByText('default', { exact: true })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Identity policy management options' })).toHaveCount(0);
  });

  await test.step('metadata tab hides the metadata menu', async () => {
    await page.goto(`access/identity/entities/${entity.id}/metadata`);
    await expect(page.getByText('team', { exact: true })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Metadata options' })).toHaveCount(0);
  });
});

test('read-only user cannot change group policies, metadata or members', async ({
  page,
  playwright,
}, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  const { data: group } = await (await api.get(`/v1/identity/group/name/${GROUP}`)).json();
  await api.dispose();

  await test.step('policies tab hides the policy management menu', async () => {
    await page.goto(`access/identity/groups/${group.id}/policies`);
    await expect(page.getByText('default', { exact: true })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Identity policy management options' })).toHaveCount(0);
  });

  await test.step('metadata tab hides the metadata menu', async () => {
    await page.goto(`access/identity/groups/${group.id}/metadata`);
    await expect(page.getByText('team', { exact: true })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Metadata options' })).toHaveCount(0);
  });

  await test.step('members tab hides the member menu', async () => {
    await page.goto(`access/identity/groups/${group.id}/members`);
    await expect(page.getByText(group.member_entity_ids[0], { exact: true })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Identity member options' })).toHaveCount(0);
  });
});
