/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';
import { BasePage } from '../../pages/base';

test('tools workflow', async ({ page }) => {
  const basePage = new BasePage(page);
  await page.goto('dashboard');
  // wrap data
  await page.getByRole('link', { name: 'Operational tools' }).click();
  await page.getByRole('navigation', { name: 'toolbar filters' }).locator('label').click();
  await page.getByRole('textbox', { name: 'Key' }).fill('foo');
  await page.getByRole('textbox', { name: 'Value' }).fill('bar');
  await page.locator('label').filter({ hasText: 'Wrap TTL Vault will use the' }).click();
  await page.getByRole('textbox', { name: 'Number of units' }).fill('20');
  await page.getByLabel('ttl-unit').selectOption('h');
  await page.getByRole('button', { name: 'Wrap data' }).click();
  await page.getByRole('button', { name: 'copy hvs.' }).click();
  // access the clipboard to get the copied token value
  let clipboardValue = await page.evaluate(() => navigator.clipboard.readText());
  // lookup token
  await page.getByRole('link', { name: 'Lookup' }).click();
  await page.getByRole('textbox', { name: 'Wrapped token' }).fill(clipboardValue);
  await page.getByRole('button', { name: 'Lookup token' }).click();
  await expect(page.getByRole('heading', { name: 'Lookup token' })).toContainText('Lookup token');
  await expect(page.locator('section')).toContainText('about 20 hours');
  // rewrap token
  await page.getByRole('link', { name: 'Rewrap' }).click();
  await page.getByRole('textbox', { name: 'Wrapped token' }).fill(clipboardValue);
  await page.getByRole('button', { name: 'Rewrap token' }).click();
  await expect(page.locator('label')).toContainText('Rewrapped token');
  await page.getByRole('button', { name: 'copy hvs.' }).click();
  // access the clipboard to get the rewrapped token value
  clipboardValue = await page.evaluate(() => navigator.clipboard.readText());
  // unwrap token
  await page.getByRole('link', { name: 'Unwrap' }).click();
  await page.getByRole('textbox', { name: 'Wrapped token' }).fill(clipboardValue);
  await page.getByRole('button', { name: 'Unwrap data' }).click();
  await expect(page.getByRole('code')).toContainText('{ "foo": "bar" }');
  // generate random bytes
  await page.getByRole('link', { name: 'Random' }).click();
  await page.getByRole('spinbutton', { name: 'Number of bytes' }).fill('64');
  await page.getByLabel('Output format').selectOption('hex');
  await page.getByRole('button', { name: 'Generate' }).click();
  await expect(page.getByRole('button', { name: 'copy' })).toBeVisible();
  // hash
  await page.getByRole('link', { name: 'Hash', exact: true }).click();
  await page.getByRole('textbox', { name: 'Input' }).fill('foobar');
  await page.getByLabel('Algorithm').selectOption('sha2-512');
  await page.getByLabel('Output format').selectOption('hex');
  // there are a stack of flash messages that are blocking the encode button from being clicked
  await basePage.dismissFlashMessages();
  await page.getByRole('button', { name: 'Encode to base64' }).click();
  await page.getByRole('button', { name: 'Hash' }).click();
  await expect(page.getByLabel('copy')).toContainText(
    '0a50261ebd1a390fed2bf326f2673c145582a6342d523204973d0219337f81616a8069b012587cf5635f6925f1b56c360230c19b273500ee013e030601bf2425'
  );
  // API explorer
  await page.getByRole('link', { name: 'API explorer' }).click();
  await expect(page.getByRole('heading', { name: 'API explorer' })).toContainText('API explorer');
  await page.getByRole('textbox', { name: 'Filter by tag' }).fill('lookup-accessor');
  await expect(page.locator('#operations-0-tokenLookUpAccessor')).toContainText(
    '/auth/token/lookup-accessor'
  );
});
