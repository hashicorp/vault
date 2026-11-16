/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect } from '@playwright/test';
import { test } from '../../fixtures/demo';
import { BasePage } from '../../pages/base';

test('console workflow', async ({ page }) => {
  const basePage = new BasePage(page);
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

    await basePage.disableEngine('console-route-test');
  });
});

test('console read, list, delete, wrap and error output', async ({ page }) => {
  const basePage = new BasePage(page);
  const mount = 'console-crud-test';
  await page.goto('dashboard');
  await page.getByRole('button', { name: 'Console toggle' }).click();
  const repl = page.getByRole('textbox', { name: 'Web R.E.P.L.' });
  const panel = page.locator('.console-ui-panel');
  const run = async (command: string) => {
    await repl.fill(command);
    await repl.press('Enter');
  };

  await test.step('write a secret to a new kv mount', async () => {
    await run(`write sys/mounts/${mount} type=kv`);
    await expect(panel.getByText(`Success! Data written to: sys/mounts/${mount}`)).toBeVisible();
    await run(`write ${mount}/my-secret foo=bar`);
    await expect(panel.getByText(`Success! Data written to: ${mount}/my-secret`)).toBeVisible();
  });

  await test.step('read returns the stored data', async () => {
    await run(`read ${mount}/my-secret`);
    await expect(panel.locator('.console-ui-output').last()).toContainText(/foo\s+bar/);
    await run(`read -field=foo ${mount}/my-secret`);
    await expect(panel.locator('.console-ui-output').last()).toHaveText('bar');
  });

  await test.step('list returns the secret key', async () => {
    await run(`list ${mount}/`);
    await expect(panel.locator('.console-ui-output').last()).toContainText('my-secret');
  });

  await test.step('read with wrap-ttl returns wrapping info', async () => {
    await run(`read -wrap-ttl=5m ${mount}/my-secret`);
    await expect(panel.locator('.console-ui-output').last()).toContainText(`creation_path`);
    await expect(panel.locator('.console-ui-output').last()).toContainText(`${mount}/my-secret`);
  });

  await test.step('delete removes the secret and reading it again reports an error', async () => {
    await run(`delete ${mount}/my-secret`);
    await expect(
      panel.getByText(`Success! Data deleted (if it existed) at: ${mount}/my-secret`)
    ).toBeVisible();
    await run(`read ${mount}/my-secret`);
    const error = panel.locator('.console-ui-alert--error').last();
    await expect(error).toContainText(`Error reading from: ${mount}/my-secret.`);
    await expect(error).toContainText('Code: 404');
  });

  await test.step('clean up the mount', async () => {
    await page.getByRole('button', { name: 'Close console' }).click();
    await basePage.disableEngine(mount);
  });
});
