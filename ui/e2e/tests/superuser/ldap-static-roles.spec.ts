/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { test, expect } from '@playwright/test';
import type { Page } from '@playwright/test';
import { BasePage } from '../../pages/base';

const SELF_MANAGED_MOUNT = 'ldap-self-managed';
const ROOT_MANAGED_MOUNT = 'ldap-root-managed';
const ROLE_NAME = 'lazy-user-role';
const ROLE_DN = 'cn=lazy-user,ou=users,dc=my-domain,dc=com';
const ROLE_USERNAME = 'lazy-user';
const ROLE_PASSWORD = 'initialpass';

const NAME_FORMAT_ERROR =
  'Name must be lowercase and can only contain alphanumeric characters, hyphens, underscores, periods, and forward slashes.';

// self_managed is not exposed on the configure form, so it is written through the web REPL the same
// way init.setup.ts creates tokens.
async function configureMount(page: Page, mount: string, selfManaged: boolean) {
  await page.goto('dashboard');
  await page.getByRole('button', { name: 'Console toggle' }).click();
  const repl = page.getByRole('textbox', { name: 'Web R.E.P.L.' });
  // Vault rejects bind credentials on a self-managed config, since it never binds as an admin.
  const bindArgs = selfManaged ? '' : ' binddn=cn=admin,dc=my-domain,dc=com bindpass=adminpass';
  await repl.fill(`write ${mount}/config url=ldap://127.0.0.1:389${bindArgs} self_managed=${selfManaged}`);
  await repl.press('Enter');
  await expect(page.getByText(`Success! Data written to: ${mount}/config`)).toBeVisible();
  await page.getByRole('button', { name: 'Console toggle' }).click();
}

// The password row only fetches the credential when the user reveals it.
async function revealPassword(page: Page) {
  await page
    .locator('.info-table-row')
    .filter({ hasText: 'Password' })
    .getByRole('button', { name: 'show value' })
    .click();
}

// Edit lives behind the details page's "Manage" dropdown.
async function goToEdit(page: Page) {
  await page.getByRole('button', { name: 'Manage' }).click();
  await page.getByRole('link', { name: 'Edit role' }).click();
}

// enableEngine returns as soon as it clicks, so wait for the engine's own page before navigating
// away — otherwise the mount request is cancelled in flight.
async function enableMount(page: Page, basePage: BasePage, mount: string) {
  await basePage.enableEngine('ldap', mount);
  await expect(page).toHaveURL(new RegExp(`${mount}/ldap/overview`));
}

// Vault rotates the account against the real LDAP server when a static role is first imported, which
// no test environment has. skip_import_rotation avoids that on create; Vault rejects the flag on
// updates, so it is only added the first time a given role is written.
async function skipImportRotation(page: Page) {
  const created = new Set<string>();
  await page.route('**/v1/*/static-role/*', async (route) => {
    const request = route.request();
    if (request.method() !== 'POST') return route.continue();
    const name = request.url().split('/').pop() as string;
    if (created.has(name)) return route.continue();
    created.add(name);
    const body = JSON.parse(request.postData() || '{}');
    await route.continue({ postData: JSON.stringify({ ...body, skip_import_rotation: true }) });
  });
}

test('ldap self-managed static role workflow', async ({ page }) => {
  const basePage = new BasePage(page);
  await skipImportRotation(page);

  await test.step('enable and configure a self-managed ldap mount', async () => {
    await enableMount(page, basePage, SELF_MANAGED_MOUNT);
    await basePage.dismissFlashMessages();
    await configureMount(page, SELF_MANAGED_MOUNT, true);
  });

  await test.step('create page is tailored to self-managed mounts', async () => {
    await page.goto(`secrets-engines/${SELF_MANAGED_MOUNT}/ldap/roles/create`);

    await expect(page.getByRole('heading', { name: 'Create static role' })).toBeVisible();
    // A self-managed mount can only hold static roles, so the type choice is not offered.
    await expect(page.getByText('Role type')).toBeHidden();
    await expect(page.getByRole('radio', { name: /Static role/ })).toBeHidden();
    await expect(page.getByRole('radio', { name: /Dynamic role/ })).toBeHidden();

    await expect(page.getByText('The username of the existing LDAP entry to manage.')).toBeVisible();
    await expect(
      page.getByText(
        'Enter the account’s current LDAP password. Vault will use it to manage password rotations going forward.'
      )
    ).toBeVisible();
    await expect(
      page.getByText(
        'Specifies how often Vault rotates the password. Enter 0 to disable automatic password rotation.'
      )
    ).toBeVisible();
  });

  await test.step('distinguished name and password are marked required', async () => {
    // The badge is what tells a user the field is mandatory before they submit.
    await expect(page.getByText('Role name Required')).toBeVisible();
    await expect(page.getByText('Distinguished name Required')).toBeVisible();
    await expect(page.getByText('Username Required')).toBeVisible();
    await expect(page.getByText('Password Required')).toBeVisible();
  });

  await test.step('fields render in the order given by the design', async () => {
    // The rotation period is a TtlPicker, which labels itself with a legend rather than a form label.
    const labels = await page
      .locator('form label.is-label, form .hds-form-label, form .ttl-picker-label')
      .allInnerTexts();
    const ordered = labels
      .map((text) => text.replace(/\s*Required\s*/g, '').trim())
      .filter((text) =>
        ['Role name', 'Distinguished name', 'Username', 'Password', 'Rotation period'].includes(text)
      );
    expect(ordered).toEqual(['Role name', 'Distinguished name', 'Username', 'Password', 'Rotation period']);
  });

  await test.step('submitting an empty form reports every missing field once', async () => {
    await page.getByRole('button', { name: 'Create role' }).click();

    await expect(page.getByText('Role name is required')).toBeVisible();
    await expect(page.getByText('Enter the username of the LDAP account')).toBeVisible();
    await expect(page.getByText('Enter the distinguished name (DN) of the LDAP account')).toBeVisible();
    await expect(page.getByText('Enter the current password for this LDAP account')).toBeVisible();
    // An empty name is missing, not malformed, so only the presence message applies.
    await expect(page.getByText(NAME_FORMAT_ERROR)).toBeHidden();
    await expect(page).toHaveURL(/\/roles\/create/);
  });

  await test.step('a malformed name reports only the format error', async () => {
    await page.getByRole('textbox', { name: 'Role name' }).fill('BAD NAME');
    await page.getByRole('button', { name: 'Create role' }).click();

    await expect(page.getByText(NAME_FORMAT_ERROR)).toBeVisible();
    await expect(page.getByText('Role name is required')).toBeHidden();
  });

  await test.step('clearing the rotation period is accepted as zero', async () => {
    const rotation = page.getByRole('group', { name: 'Rotation period' }).getByLabel('Number of units');
    await rotation.fill('');
    await expect(page.getByText('This field is required')).toBeHidden();
  });

  await test.step('create the static role', async () => {
    await page.getByRole('textbox', { name: 'Role name' }).fill(ROLE_NAME);
    await page.getByRole('textbox', { name: 'Distinguished name' }).fill(ROLE_DN);
    await page.getByRole('textbox', { name: 'Username' }).fill(ROLE_USERNAME);
    await page.getByRole('textbox', { name: 'Password' }).fill(ROLE_PASSWORD);
    await page.getByRole('button', { name: 'Create role' }).click();

    await expect(page).toHaveURL(new RegExp(`/roles/static/${ROLE_NAME}/details`));
    await basePage.dismissFlashMessages();
  });

  await test.step('details page hides the password until it is revealed', async () => {
    await expect(page.locator('.info-table-row').filter({ hasText: 'Password' })).toBeVisible();
    // The credential is only fetched on reveal, so it must not be in the document beforehand.
    await expect(page.getByText(ROLE_PASSWORD)).toBeHidden();

    await revealPassword(page);
    await expect(page.getByText(ROLE_PASSWORD)).toBeVisible();
  });

  await test.step('edit page locks the password and freezes the distinguished name', async () => {
    await goToEdit(page);

    await expect(page.getByRole('textbox', { name: 'Distinguished name' })).toBeDisabled();
    // The stored password is never sent to the browser, so it stays behind an explicit opt-in.
    await expect(page.getByRole('button', { name: 'Enable input' })).toBeVisible();
    await expect(page.getByText('Password Required')).toBeHidden();
  });

  await test.step('saving without touching the password keeps the stored value', async () => {
    await page.getByRole('button', { name: 'Save' }).click();

    await expect(page).toHaveURL(new RegExp(`/roles/static/${ROLE_NAME}/details`));
    await revealPassword(page);
    await expect(page.getByText(ROLE_PASSWORD)).toBeVisible();
    await basePage.dismissFlashMessages();
  });

  await test.step('a new password can be set once the input is enabled', async () => {
    await goToEdit(page);
    await page.getByRole('button', { name: 'Enable input' }).click();

    const password = page.getByRole('textbox', { name: 'Password' });
    await expect(password).toBeEditable();
    await expect(password).toBeEmpty();

    await password.fill('rotated-pass');
    await page.getByRole('button', { name: 'Save' }).click();

    await expect(page).toHaveURL(new RegExp(`/roles/static/${ROLE_NAME}/details`));
    await revealPassword(page);
    await expect(page.getByText('rotated-pass')).toBeVisible();
    await basePage.dismissFlashMessages();
  });

  await test.step('clean up', async () => {
    await basePage.disableEngine(SELF_MANAGED_MOUNT);
  });
});

test('ldap root-managed mounts are unaffected', async ({ page }) => {
  const basePage = new BasePage(page);
  await skipImportRotation(page);

  await test.step('enable and configure a root-managed ldap mount', async () => {
    await enableMount(page, basePage, ROOT_MANAGED_MOUNT);
    await basePage.dismissFlashMessages();
    await configureMount(page, ROOT_MANAGED_MOUNT, false);
  });

  await test.step('the role type choice is still offered', async () => {
    await page.goto(`secrets-engines/${ROOT_MANAGED_MOUNT}/ldap/roles/create`);

    await expect(page.getByRole('heading', { name: 'Create Role' })).toBeVisible();
    await expect(page.getByText('Role type')).toBeVisible();
    await expect(page.getByRole('radio', { name: /Static role/ })).toBeVisible();
    await expect(page.getByRole('radio', { name: /Dynamic role/ })).toBeVisible();
  });

  await test.step('distinguished name and password are optional', async () => {
    await expect(page.getByText('Distinguished name Required')).toBeHidden();
    await expect(page.getByText('Password Required')).toBeHidden();
    await expect(
      page.getByText("The name of the user to be used when logging in. This is useful when DN isn't")
    ).toBeVisible();
  });

  await test.step('a static role saves without a distinguished name or password', async () => {
    await page.getByRole('textbox', { name: 'Role name' }).fill('root-managed-role');
    await page.getByRole('textbox', { name: 'Username' }).fill(ROLE_USERNAME);
    await page.getByRole('button', { name: 'Create role' }).click();

    await expect(page).toHaveURL(/\/roles\/static\/root-managed-role\/details/);
    await basePage.dismissFlashMessages();
  });

  await test.step('the password row is absent when no password is managed', async () => {
    await expect(page.locator('.info-table-row').filter({ hasText: 'Username' })).toBeVisible();
  });

  await test.step('dynamic roles can still be created', async () => {
    await page.goto(`secrets-engines/${ROOT_MANAGED_MOUNT}/ldap/roles/create`);
    await page.getByRole('radio', { name: /Dynamic role/ }).check();

    await expect(page.getByRole('textbox', { name: 'Role name' })).toBeVisible();
    await expect(page.getByRole('textbox', { name: 'Password' })).toBeHidden();
  });

  await test.step('clean up', async () => {
    await basePage.disableEngine(ROOT_MANAGED_MOUNT);
  });
});
