/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';

test('rgp policy workflow', async ({ page }) => {
  await page.goto('dashboard');
  // nav to rgp policies and create a policy
  await page.getByRole('link', { name: 'Access control' }).click();
  await page.getByRole('link', { name: 'Role governing policies' }).click();
  await page.getByRole('heading', { name: 'RGP policies', exact: true }).click();
  await page.getByRole('link', { name: 'Create RGP policy' }).click();
  await page.getByRole('textbox', { name: 'Policy name' }).fill('rgp-policy');
  await page.getByRole('button', { name: 'How to write a policy' }).click();
  await page.getByText('Here is an example policy').click();
  await page.getByRole('button', { name: 'Copy' }).nth(1).click();
  // access the clipboard to get the example policy
  const clipboardValue = await page.evaluate(() => navigator.clipboard.readText());
  await page.getByRole('button', { name: 'Close' }).click();
  await page.getByRole('textbox', { name: 'Policy editor' }).fill(clipboardValue);
  await page.getByLabel('Enforcement level').selectOption('advisory');
  await page.getByRole('button', { name: 'Create policy' }).click();
  await page.getByRole('heading', { name: 'rgp-policy' }).click();
  await page.getByRole('button', { name: 'Download policy' }).click();
  await expect(page.getByRole('alert', { name: 'Info' })).toBeVisible();
  await expect(page.getByLabel('Enforcement level: advisory')).toBeVisible();

  // edit
  await page.getByRole('link', { name: 'Edit policy' }).click();
  await page.getByLabel('Policy').clear();
  const updatedValue = clipboardValue + '\n# just a comment';
  await page.getByLabel('Policy').fill(updatedValue);
  await page.getByLabel('Enforcement level', { exact: true }).selectOption('hard-mandatory');
  await page.getByRole('button', { name: 'Save' }).click();
  await expect(page.getByLabel('Enforcement level: hard-')).toBeVisible();
  await expect(page.getByRole('code')).toContainText('# just a comment');

  // delete
  await page.getByRole('link', { name: 'Edit policy' }).click();
  await page.getByRole('button', { name: 'Delete policy' }).click();
  await page.getByRole('button', { name: 'Confirm' }).click();
  await expect(page.getByRole('link', { name: 'rgp-policy', exact: true })).not.toBeVisible();
});

test('egp policy workflow', async ({ page }) => {
  await page.goto('dashboard');
  // nav to egp policies and create a policy
  await page.getByRole('link', { name: 'Access control' }).click();
  await page.getByRole('link', { name: 'Endpoint governing policies' }).click();
  await page.getByRole('heading', { name: 'EGP policies', exact: true }).click();

  await page.getByRole('link', { name: 'Create EGP policy' }).click();
  await page.getByRole('textbox', { name: 'Policy name' }).fill('egp-policy');
  await page.getByRole('button', { name: 'How to write a policy' }).click();

  await page.getByText('Here is an example policy').click();
  await page.getByRole('button', { name: 'Copy' }).nth(1).click();
  // access the clipboard to get the example policy
  const clipboardValue = await page.evaluate(() => navigator.clipboard.readText());
  await page.getByRole('button', { name: 'Close' }).click();
  await page.getByRole('textbox', { name: 'Policy editor' }).fill(clipboardValue);
  await page.getByLabel('Enforcement level').selectOption('advisory');
  await page.getByRole('textbox', { name: 'Paths list item' }).fill('foo');
  await page.getByRole('button', { name: 'Add' }).click();
  await page.getByRole('textbox', { name: 'Paths list item 1' }).fill('bar');
  await page.getByRole('button', { name: 'Add' }).click();
  await page.getByRole('button', { name: 'Create policy' }).click();
  await page.getByRole('heading', { name: 'egp-policy' }).click();
  await page.getByRole('button', { name: 'Download policy' }).click();
  await expect(page.getByRole('alert', { name: 'Info' })).toBeVisible();
  await expect(page.getByLabel('Enforcement level: advisory')).toBeVisible();
  await expect(page.getByRole('listitem').filter({ hasText: 'foo' })).toBeVisible();
  await expect(page.getByRole('listitem').filter({ hasText: 'bar' })).toBeVisible();

  // edit
  await page.getByRole('link', { name: 'Edit policy' }).click();
  await page.getByLabel('Policy').clear();
  const updatedValue = clipboardValue + '\n# just a comment';
  await page.getByLabel('Policy').fill(updatedValue);
  await page.getByLabel('Enforcement level', { exact: true }).selectOption('hard-mandatory');
  await page.getByRole('button', { name: 'delete row' }).nth(1).click();
  await page.getByRole('textbox', { name: 'Paths list item 1' }).fill('baz');
  await page.getByRole('button', { name: 'Add' }).click();
  await page.getByRole('button', { name: 'Save' }).click();
  await expect(page.getByLabel('Enforcement level: hard-')).toBeVisible();
  await page.getByRole('button', { name: 'Show more code' }).click();
  await expect(page.getByText('# just a comment')).toBeVisible();
  await expect(page.getByRole('listitem').filter({ hasText: 'foo' })).toBeVisible();
  await expect(page.getByRole('listitem').filter({ hasText: 'bar' })).not.toBeVisible();
  await expect(page.getByRole('listitem').filter({ hasText: 'baz' })).toBeVisible();

  // delete
  await page.getByRole('link', { name: 'Edit policy' }).click();
  await page.getByRole('button', { name: 'Delete policy' }).click();
  await page.getByRole('button', { name: 'Confirm' }).click();
  await expect(page.getByRole('link', { name: 'egp-policy', exact: true })).not.toBeVisible();
});
