import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn, click, findAll } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import ENV from 'vault/config/environment';
import { overrideMirageResponse } from 'vault/tests/helpers/oidc-config';

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
  });

  test('it should save new provider', async function (assert) {
    assert.expect(11);
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
      .hasText('Create provider', 'Form title renders correct text');
    assert.dom('[data-test-oidc-provider-save]').hasText('Create', 'Save button has correct text');
    assert.equal(findAll('[data-test-field]').length, 3, 'renders all input fields');
    await click('[data-test-component="search-select"]#scopesSupported .ember-basic-dropdown-trigger');
    assert.dom('li.ember-power-select-option').hasText('test-scope', 'dropdown renders scopes');

    // check validation errors
    await click('[data-test-oidc-provider-save]');
    assert
      .dom('[data-test-inline-error-message]')
      .hasText('Name is required.', 'Validation message is shown for name');
    await fillIn('[data-test-input="name"]', 'test space');
    await click('[data-test-oidc-provider-save]');
    assert
      .dom('[data-test-inline-error-message]')
      .hasText('Name cannot contain whitespace.', 'Validation message is shown whitespace');

    await click('label[for=limited]');
    assert
      .dom('[data-test-component="search-select"]#allowedClientIds')
      .exists('Limited radio button shows clients search select');
    await click('[data-test-component="search-select"]#allowedClientIds .ember-basic-dropdown-trigger');
    assert.dom('li.ember-power-select-option').hasText('some-app', 'dropdown renders clients');

    await click('label[for=allow-all]');
    assert
      .dom('[data-test-component="search-select"]#allowedClientIds')
      .doesNotExist('Allow all radio button hides search select');

    await fillIn('[data-test-input="name"]', 'test-provider');
    await click('[data-test-oidc-provider-save]');
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

    assert.dom('[data-test-oidc-provider-title]').hasText('Edit provider', 'Title renders correct text');
    assert.dom('[data-test-oidc-provider-save]').hasText('Update', 'Save button has correct text');
    assert.dom('[data-test-input="name"]').isDisabled('Name input is disabled when editing');
    assert
      .dom('[data-test-input="name"]')
      .hasValue('test-provider', 'Name input is populated with model value');
    assert
      .dom('[data-test-input="issuer"]')
      .hasValue(ISSUER_URL, 'issuer url input is populated with model value');
    assert.dom('[data-test-selected-option="true"]').hasText('test-scope', 'model scope is selected');
    assert.dom('input#allow-all').isChecked('Allow all radio button is selected');
    await click('[data-test-oidc-provider-save]');
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

    await click('[data-test-oidc-provider-cancel]');
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

    await click('label[for=limited]');
    await click('[data-test-oidc-provider-cancel]');
    assert.equal(this.model.allowed_client_ids, undefined, 'Model attributes rolled back on cancel');
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
    await click('label[for=limited]');
    assert
      .dom('[data-test-component="search-select"]#allowedClientIds [data-test-component="string-list"]')
      .exists('Radio toggle shows assignments string-list input');
  });
});
