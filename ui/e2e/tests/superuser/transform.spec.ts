/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect, test } from '@playwright/test';
import { BasePage } from '../../pages/base';

test('transform workflow', async ({ page }) => {
  const basePage = new BasePage(page);

  await test.step('enable Transform secrets engine mount', async () => {
    await basePage.enableEngine('Transform', 'transform-test');
  });

  await test.step('Transform secrets engine mount saved successfully', async () => {
    await expect(page.getByText('Success', { exact: true })).toBeVisible();
    await expect(page.getByText('Successfully mounted the')).toBeVisible();
    await page.getByRole('button', { name: 'Dismiss' }).click();
  });

  await test.step('Transformation can be created', async () => {
    await page.getByRole('link', { name: 'Create transformation' }).click();
    await page.getByRole('textbox', { name: 'Name' }).fill('test-transformation');
    await page.getByRole('checkbox', { name: 'Allow deletion' }).check();
    await page.getByRole('button', { name: 'Create transformation' }).click();
    await page.getByLabel('Template').getByText('Search').click();
    await page.getByRole('option', { name: 'builtin/socialsecuritynumber' }).click();
    await page.getByRole('button', { name: 'Create transformation' }).click();
    await expect(page.getByText('test-transformation', { exact: true })).toBeVisible();
  });

  await test.step('Role can be created', async () => {
    await page.getByRole('link', { name: 'transform-test' }).click();
    await page.getByRole('link', { name: 'Roles' }).click();
    await expect(page.getByRole('heading', { name: 'No roles in this backend' })).toBeVisible();
    await expect(page.getByText('Roles in this backend will be')).toBeVisible();
    await page.getByRole('link', { name: 'Create role' }).click();
    await page.getByRole('textbox', { name: 'Name' }).fill('test-role');
    await page.getByText('Search').click();
    await page.getByRole('option', { name: 'test-transformation' }).click();
    await page.getByRole('button', { name: 'Create role' }).click();
  });

  await test.step('Template can be created', async () => {
    await page.getByRole('link', { name: 'transform-test' }).click();
    await page.getByRole('link', { name: 'Templates' }).click();
    await page.getByRole('link', { name: 'Create template' }).click();
    await page.getByRole('textbox', { name: 'Name' }).fill('test-template');
    await page.getByRole('textbox', { name: 'Pattern' }).fill('`^(19)');
    await page.getByText('Search').click();
    await page.getByRole('option', { name: 'builtin/alphalower' }).click();
    await page.getByRole('button', { name: 'Create template' }).click();
  });

  await test.step('Template saved successfully', async () => {
    await expect(page.getByText('Success')).toBeVisible();
    await expect(page.getByText('Transform template saved.')).toBeVisible();
    await page.getByRole('button', { name: 'Dismiss' }).click();
  });

  await test.step('Alphabet can be created', async () => {
    await page.getByRole('link', { name: 'transform-test' }).click();
    await page.getByRole('link', { name: 'Alphabets' }).click();
    await page.getByRole('link', { name: 'Create alphabet' }).click();
    await page.getByRole('textbox', { name: 'Name' }).click();
    await page.getByRole('textbox', { name: 'Name' }).fill('test-alphabet');
    await page.getByRole('textbox', { name: 'Alphabet' }).click();
    await page.getByRole('textbox', { name: 'Alphabet' }).fill('abc');
    await page.getByRole('button', { name: 'Create alphabet' }).click();
  });

  await test.step('Alphabet saved successfully', async () => {
    await expect(page.getByText('Success')).toBeVisible();
    await expect(page.getByText('Alphabet saved.')).toBeVisible();
    await page.getByRole('button', { name: 'Dismiss' }).click();
  });

  await test.step('Transform mount can be configured/updated', async () => {
    await page.getByRole('link', { name: 'transform-test' }).click();
    await page.getByRole('button', { name: 'Manage', exact: true }).click();
    await page.getByRole('link', { name: 'Configure' }).click();
    await page.getByRole('textbox', { name: 'Description' }).click();
    await page.getByRole('textbox', { name: 'Description' }).fill('My transform secrets engine. test');
    await page.getByRole('button', { name: 'Save changes' }).click();
  });

  await test.step('Transform mount updated successfully', async () => {
    await expect(page.getByText('Configuration saved')).toBeVisible();
    await expect(page.getByText('Engine settings successfully')).toBeVisible();
    await page.getByRole('button', { name: 'Dismiss' }).click();
  });
});
