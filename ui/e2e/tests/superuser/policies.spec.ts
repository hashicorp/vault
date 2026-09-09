/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';

test('acl policy workflow', async ({ page }) => {
  await page.goto('dashboard');
  // nav to acl policies and create a policy
  await page.getByRole('link', { name: 'Access control' }).click();
  await page.getByRole('link', { name: 'Create ACL policy' }).click();
  await page.getByRole('textbox', { name: 'Policy name' }).fill('acl-policy');
  await page.getByRole('textbox', { name: 'Resource path' }).fill('kv');
  await page.getByRole('checkbox', { name: 'read' }).check();
  await page.getByRole('checkbox', { name: 'update' }).check();
  await page.getByRole('switch', { name: 'Show preview' }).click();
  await expect(page.getByRole('code')).toContainText('path "kv" { capabilities = ["read", "update"] }');
  await page.getByRole('switch', { name: 'Hide preview' }).click();
  await page.getByRole('button', { name: 'Add rule' }).click();
  await page.getByRole('textbox', { name: 'Resource path' }).nth(1).fill('pki');
  await page.getByRole('checkbox', { name: 'list' }).nth(1).check();
  await page.getByRole('checkbox', { name: 'patch' }).nth(1).check();
  await page.getByRole('button', { name: 'Add rule' }).click();
  await page.getByRole('textbox', { name: 'Resource path' }).nth(2).fill('totp');
  await page.getByRole('checkbox', { name: 'create' }).nth(2).check();
  // check snippets
  await page.getByRole('button', { name: 'Automation snippets' }).click();
  await expect(page.getByRole('code')).toContainText(
    'resource "vault_policy" "<local identifier>" { name = "acl-policy" policy = <<EOT path "kv" { capabilities = ["read", "update"] } path "pki" { capabilities = ["list", "patch"] } path "totp" { capabilities = ["create"] } EOT }'
  );
  await page.getByRole('tab', { name: 'CLI' }).click();
  await expect(page.getByText('vault policy write acl-policy')).toBeVisible();
  await page.getByRole('button', { name: 'Delete' }).nth(2).click();
  await expect(page.getByRole('code')).toContainText(
    'vault policy write acl-policy - <<EOT path "kv" { capabilities = ["read", "update"] } path "pki" { capabilities = ["list", "patch"] } EOT'
  );
  // check change detection
  await page.getByRole('radio', { name: 'Code editor' }).check();
  await expect(page.getByRole('heading', { name: 'Policy editor' })).toBeVisible();
  await expect(page.getByRole('button', { name: 'Switch and discard changes' })).not.toBeVisible();
  await page.getByRole('radio', { name: 'Visual editor' }).check();
  await page.getByRole('radio', { name: 'Code editor' }).check();
  await page.getByRole('textbox', { name: 'Policy editor' }).clear();
  await page
    .getByRole('textbox', { name: 'Policy editor' })
    .fill(
      'vault policy write acl-policy - <<EOT path "kv" { capabilities = ["read"] } path "pki" { capabilities = ["list", "patch"] } EOT'
    );
  await page.getByRole('radio', { name: 'Visual editor' }).click();
  await page.getByRole('button', { name: 'Switch and discard changes' }).click();
  await expect(page.getByRole('code')).toContainText(
    'EOT path "kv" { capabilities = ["read", "update"] } path "pki" { capabilities = ["list", "patch"] } EOT'
  );
  await page.getByRole('button', { name: 'Create policy' }).click();
  await expect(page.getByRole('heading', { name: 'acl-policy' })).toBeVisible();
  await page.getByRole('button', { name: 'Download policy' }).click();
  await expect(page.getByRole('alert', { name: 'Info' })).toBeVisible();
  await expect(page.getByText('Policy HCL format')).toBeVisible();
  await expect(page.getByRole('button', { name: 'Automation snippets' })).toBeVisible();

  // edit
  await page.getByRole('link', { name: 'Edit policy' }).click();
  await page.getByLabel('Policy').clear();
  await page.getByLabel('Policy').fill('path "kv" { capabilities = ["read", "update"] }');
  await page.getByRole('button', { name: 'Save' }).click();
  await expect(page.getByLabel('Policy')).toContainText('path "kv" { capabilities = ["read", "update"] }');

  // delete
  await page.getByRole('link', { name: 'Edit policy' }).click();
  await page.getByRole('button', { name: 'Delete policy' }).click();
  await page.getByRole('button', { name: 'Confirm' }).click();
  await expect(page.getByRole('link', { name: 'acl-policy', exact: true })).not.toBeVisible();
});
