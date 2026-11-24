/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { expect, test } from '@playwright/test';
import { ConfigurationSettingsPage } from '../../pages/configuration-settings';

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

const KEYMGMT_EXTERNAL_MOUNT_DATA = {
  request_id: 'request_id',
  lease_id: '',
  renewable: false,
  lease_duration: 0,
  data: {
    accessor: 'vault-plugin-secrets-keymgmt_accessor',
    config: {
      default_lease_ttl: 2764800,
      force_no_cache: false,
      listing_visibility: 'hidden',
      max_lease_ttl: 2764800,
    },
    description: '',
    external_entropy_access: false,
    local: false,
    options: {},
    path: 'keymgmt-external/',
    plugin_version: 'v0.17.0+ent',
    running_plugin_version: 'v0.17.0+ent',
    running_sha256: 'sha256',
    seal_wrap: false,
    type: 'vault-plugin-secrets-keymgmt',
    uuid: 'uuid',
  },
  wrap_info: null,
  warnings: null,
  auth: null,
  mount_type: '',
};

const UPDATED_KEYMGMT_EXTERNAL_MOUNT_DATA = {
  ...KEYMGMT_EXTERNAL_MOUNT_DATA,
  data: {
    ...KEYMGMT_EXTERNAL_MOUNT_DATA.data,
    plugin_version: 'v0.18.0+ent',
    running_plugin_version: 'v0.18.0+ent',
  },
};

test('tune external keymgmt workflow', async ({ page }) => {
  const configurationSettingsPage = new ConfigurationSettingsPage(page);

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

  await test.step('mock the keymgmt external mount response', async () => {
    await page.route('**/v1/sys/internal/ui/mounts/keymgmt-external', async (route) => {
      if (route.request().method() === 'GET') {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify(KEYMGMT_EXTERNAL_MOUNT_DATA),
        });
      } else {
        await route.continue();
      }
    });
  });

  await test.step('mock the initial keymgmt external tune response', async () => {
    await page.route('**/v1/sys/mounts/keymgmt-external/tune', async (route) => {
      if (route.request().method() === 'POST') {
        await route.fulfill({
          status: 204,
          contentType: 'application/json',
        });
      } else {
        await route.continue();
      }
    });
  });

  await page.goto('dashboard');

  await test.step("navigate to the external keymgmt mount's general settings page", async () => {
    await page.goto('secrets-engines/keymgmt-external/list');
    const path = 'keymgmt-external';
    const engineType = 'vault-plugin-secrets-keymgmt';
    await configurationSettingsPage.navigateToConfiguration(path);
    await configurationSettingsPage.editAndVerifyGeneralSettings(path, engineType, true);
    await expect(page.getByRole('paragraph').nth(3)).toContainText('v0.17.0+ent (Pinned)');
  });

  await test.step('verify that selecting an unpinned version shows the override message', async () => {
    await page.getByLabel('Update version to:').selectOption('v0.16.0+ent');
    await expect(page.getByLabel('Override pinned version')).toContainText(
      'You have selected v0.16.0+ent, but version v0.17.0+ent is pinned for this plugin. Updating to this version will override the pinned version for this mount.'
    );
  });

  await test.step('reset the version selection and verify that the override message goes away', async () => {
    await page.getByLabel('Update version to:').selectOption('');
    await expect(page.getByLabel('Override pinned version')).not.toBeVisible();
  });

  await test.step('verify that selecting an unpinned version shows the override message', async () => {
    await page.getByLabel('Update version to:').selectOption('v0.18.0+ent');
    await expect(page.getByLabel('Override pinned version')).toContainText(
      'You have selected v0.18.0+ent, but version v0.17.0+ent is pinned for this plugin. Updating to this version will override the pinned version for this mount.'
    );
  });

  await test.step('mock updated mount response after tuning with a new plugin version', async () => {
    await page.route('**/v1/sys/internal/ui/mounts/keymgmt-external', async (route) => {
      if (route.request().method() === 'GET') {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify(UPDATED_KEYMGMT_EXTERNAL_MOUNT_DATA),
        });
      } else {
        await route.continue();
      }
    });
  });

  await page.getByRole('button', { name: 'Save changes' }).click();
});
