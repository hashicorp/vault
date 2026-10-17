/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect, type Page } from '@playwright/test';
import { test } from '../../fixtures/demo';

const readToken = (page: Page) =>
  page.evaluate(() => {
    const key = Object.keys(localStorage).find((key) => key.startsWith('vault-token'));
    return key ? (JSON.parse(localStorage.getItem(key) as string).token as string) : '';
  });

const readOidcData = async (page: Page, type: 'assignment' | 'key', name: string) => {
  const response = await page.request.get(`/v1/identity/oidc/${type}/${name}`, {
    headers: { 'X-Vault-Token': await readToken(page) },
  });
  await expect(response).toBeOK();
  const { data } = await response.json();
  return data;
};

const readOidcNames = async (page: Page, type: 'assignment' | 'key') => {
  const response = await page.request.get(`/v1/identity/oidc/${type}?list=true`, {
    headers: { 'X-Vault-Token': await readToken(page) },
  });
  await expect(response).toBeOK();
  const { data } = await response.json();
  return (data?.keys ?? []) as string[];
};

test('oidc workflow', async ({ page }) => {
  await page.goto('dashboard');

  await test.step('navigate to OIDC provider page', async () => {
    await page.getByRole('link', { name: 'Access control' }).click();
    await page.getByRole('link', { name: 'OIDC provider' }).click();
    await expect(page.getByRole('heading', { name: 'OIDC provider', level: 1 })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Create your first app' })).toBeVisible();
  });

  await test.step('create application', async () => {
    await page.getByRole('link', { name: 'Create your first app' }).click();
    await page.getByRole('textbox', { name: 'Application name' }).fill('test-oidc-app');
    await page.getByRole('button', { name: 'More options' }).click();
    await page
      .getByRole('group', { name: 'ID Token TTL Lease will' })
      .getByLabel('Number of units')
      .fill('30');
    await page.getByLabel('TTL unit for ID Token TTL').selectOption('m');
    await page
      .getByRole('group', { name: 'Access Token TTL Lease will' })
      .getByLabel('Number of units')
      .fill('30');
    await page.getByLabel('TTL unit for Access Token TTL').selectOption('m');
    await page.getByRole('button', { name: 'Create' }).click();

    await expect(page.getByRole('heading', { name: 'test-oidc-app' })).toBeVisible();
    await expect(page.getByText('ID Token TTL 30 minutes')).toBeVisible();
    await expect(page.getByText('Access Token TTL 30 minutes')).toBeVisible();
    await page.getByRole('link', { name: 'Available providers' }).click();
    await expect(page.getByRole('link', { name: 'default Issuer: /v1/identity/' })).toBeVisible();
    await page.getByRole('link', { name: 'Applications' }).click();
    await expect(page.getByRole('link', { name: 'test-oidc-app Client ID:' })).toBeVisible();
  });

  await test.step('create key', async () => {
    await page.getByRole('link', { name: 'Keys' }).click();
    await expect(page.getByRole('link', { name: 'default Key nav options' })).toBeVisible();
    await page.getByRole('link', { name: 'Create key' }).click();
    await page.getByRole('textbox', { name: 'Name' }).fill('test-oidc-key');
    await page.getByLabel('Algorithm').selectOption('ES256');
    await page
      .getByRole('group', { name: 'Rotation period Lease will' })
      .getByLabel('Number of units')
      .fill('30');
    await page.getByLabel('TTL unit for Rotation period').selectOption('m');
    await page
      .getByRole('group', { name: 'Verification TTL Lease will' })
      .getByLabel('Number of units')
      .fill('30');
    await page.getByLabel('TTL unit for Verification TTL').selectOption('m');
    await expect(page.locator('.radio-card').nth(1)).toHaveClass(/is-disabled/);
    await page.getByRole('button', { name: 'Create' }).click();

    await expect(page.getByRole('heading', { name: 'test-oidc-key' })).toBeVisible();
    await expect(page.getByText('Algorithm ES256')).toBeVisible();
    await expect(page.getByText('Rotation period 30 minutes')).toBeVisible();
    await expect(page.getByText('Verification TTL 30 minutes')).toBeVisible();
    const key = await readOidcData(page, 'key', 'test-oidc-key');
    expect(key.algorithm).toBe('ES256');
    expect(key.rotation_period).toBe(1800);
    expect(key.verification_ttl).toBe(1800);
    await page.reload();
    await expect(page.getByText('Algorithm ES256')).toBeVisible();
  });

  await test.step('rotate key', async () => {
    const rotateResponse = page.waitForResponse(
      (response) =>
        response.url().endsWith('/v1/identity/oidc/key/test-oidc-key/rotate') &&
        response.request().method() === 'POST'
    );
    await page.getByRole('button', { name: 'Rotate key' }).click();
    await page.getByRole('button', { name: 'Confirm' }).click();
    expect((await rotateResponse).ok()).toBe(true);
    await expect(page.getByText('Success: test-oidc-key connection was rotated.')).toBeVisible();
  });

  await test.step('edit key', async () => {
    await page.getByRole('link', { name: 'Edit key' }).click();
    await page.getByLabel('Algorithm').selectOption('EdDSA');
    await page.getByRole('button', { name: 'Update' }).click();
    await expect(page.getByText('Algorithm EdDSA')).toBeVisible();
    expect((await readOidcData(page, 'key', 'test-oidc-key')).algorithm).toBe('EdDSA');
    await page.reload();
    await expect(page.getByText('Algorithm EdDSA')).toBeVisible();
    await page.getByRole('link', { name: 'Keys' }).click();
  });

  await test.step('delete key', async () => {
    await page.getByRole('link', { name: 'test-oidc-key Key nav options' }).click();
    const deleteResponse = page.waitForResponse(
      (response) =>
        response.url().endsWith('/v1/identity/oidc/key/test-oidc-key') &&
        response.request().method() === 'DELETE'
    );
    await page.getByRole('button', { name: 'Delete key' }).click();
    await page.getByRole('button', { name: 'Confirm' }).click();
    expect((await deleteResponse).ok()).toBe(true);
    await expect
      .poll(() => readOidcNames(page, 'key'), { message: 'the deleted key is absent from the backend list' })
      .not.toContain('test-oidc-key');
    await page.reload();
    await expect(page.getByRole('link', { name: 'test-oidc-key Key nav options' })).toHaveCount(0);
  });

  await test.step('create assignment', async () => {
    // create a group and entity for the assignment
    await page.getByRole('link', { name: 'Groups' }).click();
    await page.getByRole('link', { name: 'Create group' }).click();
    await page.getByRole('textbox', { name: 'Name' }).fill('oidc-group');
    await page.getByRole('button', { name: 'Create' }).click();
    await page.getByRole('link', { name: 'Entities' }).click();
    await page.getByRole('link', { name: 'Create new entity' }).click();
    await page.getByRole('textbox', { name: 'Name' }).fill('oidc-entity');
    await page.getByRole('button', { name: 'Create' }).click();

    await page.getByRole('link', { name: 'OIDC provider' }).click();
    await page.getByRole('link', { name: 'Assignments' }).click();
    await expect(page.locator('.list-item-row')).toHaveClass(/is-disabled/);

    await page.getByRole('link', { name: 'Create assignment' }).click();
    await page.getByRole('textbox', { name: 'Name' }).fill('oidc-assignment');
    await page.getByLabel('Entities').getByText('Search').click();
    // close the dropdown before submitting - once enough entities exist the open option list
    // overlays the Create button and intercepts the click
    await page.keyboard.press('Escape');
    await page.getByRole('button', { name: 'Create' }).click();
    await expect(page.getByText('At least one entity or group')).toBeVisible();
    await page.getByLabel('Groups').getByText('Search').click();
    await page.getByRole('option', { name: 'oidc-group' }).click();
    await page.getByRole('button', { name: 'Create' }).click();
    const assignment = await readOidcData(page, 'assignment', 'oidc-assignment');
    expect(assignment.group_ids).toHaveLength(1);
    expect(assignment.entity_ids).toHaveLength(0);
    await expect(page.getByRole('link', { name: 'oidc-group' })).toBeVisible();
  });

  await test.step('edit assignment', async () => {
    await page.getByRole('link', { name: 'Edit assignment' }).click();
    await page.getByLabel('Entities').getByText('Search').click();
    await page.getByRole('option', { name: 'oidc-entity' }).click();
    await page.getByRole('button', { name: 'Update' }).click();
    await expect(page.getByRole('link', { name: 'oidc-entity' })).toBeVisible();
    const assignment = await readOidcData(page, 'assignment', 'oidc-assignment');
    expect(assignment.entity_ids).toHaveLength(1);
    expect(assignment.group_ids).toHaveLength(1);
    await page.reload();
    await expect(page.getByRole('link', { name: 'oidc-entity' })).toBeVisible();
    await page.getByRole('link', { name: 'Assignments' }).click();
    await expect(page.getByRole('link', { name: 'oidc-assignment' })).toBeVisible();
  });

  await test.step('delete assignment', async () => {
    await page.getByRole('link', { name: 'oidc-assignment' }).click();
    const deleteResponse = page.waitForResponse(
      (response) =>
        response.url().endsWith('/v1/identity/oidc/assignment/oidc-assignment') &&
        response.request().method() === 'DELETE'
    );
    await page.getByRole('button', { name: 'Delete assignment' }).click();
    await page.getByRole('button', { name: 'Confirm' }).click();
    expect((await deleteResponse).ok()).toBe(true);
    await expect(page).toHaveURL(/\/access\/oidc\/assignments\/?$/);
    await expect
      .poll(() => readOidcNames(page, 'assignment'), {
        message: 'the deleted assignment is absent from the backend list',
      })
      .not.toContain('oidc-assignment');
    await page.reload();
    await expect(page.getByRole('link', { name: 'oidc-assignment' })).toHaveCount(0);
  });

  await test.step('create provider', async () => {
    await page.getByRole('link', { name: 'Providers' }).click();
    await page.getByRole('link', { name: 'Create provider' }).click();
    await page.getByRole('textbox', { name: 'Name' }).fill('oidc-provider');
    await page.getByRole('radio', { name: 'Limit access to selected' }).check();
    await page.getByLabel('Application name').getByText('Search').click();
    await page.getByRole('option', { name: 'test-oidc-app' }).click();
    await page.getByRole('button', { name: 'Create' }).click();
    await expect(page.getByText('/v1/identity/oidc/provider/')).toBeVisible();
  });

  await test.step('create scope', async () => {
    await page.getByRole('link', { name: 'Providers' }).click();
    await page.getByRole('link', { name: 'Scopes' }).click();
    await expect(page.getByRole('heading', { name: 'No scopes yet' })).toBeVisible();
    await page.getByRole('link', { name: 'Create scope' }).click();
    await page.getByRole('textbox', { name: 'Name' }).fill('oidc-scope');
    await page.getByRole('textbox', { name: 'Description' }).fill('oidc scope description');
    await page.getByRole('textbox', { name: 'JSON Template' }).fill(`{
      "username": {{identity.entity.aliases.$MOUNT_ACCESSOR.name}},
      "contact": {
        "email": {{identity.entity.metadata.email}},
        "phone_number": {{identity.entity.metadata.phone_number}}
      },
      "groups": {{identity.entity.groups.names}}
    }`);
    await page.getByRole('button', { name: 'Create' }).click();
  });

  await test.step('edit scope', async () => {
    await page.getByRole('link', { name: 'Edit scope' }).click();
    await page.getByRole('textbox', { name: 'Description' }).fill('updated description');
    await page.getByRole('textbox', { name: 'JSON Template' }).fill(`{
      "username": {{identity.entity.aliases.$MOUNT_ACCESSOR.name}},
      "contact": {
        "email": {{identity.entity.metadata.email}}
      },
      "groups": {{identity.entity.groups.names}}
    }`);

    await page.getByRole('button', { name: 'Update' }).click();
    await expect(page.getByText('updated description')).toBeVisible();
    await expect(page.getByText('"phone_number"')).not.toBeVisible();
    await page.getByRole('link', { name: 'Scopes' }).click();
    await expect(page.getByRole('link', { name: 'oidc-scope' })).toBeVisible();
  });
});
