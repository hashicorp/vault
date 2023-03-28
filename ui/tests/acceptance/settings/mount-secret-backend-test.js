/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { currentRouteName, currentURL, settled } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import { create } from 'ember-cli-page-object';
import page from 'vault/tests/pages/settings/mount-secret-backend';
import configPage from 'vault/tests/pages/secrets/backend/configuration';
import authPage from 'vault/tests/pages/auth';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import logout from 'vault/tests/pages/logout';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import { allEngines } from 'vault/helpers/mountable-secret-engines';

const consoleComponent = create(consoleClass);

module('Acceptance | settings/mount-secret-backend', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    return authPage.login();
  });

  test('it sets the ttl correctly when mounting', async function (assert) {
    // always force the new mount to the top of the list
    const path = `mount-kv-${this.uid}`;
    const defaultTTLHours = 100;
    const maxTTLHours = 300;
    const defaultTTLSeconds = (defaultTTLHours * 60 * 60).toString();
    const maxTTLSeconds = (maxTTLHours * 60 * 60).toString();

    await page.visit();

    assert.strictEqual(currentRouteName(), 'vault.cluster.settings.mount-secret-backend');
    await page.selectType('kv');
    await page
      .next()
      .path(path)
      .toggleOptions()
      .enableDefaultTtl()
      .defaultTTLUnit('h')
      .defaultTTLVal(defaultTTLHours)
      .enableMaxTtl()
      .maxTTLUnit('h')
      .maxTTLVal(maxTTLHours)
      .submit();
    await configPage.visit({ backend: path });
    assert.strictEqual(configPage.defaultTTL, defaultTTLSeconds, 'shows the proper TTL');
    assert.strictEqual(configPage.maxTTL, maxTTLSeconds, 'shows the proper max TTL');
  });

  test('it sets the ttl when enabled then disabled', async function (assert) {
    // always force the new mount to the top of the list
    const path = `mount-kv-${this.uid}`;
    const maxTTLHours = 300;
    const maxTTLSeconds = (maxTTLHours * 60 * 60).toString();

    await page.visit();

    assert.strictEqual(currentRouteName(), 'vault.cluster.settings.mount-secret-backend');
    await page.selectType('kv');
    await page
      .next()
      .path(path)
      .toggleOptions()
      .enableDefaultTtl()
      .enableDefaultTtl()
      .enableMaxTtl()
      .maxTTLUnit('h')
      .maxTTLVal(maxTTLHours)
      .submit();
    await configPage.visit({ backend: path });
    assert.strictEqual(configPage.defaultTTL, '0', 'shows the proper TTL');
    assert.strictEqual(configPage.maxTTL, maxTTLSeconds, 'shows the proper max TTL');
  });

  test('it throws error if setting duplicate path name', async function (assert) {
    const path = `kv-duplicate`;

    await consoleComponent.runCommands([
      // delete any kv-duplicate previously written here so that tests can be re-run
      `delete sys/mounts/${path}`,
    ]);

    await page.visit();

    assert.strictEqual(currentRouteName(), 'vault.cluster.settings.mount-secret-backend');
    await page.selectType('kv');
    await page.next().path(path).submit();
    await page.secretList();
    await settled();
    await page.enableEngine();
    await page.selectType('kv');
    await page.next().path(path).submit();
    assert.dom('[data-test-alert-banner="alert"]').containsText(`path is already in use at ${path}`);
    assert.strictEqual(currentRouteName(), 'vault.cluster.settings.mount-secret-backend');

    await page.secretList();
    await settled();
    assert
      .dom(`[data-test-secret-backend-row=${path}]`)
      .exists({ count: 1 }, 'renders only one instance of the engine');
  });

  test('version 2 with no update to config endpoint still allows mount of secret engine', async function (assert) {
    const enginePath = `kv-noUpdate-${this.uid}`;
    const V2_POLICY = `
      path "${enginePath}/*" {
        capabilities = ["list","create","read","sudo","delete"]
      }
      path "sys/mounts/*"
      {
        capabilities = ["create", "read", "update", "delete", "list", "sudo"]
      }

      # List existing secrets engines.
      path "sys/mounts"
      {
        capabilities = ["read"]
      }
      # Allow page to load after mount
      path "sys/internal/ui/mounts/${enginePath}" {
        capabilities = ["read"]
      }
    `;
    await consoleComponent.runCommands([
      // delete any previous mount with same name
      `delete sys/mounts/${enginePath}`,
      `write sys/policies/acl/kv-v2-degrade policy=${btoa(V2_POLICY)}`,
      'write -field=client_token auth/token/create policies=kv-v2-degrade',
    ]);
    await settled();
    const userToken = consoleComponent.lastLogOutput;
    await logout.visit();
    await authPage.login(userToken);
    // create the engine
    await mountSecrets.visit();
    await mountSecrets.selectType('kv');
    await mountSecrets.next().path(enginePath).setMaxVersion(101).submit();
    await settled();
    assert
      .dom('[data-test-flash-message]')
      .containsText(
        `You do not have access to the config endpoint. The secret engine was mounted, but the configuration settings were not saved.`
      );
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${enginePath}/list`,
      'After mounting, redirects to secrets list page'
    );
    await configPage.visit({ backend: enginePath });
    await settled();
    assert.dom('[data-test-row-value="Maximum number of versions"]').hasText('Not set');
  });

  test('it should transition to engine route on success if defined in mount config', async function (assert) {
    await consoleComponent.runCommands([
      // delete any previous mount with same name
      `delete sys/mounts/kubernetes`,
    ]);
    await mountSecrets.visit();
    await mountSecrets.selectType('kubernetes');
    await mountSecrets.next().path('kubernetes').submit();
    const { engineRoute } = allEngines().findBy('type', 'kubernetes');
    assert.strictEqual(
      currentRouteName(),
      `vault.cluster.secrets.backend.${engineRoute}`,
      'Transitions to engine route on mount success'
    );
  });
});
