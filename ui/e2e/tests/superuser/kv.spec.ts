/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';

test('kvv2 workflow', async ({ page }) => {
  await page.goto('dashboard');
  // enable kv secrets engine
  await page.getByRole('link', { name: 'Secrets Engines' }).click();
  await page.getByRole('link', { name: 'Enable new engine' }).click();
  await page.locator('div').filter({ hasText: 'KV' }).nth(4).click();
  await page.getByRole('textbox', { name: 'Path' }).click();
  await page.getByRole('textbox', { name: 'Path' }).fill('kv-test');
  await page.getByRole('button', { name: 'Enable engine' }).click();
  // once enabled it should navigate to the secrets engine overview page
  await expect(page.locator('section')).toContainText('kv-test version 2');
  await expect(page.locator('section')).toContainText(
    'No secrets yet When created, secrets will be listed here. Create a secret to get started.'
  );
  // verify that the kv engine appears in the list view
  await page.getByRole('link', { name: 'Secrets Engines' }).click();
  await page.getByRole('link', { name: 'kv-test/' }).click();
  // create a secret
  await page.getByRole('link', { name: 'Create secret' }).click();
  await page.getByRole('textbox', { name: 'Path for this secret' }).fill('foo');
  await page.getByRole('textbox', { name: 'key' }).fill('bar');
  await page.getByRole('textbox', { name: 'bar' }).fill('baz');
  await page.getByRole('button', { name: 'Save' }).click();
  // it should navigate to the overview page for the new secret
  await expect(page.locator('section')).toContainText('foo');
  await expect(page.locator('section')).toContainText(
    'Current version Create new The current version of this secret. 1'
  );
  // verify secret details
  await page.getByRole('link', { name: 'Secret', exact: true }).click();
  await expect(page.locator('section')).toContainText('bar');
  await page.getByRole('button', { name: 'show value' }).click();
  await expect(page.locator('pre')).toContainText('baz');
  await page.locator('label').click();
  await expect(page.getByRole('code')).toContainText('{ "bar": "baz" }');
  // create metadata for the secret
  await page.getByRole('link', { name: 'Metadata', exact: true }).click();
  await expect(page.locator('#app-main-content')).toContainText(
    'No custom metadata This data is version-agnostic and is usually used to describe the secret being stored. Add metadata'
  );
  await page.getByRole('link', { name: 'Edit metadata' }).click();
  await page.getByRole('textbox', { name: 'key' }).fill('meta');
  await page.getByRole('textbox', { name: 'value' }).fill('data');
  await page.getByRole('button', { name: 'Update' }).click();
  await expect(page.locator('#app-main-content')).toContainText('meta data');
  // create new version
  await page.getByRole('link', { name: 'Version History' }).click();
  await expect(page.locator('section')).toContainText('Version 1');
  await expect(page.locator('section')).toContainText('Current');
  await page.getByRole('button', { name: 'Manage version' }).click();
  await page.getByRole('link', { name: 'Create new version from 1', exact: true }).click();
  await page.getByRole('textbox', { name: 'key' }).first().fill('bar-v2');
  await page.getByRole('textbox', { name: 'bar-v2' }).fill('baz-v2');
  await page.getByRole('button', { name: 'Save' }).click();
  await expect(page.locator('section')).toContainText(
    'Current version Create new The current version of this secret. 2'
  );
  await page.getByRole('link', { name: 'Version History' }).click();
  await expect(page.locator('section')).toContainText('Version 2');
  await expect(page.locator('section')).toContainText('Current');
  await page.getByRole('link', { name: 'Version diff' }).click();
  await expect(page.locator('section')).toContainText('bar"baz"bar-v2"baz-v2"');
  // delete version 2
  await page.goto('secrets-engines/kv-test/kv/foo');
  await page.getByRole('link', { name: 'Secret', exact: true }).click();
  await page.getByRole('button', { name: 'Delete' }).click();
  await page.getByRole('radio', { name: 'Delete this version This' }).check();
  await page.getByRole('button', { name: 'Confirm' }).click();
  await expect(page.locator('section')).toContainText(
    'Current version Deleted Create new The current version of this secret was deleted'
  );
  await page.getByRole('link', { name: 'Secret', exact: true }).click();
  await expect(page.locator('section')).toContainText(
    'Version 2 of this secret has been deleted This version has been deleted but can be undeleted. View other versions of this secret by clicking the Version History tab above. KV v2 API docs'
  );
  // undelete version
  await page.getByRole('button', { name: 'Undelete' }).click();
  await expect(page.locator('section')).toContainText(
    'Current version Create new The current version of this secret. 2'
  );
  // delete latest version
  await page.getByRole('link', { name: 'Secret', exact: true }).click();
  await page.getByRole('button', { name: 'Version' }).click();
  await page.getByRole('link', { name: 'Version 1' }).click();
  await page.getByRole('button', { name: 'Delete' }).click();
  await page.getByRole('radio', { name: 'Delete latest version This' }).check();
  await page.getByRole('button', { name: 'Confirm' }).click();
  await expect(page.locator('section')).toContainText(
    'Current version Deleted Create new The current version of this secret was deleted'
  );
  // destroy version 2
  await page.getByRole('link', { name: 'Secret', exact: true }).click();
  await page.getByRole('button', { name: 'Destroy' }).click();
  await page.getByRole('button', { name: 'Confirm' }).click();
  await expect(page.locator('section')).toContainText(
    'Current version Destroyed Create new The current version of this secret has been permanently deleted and cannot be restored. 2'
  );
  await page.getByRole('link', { name: 'Secret', exact: true }).click();
  await expect(page.locator('section')).toContainText(
    'Version 2 of this secret has been permanently destroyed A version that has been permanently deleted cannot be restored. You can view other versions of this secret in the Version History tab above. KV v2 API docs'
  );
  // destroy version 1
  await page.getByRole('button', { name: 'Version' }).click();
  await page.getByRole('link', { name: 'Version 1' }).click();
  await page.getByRole('button', { name: 'Destroy' }).click();
  await page.getByRole('button', { name: 'Confirm' }).click();
  await page.getByRole('link', { name: 'Secret', exact: true }).click();
  await page.getByRole('button', { name: 'Version' }).click();
  await page.getByRole('link', { name: 'Version 1' }).click();
  await expect(page.locator('section')).toContainText(
    'Version 1 of this secret has been permanently destroyed A version that has been permanently deleted cannot be restored. You can view other versions of this secret in the Version History tab above. KV v2 API docs'
  );
});
