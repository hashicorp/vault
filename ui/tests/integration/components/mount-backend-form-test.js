/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, render } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { setupRenderingTest } from 'ember-qunit';
import { module, test } from 'qunit';
import { mountBackend } from 'vault/tests/helpers/components/mount-backend-form-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { allowAllCapabilitiesStub, noopStub } from 'vault/tests/helpers/stubs';
import { filterEnginesByMountCategory } from 'vault/utils/all-engines-metadata';

import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import AuthMethodForm from 'vault/forms/auth/method';
import SecretsEngineForm from 'vault/forms/secrets/engine';

module('Integration | Component | mount backend form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.flashMessages = this.owner.lookup('service:flash-messages');
    this.flashMessages.registerTypes(['success', 'danger']);
    this.flashSuccessSpy = sinon.spy(this.flashMessages, 'success');
    this.store = this.owner.lookup('service:store');
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    this.server.post('/sys/auth/foo', noopStub());
    this.onMountSuccess = sinon.spy();
  });

  module('auth method', function (hooks) {
    hooks.beforeEach(function () {
      const defaults = {
        config: { listing_visibility: false },
      };
      this.model = new AuthMethodForm(defaults, { isNew: true });
    });

    test('it renders default state', async function (assert) {
      await render(
        hbs`<MountBackendForm @mountModel={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );
      assert
        .dom(GENERAL.hdsPageHeaderTitle)
        .hasText('Enable an Authentication Method', 'renders auth header in default state');

      for (const method of filterEnginesByMountCategory({
        mountCategory: 'auth',
        isEnterprise: false,
      }).filter((engine) => engine.type !== 'token')) {
        assert
          .dom(GENERAL.cardContainer(method.type))
          .hasText(method.displayName, `renders type:${method.displayName} picker`);
      }
    });

    test('it changes path when type is changed', async function (assert) {
      await render(
        hbs`<MountBackendForm @mountModel={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );

      await click(GENERAL.cardContainer('aws'));
      assert.dom(GENERAL.inputByAttr('path')).hasValue('aws', 'sets the value of the type');
      await click(GENERAL.backButton);
      await click(GENERAL.cardContainer('approle'));
      assert.dom(GENERAL.inputByAttr('path')).hasValue('approle', 'updates the value of the type');
    });

    test('it keeps path value if the user has changed it', async function (assert) {
      await render(
        hbs`<MountBackendForm @mountModel={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );
      await click(GENERAL.cardContainer('approle'));
      assert.strictEqual(this.model.type, 'approle', 'Updates type on model');
      assert.dom(GENERAL.inputByAttr('path')).hasValue('approle', 'defaults to approle (first in the list)');
      await fillIn(GENERAL.inputByAttr('path'), 'newpath');
      assert.strictEqual(this.model.path, 'newpath', 'Updates path on model');
      await click(GENERAL.backButton);
      assert.strictEqual(this.model.type, '', 'Clears type on back');
      assert.strictEqual(this.model.path, 'newpath', 'Path is still newPath');
      await click(GENERAL.cardContainer('aws'));
      assert.strictEqual(this.model.type, 'aws', 'Updates type on model');
      assert.dom(GENERAL.inputByAttr('path')).hasValue('newpath', 'keeps custom path value');
    });

    test('it does not show a selected token type when first mounting an auth method', async function (assert) {
      await render(
        hbs`<MountBackendForm @mountModel={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );
      await click(GENERAL.cardContainer('github'));
      await click(GENERAL.button('Method Options'));
      assert
        .dom(GENERAL.inputByAttr('config.token_type'))
        .hasValue('', 'token type does not have a default value.');
      const selectOptions = document.querySelector(GENERAL.inputByAttr('config.token_type')).options;
      assert.strictEqual(selectOptions[1].text, 'default-service', 'first option is default-service');
      assert.strictEqual(selectOptions[2].text, 'default-batch', 'second option is default-batch');
      assert.strictEqual(selectOptions[3].text, 'batch', 'third option is batch');
      assert.strictEqual(selectOptions[4].text, 'service', 'fourth option is service');
    });

    test('it calls mount success', async function (assert) {
      this.server.post('/sys/auth/foo', () => {
        assert.ok(true, 'it calls enable on an auth method');
        return [204, { 'Content-Type': 'application/json' }];
      });
      const spy = sinon.spy();
      this.set('onMountSuccess', spy);

      await render(
        hbs`<MountBackendForm @mountModel={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );
      await mountBackend('approle', 'foo');

      assert.true(spy.calledOnce, 'calls the passed success method');
      assert.true(
        this.flashSuccessSpy.calledWith('Successfully mounted the approle auth method at foo.'),
        'Renders correct flash message'
      );
    });
  });

  module('Plugin Version Selection Integration (Community)', function (hooks) {
    hooks.beforeEach(function () {
      // Get version service for mocking in individual tests
      this.version = this.owner.lookup('service:version');

      // Mock plugin pins API endpoint
      this.server.get('/sys/plugins/pins', () => {
        return {
          data: {
            pinned_versions: [],
          },
        };
      });

      // Set up secrets engine form with KV type already selected
      const defaults = {
        config: {},
        kv_config: {
          max_versions: 0,
          cas_required: false,
          delete_version_after: 0,
        },
        options: { version: 2 },
      };
      this.form = new SecretsEngineForm(defaults, { isNew: true });
      this.form.type = 'kv'; // Pre-select KV type to skip type selection
      this.form.data.path = 'test-path'; // Set a test path

      // Mock mount success handler
      this.onMountSuccess = sinon.spy();

      // Mock available versions (as would be passed from route)
      this.availableVersions = [
        {
          version: 'v1.16.1+builtin',
          pluginName: 'vault-plugin-secrets-kv',
          isBuiltin: true,
          sha256: 'abc123',
        },
        {
          version: 'v0.25.0',
          pluginName: 'vault-plugin-secrets-kv',
          isBuiltin: false,
          sha256: 'def456',
        },
      ];

      this.renderComponent = () =>
        render(hbs`
          <Mount::SecretsEngineForm
            @model={{this.model}}
            @onMountSuccess={{this.onMountSuccess}}
          />
        `);

      // Helper function to create a fresh model with new available versions
      this.createFreshModel = (type = 'kv', availableVersions = this.availableVersions) => {
        const defaults = {
          config: {},
          kv_config: {
            max_versions: 0,
            cas_required: false,
            delete_version_after: 0,
          },
          options: { version: 2 },
        };
        this.form = new SecretsEngineForm(defaults, { isNew: true });
        this.form.type = type;
        this.form.data.path = 'test-path';

        // Update the model structure with new data
        this.model = {
          form: this.form,
          availableVersions: availableVersions,
          hasUnversionedPlugins: false,
          pinnedVersion: null, // Will be set per test as needed
        };
      };

      // Initialize with default model
      this.createFreshModel();
    });

    test('plugin version field is hidden when only builtin versions available', async function (assert) {
      // Mock enterprise mode (even with enterprise, no external versions means no version field)
      this.version.type = 'enterprise';

      // Mock single builtin version response
      this.availableVersions = [
        {
          version: 'v1.16.1+builtin',
          pluginName: 'vault-plugin-secrets-kv',
          isBuiltin: true,
          sha256: 'abc123',
        },
      ];

      // Create a fresh model for this test with only builtin versions
      this.createFreshModel('kv', this.availableVersions);

      await this.renderComponent();

      // External radio card should be disabled when no external versions are available
      assert
        .dom(`input${GENERAL.radioCardByAttr('external')}`)
        .isDisabled('external radio card is disabled when only builtin versions available');

      // With only builtin versions, external radio should be disabled and version field hidden
      assert
        .dom(GENERAL.fieldByAttr('config.plugin_version'))
        .doesNotExist('plugin version field is hidden when only builtin versions available');

      // Check for user messaging about why external option is disabled
      assert
        .dom(GENERAL.inlineAlert)
        .exists('info message explains why external option is disabled when no external versions available');
    });

    test('external radio card is disabled for community version', async function (assert) {
      // Mock version service to simulate community mode
      this.version.type = 'community';

      await this.renderComponent();

      // External radio card should be disabled in community mode
      assert
        .dom(`input${GENERAL.radioCardByAttr('external')}`)
        .isDisabled('external radio card is disabled for community version');

      // Enterprise badge should be visible for community users
      assert
        .dom(GENERAL.badge('external-enterprise'))
        .hasText('Enterprise', 'Enterprise badge is shown for community users');

      // Plugin version field should not be visible
      assert
        .dom(GENERAL.fieldByAttr('config.plugin_version'))
        .doesNotExist('plugin version field is hidden when external card is disabled');
    });

    module('Plugin Version Selection Integration (Ent)', function (hooks) {
      hooks.beforeEach(function () {
        // Set enterprise mode for all tests in this module
        this.version.type = 'enterprise';
      });

      test('plugin version field shows when multiple versions available', async function (assert) {
        await this.renderComponent();

        // External radio card should not be disabled in enterprise mode
        assert
          .dom(`input${GENERAL.radioCardByAttr('external')}`)
          .isNotDisabled('external radio card is enabled for enterprise version');

        // Initially, version field should not be visible (builtin is default)
        assert
          .dom(GENERAL.fieldByAttr('config.plugin_version'))
          .doesNotExist('plugin version field is hidden initially with builtin selection');

        // Click External radio card to enable version selection
        await click(`input${GENERAL.radioCardByAttr('external')}`);

        // Now version field should appear
        assert
          .dom(GENERAL.fieldByAttr('config.plugin_version'))
          .exists('plugin version field appears when external is selected');

        // HDS Select component generates a different DOM structure
        assert
          .dom(`${GENERAL.fieldByAttr('config.plugin_version')} .hds-form-select`)
          .exists('plugin version field uses HDS select component');
      });

      test('selecting default version omits plugin_version from payload', async function (assert) {
        // Mock successful mount request
        this.server.post('/sys/mounts/test-path', (schema, request) => {
          const payload = JSON.parse(request.requestBody);
          const hasPluginVersion =
            Object.prototype.hasOwnProperty.call(payload, 'config') &&
            Object.prototype.hasOwnProperty.call(payload.config, 'plugin_version');
          assert.notOk(
            hasPluginVersion,
            'plugin_version is not included in payload when default is selected'
          );
          assert.strictEqual(payload.type, 'kv', 'correct engine type is sent');
          return {};
        });

        await this.renderComponent();

        // Builtin is default selection (no plugin version field visible), so just submit
        await click(GENERAL.submitButton);
      });

      test('builtin plugin type selected by default sends correct payload', async function (assert) {
        // Mock successful mount request
        this.server.post('/sys/mounts/test-path', (schema, request) => {
          const payload = JSON.parse(request.requestBody);
          // With builtin selected (default), no plugin_version should be sent
          const hasPluginVersion =
            Object.prototype.hasOwnProperty.call(payload, 'config') &&
            Object.prototype.hasOwnProperty.call(payload.config, 'plugin_version');
          assert.notOk(hasPluginVersion, 'plugin_version is not included for builtin selection');
          assert.strictEqual(payload.type, 'kv', 'type remains builtin type for builtin plugins');
          return {};
        });

        await this.renderComponent();

        // Builtin is selected by default, no version field shown, just submit
        await click(GENERAL.submitButton);
      });

      test('selecting external version includes plugin_version in payload with external type', async function (assert) {
        // Mock successful mount request
        this.server.post('/sys/mounts/test-path', (schema, request) => {
          const payload = JSON.parse(request.requestBody);
          assert.strictEqual(
            payload.config.plugin_version,
            'v0.25.0',
            'plugin_version is included for external version'
          );
          assert.strictEqual(
            payload.type,
            'vault-plugin-secrets-kv',
            'type is external plugin name for external plugins'
          );
          return {};
        });

        await this.renderComponent();

        // Click External radio card to enable version selection
        await click(`input${GENERAL.radioCardByAttr('external')}`);

        // Select the external version from the dropdown
        await fillIn(GENERAL.selectByAttr('plugin-version'), 'v0.25.0');

        await click(GENERAL.submitButton);
      });

      test('external radio card is enabled but version field is hidden when plugin has empty version', async function (assert) {
        // Create availableVersions with a plugin that has an empty version (registered without version)
        this.availableVersions = [
          {
            version: '', // Empty version when plugin registered without version
            pluginName: 'vault-plugin-secrets-keymgmt',
            isBuiltin: false,
          },
        ];

        // Create a fresh model for this test with empty version plugins
        this.createFreshModel('keymgmt', this.availableVersions);

        await this.renderComponent();

        // External radio card should NOT be disabled
        assert
          .dom(`input${GENERAL.radioCardByAttr('external')}`)
          .isNotDisabled('external radio card is enabled when plugin has empty version');

        // Click External radio card to enable version selection
        await click(`input${GENERAL.radioCardByAttr('external')}`);

        // Version field should NOT appear since only empty versions exist
        assert
          .dom(GENERAL.fieldByAttr('config.plugin_version'))
          .doesNotExist('plugin version field is hidden when only empty version plugins exist');
      });

      test('version field shows with filtered options when plugin has both empty and non-empty versions', async function (assert) {
        // Create availableVersions with both empty and non-empty versions
        this.availableVersions = [
          {
            version: 'v1.16.1+builtin',
            pluginName: 'vault-plugin-secrets-keymgmt',
            isBuiltin: true,
          },
          {
            version: '', // Empty version (should be filtered out)
            pluginName: 'vault-plugin-secrets-keymgmt',
            isBuiltin: false,
          },
          {
            version: 'v1.0.0', // Non-empty external version
            pluginName: 'vault-plugin-secrets-keymgmt',
            isBuiltin: false,
          },
        ];

        // Create a fresh model for this test with mixed versions
        this.createFreshModel('keymgmt', this.availableVersions);

        await this.renderComponent();

        // External radio card should NOT be disabled
        assert
          .dom(`input${GENERAL.radioCardByAttr('external')}`)
          .isNotDisabled('external radio card is enabled when plugin has valid external versions');

        // Click External radio card to enable version selection
        await click(`input${GENERAL.radioCardByAttr('external')}`);

        // Version field should appear since we have non-empty external versions
        assert
          .dom(GENERAL.fieldByAttr('config.plugin_version'))
          .exists('plugin version field appears when non-empty external versions exist');

        assert
          .dom(`${GENERAL.fieldByAttr('config.plugin_version')} option[value="v1.0.0"]`)
          .exists('non-empty external version option exists');
      });
    });
  });
});
