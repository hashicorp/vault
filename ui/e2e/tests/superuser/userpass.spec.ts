/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';

test('userpass workflow', async ({ page }) => {
  // nav to access control and enable userpass auth method
  await page.goto('dashboard');
  await page.getByRole('link', { name: 'Access control' }).click();
  await page.getByRole('link', { name: 'Authentication methods' }).click();

  // if intro page is visible, click enable method there otherwise click enable method in toolbar
  if (await page.getByRole('link', { name: 'Enable a new method' }).isVisible()) {
    await page.getByRole('link', { name: 'Enable a new method' }).click();
  } else {
    await page.getByRole('link', { name: 'Enable new method' }).click();
  }

  // enable userpass auth method
  await page.getByLabel('Userpass - enabled engine type').click();
  await page.getByRole('button', { name: 'Enable method' }).click();
  await page.getByRole('button', { name: 'Update options' }).click();
  await expect(page.getByRole('link', { name: 'Type of auth mount userpass/' })).toBeVisible();

  // create a test user
  await page.getByRole('link', { name: 'Type of auth mount userpass/' }).click();
  await page.getByLabel('toolbar actions').getByRole('link', { name: 'Create user' }).click();
  await page.getByRole('textbox', { name: 'Username' }).fill('testUser');
  await page.getByRole('textbox', { name: 'password', exact: true }).fill('test');
  await page.getByRole('button', { name: 'Save' }).click();
  await expect(page.getByRole('link', { name: 'testuser', exact: true })).toBeVisible();

  // log out and log in with the new user
  await page.getByRole('button', { name: 'User menu' }).click();
  await page.getByRole('link', { name: 'Log out' }).click();
  await page.getByLabel('Method').selectOption('userpass');
  await page.getByRole('textbox', { name: 'Username' }).fill('testUser');
  await page.getByRole('textbox', { name: 'Password' }).fill('test');
  await page.getByRole('button', { name: 'Sign in' }).click();

  // verify login was successful by verifying the user menu is visible and contains the username
  await expect(page.getByRole('button', { name: 'User menu' })).toBeVisible();
  await page.getByRole('button', { name: 'User menu' }).click();
  await expect(page.getByText('Testuser')).toBeVisible();
});
