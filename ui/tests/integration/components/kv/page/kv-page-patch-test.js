/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { blur, click, fillIn, find, render, waitUntil } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { FORM, PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import codemirror, { setCodeEditorValue } from 'vault/tests/helpers/codemirror';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

module('Integration | Component | kv-v2 | Page::Secret::Patch', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');

  hooks.beforeEach(async function () {
    this.transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
    this.backend = 'kv-engine';
    this.path = 'my-secret';
    this.metadata = { current_version: 4 };
    this.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.backend, route: 'list' },
      { label: this.path, route: 'index' },
      { label: 'Patch' },
    ];
    this.subkeys = {
      foo: null,
      bar: {
        baz: null,
      },
      quux: null,
    };
    this.subkeysMeta = {
      created_time: '2021-12-14T20:28:00.773477Z',
      custom_metadata: null,
      deletion_time: '',
      destroyed: false,
      version: 1,
    };

    this.renderComponent = async () => {
      return render(
        hbs`
        <Page::Secret::Patch
          @backend={{this.backend}}
          @breadcrumbs={{this.breadcrumbs}}
          @metadata={{this.metadata}}
          @path={{this.path}}
          @subkeys={{this.subkeys}}
          @subkeysMeta={{this.subkeysMeta}}
        />`,
        { owner: this.engine }
      );
    };
  });

  test('it renders', async function (assert) {
    await this.renderComponent();
    assert.dom(PAGE.breadcrumbs).hasText(`Secrets ${this.backend} ${this.path} Patch`);
    assert.dom(PAGE.title).hasText('Patch Secret to New Version');
    assert.dom(GENERAL.fieldByAttr('Path')).isDisabled();
    assert.dom(GENERAL.fieldByAttr('Path')).hasValue(this.path);
    assert.dom(GENERAL.inputByAttr('JSON')).isNotChecked();
    assert.dom(GENERAL.inputByAttr('UI')).isChecked();
    assert.dom(FORM.patchEditorForm).exists('it renders editor form by default');
    assert.dom(GENERAL.codemirror).doesNotExist();
    Object.keys(this.subkeys).forEach((key, idx) => {
      assert.dom(FORM.keyInput(idx)).hasValue(key);
      assert.dom(FORM.keyInput(idx)).isDisabled();
    });
  });

  test('it selects JSON as an edit option', async function (assert) {
    await this.renderComponent();
    assert.dom(FORM.patchEditorForm).exists();
    await click(GENERAL.inputByAttr('JSON'));
    assert.dom(GENERAL.inputByAttr('JSON')).isChecked();
    assert.dom(GENERAL.inputByAttr('UI')).isNotChecked();
    assert.dom(FORM.patchEditorForm).doesNotExist();
    assert.dom(GENERAL.codemirror).exists();
  });

  test('it transitions on cancel', async function (assert) {
    await this.renderComponent();
    await click(FORM.cancelBtn);
    const [route] = this.transitionStub.lastCall.args;
    assert.strictEqual(
      route,
      'vault.cluster.secrets.backend.kv.secret.index',
      `it transitions on cancel to: ${route}`
    );
  });

  module('it submits', function (hooks) {
    const EXAMPLE_KV_DATA_CREATE_RESPONSE = {
      request_id: 'foobar',
      data: {
        created_time: '2023-06-21T16:18:31.479993Z',
        custom_metadata: null,
        deletion_time: '',
        destroyed: false,
        version: 1,
      },
    };

    hooks.beforeEach(async function () {
      this.endpoint = `${encodePath(this.backend)}/data/${encodePath(this.path)}`;
      this.patchStub = sinon
        .stub(this.owner.lookup('service:api').secrets, 'kvV2Patch')
        .resolves(EXAMPLE_KV_DATA_CREATE_RESPONSE);
    });

    test('patch data from kv editor form', async function (assert) {
      assert.expect(2);

      await this.renderComponent();
      // patch existing, delete and create a new key key
      await click(FORM.patchEdit());
      await fillIn(FORM.valueInput(), 'foovalue');
      await blur(FORM.valueInput());
      await click(FORM.patchDelete(1));
      await fillIn(FORM.keyInput('new'), 'aKey');
      await fillIn(FORM.valueInput('new'), '1');
      await click(FORM.patchAdd);
      // add new key and do NOT click add
      await fillIn(FORM.keyInput('new'), 'bKey');
      await fillIn(FORM.valueInput('new'), 'null');
      await click(FORM.saveBtn);

      const payload = {
        data: { bar: null, foo: 'foovalue', aKey: '1', bKey: 'null' },
        options: {
          cas: this.metadata.current_version,
        },
      };
      assert.true(
        this.patchStub.calledWith(this.path, this.backend, payload),
        'Patch request made with correct args'
      );
      assert.true(
        this.transitionStub.calledWith('vault.cluster.secrets.backend.kv.secret.index'),
        'transitions to overview route on save'
      );
    });

    test('patch data from json form', async function (assert) {
      assert.expect(2);

      await this.renderComponent();
      await click(GENERAL.inputByAttr('JSON'));
      await waitUntil(() => find('.cm-editor'));
      const editor = codemirror();
      setCodeEditorValue(editor, '{ "foo": "foovalue", "bar":null, "number":1 }');
      await click(FORM.saveBtn);

      const payload = {
        data: { foo: 'foovalue', bar: null, number: 1 },
        options: {
          cas: this.metadata.current_version,
        },
      };
      assert.true(
        this.patchStub.calledWith(this.path, this.backend, payload),
        'Patch request made with correct args'
      );
      assert.true(
        this.transitionStub.calledWith('vault.cluster.secrets.backend.kv.secret.index'),
        'transitions to overview route on save'
      );
    });

    // this assertion confirms submit allows empty values
    test('empty string values from kv editor form', async function (assert) {
      assert.expect(1);

      await this.renderComponent();
      await click(FORM.patchEdit());
      // edit existing key's value
      await fillIn(FORM.valueInput(), '');
      // add a new key with empty value, click add
      await fillIn(FORM.keyInput('new'), 'aKey');
      await fillIn(FORM.valueInput('new'), '');
      await click(FORM.patchAdd);
      // add new key and do NOT click add
      await fillIn(FORM.keyInput('new'), 'bKey');
      await fillIn(FORM.valueInput('new'), '');
      await click(FORM.saveBtn);

      const payload = {
        data: { foo: '', aKey: '', bKey: '' },
        options: {
          cas: this.metadata.current_version,
        },
      };
      assert.true(
        this.patchStub.calledWith(this.path, this.backend, payload),
        'Patch request made with correct args'
      );
    });

    // this assertion confirms submit allows empty values
    test('empty string value from json form', async function (assert) {
      assert.expect(1);

      await this.renderComponent();
      await click(GENERAL.inputByAttr('JSON'));
      await waitUntil(() => find('.cm-editor'));
      const editor = codemirror();
      setCodeEditorValue(editor, '{ "foo": "" }');
      await click(FORM.saveBtn);

      const payload = {
        data: { foo: '' },
        options: {
          cas: this.metadata.current_version,
        },
      };
      assert.true(
        this.patchStub.calledWith(this.path, this.backend, payload),
        'Patch request made with correct args'
      );
    });

    test('patch data without metadata permissions', async function (assert) {
      assert.expect(2);
      this.metadata = null;

      await this.renderComponent();
      await fillIn(FORM.keyInput('new'), 'aKey');
      await fillIn(FORM.valueInput('new'), '1');
      await click(FORM.saveBtn);

      const payload = {
        data: { aKey: '1' },
        options: {
          cas: this.subkeysMeta.version,
        },
      };
      assert.true(
        this.patchStub.calledWith(this.path, this.backend, payload),
        'Patch request made with correct args'
      );
      assert.true(
        this.transitionStub.calledWith('vault.cluster.secrets.backend.kv.secret.index'),
        'transitions to overview route on save'
      );
    });
  });

  module('it does not submit', function (hooks) {
    hooks.beforeEach(async function () {
      this.endpoint = `${encodePath(this.backend)}/data/${encodePath(this.path)}`;
      this.flashSpy = sinon.spy(this.owner.lookup('service:flash-messages'), 'info');
      const errors = { errors: ['Something went wrong. This should not have happened!'] };
      this.patchStub = sinon
        .stub(this.owner.lookup('service:api').secrets, 'kvV2Patch')
        .rejects(getErrorResponse(errors, 500));
    });

    test('if no changes from kv editor form', async function (assert) {
      assert.expect(3);

      await this.renderComponent();
      await click(FORM.saveBtn);
      assert.dom(GENERAL.messageError).doesNotExist('PATCH request is not made');
      const route = this.transitionStub.lastCall?.args[0] || '';
      const flash = this.flashSpy.lastCall?.args[0] || '';
      assert.strictEqual(
        route,
        'vault.cluster.secrets.backend.kv.secret.index',
        `it transitions to overview route: ${route}`
      );
      assert.strictEqual(
        flash,
        `No changes to submit. No updates made to "${this.path}".`,
        `flash message has message: "${flash}"`
      );
    });

    test('if no changes from json form', async function (assert) {
      assert.expect(3);

      await this.renderComponent();
      await click(GENERAL.inputByAttr('JSON'));
      await waitUntil(() => find('.cm-editor'));
      await click(FORM.saveBtn);
      assert.dom(GENERAL.messageError).doesNotExist('PATCH request is not made');
      const route = this.transitionStub.lastCall?.args[0] || '';
      const flash = this.flashSpy.lastCall?.args[0] || '';
      assert.strictEqual(
        route,
        'vault.cluster.secrets.backend.kv.secret.index',
        `it transitions to overview route: ${route}`
      );
      assert.strictEqual(
        flash,
        `No changes to submit. No updates made to "${this.path}".`,
        `flash message has message: "${flash}"`
      );
    });
  });

  module('it passes error', function (hooks) {
    hooks.beforeEach(async function () {
      const errors = { errors: ['permission denied'] };
      this.patchStub = sinon
        .stub(this.owner.lookup('service:api').secrets, 'kvV2Patch')
        .rejects(getErrorResponse(errors, 403));
    });

    test('to kv editor form', async function (assert) {
      assert.expect(2);

      await this.renderComponent();
      // patch existing, delete and create a new key key
      await click(FORM.patchEdit());
      await fillIn(FORM.valueInput(), 'foovalue');
      await blur(FORM.valueInput());
      await click(FORM.patchDelete(1));
      await fillIn(FORM.keyInput('new'), 'aKey');
      await fillIn(FORM.valueInput('new'), 'aValue');
      await click(FORM.patchAdd);
      // add new key and do NOT click add
      await fillIn(FORM.keyInput('new'), 'bKey');
      await fillIn(FORM.valueInput('new'), 'bValue');
      await click(FORM.saveBtn);
      assert.dom(GENERAL.messageError).hasText('Error permission denied');
      assert.dom(GENERAL.inlineError).hasText('There was an error submitting this form.');
    });

    test('to json form', async function (assert) {
      assert.expect(2);
      await this.renderComponent();
      await click(GENERAL.inputByAttr('JSON'));
      await waitUntil(() => find('.cm-editor'));
      const editor = codemirror();
      setCodeEditorValue(editor, '{ "foo": "foovalue", "bar":null, "number":1 }');
      await click(FORM.saveBtn);
      await click(FORM.saveBtn);
      assert.dom(GENERAL.messageError).hasText('Error permission denied');
      assert.dom(GENERAL.inlineError).hasText('There was an error submitting this form.');
    });
  });
});
