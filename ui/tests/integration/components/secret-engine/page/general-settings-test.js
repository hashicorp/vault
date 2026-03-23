/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const keyManagementMockModel = {
  secretsEngine: {
    accessor: 'keymgmt_accessor',
    config: {
      default_lease_ttl: 2073600,
      force_no_cache: false,
      listing_visibility: 'hidden',
      max_lease_ttl: 4320000,
    },
    description: 'hello',
    external_entropy_access: false,
    local: true,
    options: {},
    path: 'keymgmt/',
    plugin_version: '',
    running_plugin_version: 'v0.17.1+builtin',
    running_sha256: '',
    seal_wrap: false,
    type: 'keymgmt',
    uuid: '4ea92618-5b52-f89a-9cbe-b65dc7e65689',
    id: 'keymgmt',
    backendConfigurationLink: `vault.cluster.secrets.backend.configuration`,
  },
  pinnedVersion: null,
  versions: ['v0.17.1+builtin'],
};

module('Integration | Component | SecretEngine::Page::GeneralSettings', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.model = keyManagementMockModel;
    this.breadcrumbs = [
      { label: 'Secrets', route: 'vault.cluster.secrets' },
      {
        label: this.model.secretsEngine.id,
        route: 'vault.cluster.secrets.backend.list-root',
        model: this.model.secretsEngine.id,
      },
      { label: 'Configuration' },
    ];

    this.server.get('/sys/internal/ui/mounts/:path', () => {
      return {
        data: {
          plugin_version: '',
          running_plugin_version: this.model.secretsEngine.running_plugin_version,
          config: {
            override_pinned_version: false,
          },
        },
      };
    });
  });

  test('it shows general settings form', async function (assert) {
    assert.expect(4);

    await render(hbs`
      <SecretEngine::Page::GeneralSettings @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />
    `);
    assert.dom(GENERAL.cardContainer('secrets duration')).exists(`Lease duration card exists`);
    assert.dom(GENERAL.cardContainer('security')).exists(`Security card exists`);
    assert.dom(GENERAL.cardContainer('version')).exists(`Version card exists`);
    assert.dom(GENERAL.cardContainer('metadata')).exists(`Metadata card exists`);
  });

  test('it sends override_pinned_version=true when selecting version different from pinned', async function (assert) {
    assert.expect(1);

    // Set up model with multiple versions and pinned version
    this.model.versions = ['v0.16.0+builtin', 'v0.17.1+builtin', 'v0.18.0+builtin'];
    this.model.pinnedVersion = 'v0.16.0+builtin';

    let tuneRequest = null;
    this.server.post('/sys/mounts/:path/tune', (schema, request) => {
      tuneRequest = JSON.parse(request.requestBody);
      return {};
    });

    await render(hbs`
      <SecretEngine::Page::GeneralSettings @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />
    `);
    // Wait for version data to load

    // Select a version different from pinned version
    await fillIn(GENERAL.inputByAttr('plugin-version'), 'v0.18.0+builtin');
    await click(GENERAL.submitButton);

    assert.true(
      tuneRequest?.override_pinned_version,
      'Should send override_pinned_version=true when selecting different version'
    );
  });

  test('it sends override_pinned_version=false and excludes plugin_version when selecting the pinned version', async function (assert) {
    assert.expect(2);

    // Set up model with multiple versions
    this.model.versions = ['v0.16.0+builtin', 'v0.17.1+builtin', 'v0.18.0+builtin'];
    this.model.pinnedVersion = 'v0.16.0+builtin';

    let tuneRequest = null;
    this.server.post('/sys/mounts/:path/tune', (schema, request) => {
      tuneRequest = JSON.parse(request.requestBody);
      return {};
    });

    await render(hbs`
      <SecretEngine::Page::GeneralSettings @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />
    `);

    // Select the pinned version (should have "(Pinned)" label)
    await fillIn(GENERAL.inputByAttr('plugin-version'), 'v0.16.0+builtin (Pinned)');
    await click(GENERAL.submitButton);

    /*
     * This assertion is intentional and prevents a regression.
     *
     * This test verifies that override_pinned_version is explicitly sent in the request.
     * If the client API methods are used instead of raw request, they omit this parameter
     * from the request body.
     *
     * DO NOT change this to assert.false() - it must be strictEqual to catch empty payloads!
     */
    // eslint-disable-next-line qunit/no-assert-equal-boolean
    assert.strictEqual(
      tuneRequest?.override_pinned_version,
      false,
      'Should send override_pinned_version=false when selecting pinned version'
    );

    assert.false(
      'plugin_version' in tuneRequest,
      'Should exclude plugin_version when selecting pinned version'
    );
  });

  test('it does not send override_pinned_version when no pinned version exists', async function (assert) {
    assert.expect(1);

    // Set up model with multiple versions
    this.model.versions = ['v0.17.1+builtin', 'v0.18.0+builtin'];
    this.model.pinnedVersion = null;

    let tuneRequest = null;
    this.server.post('/sys/mounts/:path/tune', (schema, request) => {
      tuneRequest = JSON.parse(request.requestBody);
      return {};
    });

    await render(hbs`
      <SecretEngine::Page::GeneralSettings @model={{this.model}} @breadcrumbs={{this.breadcrumbs}} />
    `);
    // Wait for version data to load

    // Select a version when no pinned version exists
    await fillIn(GENERAL.inputByAttr('plugin-version'), 'v0.18.0+builtin');
    await click(GENERAL.submitButton);

    assert.false(
      'override_pinned_version' in tuneRequest,
      'Should not send override_pinned_version flag when no pinned version exists'
    );
  });
});
