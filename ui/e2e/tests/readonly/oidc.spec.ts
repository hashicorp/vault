/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect, type Page } from '@playwright/test';
import { test } from '../../fixtures/demo';
import { rootApi } from '../../fixtures/root-api';

const KEY_NAME = 'readonly-oidc-key';
const ASSIGNMENT_NAME = 'readonly-oidc-assignment';
const ENTITY_NAME = 'readonly-oidc-entity';
let entityId = '';

const readToken = (page: Page) =>
  page.evaluate(() => {
    const key = Object.keys(localStorage).find((key) => key.startsWith('vault-token'));
    return key ? (JSON.parse(localStorage.getItem(key) as string).token as string) : '';
  });

const readOidcKeys = async (page: Page) => {
  const response = await page.request.get('/v1/identity/oidc/key?list=true', {
    headers: { 'X-Vault-Token': await readToken(page) },
  });
  await expect(response).toBeOK();
  const { data } = await response.json();
  return (data?.keys ?? []) as string[];
};

test.beforeAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  try {
    const existingKey = await api.get(`/v1/identity/oidc/key/${KEY_NAME}`);
    if (!existingKey.ok()) {
      await expect(
        await api.post(`/v1/identity/oidc/key/${KEY_NAME}`, {
          data: {
            algorithm: 'RS256',
            rotation_period: '24h',
            verification_ttl: '24h',
            allowed_client_ids: ['*'],
          },
        })
      ).toBeOK();
    }

    const existingEntity = await api.get(`/v1/identity/entity/name/${ENTITY_NAME}`);
    if (existingEntity.ok()) {
      entityId = (await existingEntity.json()).data.id;
    } else {
      const response = await api.post('/v1/identity/entity', { data: { name: ENTITY_NAME } });
      await expect(response).toBeOK();
      entityId = (await response.json()).data.id;
    }

    const existingAssignment = await api.get(`/v1/identity/oidc/assignment/${ASSIGNMENT_NAME}`);
    if (!existingAssignment.ok()) {
      await expect(
        await api.post(`/v1/identity/oidc/assignment/${ASSIGNMENT_NAME}`, {
          data: { entity_ids: [entityId] },
        })
      ).toBeOK();
    }
  } finally {
    await api.dispose();
  }
});

test.afterAll(async ({ playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  try {
    const existingAssignment = await api.get(`/v1/identity/oidc/assignment/${ASSIGNMENT_NAME}`);
    if (existingAssignment.ok()) {
      await expect(await api.delete(`/v1/identity/oidc/assignment/${ASSIGNMENT_NAME}`)).toBeOK();
    }
    const existingKey = await api.get(`/v1/identity/oidc/key/${KEY_NAME}`);
    if (existingKey.ok()) await expect(await api.delete(`/v1/identity/oidc/key/${KEY_NAME}`)).toBeOK();
    if (entityId) await expect(await api.delete(`/v1/identity/entity/id/${entityId}`)).toBeOK();
  } finally {
    await api.dispose();
  }
});

test('read-only user can view an OIDC key but cannot mutate it', async ({ page }) => {
  await page.goto(`access/oidc/keys/${KEY_NAME}`);
  await page.getByRole('link', { name: 'Details' }).click();
  await expect(page.getByRole('heading', { name: KEY_NAME })).toBeVisible();
  await expect(page.getByText('Algorithm RS256')).toBeVisible();
  await expect(page.getByRole('button', { name: 'Delete key' })).toHaveCount(0);
  await expect(page.getByRole('button', { name: 'Rotate key' })).toHaveCount(0);
  await expect(page.getByRole('link', { name: 'Edit key' })).toHaveCount(0);

  const deniedWrite = await page.request.post(`/v1/identity/oidc/key/${KEY_NAME}`, {
    headers: { 'X-Vault-Token': await readToken(page) },
    data: { algorithm: 'EdDSA' },
  });
  expect(deniedWrite.status()).toBe(403);
  const deniedRotation = await page.request.post(`/v1/identity/oidc/key/${KEY_NAME}/rotate`, {
    headers: { 'X-Vault-Token': await readToken(page) },
    data: { verification_ttl: '24h' },
  });
  expect(deniedRotation.status()).toBe(403);
  const deniedDelete = await page.request.delete(`/v1/identity/oidc/key/${KEY_NAME}`, {
    headers: { 'X-Vault-Token': await readToken(page) },
  });
  expect(deniedDelete.status()).toBe(403);
  await expect
    .poll(() => readOidcKeys(page), { message: 'the denied update preserves the existing key' })
    .toContain(KEY_NAME);
});

test('read-only user can view an OIDC assignment but cannot mutate it', async ({ page }) => {
  await page.goto('access/oidc/assignments');
  await page.getByRole('link', { name: ASSIGNMENT_NAME }).click();
  await expect(page.getByRole('heading', { name: ASSIGNMENT_NAME })).toBeVisible();
  await expect(page.getByRole('link', { name: ENTITY_NAME })).toBeVisible();
  await expect(page.getByRole('link', { name: 'Edit assignment' })).toHaveCount(0);
  await expect(page.getByRole('button', { name: 'Delete assignment' })).toHaveCount(0);

  const assignmentBefore = await page.request.get(`/v1/identity/oidc/assignment/${ASSIGNMENT_NAME}`, {
    headers: { 'X-Vault-Token': await readToken(page) },
  });
  await expect(assignmentBefore).toBeOK();
  const originalData = (await assignmentBefore.json()).data;

  const deniedWrite = await page.request.post(`/v1/identity/oidc/assignment/${ASSIGNMENT_NAME}`, {
    headers: { 'X-Vault-Token': await readToken(page) },
    data: { entity_ids: [] },
  });
  expect(deniedWrite.status()).toBe(403);
  const deniedDelete = await page.request.delete(`/v1/identity/oidc/assignment/${ASSIGNMENT_NAME}`, {
    headers: { 'X-Vault-Token': await readToken(page) },
  });
  expect(deniedDelete.status()).toBe(403);

  const assignmentAfter = await page.request.get(`/v1/identity/oidc/assignment/${ASSIGNMENT_NAME}`, {
    headers: { 'X-Vault-Token': await readToken(page) },
  });
  await expect(assignmentAfter).toBeOK();
  expect((await assignmentAfter.json()).data.entity_ids).toEqual(originalData.entity_ids);
});
