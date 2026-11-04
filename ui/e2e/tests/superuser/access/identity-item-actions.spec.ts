/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { test } from '../../../fixtures/demo';
import { rootApi } from '../../../fixtures/root-api';

import type { APIRequestContext, Page } from '@playwright/test';

// Each test seeds its own entity/group so the tests share no state. Writing by name is an upsert,
// which keeps seeding idempotent when Playwright re-runs a test in a fresh worker.
const seed = async (api: APIRequestContext, type: 'entity' | 'group', name: string, data = {}) => {
  await expect(await api.post(`/v1/identity/${type}/name/${name}`, { data })).toBeOK();
  return read(api, type, name);
};

const read = async (api: APIRequestContext, type: 'entity' | 'group', name: string) => {
  const res = await api.get(`/v1/identity/${type}/name/${name}`);
  await expect(res).toBeOK();
  return (await res.json()).data;
};

// Opens the row action menu, picks the named item and confirms the modal.
const removeFromRow = async (page: Page, rowText: string, menuName: string, itemName: string) => {
  const row = page.locator('.list-item-row, .info-table-row', { hasText: rowText });
  await row.getByRole('button', { name: menuName }).click();
  await page.getByRole('button', { name: itemName }).click();
  await page.getByRole('button', { name: 'Confirm' }).click();
};

test.describe('identity item actions', () => {
  let api: APIRequestContext;

  test.beforeAll(async ({ playwright }, testInfo) => {
    api = await rootApi(playwright, testInfo);
  });

  test.afterAll(async () => {
    await api.dispose();
  });

  test('removing a policy from an entity persists and keeps other fields', async ({ page }) => {
    const entity = await seed(api, 'entity', 'actions-entity-policy', {
      policies: ['default', 'remove-me'],
      metadata: { team: 'ui' },
    });

    await page.goto(`access/identity/entities/${entity.id}/policies`);
    await expect(page.getByText('remove-me', { exact: true })).toBeVisible();

    await removeFromRow(page, 'remove-me', 'Identity policy management options', 'Remove from entity');
    await expect(page.getByText(`Successfully removed 'remove-me' policy from ${entity.id}`)).toBeVisible();
    await expect(page.getByText('remove-me', { exact: true })).toHaveCount(0);

    await page.reload();
    await expect(page.getByText('default', { exact: true })).toBeVisible();
    await expect(page.getByText('remove-me', { exact: true })).toHaveCount(0);

    const after = await read(api, 'entity', 'actions-entity-policy');
    expect(after.policies).toEqual(['default']);
    expect(after.metadata).toEqual({ team: 'ui' });
    expect(after.disabled).toBe(false);
  });

  test('removing a metadata key from an entity persists and keeps other fields', async ({ page }) => {
    const entity = await seed(api, 'entity', 'actions-entity-metadata', {
      policies: ['default'],
      metadata: { keep: 'yes', drop: 'no' },
    });

    await page.goto(`access/identity/entities/${entity.id}/metadata`);
    await expect(page.getByText('drop', { exact: true })).toBeVisible();

    await removeFromRow(page, 'drop', 'Metadata options', 'Remove');
    await expect(page.getByText("Successfully removed 'drop' from metadata")).toBeVisible();

    await page.reload();
    await expect(page.getByText('keep', { exact: true })).toBeVisible();
    await expect(page.getByText('drop', { exact: true })).toHaveCount(0);

    const after = await read(api, 'entity', 'actions-entity-metadata');
    expect(after.metadata).toEqual({ keep: 'yes' });
    expect(after.policies).toEqual(['default']);
  });

  test('enabling a disabled entity persists and keeps other fields', async ({ page }) => {
    const entity = await seed(api, 'entity', 'actions-entity-enable', {
      disabled: true,
      policies: ['default'],
      metadata: { team: 'ui' },
    });

    await page.goto(`access/identity/entities/${entity.id}/details`);
    await expect(page.getByText('is disabled. All associated tokens cannot be used')).toBeVisible();

    await page.getByRole('button', { name: 'Enable' }).click();
    await expect(page.getByText(`Successfully enabled entity: ${entity.id}`)).toBeVisible();
    await expect(page.getByText('is disabled. All associated tokens cannot be used')).toHaveCount(0);

    await page.reload();
    await expect(page.getByText('is disabled. All associated tokens cannot be used')).toHaveCount(0);

    const after = await read(api, 'entity', 'actions-entity-enable');
    expect(after.disabled).toBe(false);
    expect(after.policies).toEqual(['default']);
    expect(after.metadata).toEqual({ team: 'ui' });
  });

  test('removing a policy from a group persists and keeps other fields', async ({ page }) => {
    const group = await seed(api, 'group', 'actions-group-policy', {
      policies: ['default', 'remove-me'],
      metadata: { team: 'ui' },
    });

    await page.goto(`access/identity/groups/${group.id}/policies`);
    await expect(page.getByText('remove-me', { exact: true })).toBeVisible();

    await removeFromRow(page, 'remove-me', 'Identity policy management options', 'Remove from group');
    await expect(page.getByText(`Successfully removed 'remove-me' policy from ${group.id}`)).toBeVisible();

    await page.reload();
    await expect(page.getByText('default', { exact: true })).toBeVisible();
    await expect(page.getByText('remove-me', { exact: true })).toHaveCount(0);

    const after = await read(api, 'group', 'actions-group-policy');
    expect(after.policies).toEqual(['default']);
    expect(after.metadata).toEqual({ team: 'ui' });
  });

  test('removing a metadata key from a group persists and keeps other fields', async ({ page }) => {
    const group = await seed(api, 'group', 'actions-group-metadata', {
      policies: ['default'],
      metadata: { keep: 'yes', drop: 'no' },
    });

    await page.goto(`access/identity/groups/${group.id}/metadata`);
    await expect(page.getByText('drop', { exact: true })).toBeVisible();

    await removeFromRow(page, 'drop', 'Metadata options', 'Remove');
    await expect(page.getByText("Successfully removed 'drop' from metadata")).toBeVisible();

    await page.reload();
    await expect(page.getByText('keep', { exact: true })).toBeVisible();
    await expect(page.getByText('drop', { exact: true })).toHaveCount(0);

    const after = await read(api, 'group', 'actions-group-metadata');
    expect(after.metadata).toEqual({ keep: 'yes' });
    expect(after.policies).toEqual(['default']);
  });

  test('removing a member from a group persists and keeps other members and fields', async ({ page }) => {
    const keep = await seed(api, 'entity', 'actions-member-keep');
    const drop = await seed(api, 'entity', 'actions-member-drop');
    const group = await seed(api, 'group', 'actions-group-members', {
      policies: ['default'],
      member_entity_ids: [keep.id, drop.id],
    });

    await page.goto(`access/identity/groups/${group.id}/members`);
    await expect(page.getByText(drop.id, { exact: true })).toBeVisible();

    await removeFromRow(page, drop.id, 'Identity member options', 'Remove');
    await expect(page.getByText(`Successfully removed '${drop.id}' from the group`)).toBeVisible();

    await page.reload();
    await expect(page.getByText(keep.id, { exact: true })).toBeVisible();
    await expect(page.getByText(drop.id, { exact: true })).toHaveCount(0);

    const after = await read(api, 'group', 'actions-group-members');
    expect(after.member_entity_ids).toEqual([keep.id]);
    expect(after.policies).toEqual(['default']);
  });
});
