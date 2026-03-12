/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect, test } from '@playwright/test';

test('keymgmt workflow', async ({ page }) => {
  await test.step('mock the distribution response', async () => {
    await page.route('**/v1/keymgmt-builtin/kms/test-provider/key/test-key', async (route) => {
      if (route.request().method() === 'PUT') {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
        });
      } else {
        await route.continue();
      }
    });
  });

  await page.goto('dashboard');

  await test.step('enable Key Management secrets engine', async () => {
    await page.getByRole('link', { name: 'Secrets', exact: true }).click();
    await page.getByRole('link', { name: 'Enable new engine' }).click();
    await page.getByLabel('Key Management - enabled').click();
    await page.getByRole('textbox', { name: 'Path' }).fill('keymgmt-builtin');
    await page.getByRole('button', { name: 'Method Options' }).click();
    await page.getByRole('textbox', { name: 'Description' }).fill('This is a keymgmt mount.');
    await page.getByRole('checkbox', { name: 'Local' }).check();
    await page.getByRole('button', { name: 'Enable engine' }).click();

    await expect(page.getByText('Success', { exact: true })).toBeVisible();
    await expect(
      page.getByText('Successfully mounted the keymgmt secrets engine at keymgmt-builtin')
    ).toBeVisible();
    await page.getByRole('button', { name: 'Dismiss' }).click();
  });

  await test.step('create a provider', async () => {
    await page.getByRole('link', { name: 'Create provider' }).click();
    await page.getByLabel('Type').selectOption('azurekeyvault');
    await page.getByRole('textbox', { name: 'Provider name' }).fill('test-provider');
    await page.getByRole('textbox', { name: 'Key Vault instance name' }).fill('keyvault-name');
    await page.getByRole('textbox', { name: 'client_id' }).fill('a0454cd1-e28e-405e-bc50-7477fa8a00b7');
    await page.getByRole('textbox', { name: 'client_secret' }).fill('eR%HizuCVEpAKgeaUEx');
    await page.getByRole('textbox', { name: 'tenant_id' }).fill('cd4bf224-d114-4f96-9bbc-b8f45751c43f');
    await page.getByRole('button', { name: 'Create provider' }).click();

    await expect(page.locator('span').filter({ hasText: 'test-provider' })).toBeVisible();
    await expect(page.getByText('Azure Key Vault')).toBeVisible();
    await expect(page.getByText('keyvault-name')).toBeVisible();
    await expect(page.getByText('None')).toBeVisible();
  });

  await test.step('create a key', async () => {
    await page.getByRole('link', { name: 'Keys' }).click();
    await page.getByRole('link', { name: 'Create key' }).click();
    await page.getByRole('textbox', { name: 'Key name' }).fill('test-key');
    await page.getByRole('checkbox', { name: 'Allow deletion' }).check();
    await page.getByRole('button', { name: 'Create key' }).click();

    await expect(page.locator('section')).toContainText('test-key');
    await expect(page.locator('section')).toContainText('rsa-2048');
    await expect(page.locator('section')).toContainText('Yes');
  });

  await test.step('distribute key to provider', async () => {
    await page.getByLabel('toolbar actions').getByRole('button', { name: 'Distribute key' }).click();
    await page.getByText('Search').click();
    await page.getByRole('option', { name: 'test-provider' }).click();
    await page.getByText('Encrypt').click();
    await page.getByText('Decrypt').click();
    await page.getByText('Sign').click();
    await page.getByText('Verify').click();
    await page.getByText('Wrap', { exact: true }).click();
    await page.getByText('Unwrap').click();
    await page.getByRole('radio', { name: 'HSM' }).check();
    await page.getByRole('button', { name: 'Distribute key' }).click();

    await expect(page.getByText('Success', { exact: true })).toBeVisible();
    await expect(page.getByText('Successfully distributed key test-key to test-provider')).toBeVisible();
    await page.getByRole('button', { name: 'Dismiss' }).click();
  });
});
