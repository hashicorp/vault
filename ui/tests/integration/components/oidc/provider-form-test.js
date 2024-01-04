/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn, click, findAll } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import ENV from 'vault/config/environment';
import {
  SELECTORS,
  OIDC_BASE_URL,
  CLIENT_LIST_RESPONSE,
  overrideMirageResponse,
  overrideCapabilities,
} from 'vault/tests/helpers/oidc-config';
import parseURL from 'core/utils/parse-url';
import { setRunOptions } from 'ember-a11y-testing/test-support';

const ISSUER_URL = 'http://127.0.0.1:8200/v1/identity/oidc/provider/test-provider';

module('Integration | Component | oidc/provider-form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.before(function () {
    ENV['ember-cli-mirage'].handler = 'oidcConfig';
  });

  hooks.after(function () {
    ENV['ember-cli-mirage'].handler = null;
  });

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.server.get('/identity/oidc/scope', () => {
      return {
        request_id: 'scope-list-id',
        lease_id: '',
        renewable: false,
        lease_duration: 0,
        data: {
          keys: ['test-scope'],
        },
        wrap_info: null,
        warnings: null,
        auth: null,
      };
    });
    this.server.get('/identity/oidc/client', () => overrideMirageResponse(null, CLIENT_LIST_RESPONSE));
    setRunOptions({
      rules: {
        // TODO: Fix SearchSelect component
        'aria-required-attr': { enabled: false },
        label: { enabled: false },
        // TODO: fix RadioCard component (replace with HDS)
        'aria-valid-attr-value': { enabled: false },
        'nested-interactive': { enabled: false },
      },
    });
  });

  test('it should save new provider', async function (assert) {
    assert.expect(13);
    this.server.post('/identity/oidc/provider/test-provider', (schema, req) => {
      assert.ok(true, 'Request made to save provider');
      return JSON.parse(req.requestBody);
    });
    this.model = this.store.createRecord('oidc/provider');
    this.onSave = () => assert.ok(true, 'onSave callback fires on save success');
    await render(hbs`
    <Oidc::ProviderForm
    @model={{this.model}}
    @onCancel={{this.onCancel}}
    @onSave={{this.onSave}}
    />
    `);

    assert
      .dom('[data-test-oidc-provider-title]')
      .hasText('Create Provider', 'Form title renders correct text');
    assert.dom(SELECTORS.providerSaveButton).hasText('Create', 'Save button has correct text');
    assert
      .dom('[data-test-input="issuer"]')
      .hasAttribute('placeholder', 'e.g. https://example.com:8200', 'issuer placeholder text is correct');
    assert.strictEqual(findAll('[data-test-field]').length, 3, 'renders all input fields');
    await click('[data-test-component="search-select"]#scopesSupported .ember-basic-dropdown-trigger');
    assert.dom('li.ember-power-select-option').hasText('test-scope', 'dropdown renders scopes');

    // check validation errors
    await fillIn('[data-test-input="name"]', ' ');
    await click(SELECTORS.providerSaveButton);

    const validationErrors = findAll(SELECTORS.inlineAlert);
    assert
      .dom(validationErrors[0])
      .hasText('Name is required. Name cannot contain whitespace.', 'Validation messages are shown for name');
    assert.dom(validationErrors[1]).hasText('There are 2 errors with this form.', 'Renders form error count');

    await click('[data-test-oidc-radio="limited"]');
    assert
      .dom('[data-test-component="search-select"]#allowedClientIds')
      .exists('Limited radio button shows clients search select');
    await click('[data-test-component="search-select"]#allowedClientIds .ember-basic-dropdown-trigger');
    assert.dom('li.ember-power-select-option').hasTextContaining('test-app', 'dropdown renders client name');
    assert.dom('[data-test-smaller-id]').exists('renders smaller client id in dropdown');

    await click('[data-test-oidc-radio="allow-all"]');
    assert
      .dom('[data-test-component="search-select"]#allowedClientIds')
      .doesNotExist('Allow all radio button hides search select');

    await fillIn('[data-test-input="name"]', 'test-provider');
    await click(SELECTORS.providerSaveButton);
  });

  test('it should update provider', async function (assert) {
    assert.expect(9);

    this.server.post('/identity/oidc/provider/test-provider', (schema, req) => {
      assert.ok(true, 'Request made to save provider');
      return JSON.parse(req.requestBody);
    });

    this.store.pushPayload('oidc/provider', {
      modelName: 'oidc/provider',
      name: 'test-provider',
      allowed_client_ids: ['*'],
      issuer: ISSUER_URL,
      scopes_supported: ['test-scope'],
    });

    this.model = this.store.peekRecord('oidc/provider', 'test-provider');
    this.onSave = () => assert.ok(true, 'onSave callback fires on save success');

    await render(hbs`
      <Oidc::ProviderForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);

    assert.dom('[data-test-oidc-provider-title]').hasText('Edit Provider', 'Title renders correct text');
    assert.dom(SELECTORS.providerSaveButton).hasText('Update', 'Save button has correct text');
    assert.dom('[data-test-input="name"]').isDisabled('Name input is disabled when editing');
    assert
      .dom('[data-test-input="name"]')
      .hasValue('test-provider', 'Name input is populated with model value');
    assert
      .dom('[data-test-input="issuer"]')
      .hasValue(parseURL(ISSUER_URL).origin, 'issuer value is just scheme://host:port portion of full URL');

    assert.dom('[data-test-selected-option]').hasText('test-scope', 'model scope is selected');
    assert.dom('[data-test-oidc-radio="allow-all"] input').isChecked('Allow all radio button is selected');
    await click(SELECTORS.providerSaveButton);
  });

  test('it should rollback attributes or unload record on cancel', async function (assert) {
    assert.expect(4);
    this.model = this.store.createRecord('oidc/provider');
    this.onCancel = () => assert.ok(true, 'onCancel callback fires');

    await render(hbs`
      <Oidc::ProviderForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);

    await click(SELECTORS.providerCancelButton);
    assert.true(this.model.isDestroyed, 'New model is unloaded on cancel');

    this.store.pushPayload('oidc/provider', {
      modelName: 'oidc/provider',
      name: 'test-provider',
      allowed_client_ids: ['*'],
      issuer: ISSUER_URL,
      scopes_supported: ['test-scope'],
    });

    this.model = this.store.peekRecord('oidc/provider', 'test-provider');

    await render(hbs`
      <Oidc::ProviderForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);

    await click('[data-test-oidc-radio="limited"]');
    await click(SELECTORS.providerCancelButton);
    assert.strictEqual(this.model.allowed_client_ids, undefined, 'Model attributes rolled back on cancel');
  });

  test('it should render fallback for search select', async function (assert) {
    assert.expect(2);
    this.model = this.store.createRecord('oidc/provider');
    this.server.get('/identity/oidc/scope', () => overrideMirageResponse(403));
    this.server.get('/identity/oidc/client', () => overrideMirageResponse(403));
    await render(hbs`
      <Oidc::ProviderForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);

    assert
      .dom('[data-test-component="search-select"]#scopesSupported [data-test-component="string-list"]')
      .exists('renders fall back for scopes search select');
    await click('[data-test-oidc-radio="limited"]');
    assert
      .dom('[data-test-component="search-select"]#allowedClientIds [data-test-component="string-list"]')
      .exists('Radio toggle shows assignments string-list input');
  });

  test('it should render error alerts when API returns an error', async function (assert) {
    assert.expect(2);
    this.model = this.store.createRecord('oidc/provider');
    this.server.post('/sys/capabilities-self', () => overrideCapabilities(OIDC_BASE_URL + '/providers'));
    await render(hbs`
      <Oidc::ProviderForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);
    await fillIn('[data-test-input="name"]', 'some-provider');
    await click(SELECTORS.providerSaveButton);
    assert
      .dom(SELECTORS.inlineAlert)
      .hasText('There was an error submitting this form.', 'form error alert renders ');
    assert.dom('[data-test-message-error]').exists('alert banner renders');
  });
});
