/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { test } from '../../fixtures/demo';
import { rootApi } from '../../fixtures/root-api';
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
    await configurationSettingsPage.navigateToConfiguration(path);
    await configurationSettingsPage.navigateToGeneralSettings(engineType);
    await configurationSettingsPage.verifyUnsavedChangesModalOnNavigateAway(path);
  });

  await test.step('clean up and disable engine', async () => {
    await basePage.disableEngine(path);
  });
});

test('transit key configuration persists and deletion requires saved permission', async ({
  page,
  playwright,
}, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  const mount = 'transit-key-config';
  const name = 'config-key';
  const keyPath = `/v1/${mount}/keys/${name}`;
  try {
    await test.step('create a symmetric key and confirm its persisted defaults', async () => {
      await new BasePage(page).enableEngine('Transit', mount);
      await page.getByRole('link', { name: 'Create key' }).click();
      await page.getByRole('textbox', { name: 'Name' }).fill(name);
      await page.getByRole('button', { name: 'Create key', exact: true }).click();
      await expect(page.getByLabel('breadcrumbs').getByText(name, { exact: true })).toBeVisible();
      const created = await api.get(keyPath);
      await expect(created).toBeOK();
      expect((await created.json()).data).toMatchObject({
        name,
        type: 'aes256-gcm96',
        deletion_allowed: false,
        latest_version: 1,
      });
    });

    await test.step('return to the key list, reload, filter, and reopen the key', async () => {
      await page.getByLabel('breadcrumbs').getByRole('link', { name: mount, exact: true }).click();
      await expect(page.getByRole('link', { name, exact: true })).toBeVisible();
      await page.reload();
      await expect(page.getByRole('link', { name, exact: true })).toBeVisible();
      await page.getByRole('searchbox', { name: 'Filter keys' }).fill(name);
      await expect(page).toHaveURL((url) => url.searchParams.get('pageFilter') === name);
      await expect(page.getByRole('link', { name, exact: true })).toBeVisible();
      await page.getByRole('link', { name, exact: true }).click();
    });

    await test.step('rotate the key and retain the new version after reload', async () => {
      await page.getByRole('link', { name: 'Versions', exact: true }).click();
      await page.getByRole('button', { name: 'Rotate key', exact: true }).click();
      await page.getByRole('button', { name: 'Confirm', exact: true }).click();
      await expect(page.getByText('Version 2', { exact: true })).toBeVisible();
      const rotated = await api.get(keyPath);
      await expect(rotated).toBeOK();
      expect((await rotated.json()).data.latest_version).toBe(2);
      await page.reload();
      await expect(page.getByText('Version 2', { exact: true })).toBeVisible();
    });

    await test.step('save configuration before enabling deletion', async () => {
      await page.getByRole('link', { name: 'Edit key', exact: true }).click();
      await expect(page.getByRole('button', { name: 'Delete transit key' })).toHaveCount(0);
      await page.getByLabel('Allow deletion', { exact: true }).check();
      await expect(page.getByRole('button', { name: 'Delete transit key' })).toHaveCount(0);
      await page.locator('[data-test-toggle-label="Auto-rotation period"]').click();
      await page.locator('[data-test-ttl-value="Auto-rotation period"]').fill('2');
      await page.locator('#min_decryption_version').selectOption('2');
      await page.locator('#min_encryption_version').selectOption('2');
      await page.getByRole('button', { name: 'Update transit key' }).click();
      await expect(page.getByRole('link', { name: 'Edit key', exact: true })).toBeVisible();
      const configured = await api.get(keyPath);
      await expect(configured).toBeOK();
      expect((await configured.json()).data).toMatchObject({
        deletion_allowed: true,
        auto_rotate_period: 172800,
        min_decryption_version: 2,
        min_encryption_version: 2,
      });
      await page.reload();
      await expect(page.locator('[data-test-row-value="Deletion allowed"]').last()).toHaveText('true');
      await page.getByRole('link', { name: 'Edit key', exact: true }).click();
      await expect(page.locator('#min_decryption_version')).toHaveValue('2');
      await expect(page.locator('#min_encryption_version')).toHaveValue('2');
    });

    await test.step('delete the key and confirm its removal from storage and the list', async () => {
      await page.getByRole('button', { name: 'Delete transit key', exact: true }).click();
      await page.getByRole('button', { name: 'Confirm', exact: true }).click();
      await expect(page.getByRole('link', { name: 'Create key', exact: true })).toBeVisible();
      expect((await api.get(keyPath)).status()).toBe(404);
      await page.reload();
      await expect(page.getByRole('link', { name, exact: true })).toHaveCount(0);
    });
  } finally {
    await expect(await api.delete(`/v1/sys/mounts/${mount}`)).toBeOK();
    await api.dispose();
  }
});

test('transit signing key exposes signing and export actions', async ({ page, playwright }, testInfo) => {
  const api = await rootApi(playwright, testInfo);
  const mount = 'transit-signing-key';
  const name = 'signing-key';
  try {
    await test.step('create an exportable signing key and display its public key options', async () => {
      await new BasePage(page).enableEngine('Transit', mount);
      await page.getByRole('link', { name: 'Create key' }).click();
      await page.getByRole('textbox', { name: 'Name' }).fill(name);
      await page.getByLabel('Type', { exact: true }).selectOption('ed25519');
      await page.getByLabel('Exportable', { exact: true }).check();
      await page.getByRole('button', { name: 'Create key', exact: true }).click();
      await expect(page.getByLabel('breadcrumbs').getByText(name, { exact: true })).toBeVisible();
      const response = await api.get(`/v1/${mount}/keys/${name}`);
      await expect(response).toBeOK();
      expect((await response.json()).data).toMatchObject({
        type: 'ed25519',
        exportable: true,
        supports_signing: true,
        supports_encryption: false,
      });
      await page.getByRole('link', { name: 'Versions', exact: true }).click();
      await expect(page.getByText('Version 1', { exact: true })).toBeVisible();
      await expect(page.getByRole('button', { name: 'Public key options' })).toBeVisible();
    });

    await test.step('sign a message and carry its visible input and signature into Verify', async () => {
      await page.getByRole('link', { name: 'Key Actions', exact: true }).click();
      await expect(page.getByRole('link', { name: /^Encrypt / })).toHaveCount(0);
      await page.getByRole('link', { name: /^Sign / }).click();
      await page.getByRole('textbox', { name: 'Input', exact: true }).fill('signed message');
      await page.getByRole('button', { name: 'Sign', exact: true }).click();
      const signature = await page
        .getByRole('dialog')
        .getByRole('button', { name: /^copy vault:v1:/ })
        .innerText();
      expect(signature).toContain('vault:v1:');
      await page.getByRole('button', { name: 'Close', exact: true }).click();
      await page.getByRole('link', { name: 'Verify', exact: true }).click();
      await expect(page.getByLabel('Input', { exact: true })).toContainText('signed message');
      const signatureEditor = page.getByLabel('Signature', { exact: true });
      await expect(signatureEditor).toBeVisible();
      await expect(signatureEditor).toContainText(signature);
      await page.getByRole('button', { name: 'Verify', exact: true }).click();
      await expect(page.getByText('Results Valid')).toBeVisible();
      await page.getByRole('button', { name: 'Close', exact: true }).click();
    });

    await test.step('export using the displayed default without changing the dropdown', async () => {
      await page.getByRole('link', { name: 'Export Key', exact: true }).click();
      await expect(page.getByLabel('Key type', { exact: true })).toHaveValue('signing');
      await page.getByRole('button', { name: 'Export key', exact: true }).click();
      await expect(page.getByRole('dialog')).toBeVisible();
    });
  } finally {
    await expect(await api.delete(`/v1/sys/mounts/${mount}`)).toBeOK();
    await api.dispose();
  }
});
