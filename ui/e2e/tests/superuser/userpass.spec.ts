/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';

test('userpass workflow', async ({ page }) => {
  await test.step('enable userpass auth method', async () => {
    await page.goto('dashboard');
    await page.getByRole('link', { name: 'Access control' }).click();
    await page.getByRole('link', { name: 'Authentication methods' }).click();
    await page.getByRole('link', { name: 'Enable new method' }).click();
    await page.getByLabel('Userpass - enabled engine type').click();
    await page.getByRole('textbox', { name: 'Path' }).fill('test-userpass');
    await page.getByRole('button', { name: 'Enable method' }).click();
  });

  await test.step('create a test user', async () => {
    await page.getByRole('link', { name: 'View method' }).click();
    await page.getByLabel('toolbar actions').getByRole('link', { name: 'Create user' }).click();
    await page.getByRole('textbox', { name: 'Username' }).fill('testUser');
    await page.getByRole('textbox', { name: 'password', exact: true }).fill('test');
    await page.getByRole('button', { name: 'Save' }).click();
    await expect(page.getByRole('link', { name: 'testuser', exact: true })).toBeVisible();
  });

  await test.step('log in with new user', async () => {
    await page.getByRole('button', { name: 'User menu' }).click();
    await page.getByRole('button', { name: 'Copy token' }).click();
    await page.getByRole('link', { name: 'Log out' }).click();
    await page.getByLabel('Method').selectOption('userpass');
    await page.getByRole('textbox', { name: 'Username' }).fill('testUser');
    await page.getByRole('textbox', { name: 'Password' }).fill('test');
    await page.getByRole('button', { name: 'Advanced settings' }).click();
    await page.getByRole('textbox', { name: 'Mount path' }).fill('test-userpass');
    await page.getByRole('button', { name: 'Sign in' }).click();
  });

  await test.step('verify login was successful', async () => {
    await page.getByRole('button', { name: 'User menu' }).click();
    await expect(page.getByText('Testuser')).toBeVisible();
  });

  await test.step('disable test-userpass auth method', async () => {
    await page.getByRole('link', { name: 'Log out' }).click();
    await page.getByLabel('Method').selectOption('token');
    // get token value from clipboard that we copied before logging out
    const suToken = await page.evaluate(() => navigator.clipboard.readText());
    await page.getByRole('textbox', { name: 'Token' }).fill(suToken);
    await page.getByRole('button', { name: 'Sign in' }).click();
    await page.getByRole('link', { name: 'Access control' }).click();
    await page.getByRole('link', { name: 'Authentication methods' }).click();
    await page.getByRole('link', { name: 'Type of auth mount test-' }).getByLabel('Overflow options').click();
    await page.getByRole('button', { name: 'Disable' }).click();
    await page.getByRole('button', { name: 'Confirm' }).click();
  });
});
