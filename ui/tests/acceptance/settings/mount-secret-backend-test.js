/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { currentRouteName, currentURL, click, findAll, fillIn, visit, typeIn } from '@ember/test-helpers';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupAuth } from 'vault/tests/helpers/auth/setup-auth';
import sinon from 'sinon';

import page from 'vault/tests/pages/settings/mount-secret-backend';
import configPage from 'vault/tests/pages/secrets/backend/configuration';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import { MOUNT_BACKEND_FORM } from 'vault/tests/helpers/components/mount-backend-form-selectors';
import { mountBackend } from 'vault/tests/helpers/components/mount-backend-form-helpers';
import { SELECTORS as OIDC } from 'vault/tests/helpers/oidc-config';
import { filterEnginesByMountCategory } from 'vault/utils/all-engines-metadata';
import engineDisplayData from 'vault/helpers/engines-display-data';

// enterprise backends are tested separately
const BACKENDS_WITH_ENGINES = ['kv', 'pki', 'ldap', 'kubernetes'];

module('Acceptance | settings/mount-secret-backend', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);
  setupAuth(hooks);

  hooks.beforeEach(function () {
    this.calcDays = (hours) => {
      const days = Math.floor(hours / 24);
      const remainder = hours % 24;
      return `${days} days ${remainder} hours`;
    };
    this.kvPath = 'kv';
    this.mount = this.server.create('mount', 'isKv');

    this.api = this.owner.lookup('service:api');
    this.capabilities = this.owner.lookup('service:capabilities');

    this.mountStub = sinon.stub(this.api.sys, 'mountsEnableSecretsEngine').resolves();
    this.mountReadStub = sinon.stub(this.api.sys, 'internalUiReadMountInformation');
    this.mountReadStub.callsFake(() => this.mount);
    this.capabilitiesForStub = sinon.stub(this.capabilities, 'for');

    // handle requests but ignore responses since they are not needed in the scope of these tests
    ['metadata', 'config', 'roles', 'keys'].forEach((path) => {
      this.server.get(`/:path/${path}`, {}, 404);
    });
    this.server.get('/kv', {}, 404);

    // pki/config/acme uses hydrateModel which is not necessary for these tests
    sinon.stub(this.owner.lookup('service:pathHelp'), 'hydrateModel').resolves({});
  });

  test('it sets the ttl correctly when mounting', async function (assert) {
    const defaultTTLHours = 100;
    const maxTTLHours = 300;
    this.mount.config.default_lease_ttl = 360000;
    this.mount.config.max_lease_ttl = 1080000;

    await page.visit();

    assert.strictEqual(currentRouteName(), 'vault.cluster.settings.mount-secret-backend');
    await click(MOUNT_BACKEND_FORM.mountType('kv'));
    await fillIn(GENERAL.inputByAttr('path'), this.kvPath);
    await click(GENERAL.button('Method Options'));
    await click(GENERAL.toggleInput('Default Lease TTL'));
    await page.defaultTTLUnit('h').defaultTTLVal(defaultTTLHours);
    await click(GENERAL.toggleInput('Max Lease TTL'));
    await page.maxTTLUnit('h').maxTTLVal(maxTTLHours);
    await click(GENERAL.submitButton);
    await configPage.visit({ backend: this.kvPath });
    assert.strictEqual(configPage.defaultTTL, `${this.calcDays(defaultTTLHours)}`, 'shows the proper TTL');
    assert.strictEqual(configPage.maxTTL, `${this.calcDays(maxTTLHours)}`, 'shows the proper max TTL');
  });

  test('it sets the ttl when enabled then disabled', async function (assert) {
    const maxTTLHours = 300;
    this.mount.config.max_lease_ttl = 1080000;

    await page.visit();

    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.settings.mount-secret-backend',
      'navigates to mount page'
    );
    await click(MOUNT_BACKEND_FORM.mountType('kv'));
    await fillIn(GENERAL.inputByAttr('path'), this.kvPath);
    await click(GENERAL.button('Method Options'));
    await click(GENERAL.toggleInput('Default Lease TTL'));
    await click(GENERAL.toggleInput('Max Lease TTL'));
    await page.maxTTLUnit('h').maxTTLVal(maxTTLHours);
    await click(GENERAL.submitButton);
    await configPage.visit({ backend: this.kvPath });
    assert.strictEqual(configPage.defaultTTL, '1 month 1 day', 'shows system default TTL');
    assert.strictEqual(configPage.maxTTL, `${this.calcDays(maxTTLHours)}`, 'shows the proper max TTL');
  });

  test('it sets the max ttl after pki chosen, resets after', async function (assert) {
    await page.visit();
    assert.strictEqual(currentRouteName(), 'vault.cluster.settings.mount-secret-backend');
    await click(MOUNT_BACKEND_FORM.mountType('pki'));
    assert.dom('[data-test-input="config.max_lease_ttl"]').exists();
    assert
      .dom('[data-test-input="config.max_lease_ttl"] [data-test-ttl-toggle]')
      .isChecked('Toggle is checked by default');
    assert.dom('[data-test-input="config.max_lease_ttl"] [data-test-ttl-value]').hasValue('3650');
    assert.dom('[data-test-input="config.max_lease_ttl"] [data-test-select="ttl-unit"]').hasValue('d');

    // Go back and choose a different type
    await click(GENERAL.backButton);
    await click(MOUNT_BACKEND_FORM.mountType('database'));
    assert.dom('[data-test-input="config.max_lease_ttl"]').exists('3650');
    assert
      .dom('[data-test-input="config.max_lease_ttl"] [data-test-ttl-toggle]')
      .isNotChecked('Toggle is unchecked by default');
    await click(GENERAL.toggleInput('Max Lease TTL'));
    assert.dom('[data-test-input="config.max_lease_ttl"] [data-test-ttl-value]').hasValue('');
    assert.dom('[data-test-input="config.max_lease_ttl"] [data-test-select="ttl-unit"]').hasValue('s');
  });

  test('it throws error if setting duplicate path name', async function (assert) {
    const error = `path is already in use at ${this.kvPath}`;
    this.mountStub.rejects({ errors: [error] });

    await page.visit();

    assert.strictEqual(currentRouteName(), 'vault.cluster.settings.mount-secret-backend');
    await mountBackend('kv', this.kvPath);

    assert.strictEqual(currentRouteName(), 'vault.cluster.settings.mount-secret-backend');
    assert.dom('[data-test-message-error-description]').containsText(error);
  });

  test('version 2 with no update to config endpoint still allows mount of secret engine', async function (assert) {
    this.capabilitiesForStub.withArgs('kvConfig', { path: this.kvPath }).returns({ canUpdate: false });
    this.server.get(`/${this.kvPath}/metadata`, {}, 404);

    await mountSecrets.visit();
    await click(MOUNT_BACKEND_FORM.mountType('kv'));
    await fillIn(GENERAL.inputByAttr('path'), this.kvPath);
    await mountSecrets.setMaxVersion(101);
    await click(GENERAL.submitButton);

    assert
      .dom('[data-test-flash-message]')
      .containsText(
        `You do not have access to the config endpoint. The secret engine was mounted, but the configuration settings were not saved.`
      );
    assert.strictEqual(
      currentURL(),
      `/vault/secrets/${this.kvPath}/kv/list`,
      'After mounting, redirects to secrets list page'
    );
  });

  test('it should transition to mountable addon engine after mount success', async function (assert) {
    // test supported backends that ARE ember engines (enterprise only engines are tested individually)
    const addons = filterEnginesByMountCategory({ mountCategory: 'secret', isEnterprise: false }).filter(
      (e) => BACKENDS_WITH_ENGINES.includes(e.type)
    );
    assert.expect(addons.length);

    for (const engine of addons) {
      await mountSecrets.visit();
      await mountBackend(engine.type, engine.type);

      assert.strictEqual(
        currentRouteName(),
        `vault.cluster.secrets.backend.${engine.engineRoute}`,
        `Transitions to ${engine.displayName} route on mount success`
      );
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
      this.mount = this.server.create('mount', { type: engine.type });
      await mountSecrets.visit();
      await click(MOUNT_BACKEND_FORM.mountType(engine.type));
      await fillIn(GENERAL.inputByAttr('path'), engine.type);
      if (engine.type === 'kv') {
        await click(GENERAL.button('Method Options'));
        await mountSecrets.version(1);
      }
      await click(GENERAL.submitButton);

      const route = engineDisplayData(engine.type)?.isOnlyMountable ? 'configuration.index' : 'list-root';
      assert.strictEqual(
        currentRouteName(),
        `vault.cluster.secrets.backend.${route}`,
        `${engine.type} navigates to the correct view (either list if not configuration only or configuration if it is).`
      );
    }
  });

  test('it should transition back to backend list for unsupported backends', async function (assert) {
    const unsupported = filterEnginesByMountCategory({ mountCategory: 'secret', isEnterprise: false }).filter(
      (e) => !supportedSecretBackends().includes(e.type)
    );
    assert.expect(unsupported.length);

    for (const engine of unsupported) {
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
    await mountSecrets.visit();
    await mountBackend('kv', v2);
    assert.strictEqual(currentURL(), `/vault/secrets/${v2}/kv/list`, `${v2} navigates to list url`);
    assert.strictEqual(
      currentRouteName(),
      `vault.cluster.secrets.backend.kv.list`,
      `${v2} navigates to list url`
    );

    const v1 = 'kv';
    this.mount = this.server.create('mount', 'isKvV1');
    await mountSecrets.visit();
    await click(MOUNT_BACKEND_FORM.mountType('kv'));
    await fillIn(GENERAL.inputByAttr('path'), v1);
    await click(GENERAL.button('Method Options'));
    await mountSecrets.version(1);
    await click(GENERAL.submitButton);

    assert.strictEqual(currentURL(), `/vault/secrets/${v1}/list`, `${v1} navigates to list url`);
    assert.strictEqual(
      currentRouteName(),
      `vault.cluster.secrets.backend.list-root`,
      `${v1} navigates to list route`
    );
  });

  module('WIF secret engines', function (hooks) {
    hooks.beforeEach(function () {
      this.type = 'aws';
      this.newKey = 'test-oidc-key';
      this.keys = ['default'];

      this.server.get('/identity/oidc/key', () => {
        return {
          data: { keys: this.keys },
        };
      });
      this.server.post('/identity/oidc/key/:name', {}, 204);
      this.server.get('/aws/roles', {}, 404);

      this.mountReadStub.resolves(
        this.server.create('mount', 'isAws', { config: { identity_token_key: this.newKey } })
      );
      sinon.stub(this.api.secrets, 'awsReadRootIamCredentialsConfiguration').resolves({
        data: {
          identity_token_key: this.newKey,
          identity_token_ttl: 3600,
          role_arn: 'arn:aws:iam::123456789012:role/VaultRole',
        },
      });
      sinon.stub(this.api.secrets, 'awsReadLeaseConfiguration').resolves({ data: {} });
      sinon.stub(this.api.identity, 'oidcReadConfiguration').resolves({ data: {} });
    });

    test('it sets identity_token_key on mount config using search select list, resets after', async function (assert) {
      await page.visit();
      await click(MOUNT_BACKEND_FORM.mountType(this.type));
      await click(GENERAL.button('Method Options'));
      assert.dom('[data-test-search-select-with-modal]').exists('Search select with modal component renders');
      await clickTrigger('#key');
      const dropdownOptions = findAll('[data-option-index]').map((o) => o.innerText);
      assert.ok(dropdownOptions.includes('default'), 'search select options show default');
      await click(GENERAL.searchSelect.option(GENERAL.searchSelect.optionIndex('default')));
      assert
        .dom(GENERAL.searchSelect.selectedOption())
        .hasText('default', 'default was selected and displays in the search select');
      await click(GENERAL.backButton);
      // Choose a non-wif engine
      await click(MOUNT_BACKEND_FORM.mountType('ssh'));
      assert
        .dom('[data-test-search-select-with-modal]')
        .doesNotExist('for type ssh, the modal field does not render.');
    });

    test('it allows a user with permissions to oidc/key to create an identity_token_key', async function (assert) {
      await visit('/vault/settings/mount-secret-backend');
      await click(MOUNT_BACKEND_FORM.mountType(this.type));
      await fillIn(GENERAL.inputByAttr('path'), this.type);
      await click(GENERAL.button('Method Options'));
      await clickTrigger('#key');
      // create new key
      await fillIn(GENERAL.searchSelect.searchInput, this.newKey);
      await click(GENERAL.searchSelect.options);
      assert.dom('#search-select-modal').exists(`modal with form opens for engine ${this.type}`);
      assert
        .dom('[data-test-modal-title]')
        .hasText('Create new key', `Create key modal renders for engine: ${this.type}`);

      await click(OIDC.keySaveButton);
      assert.dom('#search-select-modal').doesNotExist(`modal disappears onSave for engine ${this.type}`);
      assert
        .dom(GENERAL.searchSelect.selectedOption())
        .hasText(this.newKey, `${this.newKey} is now selected`);

      await click(GENERAL.submitButton);
      await visit(`/vault/secrets/${this.type}/configuration`);
      await click(SES.configurationToggle);
      assert
        .dom(GENERAL.infoRowValue('Identity token key'))
        .hasText(this.newKey, `shows identity token key on configuration page for engine: ${this.type}`);
    });

    test('it allows user with NO access to oidc/key to manually input an identity_token_key', async function (assert) {
      this.server.get('/identity/oidc/key', {}, 403);

      await page.visit();
      await click(MOUNT_BACKEND_FORM.mountType(this.type));
      await fillIn(GENERAL.inputByAttr('path'), this.type);
      await click(GENERAL.button('Method Options'));
      // type-in fallback component to create new key
      await typeIn(GENERAL.inputSearch('key'), this.newKey);
      await click(GENERAL.submitButton);
      assert
        .dom(GENERAL.latestFlashContent)
        .hasText(`Successfully mounted the ${this.type} secrets engine at ${this.type}.`);

      await visit(`/vault/secrets/${this.type}/configuration`);

      await click(SES.configurationToggle);
      assert
        .dom(GENERAL.infoRowValue('Identity token key'))
        .hasText(this.newKey, `shows identity token key on configuration page for engine: ${this.type}`);
    });
  });
});
