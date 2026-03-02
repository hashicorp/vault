/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';
import fs from 'fs';
import path from 'path';

const keysPath = path.resolve(__dirname, '../../tmp/superuser-keys.json');

test('sealing/unsealing workflow', async ({ page }) => {
  await page.goto('dashboard');
  await page.getByRole('link', { name: 'Resilience and recovery' }).click();
  await page.getByRole('link', { name: 'Seal Vault' }).click();
  await page.getByRole('button', { name: 'Seal' }).click();
  await page.getByRole('button', { name: 'Confirm' }).click();
  await expect(page.getByText('Vault is sealed')).toBeVisible();

  // unseal vault for sequential tests
  const { keys, root_token } = JSON.parse(fs.readFileSync(keysPath, 'utf-8'));
  await page.getByRole('textbox', { name: 'Unseal Key Portion' }).fill(keys[0]);
  await page.getByRole('button', { name: 'Unseal' }).click();
  await page.getByRole('textbox', { name: 'Token' }).fill(root_token);
  await page.getByRole('button', { name: 'Sign in' }).click();
  await expect(page.getByRole('button', { name: 'root' })).toBeVisible();
  await expect(page.getByRole('link', { name: 'Dashboard' })).toBeVisible();
});
