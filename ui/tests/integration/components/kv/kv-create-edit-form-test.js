/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { hbs } from 'ember-cli-htmlbars';
import { fillIn, typeIn, render, click, waitFor, findAll } from '@ember/test-helpers';
import codemirror, { getCodeEditorValue, setCodeEditorValue } from 'vault/tests/helpers/codemirror';
import { FORM } from 'vault/tests/helpers/kv/kv-selectors';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import KvForm from 'vault/forms/secrets/kv';
import sinon from 'sinon';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

module('Integration | Component | kv-v2 | KvCreateEditForm', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');

  hooks.beforeEach(function () {
    this.backend = 'my-kv-engine';
    this.showJson = false;
    this.onChange = sinon.stub();

    this.renderComponent = (isNew = false) => {
      this.path = isNew ? undefined : 'my-secret';
      const options = isNew ? undefined : { cas: 2 };
      const secretData = isNew ? undefined : { foo: 'bar' };
      this.form = new KvForm(
        {
          path: this.path,
          secretData,
          max_versions: 0,
          delete_version_after: '0s',
          cas_required: false,
          options,
        },
        { isNew }
      );
      return render(
        hbs`
          <KvCreateEditForm
            @form={{this.form}}
            @path={{this.path}}
            @backend={{this.backend}}
            @showJson={{this.showJson}}
            @onChange={{this.onChange}}
            as |modelValidations|
          >
            <span data-test-yield-block>{{modelValidations.invalidFormMessage}}</span>
          </KvCreateEditForm>
        `,
        { owner: this.engine }
      );
    };

    const api = this.owner.lookup('service:api');
    this.dataWrite = sinon.stub(api.secrets, 'kvV2Write').resolves();
    this.metadataWrite = sinon.stub(api.secrets, 'kvV2WriteMetadata').resolves();

    this.transitionTo = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    setRunOptions({
      rules: {
        // failing on .CodeMirror-scroll
        'scrollable-region-focusable': { enabled: false },
      },
    });
  });

  test('it should render namespace reminder', async function (assert) {
    this.owner.lookup('service:namespace').path = 'test';
    await this.renderComponent();
    assert.dom('#namespace-reminder').hasText('This secret will be created in the test/ namespace.');
  });

  test('it should render validation errors and warnings', async function (assert) {
    await this.renderComponent(true);

    await click(FORM.saveBtn);

    assert.dom(FORM.validationError('path')).hasText(`Path can't be blank.`);
    assert.dom(FORM.invalidFormAlert).hasText('There is an error with this form.');

    await typeIn(FORM.inputByAttr('path'), 'my secret');
    assert
      .dom(FORM.validationWarning)
      .hasText(
        `Path contains whitespace. If this is desired, you'll need to encode it with %20 in API requests.`
      );
  });

  test('it should render JSON editor and add new secretData', async function (assert) {
    assert.expect(4);

    this.showJson = true;
    await this.renderComponent(true);
    await waitFor('.cm-editor');

    const editor = codemirror();
    const editorValue = getCodeEditorValue(editor);
    assert.strictEqual(
      editorValue,
      `{\n  "": ""\n}`,
      'json editor initializes with empty object that includes whitespace'
    );
    setCodeEditorValue(editor, 'blah');

    await waitFor('.cm-lint-marker');
    const lintMarkers = findAll('.cm-lint-marker');

    assert.strictEqual(lintMarkers.length, 1, 'codemirror lints input');
    setCodeEditorValue(editor, `{ "hello": "there"}`);
    assert.propEqual(this.form.data.secretData, { hello: 'there' }, 'json editor updates secret data');
    assert.true(
      this.onChange.calledWith({ hello: 'there' }),
      'onChange is called with secretData on json editor change'
    );
  });

  test('it should create new secret', async function (assert) {
    const path = 'new-secret';
    const secretData = { foo: 'bar' };

    await this.renderComponent(true);
    await fillIn(FORM.inputByAttr('path'), path);
    await fillIn(FORM.keyInput(), 'foo');
    await fillIn(FORM.maskedValueInput(), 'bar');

    assert.dom(FORM.dataInputLabel({})).hasText('Secret data', 'Correct data label renders');
    assert.true(
      this.onChange.calledWith(secretData),
      'onChange is called with secretData on kv object change'
    );

    await click(FORM.saveBtn);
    assert.true(this.dataWrite.calledWith(path, this.backend, { data: secretData }), 'secret data is saved');
    assert.true(this.metadataWrite.notCalled, 'metadata is not saved when there are no changes');
    assert.true(
      this.transitionTo.calledWith('vault.cluster.secrets.backend.kv.secret.index', path),
      'transitions to secret on success'
    );

    // metadata is updated outside of component
    // simulate metadata change by updating form data
    this.form.data.max_versions = 5;
    await click(FORM.saveBtn);

    assert.true(this.dataWrite.calledWith(path, this.backend, { data: secretData }), 'secret data is saved');
    assert.true(
      this.metadataWrite.calledWith(path, this.backend, {
        max_versions: 5,
        delete_version_after: '0s',
        cas_required: false,
      }),
      'metadata is saved when there are changes'
    );
    assert.true(
      this.transitionTo.calledWith('vault.cluster.secrets.backend.kv.secret.index', path),
      'transitions to secret on success'
    );
  });

  test('it should create new secret version', async function (assert) {
    await this.renderComponent();

    assert.dom(FORM.inputByAttr('path')).isDisabled();
    assert.dom(FORM.inputByAttr('path')).hasValue(this.path);
    assert.dom(FORM.dataInputLabel({})).hasText('Version data', 'Correct data label renders');
    assert.dom(FORM.keyInput()).hasValue('foo');
    await click(`${FORM.valueInput()} button`); // reveal value
    assert.dom(FORM.maskedValueInput()).hasValue('bar');

    await fillIn(FORM.keyInput(1), 'bar');
    await fillIn(FORM.maskedValueInput(1), 'baz');
    await click(FORM.saveBtn);

    assert.true(
      this.dataWrite.calledWith(this.path, this.backend, {
        data: { foo: 'bar', bar: 'baz' },
        options: { cas: 2 },
      }),
      'secret data is saved'
    );
    assert.true(this.metadataWrite.notCalled, 'metadata is not saved when there are no changes');
    assert.true(
      this.transitionTo.calledWith('vault.cluster.secrets.backend.kv.secret.index', this.path),
      'transitions to secret on success'
    );
  });

  test('it should handle save errors', async function (assert) {
    await this.renderComponent();
    this.form.data.max_versions = 5;

    // data save failure
    const dataError = getErrorResponse({ errors: ['error saving secret data'] }, 400);
    this.dataWrite.rejects(dataError);
    await click(FORM.saveBtn);

    assert.dom(FORM.messageError).includesText('error saving secret data', 'data error message renders');
    assert.true(this.metadataWrite.notCalled, 'metadata is not saved on data save failure');

    // data control group error
    sinon
      .stub(this.owner.lookup('service:control-group'), 'logFromError')
      .returns({ content: 'A Control Group was encountered' });
    const ctrlError = getErrorResponse({ accessor: 'foobar', isControlGroupError: true }, 500);
    this.dataWrite.rejects(ctrlError);
    await click(FORM.saveBtn);

    assert
      .dom(FORM.messageError)
      .includesText('A Control Group was encountered', 'control group error message renders');

    // data success and metadata failure
    this.dataWrite.resolves();
    const metaError = getErrorResponse({ errors: ['error saving secret metadata'] }, 400);
    this.metadataWrite.rejects(metaError);
    const flashStub = sinon.stub(this.owner.lookup('service:flash-messages'), 'danger');
    await click(FORM.saveBtn);

    const flashMessage = 'Secret data was saved but metadata was not: error saving secret metadata';
    assert.true(flashStub.calledWith(flashMessage), 'flash message displays with metadata save error');
  });
});
