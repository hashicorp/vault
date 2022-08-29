import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn, click, findAll } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import ENV from 'vault/config/environment';
import { overrideMirageResponse, CLIENT_LIST_RESPONSE } from 'vault/tests/helpers/oidc-config';

module('Integration | Component | oidc/key-form', function (hooks) {
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
    this.server.get('/identity/oidc/client', () => overrideMirageResponse(null, CLIENT_LIST_RESPONSE));
  });

  test('it should save new key', async function (assert) {
    assert.expect(9);
    this.server.post('/identity/oidc/key/test-key', (schema, req) => {
      assert.ok(true, 'Request made to save key');
      return JSON.parse(req.requestBody);
    });
    this.model = this.store.createRecord('oidc/key');
    this.onSave = () => assert.ok(true, 'onSave callback fires on save success');
    await render(hbs`
    <Oidc::KeyForm
    @model={{this.model}}
    @onCancel={{this.onCancel}}
    @onSave={{this.onSave}}
    />
    `);

    assert.dom('[data-test-oidc-key-title]').hasText('Create key', 'Form title renders correct text');
    assert.dom('[data-test-oidc-key-save]').hasText('Create', 'Save button has correct text');
    assert.dom('[data-test-input="algorithm"]').hasValue('RS256', 'default algorithm is correct');
    assert.equal(findAll('[data-test-field]').length, 4, 'renders all input fields');

    // check validation errors
    await click('[data-test-oidc-key-save]');
    assert
      .dom('[data-test-inline-error-message]')
      .hasText('Name is required.', 'Validation message is shown for name');
    await fillIn('[data-test-input="name"]', 'test space');
    await click('[data-test-oidc-key-save]');
    assert
      .dom('[data-test-inline-error-message]')
      .hasText('Name cannot contain whitespace.', 'Validation message is shown whitespace');

    assert.dom('label[for=limited] input').isDisabled('limit radio button disabled on create');
    await fillIn('[data-test-input="name"]', 'test-key');
    await click('[data-test-oidc-key-save]');
  });

  test('it should update key and limit access to selected applications', async function (assert) {
    assert.expect(12);

    this.server.post('/identity/oidc/key/test-key', (schema, req) => {
      assert.ok(true, 'Request made to update key');
      return JSON.parse(req.requestBody);
    });

    this.store.pushPayload('oidc/key', {
      modelName: 'oidc/key',
      name: 'test-key',
      allowed_client_ids: ['*'],
    });

    this.model = this.store.peekRecord('oidc/key', 'test-key');
    this.onSave = () => assert.ok(true, 'onSave callback fires on update success');

    await render(hbs`
      <Oidc::KeyForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);

    assert.dom('[data-test-oidc-key-title]').hasText('Edit key', 'Title renders correct text');
    assert.dom('[data-test-oidc-key-save]').hasText('Update', 'Save button has correct text');
    assert.dom('[data-test-input="name"]').isDisabled('Name input is disabled when editing');
    assert.dom('[data-test-input="name"]').hasValue('test-key', 'Name input is populated with model value');
    assert.dom('input#allow-all').isChecked('Allow all radio button is selected');

    await click('label[for=limited]');
    assert
      .dom('[data-test-component="search-select"]#allowedClientIds')
      .exists('Limited radio button shows clients search select');
    await click('[data-test-component="search-select"]#allowedClientIds .ember-basic-dropdown-trigger');
    assert.equal(findAll('li.ember-power-select-option').length, 1, 'dropdown only renders one option');
    assert
      .dom('li.ember-power-select-option')
      .hasTextContaining('app-1', 'dropdown contains client that references key');
    assert.dom('[data-test-smaller-id]').exists('renders smaller client id in dropdown');

    await click('label[for=allow-all]');
    assert
      .dom('[data-test-component="search-select"]#allowedClientIds')
      .doesNotExist('Allow all radio button hides search select');

    await click('[data-test-oidc-key-save]');
  });

  test('it should rollback attributes or unload record on cancel', async function (assert) {
    assert.expect(4);
    this.model = this.store.createRecord('oidc/key');
    this.onCancel = () => assert.ok(true, 'onCancel callback fires');

    await render(hbs`
      <Oidc::KeyForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);

    await click('[data-test-oidc-key-cancel]');
    assert.true(this.model.isDestroyed, 'New model is unloaded on cancel');

    this.store.pushPayload('oidc/key', {
      modelName: 'oidc/key',
      name: 'test-key',
      allowed_client_ids: ['*'],
    });

    this.model = this.store.peekRecord('oidc/key', 'test-key');

    await render(hbs`
      <Oidc::KeyForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);

    await click('label[for=limited]');
    await click('[data-test-oidc-key-cancel]');
    assert.equal(this.model.allowed_client_ids, undefined, 'Model attributes rolled back on cancel');
  });

  test('it should render fallback for search select', async function (assert) {
    assert.expect(1);

    this.server.post('/identity/oidc/key/test-key', (schema, req) => {
      assert.ok(true, 'Request made to update key');
      return JSON.parse(req.requestBody);
    });

    this.store.pushPayload('oidc/key', {
      modelName: 'oidc/key',
      name: 'test-key',
      allowed_client_ids: ['*'],
    });

    this.model = this.store.peekRecord('oidc/key', 'test-key');

    this.server.get('/identity/oidc/client', () => overrideMirageResponse(403));
    await render(hbs`
      <Oidc::KeyForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);

    await click('label[for=limited]');
    assert
      .dom('[data-test-component="search-select"]#allowedClientIds [data-test-component="string-list"]')
      .exists('Radio toggle shows client string-list input');
  });
});
