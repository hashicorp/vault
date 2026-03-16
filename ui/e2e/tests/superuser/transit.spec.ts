/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';
import { BasePage } from '../../pages/base';
import { ConfigurationSettingsPage } from '../../pages/configuration-settings';

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

test('transit tune workflow', async ({ page }) => {
  const basePage = new BasePage(page);
  const configurationSettingsPage = new ConfigurationSettingsPage(page);

  const path = 'transit-tune';
  const engineType = 'transit';

  await test.step('enable Transit secrets engine mount', async () => {
    await basePage.enableEngine(engineType, path);
  });

  await test.step('navigate to configuration page from manage dropdown ', async () => {
    await configurationSettingsPage.navigateToConfiguration(path);
  });

  // Transit does not have plugin settings, so we only need to test for general settings

  await test.step('navigate and verify general settings form', async () => {
    await configurationSettingsPage.navigateToGeneralSettings(engineType);
    await configurationSettingsPage.editAndVerifyGeneralSettings(path, engineType);
    await page.getByRole('link', { name: 'Exit configuration' }).click();
  });

  await test.step('ensure that we navigate back to the transit overview page when Exit configuration is clicked', async () => {
    await expect(
      page
        .locator('div')
        .filter({ hasText: `${path} Manage Create key` })
        .nth(3)
    ).toBeVisible();
  });

  await test.step('verify unsaved changes modal works in general settings', async () => {
    // Navigate back to general settings
    await configurationSettingsPage.navigateToConfiguration(path);
    await configurationSettingsPage.navigateToGeneralSettings(engineType);

    // Test Unsaved changes modal
    await configurationSettingsPage.verifyUnsavedChangesModalOnNavigateAway(path);
  });

  await test.step('clean up and disable engine', async () => {
    await basePage.disableEngine(path);
  });
});
