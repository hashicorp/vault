/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test as base, expect } from '@playwright/test';
import fs from 'fs';
import path from 'path';
import { USER_POLICY_MAP } from './policies';
import { TEST_USERS } from './test-users';
import { DISMISSED_WIZARD_KEY, WIZARD_ID_MAP } from '../app/utils/constants/wizard';

import type { TestUserType } from './test-users';

export type UserSetupOptions = {
  userType: TestUserType;
};

// Record<never, never> preserves Playwright's worker fixture-tuple overload.
export const setup = base.extend<Record<never, never>, UserSetupOptions>({
  userType: [TEST_USERS.superuser, { option: true, scope: 'worker' }],
});

// setup will run once before all tests
setup('initialize vault and setup user for testing', async ({ page, userType }) => {
  // on fresh app load navigating to the root will land us on the initialize page
  await page.goto('./');
  // manually update dismissed wizards so that they don't have to be skipped in tests before persisting storage state;
  await page.evaluate(
    ({ key, ids }) => {
      localStorage.setItem(key, JSON.stringify(ids));
    },
    { key: DISMISSED_WIZARD_KEY, ids: Object.values(WIZARD_ID_MAP) }
  );
  await setup.step('server uses the persona storage backend', async () => {
    const response = await page.request.get('/v1/sys/seal-status');
    await expect(response).toBeOK();
    expect((await response.json()).storage_type).toBe(userType.storage);
  });
  if (userType.storage === 'raft') {
    await page.getByRole('radio', { name: 'Create a new Raft cluster' }).check();
    await page.getByRole('button', { name: 'Next' }).click();
  }
  // initialize vault
  await page.getByRole('spinbutton', { name: 'Key shares' }).fill('1');
  await page.getByRole('spinbutton', { name: 'Key threshold' }).fill('1');
  await page.getByRole('button', { name: 'Initialize' }).click();
  // listen for download event so we can get the unseal key and root token
  const downloadPromise = page.waitForEvent('download');
  await page.getByRole('button', { name: 'Download keys' }).click();
  const download = await downloadPromise;
  const keysPath = path.join(__dirname, `/tmp/${userType.name}-keys.json`);
  await download.saveAs(keysPath);
  const { keys, root_token } = JSON.parse(fs.readFileSync(keysPath, 'utf-8'));
  // unseal vault
  await page.getByRole('link', { name: 'Continue to Unseal' }).click();
  await page.getByRole('textbox', { name: 'Unseal Key Portion' }).fill(keys[0]);
  await page.getByRole('button', { name: 'Unseal' }).click();
  await setup.step('waits for the unsealed node to become active before login', async () => {
    await expect
      .poll(async () => (await page.request.get('/v1/sys/health')).status(), { timeout: 30_000 })
      .toBe(200);
  });
  // use the root token to login
  await page.getByRole('textbox', { name: 'Token' }).fill(root_token);
  await page.getByRole('button', { name: 'Sign in' }).click();
  // Personas can share policy definitions while keeping separate servers and sessions.
  await page.getByRole('link', { name: 'Access control', exact: true }).click();
  // if the intro page is shown, click the create policy link there, otherwise click the create policy link in the toolbar on main page
  if (await page.getByRole('link', { name: 'Create a policy' }).isVisible()) {
    await page.getByRole('link', { name: 'Create a policy' }).click();
  } else {
    await page.getByRole('link', { name: 'Create ACL policy' }).click();
  }
  await page.getByRole('textbox', { name: 'Policy name' }).fill(userType.policy);
  await page.getByRole('radio', { name: 'Code editor' }).check();
  const policy = USER_POLICY_MAP[userType.policy];
  if (!policy) {
    throw new Error(
      `No policy "${userType.policy}" defined for persona "${userType.name}" in USER_POLICY_MAP`
    );
  }
  await page.getByRole('textbox', { name: 'Policy editor' }).fill(policy);
  await page.getByRole('button', { name: 'Create policy' }).click();
  // there is no UI workflow for creating tokens with specific policies
  // generate a token using the web REPL and assign the new policy to it
  await page.getByRole('button', { name: 'Console toggle' }).click();
  await page
    .getByRole('textbox', { name: 'web R.E.P.L.' })
    .fill(`write -field=client_token auth/token/create policies=${userType.policy} ttl=1d`);
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
  await setup.step('persona token has the selected policy rather than root privileges', async () => {
    const response = await page.request.get('/v1/auth/token/lookup-self', {
      headers: { 'X-Vault-Token': newToken },
    });
    await expect(response).toBeOK();
    const { data } = await response.json();
    expect(data.policies).toContain(userType.policy);
    expect(data.policies).not.toContain('root');
  });
  // save the localStorage state to file which includes the auth token and dismissed wizards
  // subsequent tests can then reuse the session data
  await page.context().storageState({ path: path.join(__dirname, `/tmp/${userType.name}-session.json`) });
});
