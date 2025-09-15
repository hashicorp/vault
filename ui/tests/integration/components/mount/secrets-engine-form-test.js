/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, typeIn } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { allowAllCapabilitiesStub, noopStub } from 'vault/tests/helpers/stubs';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
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
    this.flashMessages.registerTypes(['success', 'danger']);
    this.flashSuccessSpy = sinon.spy(this.flashMessages, 'success');
    this.store = this.owner.lookup('service:store');
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
    this.model = new SecretsEngineForm(defaults, { isNew: true });
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
    this.model.type = 'azure';
    this.model.data.path = 'azure'; // Set path to match type as would happen in the route
    await render(
      hbs`<Mount::SecretsEngineForm @model={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
    );
    assert.dom(GENERAL.inputByAttr('path')).hasValue('azure', 'path matches type');
  });

  test('it keeps custom path value', async function (assert) {
    this.model.type = 'kv';
    this.model.data.path = 'custom-path';
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

    this.model.type = 'ssh';
    this.model.data.path = 'foo';

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

  module('KV engine', function () {
    test('it shows KV specific fields when type is kv', async function (assert) {
      this.model.type = 'kv';
      await render(
        hbs`<Mount::SecretsEngineForm @model={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );
      assert.dom(GENERAL.inputByAttr('kv_config.max_versions')).exists('shows max versions field');
      assert.dom(GENERAL.inputByAttr('kv_config.cas_required')).exists('shows CAS required field');
      assert.dom(GENERAL.inputByAttr('kv_config.delete_version_after')).exists('shows delete after field');
    });
  });

  module('WIF secret engines', function () {
    test('it shows identity_token_key when type is a WIF engine and hides when its not', async function (assert) {
      // Test AWS (a WIF engine)
      this.model.type = 'aws';
      this.model.applyTypeSpecificDefaults();

      // Initialize config object for WIF engines
      if (!this.model.data.config) {
        this.model.data.config = {};
      }

      await render(
        hbs`<Mount::SecretsEngineForm @model={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );

      // First check if the Method Options group is being rendered at all
      assert.dom('[data-test-button="Method Options"]').exists('Method Options toggle button exists');

      // Click to expand Method Options if it's collapsed
      await click('[data-test-button="Method Options"]');

      assert
        .dom(GENERAL.fieldByAttr('config.identity_token_key'))
        .exists('Identity token key field shows for AWS engine');

      // Test KV (not a WIF engine)
      this.model.type = 'kv';
      this.model.applyTypeSpecificDefaults();

      await render(
        hbs`<Mount::SecretsEngineForm @model={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );

      assert
        .dom(GENERAL.fieldByAttr('config.identity_token_key'))
        .doesNotExist('Identity token key field hidden for KV engine');
    });

    test('it updates identity_token_key if user has changed it', async function (assert) {
      this.model.type = WIF_ENGINES[0]; // Use first WIF engine
      this.model.applyTypeSpecificDefaults();
      // Initialize config object
      if (!this.model.data.config) {
        this.model.data.config = {};
      }
      await render(
        hbs`<Mount::SecretsEngineForm @model={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );

      // Expand Method Options section to show identity_token_key field
      await click(GENERAL.button('Method Options'));

      assert.strictEqual(
        this.model.data.config.identity_token_key,
        undefined,
        'On init identity_token_key is not set on the model'
      );

      // SearchSelectWithModal likely uses fallback component when no OIDC models are found
      await typeIn(GENERAL.inputSearch('key'), 'specialKey');

      assert.strictEqual(
        this.model.data.config.identity_token_key,
        'specialKey',
        'updates model with custom identity_token_key'
      );
    });
  });

  module('PKI engine', function () {
    test('it sets default max lease TTL for PKI', async function (assert) {
      this.model.type = 'pki';
      this.model.applyTypeSpecificDefaults();

      assert.strictEqual(
        this.model.data.config.max_lease_ttl,
        '3650d',
        'sets PKI default max lease TTL to 10 years'
      );
    });
  });
});
