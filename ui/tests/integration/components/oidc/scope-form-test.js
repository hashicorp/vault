/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, fillIn, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { SELECTORS } from 'vault/tests/helpers/oidc-config';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';
import OidcScopeForm from 'vault/forms/oidc/scope';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

module('Integration | Component | oidc/scope-form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    const api = this.owner.lookup('service:api');
    this.writeStub = sinon.stub(api.identity, 'oidcWriteScope').resolves();
    this.onCancel = sinon.spy();
    this.onSave = sinon.spy();

    this.renderComponent = (scope) => {
      this.form = new OidcScopeForm(scope || {}, { isNew: !scope });
      return render(hbs`
        <Oidc::ScopeForm
          @form={{this.form}}
          @onCancel={{this.onCancel}}
          @onSave={{this.onSave}}
        />
      `);
    };
  });

  test('it should save new scope', async function (assert) {
    assert.expect(8);

    await this.renderComponent();

    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Create Scope', 'Form title renders');
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

    assert.true(this.onSave.calledOnce, 'onSave callback is called on successful save');
    assert.true(
      this.writeStub.calledWith('test', { description: 'this is a test' }),
      'API is called with correct parameters'
    );
  });

  test('it should update scope', async function (assert) {
    assert.expect(9);

    await this.renderComponent({ name: 'test', description: 'this is a test' });

    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Edit Scope', 'Form title renders');
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

    assert.true(this.onSave.calledOnce, 'onSave callback is called on successful save');
    assert.true(
      this.writeStub.calledWith('test', { description: 'this is an edit test' }),
      'API is called with correct parameters'
    );
  });

  test('it should trigger on cancel callback', async function (assert) {
    assert.expect(1);

    await this.renderComponent();
    await click(SELECTORS.scopeCancelButton);
    assert.true(this.onCancel.calledOnce, 'onCancel callback is called when cancel button is clicked');
  });

  test('it should show example template modal', async function (assert) {
    assert.expect(5);

    const MODAL = (e) => `[data-test-scope-modal="${e}"]`;

    await this.renderComponent();

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

    this.writeStub.rejects(getErrorResponse());
    await this.renderComponent();

    await fillIn(GENERAL.inputByAttr('name'), 'test-scope');
    await click(SELECTORS.scopeSaveButton);
    assert
      .dom(GENERAL.inlineAlert)
      .hasText('There was an error submitting this form.', 'form error alert renders ');
    assert.dom(GENERAL.messageError).exists('alert banner renders');
  });
});
