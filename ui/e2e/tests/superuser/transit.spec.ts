/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';
import { BasePage } from '../../pages/base';

test('transit workflow', async ({ page }) => {
  const basePage = new BasePage(page);

  // enable Transit Engine
  await basePage.enableEngine('Transit', 'transit-workflow');

  // create key
  await page.getByRole('link', { name: 'Create key' }).click();
  await page.getByRole('textbox', { name: 'Name' }).fill('testKey');
  await page.getByRole('button', { name: 'Create key' }).click();
  await expect(page.getByLabel('breadcrumbs').getByText('testKey')).toBeVisible();
  await basePage.dismissFlashMessages();

  // rotate key
  await page.getByRole('link', { name: 'Versions' }).click();
  await page.getByRole('button', { name: 'Rotate key' }).click();
  await page.getByRole('button', { name: 'Confirm' }).click();
  //verify a new version is created after rotation
  await expect(page.getByText('Version 2')).toBeVisible();

  // encrypt text
  await page.getByRole('link', { name: 'Key Actions' }).click();
  await page.getByRole('link', { name: 'Encrypt Looks up wrapping' }).click();
  await page.getByRole('textbox', { name: 'Plaintext' }).fill('testString');
  await page.getByRole('button', { name: 'Encrypt' }).click();

  // grab the generated ciphertext and copy it to clipboard
  await page.getByRole('button', { name: 'copy vault:v2:' }).click();
  await page.getByRole('button', { name: 'Close' }).click();

  // decrypt text
  await page.getByRole('link', { name: 'Decrypt' }).click();
  await page.getByRole('textbox', { name: 'Ciphertext' }).press('ControlOrMeta+v');
  await page.getByRole('button', { name: 'Decrypt' }).click();
  await expect(page.getByText('Copy your unwrapped data')).toBeVisible();
  await page.getByRole('button', { name: 'copy dGVzdFN0cmluZw==' }).click();
  await page.getByRole('button', { name: 'Close' }).click();

  // generate datakey
  await page.getByRole('link', { name: 'Datakey' }).click();
  await page.getByRole('button', { name: 'Create datakey' }).click();
  await page.locator('html').click();
  await expect(page.getByText('Copy your generated key')).toBeVisible();
  await page.getByRole('button', { name: 'copy vault:v2:' }).click();
  await page.getByRole('button', { name: 'Close' }).click();

  // rewrap key
  await page.getByRole('link', { name: 'Rewrap' }).click();
  await page.getByRole('textbox', { name: 'Ciphertext' }).press('ControlOrMeta+v');
  await page.getByRole('button', { name: 'Rewrap' }).click();
  await expect(page.getByText('Copy your token')).toBeVisible();
  await page.getByRole('button', { name: 'Close' }).click();

  // HMAC generate
  await page.getByRole('link', { name: 'HMAC' }).click();
  await page.getByRole('textbox', { name: 'Input' }).fill('test');
  await page.getByRole('button', { name: 'HMAC' }).click();
  await expect(page.getByText('Copy your unwrapped data')).toBeVisible();
  await page.getByRole('button', { name: 'copy vault:v2:' }).click();
  await page.getByRole('button', { name: 'Close' }).click();

  // HMAC verify
  await page.getByRole('link', { name: 'Verify' }).click();
  //verify test and hmac is prefilled from previous step
  await expect(
    page
      .getByLabel('Input')
      .locator('div')
      .filter({ hasText: /^test$/ })
  ).toBeVisible();
  await expect(page.getByText(/vault:v2:.*/)).toBeVisible();
  await page.getByRole('button', { name: 'Verify' }).click();

  await expect(page.getByText('Results Valid')).toBeVisible();
  await page.getByRole('button', { name: 'Close' }).click();
});
