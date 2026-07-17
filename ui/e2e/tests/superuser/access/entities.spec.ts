/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';

test('entities workflow', async ({ page }) => {
  await test.step('should display entities page', async () => {
    await page.goto('dashboard');
    await page.getByRole('link', { name: 'Access control' }).click();
    await page.getByRole('link', { name: 'Entities' }).click();
    await expect(page.getByRole('heading', { name: 'Entities', exact: true })).toContainText('Entities');
    await expect(page.getByRole('heading', { name: 'No entities yet' })).toContainText('No entities yet');
    await expect(page.getByText('A list of entities in this')).toContainText(
      'A list of entities in this namespace will be listed here. Create your first entity to get started.'
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
    await page.getByRole('link', { name: 'Entities' }).click();
  });

  await test.step('create entity', async () => {
    await page.getByRole('link', { name: 'Create entity' }).click();
    await page.getByRole('textbox', { name: 'Name' }).fill('entity-1');
    await page.getByRole('textbox', { name: 'Key' }).fill('hello');
    await page.getByRole('textbox', { name: 'Value' }).fill('world');
    await page.locator('.ember-basic-dropdown-trigger').first().click();
    await page.locator('.ember-power-select-option', { hasText: 'test-policy' }).click();
    await page.getByRole('button', { name: 'Create' }).click();
  });

  await test.step('should display entity detail page', async () => {
    await expect(page.getByRole('heading', { name: 'entity-' })).toContainText('entity-1');
    await page.locator('div:nth-child(2) > .column.is-flex-center').click();
    await expect(page.getByText('Name entity-')).toBeVisible();
  });

  await test.step('should display correct entity information on each tab', async () => {
    // Aliases tab
    await expect(page.getByRole('link', { name: 'Aliases' })).toBeVisible();
    await page.getByRole('link', { name: 'Aliases' }).click();
    await expect(page.getByRole('heading', { name: 'No entity aliases for entity-1 yet' })).toBeVisible();

    // Policies tab
    await expect(page.getByRole('link', { name: 'Policies', exact: true })).toBeVisible();
    await page.getByRole('link', { name: 'Policies', exact: true }).click();
    await expect(page.locator('section')).toContainText('test-policy');
    await page.getByRole('button', { name: 'Identity policy management' }).click();
    await page.getByRole('link', { name: 'View policy', exact: true }).click();
    await page.getByText('Vault ACL policies test-policy test-policy Download policy').click();
    await expect(page.getByText('Vault ACL policies test-policy test-policy Download policy')).toBeVisible();
    await page.getByRole('link', { name: 'Entities' }).click();
    await page.getByRole('link', { name: 'entity-1', exact: true }).click();

    // Groups tab
    await expect(page.getByRole('link', { name: 'Groups' }).nth(1)).toBeVisible();
    await page.getByRole('link', { name: 'Groups' }).nth(1).click();
    await expect(
      page.getByRole('heading', { name: 'entity-1 is not a member of any groups.' })
    ).toBeVisible();

    // Metadata tab
    await page.getByRole('link', { name: 'Metadata' }).click();
    await expect(page.getByText('hello world')).toBeVisible();
    await expect(page.getByRole('link', { name: 'Metadata' })).toBeVisible();
  });

  await test.step('create and view aliases', async () => {
    await page.getByRole('link', { name: 'Add alias' }).click();
    await expect(page.getByRole('heading', { name: 'Create Entity Alias for' })).toBeVisible();
    await page.getByRole('textbox', { name: 'Name' }).fill('alias-1');
    await page.getByRole('button', { name: 'Create' }).click();
    await expect(page.locator('.hds-page-header__title-wrapper')).toBeVisible();

    await page.getByRole('link', { name: 'Edit entity alias' }).click();
    await expect(page.locator('.hds-page-header__title-wrapper')).toBeVisible();
    await page.getByRole('link', { name: 'Entity aliases' }).click();
    await expect(page.getByRole('link', { name: 'alias-1', exact: true })).toBeVisible();
  });

  await test.step('cleanup by deleting entities if it was not deleted in previous step', async () => {
    await page.getByLabel('navigation for entities').getByRole('link', { name: 'Entities' }).click();
    await page.getByRole('button', { name: 'Identity management options' }).click();
    await page.getByRole('button', { name: 'Delete' }).click();
    await page.getByRole('button', { name: 'Confirm' }).click();
    await expect(page.getByRole('heading', { name: 'No entities yet' })).toBeVisible();
    await page.getByRole('link', { name: 'Aliases' }).click();
    await expect(page.getByRole('heading', { name: 'No entity aliases yet' })).toBeVisible();
  });
});
