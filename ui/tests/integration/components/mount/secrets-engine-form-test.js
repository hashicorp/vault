/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, render, typeIn } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupRenderingTest } from 'ember-qunit';
import { module, test } from 'qunit';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import {
  allowAllCapabilitiesStub,
  capabilitiesStub,
  noopStub,
  overrideResponse,
} from 'vault/tests/helpers/stubs';
import { ALL_ENGINES } from 'vault/utils/all-engines-metadata';

import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import SecretsEngineForm from 'vault/forms/secrets/engine';

const WIF_ENGINES = ALL_ENGINES.filter((e) => e.isWIF).map((e) => e.type);

module('Integration | Component | mount/secrets-engine-form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.flashMessages = this.owner.lookup('service:flash-messages');
    this.flashSuccessSpy = sinon.spy(this.flashMessages, 'success');
    this.flashWarningSpy = sinon.spy(this.flashMessages, 'warning');
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    this.server.post('/sys/mounts/foo', noopStub());
    this.onMountSuccess = sinon.spy();

    const defaults = {
      config: { listing_visibility: false },
      kv_config: {
        max_versions: 0,
        cas_required: false,
        delete_version_after: 0,
      },
      options: { version: 2 },
    };
    this.form = new SecretsEngineForm(defaults, { isNew: true });

    this.model = {
      form: this.form,
      availableVersions: [],
      hasUnversionedPlugins: false,
    };
  });

  test('it renders secret engine form', async function (assert) {
    await render(
      hbs`<Mount::SecretsEngineForm @model={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
    );
    assert.dom(GENERAL.breadcrumbs).exists('renders breadcrumbs');
    assert.dom(GENERAL.submitButton).hasText('Enable engine', 'renders submit button');
    assert.dom(GENERAL.backButton).hasText('Back', 'renders back button');
  });

  test('it changes path when type is set', async function (assert) {
    this.form.type = 'azure';
    this.form.data.path = 'azure'; // Set path to match type as would happen in the route
    await render(
      hbs`<Mount::SecretsEngineForm @model={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
    );
    assert.dom(GENERAL.inputByAttr('path')).hasValue('azure', 'path matches type');
  });

  test('it keeps custom path value', async function (assert) {
    this.form.type = 'kv';
    this.form.data.path = 'custom-path';
    await render(
      hbs`<Mount::SecretsEngineForm @model={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
    );
    assert.dom(GENERAL.inputByAttr('path')).hasValue('custom-path', 'keeps custom path');
  });

  test('it calls mount success', async function (assert) {
    assert.expect(3);

    this.server.post('/sys/mounts/foo', () => {
      assert.ok(true, 'it calls enable on a secrets engine');
      return [204, { 'Content-Type': 'application/json' }];
    });
    const spy = sinon.spy();
    this.set('onMountSuccess', spy);

    this.form.type = 'ssh';
    this.form.data.path = 'foo';

    await render(
      hbs`<Mount::SecretsEngineForm @model={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
    );

    await click(GENERAL.submitButton);

    assert.true(spy.calledOnce, 'calls the passed success method');
    assert.true(
      this.flashSuccessSpy.calledWith('Successfully mounted the ssh secrets engine at foo.'),
      'Renders correct flash message'
    );
  });

  module('KV engine', function (hooks) {
    hooks.beforeEach(function () {
      this.form.type = 'kv';
    });

    test('it shows KV specific fields when type is kv', async function (assert) {
      await render(
        hbs`<Mount::SecretsEngineForm @model={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );
      assert.dom(GENERAL.inputByAttr('kv_config.max_versions')).exists('shows max versions field');
      assert.dom(GENERAL.inputByAttr('kv_config.cas_required')).exists('shows CAS required field');
      assert.dom(GENERAL.inputByAttr('kv_config.delete_version_after')).exists('shows delete after field');
    });

    test('version 2 with no update to config endpoint still allows mount of secret engine', async function (assert) {
      assert.expect(6);
      this.server.post('/sys/capabilities-self', () => capabilitiesStub('my-kv-engine/config', ['deny']));
      this.server.post('/sys/mounts/my-kv-engine', (schema, req) => {
        assert.true(true, 'it makes request to mount engine');
        const payload = JSON.parse(req.requestBody);
        const expected = {
          config: { listing_visibility: 'hidden', force_no_cache: false },
          options: { version: 2 },
          type: 'kv',
        };
        assert.propEqual(payload, expected, 'mount request has expected payload');
        return overrideResponse(204);
      });

      await render(
        hbs`<Mount::SecretsEngineForm @model={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );
      await fillIn(GENERAL.inputByAttr('path'), 'my-kv-engine');
      await fillIn(GENERAL.inputByAttr('kv_config.max_versions'), '101');
      await click(GENERAL.submitButton);
      const [message] = this.flashWarningSpy.lastCall.args;
      assert.strictEqual(
        message,
        `You do not have access to the config endpoint. The secret engine was mounted, but the configuration settings were not saved.`,
        'it calls warning flash with expected message'
      );
      const [type, enginePath, useEngineRoute] = this.onMountSuccess.lastCall.args;
      assert.strictEqual(type, 'kv', 'onMountSuccess called with expected type');
      assert.strictEqual(enginePath, 'my-kv-engine', 'onMountSuccess called with expected engine path');
      assert.true(useEngineRoute, 'onMountSuccess called useEngineRoute: true');
    });
  });

  module('WIF secret engines', function () {
    test('it shows identity_token_key when type is a WIF engine and hides when its not', async function (assert) {
      // Test AWS (a WIF engine)
      this.form.type = 'aws';
      this.form.applyTypeSpecificDefaults();

      // Initialize config object for WIF engines
      if (!this.form.data.config) {
        this.form.data.config = {};
      }

      await render(
        hbs`<Mount::SecretsEngineForm @model={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );

      // First check if the Method Options group is being rendered at all
      assert.dom(GENERAL.button('Method Options')).exists('Method Options toggle button exists');

      // Click to expand Method Options if it's collapsed
      await click(GENERAL.button('Method Options'));

      assert
        .dom(GENERAL.fieldByAttr('config.identity_token_key'))
        .exists('Identity token key field shows for AWS engine');

      // Test KV (not a WIF engine)
      this.form.type = 'kv';
      this.form.applyTypeSpecificDefaults();

      await render(
        hbs`<Mount::SecretsEngineForm @model={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );

      assert
        .dom(GENERAL.fieldByAttr('config.identity_token_key'))
        .doesNotExist('Identity token key field hidden for KV engine');
    });

    test('it updates identity_token_key if user has changed it', async function (assert) {
      this.form.type = WIF_ENGINES[0]; // Use first WIF engine
      this.form.applyTypeSpecificDefaults();
      // Initialize config object
      if (!this.form.data.config) {
        this.form.data.config = {};
      }
      await render(
        hbs`<Mount::SecretsEngineForm @model={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );

      // Expand Method Options section to show identity_token_key field
      await click(GENERAL.button('Method Options'));

      assert.strictEqual(
        this.form.data.config.identity_token_key,
        undefined,
        'On init identity_token_key is not set on the model'
      );

      // SearchSelectWithModal likely uses fallback component when no OIDC models are found
      await typeIn(GENERAL.inputSearch('key'), 'specialKey');

      assert.strictEqual(
        this.form.data.config.identity_token_key,
        'specialKey',
        'updates model with custom identity_token_key'
      );
    });
  });

  module('PKI engine', function () {
    test('it sets default max lease TTL for PKI', async function (assert) {
      this.form.type = 'pki';
      this.form.applyTypeSpecificDefaults();

      assert.strictEqual(
        this.form.data.config.max_lease_ttl,
        '3650d',
        'sets PKI default max lease TTL to 10 years'
      );
    });
  });

  module('Plugin registration and versioning', function (hooks) {
    hooks.beforeEach(function () {
      this.form.type = 'keymgmt';
      this.form.data.path = 'keymgmt';

      // Mock version service for enterprise checks
      this.versionService = this.owner.lookup('service:version');
      sinon.stub(this.versionService, 'isEnterprise').value(true);

      // Setup available versions for testing and add to model structure
      const availableVersions = [
        { version: '1.0.0', pluginName: 'vault-plugin-secrets-keymgmt', isBuiltin: false },
        { version: '1.1.0', pluginName: 'vault-plugin-secrets-keymgmt', isBuiltin: false },
        { version: '2.0.0', pluginName: 'vault-plugin-secrets-keymgmt', isBuiltin: false },
        { version: '', pluginName: 'keymgmt', isBuiltin: true }, // Built-in version
      ];
      this.availableVersions = availableVersions;
      this.model.availableVersions = availableVersions;
    });

    test('it renders plugin type selection radio cards', async function (assert) {
      await render(
        hbs`<Mount::SecretsEngineForm @model={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );

      assert.dom(`input${GENERAL.radioCardByAttr('builtin')}`).exists('shows built-in plugin radio card');
      assert.dom(`input${GENERAL.radioCardByAttr('external')}`).exists('shows external plugin radio card');
      assert
        .dom(`input${GENERAL.radioCardByAttr('builtin')}`)
        .isChecked('built-in plugin is selected by default');
      assert
        .dom(`input${GENERAL.radioCardByAttr('external')}`)
        .isNotChecked('external plugin is not selected by default');
    });

    test('it defaults to built-in plugin type', async function (assert) {
      await render(
        hbs`<Mount::SecretsEngineForm 
          @model={{this.model}} 
          @onMountSuccess={{this.onMountSuccess}} 
        />`
      );

      assert
        .dom(GENERAL.fieldByAttr('config.plugin_version'))
        .doesNotExist('plugin version field is hidden for built-in');
      assert.strictEqual(this.form.type, 'keymgmt', 'model type remains as built-in name');
    });

    test('it shows plugin version field when external plugin is selected', async function (assert) {
      await render(
        hbs`<Mount::SecretsEngineForm 
          @model={{this.model}} 
          @onMountSuccess={{this.onMountSuccess}} 
        />`
      );

      await click(`input${GENERAL.radioCardByAttr('external')}`);

      assert
        .dom(GENERAL.fieldByAttr('config.plugin_version'))
        .exists('plugin version field appears for external');
      assert.dom(GENERAL.selectByAttr('plugin-version')).exists('plugin version select is rendered');
      assert.strictEqual(
        this.form.type,
        'vault-plugin-secrets-keymgmt',
        'model type updates to external plugin name'
      );
    });

    test('it populates version dropdown with sorted options', async function (assert) {
      await render(
        hbs`<Mount::SecretsEngineForm 
          @model={{this.model}} 
          @onMountSuccess={{this.onMountSuccess}} 
        />`
      );

      await click(`input${GENERAL.radioCardByAttr('external')}`);

      // Note: version option selectors may need custom data-test attributes
      assert.dom('[data-test-version-option="2.0.0"]').exists('includes version 2.0.0');
      assert.dom('[data-test-version-option="1.1.0"]').exists('includes version 1.1.0');
      assert.dom('[data-test-version-option="1.0.0"]').exists('includes version 1.0.0');
    });

    test('it disables external plugin when no enterprise license', async function (assert) {
      this.versionService.isEnterprise = false;

      await render(
        hbs`<Mount::SecretsEngineForm 
          @model={{this.model}} 
          @onMountSuccess={{this.onMountSuccess}} 
        />`
      );

      assert
        .dom(`input${GENERAL.radioCardByAttr('external')}`)
        .isDisabled('external plugin is disabled without enterprise');
      assert.dom('.hds-badge').hasText('Enterprise', 'shows enterprise badge');
    });

    test('it disables external plugin when no external versions available', async function (assert) {
      this.model.availableVersions = [{ version: '', pluginName: 'keymgmt', isBuiltin: true }];

      await render(
        hbs`<Mount::SecretsEngineForm 
          @model={{this.model}} 
          @onMountSuccess={{this.onMountSuccess}} 
        />`
      );

      assert
        .dom(`input${GENERAL.radioCardByAttr('external')}`)
        .isDisabled('external plugin is disabled when no external versions');
    });

    test('it updates plugin version when selection changes', async function (assert) {
      await render(
        hbs`<Mount::SecretsEngineForm 
          @model={{this.model}} 
          @onMountSuccess={{this.onMountSuccess}} 
        />`
      );

      await click(`input${GENERAL.radioCardByAttr('external')}`);
      await fillIn(GENERAL.selectByAttr('plugin-version'), '1.0.0');

      assert.strictEqual(
        this.form.data.config.plugin_version,
        '1.0.0',
        'updates model config with selected version'
      );
    });

    test('it clears plugin version when switching back to built-in', async function (assert) {
      await render(
        hbs`<Mount::SecretsEngineForm 
          @model={{this.model}} 
          @onMountSuccess={{this.onMountSuccess}} 
        />`
      );

      // Select external and set version
      await click(`input${GENERAL.radioCardByAttr('external')}`);
      await fillIn(GENERAL.selectByAttr('plugin-version'), '1.0.0');

      // Switch back to built-in
      await click(`input${GENERAL.radioCardByAttr('builtin')}`);

      assert.strictEqual(this.form.data.config.plugin_version, '', 'clears plugin version for built-in');
      assert.strictEqual(this.form.type, 'keymgmt', 'resets model type to built-in name');
      assert.dom(GENERAL.fieldByAttr('config.plugin_version')).doesNotExist('hides plugin version field');
    });

    test('it shows unversioned plugins warning when hasUnversionedPlugins is true', async function (assert) {
      this.model.hasUnversionedPlugins = true;

      await render(
        hbs`<Mount::SecretsEngineForm 
          @model={{this.model}} 
          @onMountSuccess={{this.onMountSuccess}} 
        />`
      );

      await click(`input${GENERAL.radioCardByAttr('external')}`);

      assert
        .dom(GENERAL.helpTextByAttr('config.plugin_version'))
        .containsText(
          'Un-versioned plugins are not supported, they must be enabled via CLI',
          'shows unversioned plugins warning'
        );
    });

    test('it hides unversioned plugins warning when hasUnversionedPlugins is false', async function (assert) {
      this.model.hasUnversionedPlugins = false;

      await render(
        hbs`<Mount::SecretsEngineForm 
          @model={{this.model}} 
          @onMountSuccess={{this.onMountSuccess}} 
        />`
      );

      await click(`input${GENERAL.radioCardByAttr('external')}`);

      assert
        .dom(GENERAL.fieldByAttr('config.plugin_version'))
        .doesNotContainText('Un-versioned plugins are not supported', 'hides unversioned plugins warning');
    });

    test('it hides unversioned plugins warning when hasUnversionedPlugins is not provided', async function (assert) {
      await render(
        hbs`<Mount::SecretsEngineForm 
          @model={{this.model}} 
          @onMountSuccess={{this.onMountSuccess}} 
        />`
      );

      await click(`input${GENERAL.radioCardByAttr('external')}`);

      assert
        .dom(GENERAL.fieldByAttr('config.plugin_version'))
        .doesNotContainText(
          'Un-versioned plugins are not supported',
          'hides unversioned plugins warning when property not provided'
        );
    });
  });

  module('Plugin pins integration', function (hooks) {
    hooks.beforeEach(function () {
      this.form.type = 'keymgmt';
      this.form.data.path = 'keymgmt';
      this.versionService = this.owner.lookup('service:version');
      sinon.stub(this.versionService, 'isEnterprise').value(true);

      const availableVersions = [
        { version: '1.0.0', pluginName: 'vault-plugin-secrets-keymgmt', isBuiltin: false },
        { version: '1.1.0', pluginName: 'vault-plugin-secrets-keymgmt', isBuiltin: false },
        { version: '2.0.0', pluginName: 'vault-plugin-secrets-keymgmt', isBuiltin: false },
      ];
      this.availableVersions = availableVersions;
      this.model.availableVersions = availableVersions;

      // Add pinned version to model data for tests
      this.model.pinnedVersion = '1.1.0';
    });

    test('it shows pinned version first in dropdown', async function (assert) {
      await render(
        hbs`<Mount::SecretsEngineForm 
          @model={{this.model}} 
          @onMountSuccess={{this.onMountSuccess}} 
        />`
      );

      await click(`input${GENERAL.radioCardByAttr('external')}`);

      // Check that pinned version is selected by default
      assert
        .dom(GENERAL.selectByAttr('plugin-version'))
        .hasValue('1.1.0', 'pinned version is selected by default');

      // Check pinned label appears
      assert
        .dom('[data-test-version-option="1.1.0"]')
        .hasText('1.1.0 (pinned)', 'shows pinned label in dropdown');
    });

    test('it shows pinned version in helper text', async function (assert) {
      await render(
        hbs`<Mount::SecretsEngineForm 
          @model={{this.model}} 
          @onMountSuccess={{this.onMountSuccess}} 
        />`
      );

      await click(`input${GENERAL.radioCardByAttr('external')}`);

      assert
        .dom(`${GENERAL.fieldByAttr('config.plugin_version')} .hds-form-helper-text`)
        .containsText('1.1.0 is pinned', 'shows pinned version in helper text');
    });

    test('it shows warning when selecting non-pinned version', async function (assert) {
      await render(
        hbs`<Mount::SecretsEngineForm 
          @model={{this.model}} 
          @onMountSuccess={{this.onMountSuccess}} 
        />`
      );

      await click(`input${GENERAL.radioCardByAttr('external')}`);
      // Wait for external plugin to be selected and version field to appear

      await fillIn(GENERAL.selectByAttr('plugin-version'), '2.0.0');
      // Wait for warning logic to process

      assert.dom('.hds-alert').exists('shows warning alert');
      assert
        .dom('.hds-alert .hds-alert__title')
        .hasText('Version differs from pinned', 'shows correct warning title');
      assert
        .dom('.hds-alert .hds-alert__description')
        .containsText(
          'You have selected 2.0.0, but version 1.1.0 is pinned',
          'shows correct warning description'
        );
    });

    test('it does not show warning when using pinned version', async function (assert) {
      await render(
        hbs`<Mount::SecretsEngineForm 
          @model={{this.model}} 
          @onMountSuccess={{this.onMountSuccess}} 
        />`
      );

      await click(`input${GENERAL.radioCardByAttr('external')}`);
      // Pinned version should be selected by default, no need to change

      assert.dom('.hds-alert--color-warning').doesNotExist('does not show warning when using pinned version');
    });

    test('it handles plugins with no pins correctly', async function (assert) {
      // Clear pinned version
      this.model.pinnedVersion = null;

      await render(
        hbs`<Mount::SecretsEngineForm 
          @model={{this.model}} 
          @onMountSuccess={{this.onMountSuccess}} 
        />`
      );

      await click(`input${GENERAL.radioCardByAttr('external')}`);

      // Should default to highest semantic version
      assert
        .dom(GENERAL.selectByAttr('plugin-version'))
        .hasValue('2.0.0', 'defaults to highest version when no pins');
      assert
        .dom(`${GENERAL.fieldByAttr('config.plugin_version')} .hds-form-helper-text`)
        .doesNotContainText('pinned', 'does not show pinned text when no pins');
      assert.dom('.hds-alert--color-warning').doesNotExist('does not show warning when no pins');
    });
  });

  module('Plugin version configuration handling', function (hooks) {
    hooks.beforeEach(function () {
      this.form.type = 'keymgmt';
      this.form.data.path = 'keymgmt';
      this.versionService = this.owner.lookup('service:version');
      sinon.stub(this.versionService, 'isEnterprise').value(true);

      const availableVersions = [
        { version: '1.0.0', pluginName: 'vault-plugin-secrets-keymgmt', isBuiltin: false },
        { version: '2.0.0', pluginName: 'vault-plugin-secrets-keymgmt', isBuiltin: false },
      ];
      this.availableVersions = availableVersions;
      this.model.availableVersions = availableVersions;

      // No pinned version for this test
      this.model.pinnedVersion = null;
    });

    test('it includes plugin_version in config for external plugins', async function (assert) {
      this.server.post('/sys/mounts/keymgmt', (schema, req) => {
        const payload = JSON.parse(req.requestBody);
        assert.strictEqual(
          payload.config.plugin_version,
          '2.0.0',
          'includes plugin_version in mount request'
        );
        assert.false(
          Object.hasOwn(payload.config, 'override_pinned_version'),
          'does not include override flag when no pins'
        );
        return [204, { 'Content-Type': 'application/json' }];
      });

      await render(
        hbs`<Mount::SecretsEngineForm 
          @model={{this.model}} 
          @onMountSuccess={{this.onMountSuccess}} 
        />`
      );

      await click(`input${GENERAL.radioCardByAttr('external')}`);
      await click(GENERAL.submitButton);
    });

    test('it includes override flag when using non-pinned version', async function (assert) {
      // Set pinned version for keymgmt plugin
      this.model.pinnedVersion = '1.0.0';

      this.server.post('/sys/mounts/keymgmt', (schema, req) => {
        const payload = JSON.parse(req.requestBody);
        assert.strictEqual(payload.config.plugin_version, '2.0.0', 'includes selected plugin_version');
        assert.true(
          payload.config.override_pinned_version,
          'includes override flag when using non-pinned version'
        );
        return [204, { 'Content-Type': 'application/json' }];
      });

      await render(
        hbs`<Mount::SecretsEngineForm 
          @model={{this.model}} 
          @onMountSuccess={{this.onMountSuccess}} 
        />`
      );

      await click(`input${GENERAL.radioCardByAttr('external')}`);
      await fillIn(GENERAL.selectByAttr('plugin-version'), '2.0.0');
      await click(GENERAL.submitButton);
    });

    test('it omits plugin_version when using pinned version', async function (assert) {
      // Set pinned version for keymgmt plugin
      this.model.pinnedVersion = '1.0.0';

      this.server.post('/sys/mounts/keymgmt', (schema, req) => {
        const payload = JSON.parse(req.requestBody);
        assert.false(
          Object.hasOwn(payload.config, 'plugin_version'),
          'omits plugin_version when using pinned version'
        );
        assert.false(
          Object.hasOwn(payload.config, 'override_pinned_version'),
          'omits override flag when using pinned version'
        );
        return [204, { 'Content-Type': 'application/json' }];
      });

      await render(
        hbs`<Mount::SecretsEngineForm 
          @model={{this.model}} 
          @onMountSuccess={{this.onMountSuccess}} 
        />`
      );

      await click(`input${GENERAL.radioCardByAttr('external')}`);
      // The pinned version (1.0.0) should be auto-selected
      await click(GENERAL.submitButton);
    });

    test('it does not include plugin_version for built-in plugins', async function (assert) {
      this.server.post('/sys/mounts/keymgmt', (schema, req) => {
        const payload = JSON.parse(req.requestBody);
        assert.false(
          Object.hasOwn(payload.config, 'plugin_version'),
          'does not include plugin_version for built-in'
        );
        assert.false(
          Object.hasOwn(payload.config, 'override_pinned_version'),
          'does not include override flag for built-in'
        );
        assert.strictEqual(payload.type, 'keymgmt', 'uses built-in type name');
        return [204, { 'Content-Type': 'application/json' }];
      });

      await render(
        hbs`<Mount::SecretsEngineForm 
          @model={{this.model}} 
          @onMountSuccess={{this.onMountSuccess}} 
        />`
      );

      // Built-in is selected by default
      await click(GENERAL.submitButton);
    });
  });

  module('Error handling and edge cases', function (hooks) {
    hooks.beforeEach(function () {
      this.form.type = 'keymgmt';
      this.form.data.path = 'keymgmt';
      this.versionService = this.owner.lookup('service:version');
      sinon.stub(this.versionService, 'isEnterprise').value(true);

      // No pinned version for error handling tests
      this.model.pinnedVersion = null;
    });

    test('it handles empty available versions gracefully', async function (assert) {
      this.model.availableVersions = [];

      await render(
        hbs`<Mount::SecretsEngineForm 
          @model={{this.model}} 
          @onMountSuccess={{this.onMountSuccess}} 
        />`
      );

      // External should be disabled
      assert
        .dom(`input${GENERAL.radioCardByAttr('external')}`)
        .isDisabled('external plugin disabled when no versions');
      assert
        .dom(GENERAL.fieldByAttr('config.plugin_version'))
        .doesNotExist('plugin version field hidden when no external versions');
    });

    test('it handles missing availableVersions argument', async function (assert) {
      await render(
        hbs`<Mount::SecretsEngineForm @model={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );

      // External should be disabled
      assert
        .dom(`input${GENERAL.radioCardByAttr('external')}`)
        .isDisabled('external plugin disabled when availableVersions not provided');
    });

    test('it shows version field immediately when pinned version available', async function (assert) {
      // Set pinned version
      this.model.pinnedVersion = '1.0.0';

      // Set up available versions for this test
      this.model.availableVersions = [
        { version: '1.0.0', pluginName: 'vault-plugin-secrets-keymgmt', isBuiltin: false },
      ];

      await render(
        hbs`<Mount::SecretsEngineForm 
          @model={{this.model}} 
          @onMountSuccess={{this.onMountSuccess}} 
        />`
      );

      await click(`input${GENERAL.radioCardByAttr('external')}`);

      // Version field should show even before pins are loaded, since versions are available
      assert
        .dom(GENERAL.fieldByAttr('config.plugin_version'))
        .exists('version field shows when external selected and versions available');

      // Field should remain visible since pinned version is available immediately
      assert
        .dom(GENERAL.fieldByAttr('config.plugin_version'))
        .exists('version field remains visible with pinned version');
    });
  });
});
