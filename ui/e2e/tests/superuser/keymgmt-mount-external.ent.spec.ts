/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect, test } from '@playwright/test';

const PINNED_PLUGIN_DATA = {
  data: {
    name: 'vault-plugin-secrets-keymgmt',
    type: 'secret',
    version: 'v0.17.0+ent',
  },
};

const PLUGIN_CATALOG_DATA = {
  request_id: 'request_id',
  lease_id: '',
  renewable: false,
  lease_duration: 0,
  data: {
    detailed: [
      {
        builtin: true,
        deprecation_status: 'supported',
        name: 'keymgmt',
        type: 'secret',
        version: 'v0.18.1+builtin',
      },
      {
        builtin: false,
        name: 'vault-plugin-secrets-keymgmt',
        sha256: 'sha256',
        type: 'secret',
        version: '',
      },
      {
        builtin: false,
        name: 'vault-plugin-secrets-keymgmt',
        sha256: 'sha256',
        type: 'secret',
        version: 'v0.16.0+ent',
      },
      {
        builtin: false,
        name: 'vault-plugin-secrets-keymgmt',
        sha256: 'sha256',
        type: 'secret',
        version: 'v0.17.0+ent',
      },
      {
        builtin: false,
        name: 'vault-plugin-secrets-keymgmt',
        sha256: 'sha256',
        type: 'secret',
        version: 'v0.18.0+ent',
      },
    ],
  },
};

test('mount external keymgmt workflow', async ({ page }) => {
  await test.step('mock the keymgmt pinned version response', async () => {
    await page.route('**v1/sys/plugins/pins/secret/vault-plugin-secrets-keymgmt', async (route) => {
      if (route.request().method() === 'GET') {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify(PINNED_PLUGIN_DATA),
        });
      } else {
        await route.continue();
      }
    });
  });

  await test.step('mock the plugin catalog response to return builtin and external keymgmt plugins', async () => {
    await page.route('**v1/sys/plugins/catalog', async (route) => {
      if (route.request().method() === 'GET') {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify(PLUGIN_CATALOG_DATA),
        });
      } else {
        await route.continue();
      }
    });
  });

  await page.goto('dashboard');

  await test.step('navigate to enable Key Management engine', async () => {
    await page.getByRole('link', { name: 'Secrets', exact: true }).click();
    await page.getByRole('link', { name: 'Enable new engine' }).click();
    await page.getByLabel('Key Management - enabled').click();
    await page.getByRole('textbox', { name: 'Path' }).fill('keymgmt-external');
  });

  await test.step('verify builtin and external plugin type options are visible', async () => {
    await expect(page.getByText('Built-in plugin Preregistered')).toBeVisible();
    await expect(page.getByText('External plugin External')).toBeVisible();
    await expect(page.getByText('Plugin version Required')).not.toBeVisible();
  });

  await test.step('selecting external plugin type shows plugin version dropdown', async () => {
    await page.locator('label:nth-child(2) > .hds-form-radio-card__control-wrapper').click();
    await expect(page.getByText('Plugin version Required')).toBeVisible();
    await expect(page.getByLabel('Plugin version Required')).toContainText(
      'v0.17.0+ent (pinned) v0.16.0+ent v0.18.0+ent'
    );
  });

  await test.step('pinned version is selected by default with no warning', async () => {
    await expect(page.getByLabel('Version differs from pinned')).not.toBeVisible();
  });

  await test.step('selecting a non-pinned version shows a warning', async () => {
    await page.getByLabel('Plugin version Required').selectOption('v0.16.0+ent');
    await expect(page.getByLabel('Version differs from pinned')).toContainText(
      'You have selected v0.16.0+ent, but version v0.17.0+ent is pinned for this plugin. Enabling the engine with this version will override the pinned version for this mount.'
    );
  });

  await test.step('re-selecting the pinned version clears the warning', async () => {
    await page.getByLabel('Plugin version Required').selectOption('v0.17.0+ent');
    await expect(page.getByLabel('Version differs from pinned')).not.toBeVisible();
  });

  await test.step('enabling engine shows error with external plugin name', async () => {
    await page.getByRole('button', { name: 'Enable engine' }).click();
    await expect(page.getByLabel('Error')).toContainText(
      'plugin not found in the catalog: vault-plugin-secrets-keymgmt'
    );
  });
});
