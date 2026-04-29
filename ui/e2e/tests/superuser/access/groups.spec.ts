/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';

test('groups workflow', async ({ page }) => {
  await test.step('should display groups page', async () => {
    await page.goto('dashboard');
    await page.getByRole('link', { name: 'Access control' }).click();
    await page.getByRole('link', { name: 'Groups' }).click();
    await expect(page.getByRole('heading', { name: 'Groups', exact: true })).toContainText('Groups');
    await expect(page.getByRole('heading', { name: 'No groups yet' })).toContainText('No groups yet');
    await expect(page.getByText('A list of groups in this')).toContainText(
      'A list of groups in this namespace will be listed here. Create your first group to get started.'
    );
  });

  await test.step('create policy', async () => {
    await page.getByRole('link', { name: 'ACL Policies' }).click();
    await page.getByRole('link', { name: 'Create ACL policy' }).click();
    await page.getByRole('textbox', { name: 'Policy name' }).click();
    await page.getByRole('textbox', { name: 'Policy name' }).fill('test-policy');
    await page.getByRole('radio', { name: 'Code editor' }).check();
    await page
      .getByRole('textbox', { name: 'Policy editor' })
      .fill('path "auth/token/lookup-self" { capabilities = ["read"]}');
    await page.getByRole('button', { name: 'Create policy' }).click();
    await page.getByRole('link', { name: 'Groups' }).click();
  });

  await test.step('create group', async () => {
    await page.getByRole('link', { name: 'Create group' }).click();
    await page.getByRole('textbox', { name: 'Name' }).fill('group-1');
    await page.getByRole('textbox', { name: 'Key' }).fill('hello');
    await page.getByRole('textbox', { name: 'Value' }).fill('world');
    await page.locator('.ember-basic-dropdown-trigger').first().click();
    await page.locator('.ember-power-select-option', { hasText: 'test-policy' }).click();
    await page.getByRole('button', { name: 'Create' }).click();
  });

  await test.step('should display group detail page', async () => {
    await expect(page.getByRole('heading', { name: 'group-' })).toContainText('group-1');
    await page.locator('div:nth-child(2) > .column.is-flex-center').click();
    await expect(page.getByText('Name group-')).toBeVisible();
  });

  await test.step('should display correct group information on each tab', async () => {
    // Policies tab
    await expect(page.getByRole('link', { name: 'Policies', exact: true })).toBeVisible();
    await page.getByRole('link', { name: 'Policies', exact: true }).click();
    await expect(page.locator('section')).toContainText('test-policy');
    await page.getByRole('button', { name: 'Identity policy management' }).click();
    await page.getByRole('link', { name: 'View policy', exact: true }).click();
    await page.getByText('Vault ACL policies test-policy test-policy Download policy').click();
    await expect(page.getByText('Vault ACL policies test-policy test-policy Download policy')).toBeVisible();
    await page.getByRole('link', { name: 'Groups' }).click();
    await page.getByRole('link', { name: 'group-1', exact: true }).click();

    // Members tab
    await expect(page.getByRole('link', { name: 'Members' })).toBeVisible();
    await page.getByRole('link', { name: 'Members' }).click();
    await expect(page.getByRole('heading', { name: 'No members in this group yet' })).toBeVisible();

    // Parent groups tab
    await expect(page.getByRole('link', { name: 'Parent groups' })).toBeVisible();
    await page.getByRole('link', { name: 'Parent groups' }).click();
    await expect(page.getByRole('heading', { name: 'This group has no parent' })).toBeVisible();

    // Metadata tab
    await page.getByRole('link', { name: 'Metadata' }).click();
    await expect(page.getByText('hello world')).toBeVisible();
    await expect(page.getByRole('link', { name: 'Metadata' })).toBeVisible();
  });

  await test.step('edit and delete group', async () => {
    await page.getByRole('link', { name: 'Edit group' }).click();
    await expect(page.getByRole('heading', { name: 'Edit Group-1' })).toContainText('Edit Group-1');
    await page.getByRole('button', { name: 'Delete group' }).click();
    await page.getByRole('button', { name: 'Confirm' }).click();
  });

  await test.step('should show empty state after group deletion', async () => {
    await expect(page.getByRole('heading', { name: 'No groups yet' })).toContainText('No groups yet');
  });
});
