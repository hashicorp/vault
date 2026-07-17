/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect, test } from '@playwright/test';

test('mfa workflow', async ({ page }) => {
  await page.goto('dashboard');

  await test.step('create userpass auth method and user', async () => {
    await page.getByRole('link', { name: 'Access control' }).click();
    await page.getByRole('link', { name: 'Authentication methods' }).click();
    await page.getByRole('link', { name: 'Enable new method' }).click();
    await page.getByLabel('Userpass').click();
    await page.getByRole('button', { name: 'Enable method' }).click();
    await page.getByRole('button', { name: 'Update options' }).click();

    await page.getByRole('link', { name: 'Type of auth mount userpass/' }).click();
    await page.getByLabel('toolbar actions').getByRole('link', { name: 'Create user' }).click();
    await page.getByRole('textbox', { name: 'Username' }).fill('bob');
    await page.getByRole('textbox', { name: 'password', exact: true }).fill('bobpassword');
    await page.getByRole('button', { name: 'Save' }).click();
    await expect(page.getByRole('link', { name: 'bob', exact: true })).toBeVisible();
  });

  await test.step('navigate to MFA page', async () => {
    await page.getByRole('link', { name: 'Multi-factor authentication' }).click();
    await expect(page.getByRole('img', { name: 'MFA configure diagram' })).toBeVisible();
  });

  await test.step('create method with enforcement', async () => {
    await page.getByRole('link', { name: 'Configure MFA' }).click();
    await page.getByRole('radio', { name: 'TOTP' }).check();
    await page.getByRole('button', { name: 'Next' }).click();
    await page.getByRole('textbox', { name: 'Issuer' }).fill('mfa-totp-issuer');
    await page.getByLabel('TTL unit for Period').selectOption('m');
    await page.getByRole('textbox', { name: 'Number of units' }).fill('5');
    await page.getByRole('radio', { name: 'SHA1' }).check();
    await page.getByRole('radio', { name: '6', exact: true }).check();
    await page.getByRole('textbox', { name: 'Name' }).fill('mfa-totp-enforcement');
    await page.getByLabel('target-type').selectOption('method');
    await page.getByLabel('Auth method').locator('select').selectOption('userpass');
    await page.getByRole('button', { name: 'Add' }).click();
    await page.getByRole('button', { name: 'Continue' }).click();

    await expect(page.getByText('Issuer mfa-totp-issuer')).toBeVisible();
    await expect(page.getByText('Period 5 minutes')).toBeVisible();
    await expect(page.getByText('Digits TOTP code length. 6')).toBeVisible();
    await expect(page.getByText('Enable self-enrollment No')).toBeVisible();
    await page.getByRole('link', { name: 'Enforcements' }).click();
    await expect(page.getByRole('link', { name: 'mfa-totp-enforcement' })).toBeVisible();
  });

  await test.step('verify mfa enforcement on login', async () => {
    await page.getByRole('button', { name: 'User menu' }).click();
    await page.getByRole('button', { name: 'Copy token' }).click();
    await page.getByRole('link', { name: 'Log out' }).click();

    await page.getByLabel('Method').selectOption('userpass');
    await page.getByRole('textbox', { name: 'Username' }).fill('bob');
    await page.getByRole('textbox', { name: 'Password' }).fill('bobpassword');
    await page.getByRole('button', { name: 'Sign in' }).click();

    await expect(page.getByRole('heading', { name: 'Multi-factor authentication' })).toBeVisible();
    await page.getByText('Enter your authentication code to log in. TOTP passcode').click();
    await page.getByRole('textbox', { name: 'TOTP passcode' }).fill('12');
    await page.getByRole('button', { name: 'Verify' }).click();
    await expect(page.getByRole('alert', { name: 'Error' })).toBeVisible();
    await page.getByRole('button', { name: 'Cancel' }).click();
    await expect(page.getByRole('heading', { name: 'Sign in to Vault' })).toBeVisible();
  });

  await test.step('delete mfa method and userpass auth method', async () => {
    await page.getByLabel('Method').selectOption('token');
    // get token value from clipboard that we copied before logging out
    const suToken = await page.evaluate(() => navigator.clipboard.readText());
    await page.getByRole('textbox', { name: 'Token' }).fill(suToken);
    await page.getByRole('button', { name: 'Sign in' }).click();
    await page.getByRole('link', { name: 'Access control' }).click();
    await page.getByRole('link', { name: 'Multi-factor authentication' }).click();
    await page.getByRole('link', { name: 'TOTP' }).click();
    await page.getByRole('link', { name: 'Enforcements' }).click();
    await page.getByRole('link', { name: 'mfa-totp-enforcement' }).click();
    await page.getByRole('button', { name: 'Delete' }).click();
    await page.getByRole('textbox', { name: 'Type mfa-totp-enforcement' }).fill('mfa-totp-enforcement');
    await page.getByLabel('Delete enforcement?').getByRole('button', { name: 'Delete' }).click();
    await page.getByRole('link', { name: 'Methods', exact: true }).click();
    await page.getByRole('link', { name: 'TOTP' }).click();
    await page.getByRole('link', { name: 'Configuration' }).click();
    await page.getByRole('button', { name: 'Delete' }).click();
    await page.getByRole('button', { name: 'Confirm' }).click();

    await page.getByRole('link', { name: 'Authentication methods' }).click();
    await page
      .getByRole('link', { name: 'Type of auth mount userpass/' })
      .getByLabel('Overflow options')
      .click();
    await page.getByRole('button', { name: 'Disable' }).click();
    await page.getByRole('button', { name: 'Confirm' }).click();
  });
});
