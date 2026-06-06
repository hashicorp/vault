/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect, test } from '@playwright/test';

test('oidc workflow', async ({ page }) => {
  await page.goto('dashboard');

  await test.step('navigate to OIDC provider page', async () => {
    await page.getByRole('link', { name: 'Access control' }).click();
    await page.getByRole('link', { name: 'OIDC provider' }).click();
    await expect(page.getByRole('img', { name: 'Example flow of a user' })).toBeVisible();
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
  });

  await test.step('rotate key', async () => {
    await page.getByRole('button', { name: 'Rotate key' }).click();
    await page.getByRole('button', { name: 'Confirm' }).click();
  });

  await test.step('edit key', async () => {
    await page.getByRole('link', { name: 'Edit key' }).click();
    await page.getByLabel('Algorithm').selectOption('EdDSA');
    await page.getByRole('button', { name: 'Update' }).click();
    await expect(page.getByText('Algorithm EdDSA')).toBeVisible();
    await page.getByRole('link', { name: 'Keys' }).click();
  });

  await test.step('create assignment', async () => {
    // create a group and entity for the assignment
    await page.getByRole('link', { name: 'Groups' }).click();
    await page.getByRole('link', { name: 'Create group' }).click();
    await page.getByRole('textbox', { name: 'Name' }).fill('oidc-group');
    await page.getByRole('button', { name: 'Create' }).click();
    await page.getByRole('link', { name: 'Entities' }).click();
    await page.getByRole('link', { name: 'Create entity' }).click();
    await page.getByRole('textbox', { name: 'Name' }).fill('oidc-entity');
    await page.getByRole('button', { name: 'Create' }).click();

    await page.getByRole('link', { name: 'OIDC provider' }).click();
    await page.getByRole('link', { name: 'Assignments' }).click();
    await expect(page.locator('.list-item-row')).toHaveClass(/is-disabled/);

    await page.getByRole('link', { name: 'Create assignment' }).click();
    await page.getByRole('textbox', { name: 'Name' }).fill('oidc-assignment');
    await page.getByLabel('Entities').getByText('Search').click();
    await page.getByRole('button', { name: 'Create' }).click();
    await expect(page.getByText('At least one entity or group')).toBeVisible();
    await page.getByLabel('Groups').getByText('Search').click();
    await page.getByRole('option', { name: 'oidc-group' }).click();
    await page.getByRole('button', { name: 'Create' }).click();
    await expect(page.getByRole('link', { name: 'oidc-group' })).toBeVisible();
  });

  await test.step('edit assignment', async () => {
    await page.getByRole('link', { name: 'Edit assignment' }).click();
    await page.getByLabel('Entities').getByText('Search').click();
    await page.getByRole('option', { name: 'oidc-entity' }).click();
    await page.getByRole('button', { name: 'Update' }).click();
    await expect(page.getByRole('link', { name: 'oidc-entity' })).toBeVisible();
    await page.getByRole('link', { name: 'Assignments' }).click();
    await expect(page.getByRole('link', { name: 'oidc-assignment' })).toBeVisible();
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
