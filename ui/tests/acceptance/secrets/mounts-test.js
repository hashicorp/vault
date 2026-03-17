/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import {
  click,
  currentRouteName,
  currentURL,
  fillIn,
  findAll,
  typeIn,
  visit,
  waitUntil,
} from '@ember/test-helpers';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import { setupApplicationTest } from 'ember-qunit';
import { module, skip, test } from 'qunit';
import { v4 as uuidv4 } from 'uuid';
import { runCmd, tokenWithPolicyCmd } from 'vault/tests/helpers/commands';

import { create } from 'ember-cli-page-object';
import engineDisplayData from 'vault/helpers/engines-display-data';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { mountBackend } from 'vault/tests/helpers/components/mount-backend-form-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SELECTORS as OIDC } from 'vault/tests/helpers/oidc-config';
import { adminOidcCreate, adminOidcCreateRead } from 'vault/tests/helpers/secret-engine/policy-generator';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import { default as mountSecrets, default as page } from 'vault/tests/pages/settings/mount-secret-backend';
import { filterEnginesByMountCategory } from 'vault/utils/all-engines-metadata';

const consoleComponent = create(consoleClass);

// enterprise backends are tested separately
const BACKENDS_WITH_ENGINES = ['kv', 'pki', 'ldap', 'kubernetes'];
module('Acceptance | secrets-engines/enable', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    this.uid = uuidv4();
    this.calcDays = (hours) => {
      const days = Math.floor(hours / 24);
      const remainder = hours % 24;
      return `${days} days ${remainder} hours`;
    };
    return login();
  });

  test('it sets the ttl correctly when mounting', async function (assert) {
    // always force the new mount to the top of the list
    const path = `mount-kv-${this.uid}`;
    const defaultTTLHours = 100;
    const maxTTLHours = 300;
    await page.visit();

    assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.enable.index');
    await click(GENERAL.cardContainer('kv'));
    await fillIn(GENERAL.inputByAttr('path'), path);
    await click(GENERAL.button('Method Options'));
    await click(GENERAL.toggleInput('Default Lease TTL'));
    await page.defaultTTLUnit('h').defaultTTLVal(defaultTTLHours);
    await click(GENERAL.toggleInput('Max Lease TTL'));
    await page.maxTTLUnit('h').maxTTLVal(maxTTLHours);
    await click(GENERAL.submitButton);

    assert.dom(GENERAL.dropdownToggle('Manage')).exists('renders manage dropdown');
    await click(GENERAL.dropdownToggle('Manage'));
    await click(GENERAL.menuItem('Configure'));
    await click(GENERAL.tabLink('general-settings'));
    assert
      .dom(GENERAL.inputByAttr('default_lease_ttl'))
      .hasValue(`${defaultTTLHours}`, 'shows the proper TTL');
    assert.dom(GENERAL.selectByAttr('default_lease_ttl')).hasValue('h', 'shows the proper TTL unit');

    assert.dom(GENERAL.inputByAttr('max_lease_ttl')).hasValue(`${maxTTLHours}`, 'shows the proper max TTL');
    assert.dom(GENERAL.selectByAttr('max_lease_ttl')).hasValue('h', 'shows the proper max TTL unit');
  });

  test('it sets the ttl when enabled then disabled', async function (assert) {
    // always force the new mount to the top of the list
    const path = `mount-kv-${this.uid}`;
    const maxTTLHours = 300;

    await page.visit();

    assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.enable.index', 'navigates to mount page');
    await click(GENERAL.cardContainer('kv'));
    await fillIn(GENERAL.inputByAttr('path'), path);
    await click(GENERAL.button('Method Options'));
    await click(GENERAL.toggleInput('Default Lease TTL'));
    await click(GENERAL.toggleInput('Max Lease TTL'));
    await page.maxTTLUnit('h').maxTTLVal(maxTTLHours);
    await click(GENERAL.submitButton);

    assert.dom(GENERAL.dropdownToggle('Manage')).exists('renders manage dropdown');
    await click(GENERAL.dropdownToggle('Manage'));
    await click(GENERAL.menuItem('Configure'));
    await click(GENERAL.tabLink('general-settings'));

    assert.dom(GENERAL.inputByAttr('default_lease_ttl')).hasValue('32', 'shows system default TTL');
    assert.dom(GENERAL.inputByAttr('max_lease_ttl')).hasValue(`${maxTTLHours}`, 'shows the proper max TTL');
    assert.dom(GENERAL.selectByAttr('max_lease_ttl')).hasValue('h', 'shows the proper max TTL unit');
  });

  test('it sets the max ttl after pki chosen, resets after', async function (assert) {
    await page.visit();
    assert.strictEqual(currentRouteName(), 'vault.cluster.secrets.enable.index');
    await click(GENERAL.cardContainer('pki'));
    assert.dom('[data-test-input="config.max_lease_ttl"]').exists();
    assert
      .dom('[data-test-input="config.max_lease_ttl"] [data-test-ttl-toggle]')
      .isChecked('Toggle is checked by default');
    assert.dom('[data-test-input="config.max_lease_ttl"] [data-test-ttl-value]').hasValue('3650');
    assert.dom('[data-test-input="config.max_lease_ttl"] [data-test-select="ttl-unit"]').hasValue('d');

    // Go back and choose a different type
    await click(GENERAL.backButton);
    await click(GENERAL.cardContainer('database'));
    assert.dom('[data-test-input="config.max_lease_ttl"]').exists('3650');
    assert
      .dom('[data-test-input="config.max_lease_ttl"] [data-test-ttl-toggle]')
      .isNotChecked('Toggle is unchecked by default');
    await click(GENERAL.toggleInput('Max Lease TTL'));
    assert.dom('[data-test-input="config.max_lease_ttl"] [data-test-ttl-value]').hasValue('');
    assert.dom('[data-test-input="config.max_lease_ttl"] [data-test-select="ttl-unit"]').hasValue('s');
  });

  test('it should transition to mountable addon engine after mount success', async function (assert) {
    // test supported backends that ARE ember engines (enterprise only engines are tested individually)
    const addons = filterEnginesByMountCategory({ mountCategory: 'secret', isEnterprise: false }).filter(
      (e) => BACKENDS_WITH_ENGINES.includes(e.type)
    );
    assert.expect(addons.length);

    for (const engine of addons) {
      await consoleComponent.runCommands([
        // delete any previous mount with same name
        `delete sys/mounts/${engine.type}`,
      ]);
      await mountSecrets.visit();
      await mountBackend(engine.type, engine.type);

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
    const engines = filterEnginesByMountCategory({ mountCategory: 'secret', isEnterprise: false }).filter(
      (e) => (nonEngineBackends.includes(e.type) || e.type === 'kv') && e.type !== 'cubbyhole'
    );
    assert.expect(engines.length);

    for (const engine of engines) {
      await consoleComponent.runCommands([
        // delete any previous mount with same name
        `delete sys/mounts/${engine.type}`,
      ]);
      await mountSecrets.visit();
      await click(GENERAL.cardContainer(engine.type));
      await fillIn(GENERAL.inputByAttr('path'), engine.type);
      if (engine.type === 'kv') {
        await click(GENERAL.button('Method Options'));
        await mountSecrets.version(1);
      }
      await click(GENERAL.submitButton);

      const route = engineDisplayData(engine.type)?.isOnlyMountable
        ? 'configuration.general-settings'
        : 'list-root';
      const expectedRoute = `vault.cluster.secrets.backend.${route}`;
      await waitUntil(() => currentRouteName() === expectedRoute);
      assert.strictEqual(
        currentRouteName(),
        expectedRoute,
        `${engine.type} navigates to the correct view (either list if not configuration only or configuration if it is).`
      );

      await consoleComponent.runCommands([
        // cleanup after
        `delete sys/mounts/${engine.type}`,
      ]);
    }
  });

  test('it should transition back to backend list for unsupported backends', async function (assert) {
    const unsupported = filterEnginesByMountCategory({ mountCategory: 'secret', isEnterprise: false }).filter(
      (e) => !supportedSecretBackends().includes(e.type)
    );
    assert.expect(unsupported.length);

    for (const engine of unsupported) {
      await consoleComponent.runCommands([
        // delete any previous mount with same name
        `delete sys/mounts/${engine.type}`,
      ]);
      await mountSecrets.visit();
      await mountBackend(engine.type, engine.type);

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
    await mountBackend('kv', v2);
    assert.strictEqual(currentURL(), `/vault/secrets-engines/${v2}/kv/list`, `${v2} navigates to list url`);
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
    await click(GENERAL.cardContainer('kv'));
    await fillIn(GENERAL.inputByAttr('path'), v1);
    await click(GENERAL.button('Method Options'));
    await mountSecrets.version(1);
    await click(GENERAL.submitButton);

    assert.strictEqual(currentURL(), `/vault/secrets-engines/${v1}/list`, `${v1} navigates to list url`);
    assert.strictEqual(
      currentRouteName(),
      `vault.cluster.secrets.backend.list-root`,
      `${v1} navigates to list route`
    );
  });

  // Condensed tests for these specific engines here as they just check if they are added to the list after mounting
  test('enable alicloud', async function (assert) {
    const enginePath = `alicloud-${this.uid}`;
    await mountSecrets.visit();
    await mountBackend('alicloud', enginePath);

    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backends',
      'redirects to the backends page'
    );
    await fillIn(GENERAL.inputSearch('secret-engine-path'), enginePath);
    assert.dom(GENERAL.listItem(`${enginePath}/`)).exists();

    // cleanup
    await runCmd(`delete sys/mounts/${enginePath}`);
  });

  test('enable gcpkms', async function (assert) {
    const enginePath = `gcpkms-${this.uid}`;
    await mountSecrets.visit();
    await mountBackend('gcpkms', enginePath);

    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backends',
      'redirects to the backends page'
    );
    await fillIn(GENERAL.inputSearch('secret-engine-path'), enginePath);
    assert.dom(GENERAL.listItem(`${enginePath}/`)).exists();
    // cleanup
    await runCmd(`delete sys/mounts/${enginePath}`);
  });

  module('WIF secret engines', function () {
    test('it sets identity_token_key on mount config using search select list, resets after', async function (assert) {
      // create an oidc/key
      await runCmd(`write identity/oidc/key/some-key allowed_client_ids="*"`);

      await page.visit();
      await click(GENERAL.cardContainer('aws')); // only testing aws of the WIF engines as the functionality for all others WIF engines in this form are the same
      await click(GENERAL.button('Method Options'));
      assert.dom('[data-test-search-select-with-modal]').exists('Search select with modal component renders');
      await clickTrigger('#key');
      const dropdownOptions = findAll('[data-option-index]').map((o) => o.innerText);
      assert.ok(dropdownOptions.includes('some-key'), 'search select options show some-key');
      await click(GENERAL.searchSelect.option(GENERAL.searchSelect.optionIndex('some-key')));
      assert
        .dom(GENERAL.searchSelect.selectedOption())
        .hasText('some-key', 'some-key was selected and displays in the search select');
      await click(GENERAL.backButton);
      // Choose a non-wif engine
      await click(GENERAL.cardContainer('ssh'));
      assert
        .dom('[data-test-search-select-with-modal]')
        .doesNotExist('for type ssh, the modal field does not render.');
      // cleanup
      await runCmd(`delete identity/oidc/key/some-key`);
    });

    // TODO: Revisit these two OIDC tests, the config form is rendering but should have an info row value? not sure why
    skip('it allows a user with permissions to oidc/key to create an identity_token_key', async function (assert) {
      const engine = 'aws'; // only testing aws of the WIF engines as the functionality for all others WIF engines in this form are the same
      await login();
      const path = `secrets-adminPolicy-${engine}`;
      const newKey = `key-${engine}-${uuidv4()}`;
      const secrets_admin_policy = adminOidcCreateRead(path);
      const secretsAdminToken = await runCmd(
        tokenWithPolicyCmd(`secrets-admin-${path}`, secrets_admin_policy)
      );

      await login(secretsAdminToken);
      await visit('/vault/secrets-engines/enable');
      await click(GENERAL.cardContainer(engine));
      await fillIn(GENERAL.inputByAttr('path'), path);
      await click(GENERAL.button('Method Options'));
      await clickTrigger('#key');
      // create new key
      await fillIn(GENERAL.searchSelect.searchInput, newKey);
      await click(GENERAL.searchSelect.options);
      assert.dom('#search-select-modal').exists(`modal with form opens for engine ${engine}`);
      assert
        .dom('[data-test-modal-title]')
        .hasText('Create new key', `Create key modal renders for engine: ${engine}`);

      await click(OIDC.keySaveButton);
      assert.dom('#search-select-modal').doesNotExist(`modal disappears onSave for engine ${engine}`);
      assert.dom(GENERAL.searchSelect.selectedOption()).hasText(newKey, `${newKey} is now selected`);

      await click(GENERAL.submitButton);
      await visit(`/vault/secrets-engines/${path}/configuration`);
      await click(GENERAL.tab('plugin-settings'));

      assert
        .dom(GENERAL.infoRowValue('Identity token key'))
        .hasText(newKey, `shows identity token key on configuration page for engine: ${engine}`);

      // cleanup
      await runCmd(`delete sys/mounts/${path}`);
      await runCmd(`delete identity/oidc/key/some-key`);
      await runCmd(`delete identity/oidc/key/${newKey}`);
    });

    skip('it allows user with NO access to oidc/key to manually input an identity_token_key', async function (assert) {
      const engine = 'aws'; // only testing aws of the WIF engines as the functionality for all others WIF engines in this form are the same
      await login();
      const path = `secrets-noOidcAdmin-${engine}`;
      const secretsNoOidcAdminPolicy = adminOidcCreate(path);
      const secretsNoOidcAdminToken = await runCmd(
        tokenWithPolicyCmd(`secrets-noOidcAdmin-${path}`, secretsNoOidcAdminPolicy)
      );
      // create an oidc/key that they can then use even if they can't read it.
      await runCmd(`write identity/oidc/key/general-key allowed_client_ids="*"`);

      await login(secretsNoOidcAdminToken);
      await page.visit();
      await click(GENERAL.cardContainer(engine));
      await fillIn(GENERAL.inputByAttr('path'), path);
      await click(GENERAL.button('Method Options'));
      // type-in fallback component to create new key
      await typeIn(GENERAL.inputSearch('key'), 'general-key');
      await click(GENERAL.submitButton);
      assert
        .dom(GENERAL.latestFlashContent)
        .hasText(`Successfully mounted the ${engine} secrets engine at ${path}.`);

      await visit(`/vault/secrets-engines/${path}/configuration`);

      await click(SES.configurationToggle);
      assert
        .dom(GENERAL.infoRowValue('Identity token key'))
        .hasText('general-key', `shows identity token key on configuration page for engine: ${engine}`);

      // cleanup
      await runCmd(`delete sys/mounts/${path}`);
    });
  });
});
