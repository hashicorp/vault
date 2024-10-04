/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import {
  currentRouteName,
  currentURL,
  settled,
  click,
  findAll,
  fillIn,
  visit,
  typeIn,
} from '@ember/test-helpers';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';
import { runCmd, tokenWithPolicyCmd } from 'vault/tests/helpers/commands';

import { create } from 'ember-cli-page-object';
import page from 'vault/tests/pages/settings/mount-secret-backend';
import configPage from 'vault/tests/pages/secrets/backend/configuration';
import authPage from 'vault/tests/pages/auth';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import logout from 'vault/tests/pages/logout';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import { mountableEngines } from 'vault/helpers/mountable-secret-engines'; // allEngines() includes enterprise engines, those are tested elsewhere
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import { SELECTORS as OIDC } from 'vault/tests/helpers/oidc-config';
import { adminOidcCreateRead, adminOidcCreate } from 'vault/tests/helpers/secret-engine/policy-generator';
import { WIF_ENGINES } from 'vault/helpers/mountable-secret-engines';

const consoleComponent = create(consoleClass);

// enterprise backends are tested separately
const BACKENDS_WITH_ENGINES = ['kv', 'pki', 'ldap', 'kubernetes'];
module('Acceptance | settings/mount-secret-backend', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    this.calcDays = (hours) => {
      const days = Math.floor(hours / 24);
      const remainder = hours % 24;
      return `${days} days ${remainder} hours`;
    };
    return authPage.login();
  });

  test('it sets the ttl correctly when mounting', async function (assert) {
    // always force the new mount to the top of the list
    const path = `mount-kv-${this.uid}`;
    const defaultTTLHours = 100;
    const maxTTLHours = 300;
    await page.visit();

    assert.strictEqual(currentRouteName(), 'vault.cluster.settings.mount-secret-backend');
    await page.selectType('kv');
    await page
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
    assert.strictEqual(configPage.defaultTTL, `${this.calcDays(defaultTTLHours)}`, 'shows the proper TTL');
    assert.strictEqual(configPage.maxTTL, `${this.calcDays(maxTTLHours)}`, 'shows the proper max TTL');
  });

  test('it sets the ttl when enabled then disabled', async function (assert) {
    // always force the new mount to the top of the list
    const path = `mount-kv-${this.uid}`;
    const maxTTLHours = 300;

    await page.visit();

    assert.strictEqual(currentRouteName(), 'vault.cluster.settings.mount-secret-backend');
    await page.selectType('kv');
    await page
      .path(path)
      .toggleOptions()
      .enableDefaultTtl()
      .enableMaxTtl()
      .maxTTLUnit('h')
      .maxTTLVal(maxTTLHours)
      .submit();
    await configPage.visit({ backend: path });
    assert.strictEqual(configPage.defaultTTL, '1 month 1 day', 'shows system default TTL');
    assert.strictEqual(configPage.maxTTL, `${this.calcDays(maxTTLHours)}`, 'shows the proper max TTL');
  });

  test('it sets the max ttl after pki chosen, resets after', async function (assert) {
    await page.visit();
    assert.strictEqual(currentRouteName(), 'vault.cluster.settings.mount-secret-backend');
    await page.selectType('pki');
    assert.dom('[data-test-input="maxLeaseTtl"]').exists();
    assert
      .dom('[data-test-input="maxLeaseTtl"] [data-test-ttl-toggle]')
      .isChecked('Toggle is checked by default');
    assert.dom('[data-test-input="maxLeaseTtl"] [data-test-ttl-value]').hasValue('3650');
    assert.dom('[data-test-input="maxLeaseTtl"] [data-test-select="ttl-unit"]').hasValue('d');

    // Go back and choose a different type
    await page.back();
    await page.selectType('database');
    assert.dom('[data-test-input="maxLeaseTtl"]').exists('3650');
    assert
      .dom('[data-test-input="maxLeaseTtl"] [data-test-ttl-toggle]')
      .isNotChecked('Toggle is unchecked by default');
    await page.enableMaxTtl();
    assert.dom('[data-test-input="maxLeaseTtl"] [data-test-ttl-value]').hasValue('');
    assert.dom('[data-test-input="maxLeaseTtl"] [data-test-select="ttl-unit"]').hasValue('s');
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
    await page.path(path).submit();
    await page.secretList();
    await settled();
    await page.enableEngine();
    await page.selectType('kv');
    await page.path(path).submit();
    assert.dom('[data-test-message-error-description]').containsText(`path is already in use at ${path}`);
    assert.strictEqual(currentRouteName(), 'vault.cluster.settings.mount-secret-backend');

    await page.secretList();
    await settled();
    assert
      .dom(`[data-test-secrets-backend-link=${path}]`)
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
    await consoleComponent.toggle();
    await consoleComponent.runCommands(
      [
        // delete any previous mount with same name
        `delete sys/mounts/${enginePath}`,
        `write sys/policies/acl/kv-v2-degrade policy=${btoa(V2_POLICY)}`,
        'write -field=client_token auth/token/create policies=kv-v2-degrade',
      ],
      false
    );
    await settled();
    const userToken = consoleComponent.lastLogOutput;
    await logout.visit();
    await authPage.login(userToken);
    // create the engine
    await mountSecrets.visit();
    await mountSecrets.selectType('kv');
    await mountSecrets.path(enginePath).setMaxVersion(101).submit();
    await settled();
    assert
      .dom('[data-test-flash-message]')
      .containsText(
        `You do not have access to the config endpoint. The secret engine was mounted, but the configuration settings were not saved.`
      );
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${enginePath}/kv/list`,
      'After mounting, redirects to secrets list page'
    );
    await configPage.visit({ backend: enginePath });
    await settled();
  });

  test('it should transition to mountable addon engine after mount success', async function (assert) {
    // test supported backends that ARE ember engines (enterprise only engines are tested individually)
    const addons = mountableEngines().filter((e) => BACKENDS_WITH_ENGINES.includes(e.type));
    assert.expect(addons.length);

    for (const engine of addons) {
      await consoleComponent.runCommands([
        // delete any previous mount with same name
        `delete sys/mounts/${engine.type}`,
      ]);
      await mountSecrets.visit();
      await mountSecrets.selectType(engine.type);
      await mountSecrets.path(engine.type).submit();
      assert.strictEqual(
        currentRouteName(),
        `vault.cluster.secrets.backend.${engine.engineRoute}`,
        `Transitions to ${engine.displayName} route on mount success`
      );
      await consoleComponent.runCommands([
        // cleanup after
        `delete sys/mounts/${engine.type}`,
      ]);
    }
  });

  test('it should transition to mountable non-addon engine after mount success', async function (assert) {
    // test supported backends that are not ember engines (enterprise only engines are tested individually)
    const nonEngineBackends = supportedSecretBackends().filter((b) => !BACKENDS_WITH_ENGINES.includes(b));
    // add back kv because we want to test v1
    const engines = mountableEngines().filter((e) => nonEngineBackends.includes(e.type) || e.type === 'kv');
    assert.expect(engines.length);

    for (const engine of engines) {
      await consoleComponent.runCommands([
        // delete any previous mount with same name
        `delete sys/mounts/${engine.type}`,
      ]);
      await mountSecrets.visit();
      await mountSecrets.selectType(engine.type);
      await mountSecrets.path(engine.type);
      if (engine.type === 'kv') {
        await mountSecrets.toggleOptions().version(1);
      }
      await mountSecrets.submit();

      assert.strictEqual(
        currentRouteName(),
        `vault.cluster.secrets.backend.list-root`,
        `${engine.type} navigates to list view`
      );
      await consoleComponent.runCommands([
        // cleanup after
        `delete sys/mounts/${engine.type}`,
      ]);
    }
  });

  test('it should transition back to backend list for unsupported backends', async function (assert) {
    const unsupported = mountableEngines().filter((e) => !supportedSecretBackends().includes(e.type));
    assert.expect(unsupported.length);

    for (const engine of unsupported) {
      await consoleComponent.runCommands([
        // delete any previous mount with same name
        `delete sys/mounts/${engine.type}`,
      ]);
      await mountSecrets.visit();
      await mountSecrets.selectType(engine.type);
      await mountSecrets.path(engine.type).submit();

      assert.strictEqual(
        currentRouteName(),
        `vault.cluster.secrets.backends`,
        `${engine.type} returns to backends list`
      );
    }
  });

  test('it should transition to different locations for kv v1 and v2', async function (assert) {
    assert.expect(4);
    const v2 = 'kv-v2';
    await consoleComponent.runCommands([
      // delete any previous mount with same name
      `delete sys/mounts/${v2}`,
    ]);
    await mountSecrets.visit();
    await mountSecrets.selectType('kv');
    await mountSecrets.path(v2).submit();

    assert.strictEqual(currentURL(), `/vault/secrets/${v2}/kv/list`, `${v2} navigates to list url`);
    assert.strictEqual(
      currentRouteName(),
      `vault.cluster.secrets.backend.kv.list`,
      `${v2} navigates to list url`
    );

    const v1 = 'kv-v1';
    await consoleComponent.runCommands([
      // delete any previous mount with same name
      `delete sys/mounts/${v1}`,
    ]);
    await mountSecrets.visit();
    await mountSecrets.selectType('kv');
    await mountSecrets.path(v1).toggleOptions().version(1).submit();

    assert.strictEqual(currentURL(), `/vault/secrets/${v1}/list`, `${v1} navigates to list url`);
    assert.strictEqual(
      currentRouteName(),
      `vault.cluster.secrets.backend.list-root`,
      `${v1} navigates to list route`
    );
  });

  module('WIF secret engines', function () {
    test('it sets identity_token_key on mount config using search select list, resets after', async function (assert) {
      // create an oidc/key
      await runCmd(`write identity/oidc/key/some-key allowed_client_ids="*"`);

      for (const engine of WIF_ENGINES) {
        await page.visit();
        await page.selectType(engine);
        await click(GENERAL.toggleGroup('Method Options'));
        assert
          .dom('[data-test-search-select-with-modal]')
          .exists('Search select with modal component renders');
        await clickTrigger('#key');
        const dropdownOptions = findAll('[data-option-index]').map((o) => o.innerText);
        assert.ok(dropdownOptions.includes('some-key'), 'search select options show some-key');
        await click(GENERAL.searchSelect.option(GENERAL.searchSelect.optionIndex('some-key')));
        assert
          .dom(GENERAL.searchSelect.selectedOption())
          .hasText('some-key', 'some-key was selected and displays in the search select');
      }
      // Go back and choose a non-wif engine type
      await page.back();
      await page.selectType('ssh');
      assert
        .dom('[data-test-search-select-with-modal]')
        .doesNotExist('for type ssh, the modal field does not render.');
      // cleanup
      await runCmd(`delete identity/oidc/key/some-key`);
    });

    test('it allows a user with permissions to oidc/key to create an identity_token_key', async function (assert) {
      for (const engine of WIF_ENGINES) {
        const path = `secrets-adminPolicy-${engine}`;
        const newKey = `key-${uuidv4()}`;
        const secrets_admin_policy = adminOidcCreateRead(path);
        const secretsAdminToken = await runCmd(
          tokenWithPolicyCmd(`secrets-admin-${path}`, secrets_admin_policy)
        );

        await logout.visit();
        await authPage.login(secretsAdminToken);
        await page.visit();
        await page.selectType(engine);
        await page.path(path);
        await click(GENERAL.toggleGroup('Method Options'));
        await clickTrigger('#key');
        // create new key
        await fillIn(GENERAL.searchSelect.searchInput, newKey);
        await click(GENERAL.searchSelect.options);
        assert.dom('#search-select-modal').exists('modal with form opens');
        assert.dom('[data-test-modal-title]').hasText('Create new key', 'Create key modal renders');

        await click(OIDC.keySaveButton);
        assert.dom('#search-select-modal').doesNotExist('modal disappears onSave');
        assert.dom(GENERAL.searchSelect.selectedOption()).hasText(newKey, `${newKey} is now selected`);

        await page.submit();
        await visit(`/vault/secrets/${path}/configuration`);
        await click(SES.configurationToggle);
        assert
          .dom(GENERAL.infoRowValue('Identity Token Key'))
          .hasText(newKey, 'shows identity token key on configuration page');
        // cleanup
        await runCmd(`delete sys/mounts/${path}`);
        await runCmd(`delete identity/oidc/key/some-key`);
        await runCmd(`delete identity/oidc/key/${newKey}`);
      }
    });

    test('it allows user with NO access to oidc/key to manually input an identity_token_key', async function (assert) {
      for (const engine of WIF_ENGINES) {
        const path = `secrets-noOidcAdmin-${engine}`;
        const secretsNoOidcAdminPolicy = adminOidcCreate(path);
        const secretsNoOidcAdminToken = await runCmd(
          tokenWithPolicyCmd(`secrets-noOidcAdmin-${path}`, secretsNoOidcAdminPolicy)
        );
        // create an oidc/key that they can then use even if they can't read it.
        await runCmd(`write identity/oidc/key/general-key allowed_client_ids="*"`);

        await logout.visit();
        await authPage.login(secretsNoOidcAdminToken);
        await page.visit();
        await page.selectType(engine);
        await page.path(path);
        await click(GENERAL.toggleGroup('Method Options'));
        // type-in fallback component to create new key
        await typeIn(GENERAL.inputSearch('key'), 'general-key');
        await page.submit();
        assert
          .dom(GENERAL.latestFlashContent)
          .hasText(`Successfully mounted the aws secrets engine at ${path}.`);

        await visit(`/vault/secrets/${path}/configuration`);
        await click(SES.configurationToggle);
        assert
          .dom(GENERAL.infoRowValue('Identity Token Key'))
          .hasText('general-key', 'shows identity token key on configuration page');
        // cleanup
        await runCmd(`delete sys/mounts/${path}`);
      }
    });
  });
});
