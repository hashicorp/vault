/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { SELECTORS, OIDC_BASE_URL } from 'vault/tests/helpers/oidc-config';
import { capabilitiesStub } from 'vault/tests/helpers/stubs';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | oidc/scope-form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
  });

  test('it should save new scope', async function (assert) {
    assert.expect(8);

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
    `);

    assert.dom('[data-test-oidc-scope-title]').hasText('Create Scope', 'Form title renders');
    assert.dom(SELECTORS.scopeSaveButton).hasText('Create', 'Save button has correct label');
    await click(SELECTORS.scopeSaveButton);

    // check validation errors
    await click(SELECTORS.scopeSaveButton);

    assert
      .dom(GENERAL.validationErrorByAttr('name'))
      .hasText('Name is required.', 'Validation messages are shown for name');
    assert
      .dom(SELECTORS.inlineAlert)
      .hasText('There is an error with this form.', 'Renders form error count');

    // json editor has test coverage so let's just confirm that it renders
    assert
      .dom(`${GENERAL.inputByAttr('template')} .hds-code-editor__header`)
      .exists('JsonEditor toolbar renders');
    assert.dom(`${GENERAL.inputByAttr('template')} ${GENERAL.codemirror}`).exists('Code mirror renders');

    await fillIn(GENERAL.inputByAttr('name'), 'test');
    await fillIn(GENERAL.inputByAttr('description'), 'this is a test');
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
    `);

    assert.dom('[data-test-oidc-scope-title]').hasText('Edit Scope', 'Form title renders');
    assert.dom(SELECTORS.scopeSaveButton).hasText('Update', 'Save button has correct label');
    assert.dom(GENERAL.inputByAttr('name')).isDisabled('Name input is disabled when editing');
    assert.dom(GENERAL.inputByAttr('name')).hasValue('test', 'Name input is populated with model value');
    assert
      .dom(GENERAL.inputByAttr('description'))
      .hasValue('this is a test', 'Description input is populated with model value');
    // json editor has test coverage so let's just confirm that it renders
    assert
      .dom(`${GENERAL.inputByAttr('template')} .hds-code-editor__header`)
      .exists('JsonEditor toolbar renders');
    assert
      .dom(`${GENERAL.inputByAttr('template')} [data-test-component="code-mirror-modifier"]`)
      .exists('Code mirror renders');

    await fillIn(GENERAL.inputByAttr('description'), 'this is an edit test');
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
      `);

    await fillIn(GENERAL.inputByAttr('description'), 'changed description attribute');
    await click(SELECTORS.scopeCancelButton);
    assert.strictEqual(
      this.model.description,
      'this is a test',
      'Model attributes are rolled back on cancel'
    );
  });

  test('it should show example template modal', async function (assert) {
    assert.expect(5);
    const MODAL = (e) => `[data-test-scope-modal="${e}"]`;
    this.model = this.store.createRecord('oidc/scope');

    await render(hbs`
      <Oidc::ScopeForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);

    await click('[data-test-oidc-scope-example]');
    assert.dom(MODAL('title')).hasText('Scope template', 'Modal title renders');
    assert.dom(MODAL('text')).hasText('Example of a JSON template for scopes:', 'Modal text renders');
    assert.dom('#scope-template-modal .hds-icon-clipboard-copy').exists('Modal copy button exists');
    assert.dom('.token .string').hasText('"username"', 'Example template json renders');
    await click('[data-test-close-modal]');
    assert.dom('.hds#scope-template-modal').doesNotExist('Modal is hidden');
  });

  test('it should render error alerts when API returns an error', async function (assert) {
    assert.expect(2);
    this.model = this.store.createRecord('oidc/scope');
    this.server.post('/sys/capabilities-self', () => capabilitiesStub(OIDC_BASE_URL + '/scopes'));
    await render(hbs`
      <Oidc::ScopeForm
        @model={{this.model}}
        @onCancel={{this.onCancel}}
        @onSave={{this.onSave}}
      />
    `);
    await fillIn(GENERAL.inputByAttr('name'), 'test-scope');
    await click(SELECTORS.scopeSaveButton);
    assert
      .dom(GENERAL.inlineAlert)
      .hasText('There was an error submitting this form.', 'form error alert renders ');
    assert.dom(GENERAL.messageError).exists('alert banner renders');
  });
});
