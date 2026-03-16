/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';

test('console workflow', async ({ page }) => {
  await page.goto('dashboard');

  await test.step('open console and verify content', async () => {
    await page.getByRole('button', { name: 'Console toggle' }).click();
    // verify console is open and has expected content
    await expect(page.getByRole('textbox', { name: 'Web R.E.P.L.' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Close console' })).toBeVisible();

    // verify clicking maximize button adds expected class to panel and clicking minimize removes it
    await page.getByRole('button', { name: 'Maximize window' }).click();
    await expect(page.locator('.console-ui-panel')).toHaveClass(/fullscreen/);

    await page.getByRole('button', { name: 'Minimize window' }).click();
    await expect(page.locator('.console-ui-panel')).not.toHaveClass(/fullscreen/);
  });

  await test.step('execute command in console and verify output', async () => {
    await page
      .getByRole('textbox', { name: 'Web R.E.P.L.' })
      .fill('write sys/mounts/console-route-test type=kv');
    await page.getByRole('textbox', { name: 'Web R.E.P.L.' }).press('Enter');
    await expect(page.getByText('Success! Data written to: sys')).toBeVisible();

    // verify clicking close button hides the console
    await page.getByRole('button', { name: 'Close console' }).click();
    await expect(page.getByRole('textbox', { name: 'Web R.E.P.L.' })).toBeHidden();

    // navigate to the secrets page and verify the new mount is visible
    await page.getByRole('link', { name: 'Secrets', exact: true }).click();
    await expect(page.getByRole('link', { name: 'console-route-test/' })).toBeVisible();
  });
});
