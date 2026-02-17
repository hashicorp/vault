/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test as base } from '@playwright/test';
import fs from 'fs';
import path from 'path';
import { USER_POLICY_MAP } from './policies';

export type UserSetupOptions = {
  userType: string;
};

// use superuser as the default policy if not provided in the config for a project
export const setup = base.extend<UserSetupOptions>({
  userType: 'superuser',
});

// setup will run once before all tests
setup('initialize vault and setup user for testing', async ({ page, userType }) => {
  // on fresh app load navigating to the root will land us on the initialize page
  await page.goto('./');
  // initialize vault
  await page.getByRole('spinbutton', { name: 'Key shares' }).fill('1');
  await page.getByRole('spinbutton', { name: 'Key threshold' }).fill('1');
  await page.getByRole('button', { name: 'Initialize' }).click();
  // listen for download event so we can get the unseal key and root token
  const downloadPromise = page.waitForEvent('download');
  await page.getByRole('button', { name: 'Download keys' }).click();
  const download = await downloadPromise;
  const keysPath = path.join(__dirname, `/tmp/${userType}-keys.json`);
  await download.saveAs(keysPath);
  const { keys, root_token } = JSON.parse(fs.readFileSync(keysPath, 'utf-8'));
  // unseal vault
  await page.getByRole('link', { name: 'Continue to Unseal' }).click();
  await page.getByRole('textbox', { name: 'Unseal Key Portion' }).fill(keys[0]);
  await page.getByRole('button', { name: 'Unseal' }).click();
  // use the root token to login
  await page.getByRole('textbox', { name: 'Token' }).fill(root_token);
  await page.getByRole('button', { name: 'Sign in' }).click();
  // create a policy for a specific user persona
  // defaults to superuser but should be passed in via the project config in playwright.config.ts
  await page.getByRole('link', { name: 'Access', exact: true }).click();
  await page.getByRole('link', { name: 'Create ACL policy' }).click();
  await page.getByRole('textbox', { name: 'Policy name' }).fill(userType);
  await page.getByRole('radio', { name: 'Code editor' }).check();
  await page.getByRole('textbox', { name: 'Policy editor' }).fill(USER_POLICY_MAP[userType]);
  await page.getByRole('button', { name: 'Create policy' }).click();
  // there is no UI workflow for creating tokens with specific policies
  // generate a token using the web REPL and assign the new policy to it
  await page.getByRole('button', { name: 'Console toggle' }).click();
  await page
    .getByRole('textbox', { name: 'web R.E.P.L.' })
    .fill(`write -field=client_token auth/token/create policies=${userType} ttl=1d`);
  await page.getByRole('textbox', { name: 'web R.E.P.L.' }).press('Enter');
  const newToken = await page.locator('.console-ui-output pre').innerText();
  await page.getByRole('button', { name: 'Console toggle' }).click();
  // log out with the root token and log in with the new token/policy
  await page.getByRole('button', { name: 'User menu' }).click();
  await page.getByRole('link', { name: 'Log out' }).click();
  await page.getByRole('textbox', { name: 'Token' }).fill(newToken);
  await page.getByRole('button', { name: 'Sign in' }).click();
  // wait for the dashboard to load to ensure login was successful
  await page.waitForURL('**/dashboard');
  // save the authenticated state to file
  // subsequent tests can then reuse this session data
  await page.context().storageState({ path: path.join(__dirname, `/tmp/${userType}-session.json`) });
});
