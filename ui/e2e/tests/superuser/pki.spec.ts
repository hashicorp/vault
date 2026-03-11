/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';

test('pki workflow', async ({ page }) => {
  await page.goto('dashboard');
  // enable PKI Engine
  await page.getByRole('link', { name: 'Secrets', exact: true }).click();
  // skip if intro page is shown
  const skipButton = page.getByRole('button', { name: 'Skip' });
  if (await skipButton.isVisible()) {
    await skipButton.click();
  }
  await page.getByRole('link', { name: 'Enable new engine' }).click();
  await page.getByLabel('PKI Certificates - enabled').click();
  await page.getByRole('textbox', { name: 'Path' }).fill('pki-engine');
  await page.locator('label').filter({ hasText: 'Default Lease TTL Vault will' }).click();
  await page.getByLabel('TTL unit for Default Lease TTL').selectOption('m');
  await page
    .getByRole('group', { name: 'Default Lease TTL Lease will' })
    .getByLabel('Number of units')
    .fill('5');
  await page.getByLabel('TTL unit for Max Lease TTL').selectOption('m');
  await page
    .getByRole('group', { name: 'Max Lease TTL Lease will' })
    .getByLabel('Number of units')
    .fill('10');
  await page.getByRole('button', { name: 'Enable engine' }).click();

  // configure PKI Engine
  await expect(page.getByRole('heading', { name: 'pki-engine' })).toContainText('pki-engine');
  await expect(page.locator('section')).toContainText('PKI not configured');
  await page.getByRole('link', { name: 'Configure PKI' }).click();
  await expect(page.locator('section')).toContainText('Import a CA');
  await expect(page.locator('section')).toContainText('Generate root');
  await expect(page.locator('section')).toContainText('Generate intermediate CSR');
  await page.locator('label').filter({ hasText: 'Generate root Generates a new' }).click();
  await page.getByLabel('Type').selectOption('internal');
  await page.getByRole('textbox', { name: 'Common name' }).fill('pki-common-name');
  await page.getByRole('textbox', { name: 'Issuer name' }).fill('pki-issuer');
  await page.getByRole('textbox', { name: 'Not valid after' }).fill('36000');
  await page.getByLabel('TTL unit for Not before').selectOption('m');
  await page.getByRole('textbox', { name: 'Number of units' }).fill('10');
  await page.getByLabel('Format', { exact: true }).selectOption('der');
  await page.getByRole('textbox', { name: 'Max path length' }).fill('16');
  await page.getByRole('button', { name: 'Key parameters' }).click();
  await page.getByRole('textbox', { name: 'Key name' }).fill('pki-key');
  await page.getByLabel('Key type').selectOption('ed25519');
  await page.getByRole('button', { name: 'Done' }).click();
  await expect(page.getByRole('heading', { name: 'pki-engine configuration' })).toBeVisible();
  await expect(page.getByRole('link', { name: 'PKI Certificates settings' })).toBeVisible();
  await expect(page.getByRole('button', { name: 'Generate policy' })).toBeVisible();
  await expect(page.getByRole('link', { name: 'Exit configuration' })).toBeVisible();
  await page.getByRole('link', { name: 'General settings' }).click();
  await expect(page.getByRole('button', { name: 'Generate policy' })).not.toBeVisible();
  await page.getByRole('link', { name: 'pki-engine', exact: true }).click();

  // create role
  await page.getByRole('link', { name: 'Roles', exact: true }).click();
  await page.getByRole('link', { name: 'Create role' }).click();
  await page.getByRole('textbox', { name: 'Role Name' }).fill('pki-role');
  await expect(page.getByLabel('issuer_ref')).toContainText('pki-issuer');
  await page.locator('label').filter({ hasText: 'Use default issuer' }).click();
  await page.getByLabel('TTL unit for TTL').selectOption('m');
  await page
    .getByText('TTL Set relative certificate')
    .getByRole('textbox', { name: 'Number of units' })
    .fill('30');
  await page.getByLabel('TTL unit for Backdate validity').selectOption('m');
  await page
    .getByRole('group', { name: 'Backdate validity Lease will' })
    .getByLabel('Number of units')
    .fill('10');
  await page.locator('label').filter({ hasText: 'Max TTL Vault will use the' }).click();
  await page.getByRole('button', { name: 'Domain handling' }).click();
  await page.getByLabel('TTL unit for Max TTL').selectOption('m');
  await page
    .getByRole('group', { name: 'Max TTL Lease will expire' })
    .getByLabel('Number of units')
    .fill('50');
  await page.getByRole('checkbox', { name: 'Allow any name' }).check();
  await page.getByRole('button', { name: 'Create' }).click();
  await page.getByRole('link', { name: 'Generate Certificate' }).click();
  await page.getByRole('textbox', { name: 'Common name' }).fill('role-cert');
  await page.getByRole('textbox', { name: 'Number of units' }).fill('1');
  await page.getByLabel('TTL unit for TTL').selectOption('m');
  await page.getByRole('button', { name: 'Generate' }).click();
  await expect(page.getByLabel('Next steps')).toBeVisible();
  // used for testing role certificate revocation as issuer certificates cannot be revoked via the UI currently
  const roleSerialNumber = await page.getByText(/^[0-9a-f]{2}(?::[0-9a-f]{2}){9,}$/i).textContent();
  await page.getByLabel('breadcrumbs').getByText('Roles').click();
  await expect(page.locator('section')).toContainText('pki-role');

  // view certificates
  await page.getByRole('link', { name: 'Certificates' }).click();
  await expect(page.getByLabel('certificate serial number')).toHaveCount(2);
  await page.getByLabel('certificate serial number').filter({ hasText: roleSerialNumber }).click();
  await expect(page.getByRole('heading', { name: 'View Certificate' })).toContainText('View Certificate');
  const downloadPromise = page.waitForEvent('download');
  await page.getByRole('button', { name: 'Download' }).click();
  await expect(page.getByText('Your download has started')).toBeVisible();
  const download = await downloadPromise;
  expect(download.suggestedFilename()).toMatch(/\.pem$/);
  await expect(page.getByText('Revocation time')).not.toBeVisible();
  await page.getByRole('button', { name: 'Revoke certificate' }).click();
  await page.getByRole('button', { name: 'Confirm' }).click();
  await expect(page.getByText('Revocation time')).toBeVisible();
  await expect(page.getByRole('button', { name: 'Revoke certificate' })).not.toBeVisible();

  // generate issuer
  await page.getByRole('link', { name: 'pki-engine' }).click();
  await page.getByRole('link', { name: 'Issuers', exact: true }).click();
  await expect(page.getByRole('button', { name: 'Generate policy' })).toBeVisible();
  await expect(page.getByRole('button', { name: 'Manage', exact: true })).toBeVisible();
  await page.getByRole('button', { name: 'Generate', exact: true }).click();
  await page.getByRole('link', { name: 'Root', exact: true }).click();
  await page.getByLabel('Type').selectOption('exported');
  await page.getByRole('textbox', { name: 'Common name' }).fill('pki-common-name-exported');
  await page.getByRole('textbox', { name: 'Issuer name' }).fill('pki-issuer-exported');
  await page.getByRole('button', { name: 'Done' }).click();
  await expect(page.getByRole('heading', { name: 'View generated root' })).toBeVisible();
  await expect(page.getByLabel('Next steps')).toBeVisible();
  await page.getByRole('button', { name: 'Done' }).click();

  // generate keys
  await page.getByRole('link', { name: 'Keys' }).click();
  await expect(page.locator('section')).toContainText(
    'Below is information about the private keys used by the issuers to sign certificates. While certificates represent a public assertion of an identity, private keys represent the private part of that identity, a secret used to prove who they are and who they trust.'
  );
  await expect(page.locator('section')).toContainText('pki-key');
  await expect(page.getByRole('link', { name: 'pki-key' })).toBeVisible();
  await page.getByRole('link', { name: 'Generate' }).click();
  await page.getByRole('textbox', { name: 'Key name' }).fill('pki-generated-key');
  await page.getByLabel('Type', { exact: true }).selectOption('internal');
  await page.getByLabel('Key type').selectOption('rsa');
  await page.getByLabel('Key bits').selectOption('3072');
  await page.getByRole('button', { name: 'Generate key' }).click();
  await expect(page.getByRole('heading', { name: 'View Key' })).toBeVisible();
  await expect(page.locator('section')).toContainText('pki-generated-key');

  // overview
  await page.getByRole('link', { name: 'pki-engine' }).click();
  await expect(page.locator('section')).toContainText(
    'Issuers View issuers The total number of issuers in this PKI mount. Includes both root and intermediate certificates. 2'
  );
  await expect(page.locator('section')).toContainText(
    'Roles View roles The total number of roles in this PKI mount that have been created to generate certificates. 1'
  );
  await expect(page.getByText('Issue certificate Begin')).toBeVisible();
  await page.getByText('Type to find a role...').click();
  await expect(page.getByRole('option', { name: 'pki-role' })).toBeVisible();
  await expect(page.getByText('View certificate Quickly view')).toBeVisible();
  await page.getByText('33:a3:').click();
  await expect(page.getByRole('option', { name: roleSerialNumber })).toBeVisible();
  await expect(page.getByText('View issuer Choose or type an')).toBeVisible();
  await page.getByText('Type to find an issuer...').click();
  await expect(page.getByRole('option').first()).toBeVisible();
});
