/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect, test } from '@playwright/test';
import { BasePage } from '../../pages/base';

test('kmip workflow', async ({ page }) => {
  const basePage = new BasePage(page);

  await test.step('enable KMIP secrets engine', async () => {
    await basePage.enableEngine('KMIP', 'kmip-test');
  });

  await test.step('KMIP secrets engine mount saved successfully', async () => {
    await expect(page.getByText('Success', { exact: true })).toBeVisible();
    await expect(page.getByText('Successfully mounted the kmip secrets engine at kmip-test')).toBeVisible();
    await page.getByRole('button', { name: 'Dismiss' }).click();
  });

  await test.step('configure KMIP engine', async () => {
    await page.getByRole('button', { name: 'Manage' }).click();
    await page.getByRole('link', { name: 'Configure' }).click();

    await page.locator('label').filter({ hasText: 'Default TLS Client TTL Lease' }).click();
    await page.getByLabel('TLS CA Key type').selectOption('rsa');
    await page.getByRole('textbox', { name: 'TLS CA Key bits' }).fill('2048');
    await page.getByRole('button', { name: 'Save' }).click();

    await expect(page.getByText('Success', { exact: true })).toBeVisible();
    await expect(page.getByText('Successfully configured KMIP engine')).toBeVisible();
    await page.getByRole('button', { name: 'Dismiss' }).click();

    await expect(page.locator('[data-test-row-value="Default TLS client TTL"]')).toContainText('0');
    await expect(page.locator('span[data-test-row-value="TLS CA key type"]')).toContainText('rsa');
    await expect(page.locator('span[data-test-row-value="TLS CA key bits"]')).toContainText('2048');
  });

  await test.step('update general settings', async () => {
    await page.getByRole('link', { name: 'General settings' }).click();
    await page.getByText('Engine type kmip').click();
    await page.getByRole('textbox', { name: 'Description' }).click();
    await page.getByRole('textbox', { name: 'Description' }).press('ControlOrMeta+a');
    await page.getByRole('textbox', { name: 'Description' }).fill('abcdefg');
    await page.getByRole('button', { name: 'Save changes' }).click();

    await expect(page.getByText('Configuration saved')).toBeVisible();
    await expect(page.getByText('Engine settings successfully updated.')).toBeVisible();
    await page.getByRole('button', { name: 'Dismiss' }).click();
  });

  await test.step('create a scope', async () => {
    await page.getByRole('link', { name: 'Exit configuration' }).click();
    await expect(page.locator('section')).toContainText('KMIP Secrets Engine');
    await expect(page.locator('section')).toContainText(
      "First, let's create a scope that our roles and credentials will belong to. A client can only access objects within their role's scope."
    );
    await page.getByRole('link', { name: 'Create scope' }).click();
    await page.getByRole('textbox', { name: 'Name' }).click();
    await page.getByRole('textbox', { name: 'Name' }).fill('kmip-scope');
    await page.getByRole('button', { name: 'Save' }).click();

    await expect(page.getByText('Success', { exact: true })).toBeVisible();
    await expect(page.getByText('Successfully created scope kmip-scope')).toBeVisible();
    await page.getByRole('button', { name: 'Dismiss' }).click();
  });

  await test.step('create a role', async () => {
    await page.getByRole('button', { name: 'More options' }).click();
    await page.getByRole('link', { name: 'View scope', exact: true }).click();

    await page.locator('div').filter({ hasText: 'No roles in this scope yet' }).nth(3).click();
    await page.getByRole('link', { name: 'Create role' }).click();
    await page.getByRole('textbox', { name: 'Name', exact: true }).click();
    await page.getByRole('textbox', { name: 'Name', exact: true }).fill('kmip-role');
    await page.getByRole('button', { name: 'Save' }).click();

    await expect(page.getByText('Success', { exact: true })).toBeVisible();
    await expect(page.getByText('Successfully saved role kmip-role')).toBeVisible();
    await page.getByRole('button', { name: 'Dismiss' }).click();
  });

  await test.step('generate credentials', async () => {
    await page.getByRole('button', { name: 'More options' }).click();
    await page.getByRole('link', { name: 'View credentials', exact: true }).click();
    await expect(page.getByText('No credentials yet for this')).toBeVisible();
    await expect(page.getByText('You can generate new')).toBeVisible();
    await page.getByLabel('toolbar actions').getByRole('link', { name: 'Generate credentials' }).click();
    await page.getByLabel('Certificate format').selectOption('pem_bundle');
    await page.getByRole('button', { name: 'Save' }).click();

    await expect(page.getByText('Success', { exact: true })).toBeVisible();
    await expect(page.getByText('Successfully generated credentials from role kmip-role.')).toBeVisible();
    await page.getByRole('button', { name: 'Dismiss' }).click();
  });

  await test.step('revoke credentials', async () => {
    await page.getByRole('button', { name: 'Revoke credentials' }).click();
    await page.getByRole('button', { name: 'Confirm' }).click();

    await expect(page.getByText('Success', { exact: true })).toBeVisible();
    await expect(page.getByText('Successfully revoked credentials.')).toBeVisible();
    await page.getByRole('button', { name: 'Dismiss' }).click();

    await expect(page.getByText('No credentials yet for this')).toBeVisible();
  });

  await test.step('delete role', async () => {
    await page.getByRole('link', { name: 'kmip-scope' }).click();
    await page.getByRole('button', { name: 'More options' }).click();
    await page.getByRole('button', { name: 'Delete role' }).click();
    await page.getByRole('button', { name: 'Confirm' }).click();

    await expect(page.getByText('No roles in this scope yet')).toBeVisible();
    await expect(page.getByText('Roles let you generate')).toBeVisible();
  });

  await test.step('delete scope', async () => {
    await page.getByRole('link', { name: 'kmip-test' }).click();
    await page.getByRole('button', { name: 'More options' }).click();
    await page.getByRole('button', { name: 'Delete scope' }).click();
    await page.getByRole('button', { name: 'Confirm' }).click();

    await expect(page.getByText('KMIP Secrets Engine')).toBeVisible();
    await expect(page.getByText("First, let's create a scope")).toBeVisible();
  });
});
