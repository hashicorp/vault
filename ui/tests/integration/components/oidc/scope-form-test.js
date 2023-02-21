import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn, click, findAll } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { SELECTORS, OIDC_BASE_URL, overrideCapabilities } from 'vault/tests/helpers/oidc-config';

module('Integration | Component | oidc/scope-form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
  });

  test('it should save new scope', async function (assert) {
    assert.expect(9);

    this.server.post('/identity/oidc/scope/test', (schema, req) => {
      assert.ok(true, 'Request made to save scope');
      return JSON.parse(req.requestBody);
    });

    this.model = this.store.createRecord('oidc/scope');
    this.onSave = () => assert.ok(true, 'onSave callback fires on save success');

    await render(hbs`
      <Oidc::ScopeForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
      <div id="modal-wormhole"></div>
    `);

    assert.dom('[data-test-oidc-scope-title]').hasText('Create scope', 'Form title renders');
    assert.dom(SELECTORS.scopeSaveButton).hasText('Create', 'Save button has correct label');
    await click(SELECTORS.scopeSaveButton);

    // check validation errors
    await click(SELECTORS.scopeSaveButton);

    const validationErrors = findAll(SELECTORS.inlineAlert);
    assert.dom(validationErrors[0]).hasText('Name is required.', 'Validation messages are shown for name');
    assert.dom(validationErrors[1]).hasText('There is an error with this form.', 'Renders form error count');

    assert
      .dom('[data-test-inline-error-message]')
      .hasText('Name is required.', 'Validation message is shown for name');
    // json editor has test coverage so let's just confirm that it renders
    assert
      .dom('[data-test-input="template"] [data-test-component="json-editor-toolbar"]')
      .exists('JsonEditor toolbar renders');
    assert
      .dom('[data-test-input="template"] [data-test-component="code-mirror-modifier"]')
      .exists('Code mirror renders');

    await fillIn('[data-test-input="name"]', 'test');
    await fillIn('[data-test-input="description"]', 'this is a test');
    await click(SELECTORS.scopeSaveButton);
  });

  test('it should update scope', async function (assert) {
    assert.expect(9);

    this.server.post('/identity/oidc/scope/test', (schema, req) => {
      assert.ok(true, 'Request made to save scope');
      return JSON.parse(req.requestBody);
    });

    this.store.pushPayload('oidc/scope', {
      modelName: 'oidc/scope',
      name: 'test',
      description: 'this is a test',
    });
    this.model = this.store.peekRecord('oidc/scope', 'test');
    this.onSave = () => assert.ok(true, 'onSave callback fires on save success');

    await render(hbs`
      <Oidc::ScopeForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
      <div id="modal-wormhole"></div>
    `);

    assert.dom('[data-test-oidc-scope-title]').hasText('Edit scope', 'Form title renders');
    assert.dom(SELECTORS.scopeSaveButton).hasText('Update', 'Save button has correct label');
    assert.dom('[data-test-input="name"]').isDisabled('Name input is disabled when editing');
    assert.dom('[data-test-input="name"]').hasValue('test', 'Name input is populated with model value');
    assert
      .dom('[data-test-input="description"]')
      .hasValue('this is a test', 'Description input is populated with model value');
    // json editor has test coverage so let's just confirm that it renders
    assert
      .dom('[data-test-input="template"] [data-test-component="json-editor-toolbar"]')
      .exists('JsonEditor toolbar renders');
    assert
      .dom('[data-test-input="template"] [data-test-component="code-mirror-modifier"]')
      .exists('Code mirror renders');

    await fillIn('[data-test-input="description"]', 'this is an edit test');
    await click(SELECTORS.scopeSaveButton);
  });

  test('it should rollback attributes or unload record on cancel', async function (assert) {
    assert.expect(4);

    this.onCancel = () => assert.ok(true, 'onCancel callback fires');

    this.model = this.store.createRecord('oidc/scope');

    await render(hbs`
      <Oidc::ScopeForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
      <div id="modal-wormhole"></div>
    `);

    await click(SELECTORS.scopeCancelButton);
    assert.true(this.model.isDestroyed, 'New model is unloaded on cancel');

    this.store.pushPayload('oidc/scope', {
      modelName: 'oidc/scope',
      name: 'test',
      description: 'this is a test',
    });
    this.model = this.store.peekRecord('oidc/scope', 'test');

    await render(hbs`
    <Oidc::ScopeForm
      @model={{this.model}}
      @onCancel={{this.onCancel}}
      @onSave={{this.onSave}}
    />
    <div id="modal-wormhole"></div>
  `);

    await fillIn('[data-test-input="description"]', 'changed description attribute');
    await click(SELECTORS.scopeCancelButton);
    assert.strictEqual(
      this.model.description,
      'this is a test',
      'Model attributes are rolled back on cancel'
    );
  });

  test('it should show example template modal', async function (assert) {
    assert.expect(6);

    this.model = this.store.createRecord('oidc/scope');

    await render(hbs`
      <Oidc::ScopeForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
      <div id="modal-wormhole"></div>
    `);

    assert.dom('[data-test-modal-div]').doesNotHaveClass('is-active', 'Modal is hidden');
    await click('[data-test-oidc-scope-example]');
    assert.dom('[data-test-modal-div]').hasClass('is-active', 'Modal is shown');
    assert.dom('[data-test-modal-title]').hasText('Scope template', 'Modal title renders');
    assert
      .dom('[data-test-modal-copy]')
      .hasText('Example of a JSON template for scopes:', 'Modal copy renders');
    assert.dom('.cm-string').hasText('"username"', 'Example template json renders');
    await click('[data-test-close-modal]');
    assert.dom('[data-test-modal-div]').doesNotHaveClass('is-active', 'Modal is hidden');
  });

  test('it should render error alerts when API returns an error', async function (assert) {
    assert.expect(2);
    this.model = this.store.createRecord('oidc/scope');
    this.server.post('/sys/capabilities-self', () => overrideCapabilities(OIDC_BASE_URL + '/scopes'));
    await render(hbs`
      <Oidc::ScopeForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
      <div id="modal-wormhole"></div>
    `);
    await fillIn('[data-test-input="name"]', 'test-scope');
    await click(SELECTORS.scopeSaveButton);
    assert
      .dom(SELECTORS.inlineAlert)
      .hasText('There was an error submitting this form.', 'form error alert renders ');
    assert.dom('[data-test-alert-banner="alert"]').exists('alert banner renders');
  });
});
