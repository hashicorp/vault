import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn, click, findAll } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { create } from 'ember-cli-page-object';
import { clickTrigger } from 'ember-power-select/test-support/helpers';
import ss from 'vault/tests/pages/components/search-select';
import { setupMirage } from 'ember-cli-mirage/test-support';
import ENV from 'vault/config/environment';
import { overrideMirageResponse } from 'vault/tests/helpers/oidc-config';

const searchSelect = create(ss);

module('Integration | Component | oidc/client-form', function (hooks) {
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

  test('it should save new client', async function (assert) {
    assert.expect(15);

    this.server.post('/identity/oidc/client/some-app', (schema, req) => {
      assert.ok(true, 'Request made to save client');
      return JSON.parse(req.requestBody);
    });
    this.model = this.store.createRecord('oidc/client');
    this.onSave = () => assert.ok(true, 'onSave callback fires on save success');

    await render(hbs`
      <Oidc::ClientForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);
    assert
      .dom('[data-test-oidc-client-title]')
      .hasText('Create application', 'Form title renders correct text');
    assert.dom('[data-test-oidc-client-save]').hasText('Create', 'Save button has correct text');

    // check validation errors
    await click('[data-test-oidc-client-save]');
    assert
      .dom('[data-test-inline-error-message]')
      .hasText('Name is required.', 'Validation message is shown for name');
    await fillIn('[data-test-input="name"]', 'test space');
    await click('[data-test-oidc-client-save]');
    assert
      .dom('[data-test-inline-error-message]')
      .hasText('Name cannot contain whitespace.', 'Validation message is shown whitespace');

    await click('[data-test-toggle-group="More options"]');
    await click('label[for=limited]');
    assert
      .dom('[data-test-selected-option="true"]')
      .hasText('default', 'Search select has default key selected');
    assert
      .dom('[data-test-search-select-with-modal]')
      .exists('Limited radio button shows assignments search select');

    assert.equal(findAll('[data-test-field]').length, 6);

    await clickTrigger();
    assert.dom('li.ember-power-select-option').hasText('assignment-1', 'dropdown renders assignments');
    await fillIn('[data-test-input="name"]', 'some-app');
    await click('[data-test-oidc-client-save]');
  });

  test('it should update client', async function (assert) {
    assert.expect(9);

    this.server.post('/identity/oidc/client/some-app', (schema, req) => {
      assert.ok(true, 'Request made to save client');
      return JSON.parse(req.requestBody);
    });

    this.store.pushPayload('oidc/client', {
      modelName: 'oidc/client',
      name: 'some-app',
      clientType: 'public',
    });

    this.model = this.store.peekRecord('oidc/client', 'some-app');
    this.onSave = () => assert.ok(true, 'onSave callback fires on save success');

    await render(hbs`
      <Oidc::ClientForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);
    await click('[data-test-toggle-group="More options"]');
    assert.dom('[data-test-oidc-client-title]').hasText('Edit application', 'Title renders correct text');
    assert.dom('[data-test-oidc-client-save]').hasText('Update', 'Save button has correct text');
    assert.dom('[data-test-input="name"]').isDisabled('Name input is disabled when editing');
    assert.dom('[data-test-input="name"]').hasValue('some-app', 'Name input is populated with model value');
    assert.dom('[data-test-input="key"]').isDisabled('Signing key input is disabled');
    assert.dom('[data-test-input="key"]').hasValue('default', 'Key input populated with default');
    assert.dom('[data-test-input="clientType"]').isDisabled('client type input is disabled on edit');
    assert
      .dom('[data-test-input="clientType"] input#confidential')
      .isChecked('Correct radio button is selected');
    await click('[data-test-oidc-client-save]');
  });

  test('it should rollback attributes or unload record on cancel', async function (assert) {
    assert.expect(4);
    this.model = this.store.createRecord('oidc/client');
    this.onCancel = () => assert.ok(true, 'onCancel callback fires');

    await render(hbs`
      <Oidc::ClientForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);

    await click('[data-test-oidc-client-cancel]');
    assert.true(this.model.isDestroyed, 'New model is unloaded on cancel');

    this.store.pushPayload('oidc/client', {
      modelName: 'oidc/client',
      name: 'some-app',
      assignments: ['allow_all'],
    });
    this.model = this.store.peekRecord('oidc/client', 'some-app');

    await render(hbs`
      <Oidc::ClientForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);

    await click('[data-test-input="clientType"] input[id="public"]');
    await click('[data-test-oidc-client-cancel]');
    assert.equal(this.model.clientType, 'confidential', 'Model attributes rolled back on cancels');
  });

  test('it should show create assignment modal', async function (assert) {
    assert.expect(2);
    this.model = this.store.createRecord('oidc/client');

    await render(hbs`
      <Oidc::ClientForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
      <div id="modal-wormhole"></div>
    `);
    await click('label[for=limited]');
    await clickTrigger();
    await fillIn('.ember-power-select-search input', 'test-new');
    await searchSelect.options.objectAt(0).click();
    assert.dom('[data-test-modal-title]').hasText('Create new assignment', 'Create assignment modal renders');
    await click('[data-test-oidc-assignment-cancel]');
    assert.dom('[data-test-modal-div]').doesNotExist('Modal disappears after clicking cancel');
  });

  test('it should render fallback for search select', async function (assert) {
    assert.expect(1);
    this.model = this.store.createRecord('oidc/client');
    this.server.get('/identity/oidc/assignment', () => overrideMirageResponse(403));
    await render(hbs`
      <Oidc::ClientForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);

    await click('label[for=limited]');
    assert
      .dom('[data-test-component="string-list"]')
      .exists('Radio toggle shows assignments string-list input');
  });
});
