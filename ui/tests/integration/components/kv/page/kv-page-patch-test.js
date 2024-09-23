/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { blur, click, fillIn, find, render, waitUntil } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { FORM, PAGE } from 'vault/tests/helpers/kv/kv-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { baseSetup } from 'vault/tests/helpers/kv/kv-run-commands';
import codemirror from 'vault/tests/helpers/codemirror';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { overrideResponse } from 'vault/tests/helpers/stubs';

module('Integration | Component | kv-v2 | Page::Secret::Patch', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    baseSetup(this);
    this.transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
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
    });

    test('patch data from kv editor form', async function (assert) {
      assert.expect(3);
      this.server.patch(this.endpoint, (schema, req) => {
        const payload = JSON.parse(req.requestBody);
        const expected = {
          data: { bar: null, foo: 'foovalue', aKey: '1', bKey: 'null' },
          options: {
            cas: this.metadata.currentVersion,
          },
        };
        assert.true(true, `PATCH request made to ${this.endpoint}`);
        assert.propEqual(
          payload,
          expected,
          `payload: ${JSON.stringify(payload)} matches expected: ${JSON.stringify(payload)}`
        );
        return EXAMPLE_KV_DATA_CREATE_RESPONSE;
      });

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
      const [route] = this.transitionStub.lastCall.args;
      assert.strictEqual(
        route,
        'vault.cluster.secrets.backend.kv.secret.index',
        `it transitions on save to: ${route}`
      );
    });

    test('patch data from json form', async function (assert) {
      assert.expect(3);
      this.server.patch(this.endpoint, (schema, req) => {
        const payload = JSON.parse(req.requestBody);
        const expected = {
          data: { foo: 'foovalue', bar: null, number: 1 },
          options: {
            cas: 4,
          },
        };
        assert.true(true, `PATCH request made to ${this.endpoint}`);
        assert.propEqual(
          payload,
          expected,
          `payload: ${JSON.stringify(payload)} matches expected: ${JSON.stringify(payload)}`
        );
        return EXAMPLE_KV_DATA_CREATE_RESPONSE;
      });
      await this.renderComponent();
      await click(GENERAL.inputByAttr('JSON'));
      await waitUntil(() => find('.CodeMirror'));
      await codemirror().setValue('{ "foo": "foovalue", "bar":null, "number":1 }');
      await click(FORM.saveBtn);
      const [route] = this.transitionStub.lastCall.args;
      assert.strictEqual(
        route,
        'vault.cluster.secrets.backend.kv.secret.index',
        `it transitions on save to: ${route}`
      );
    });

    // this assertion confirms submit allows empty values
    test('empty string values from kv editor form', async function (assert) {
      assert.expect(1);
      this.server.patch(this.endpoint, (schema, req) => {
        const payload = JSON.parse(req.requestBody);
        const expected = {
          data: { foo: '', aKey: '', bKey: '' },
          options: {
            cas: this.metadata.currentVersion,
          },
        };
        assert.propEqual(
          payload,
          expected,
          `payload: ${JSON.stringify(payload)} matches expected: ${JSON.stringify(payload)}`
        );
        return EXAMPLE_KV_DATA_CREATE_RESPONSE;
      });

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
    });

    // this assertion confirms submit allows empty values
    test('empty string value from json form', async function (assert) {
      assert.expect(1);
      this.server.patch(this.endpoint, (schema, req) => {
        const payload = JSON.parse(req.requestBody);
        const expected = {
          data: { foo: '' },
          options: {
            cas: this.metadata.currentVersion,
          },
        };
        assert.propEqual(
          payload,
          expected,
          `payload: ${JSON.stringify(payload)} matches expected: ${JSON.stringify(payload)}`
        );
        return EXAMPLE_KV_DATA_CREATE_RESPONSE;
      });

      await this.renderComponent();
      await click(GENERAL.inputByAttr('JSON'));
      await waitUntil(() => find('.CodeMirror'));
      await codemirror().setValue('{ "foo": "" }');
      await click(FORM.saveBtn);
    });

    test('patch data without metadata permissions', async function (assert) {
      assert.expect(3);
      this.metadata = null;
      this.server.patch(this.endpoint, (schema, req) => {
        const payload = JSON.parse(req.requestBody);
        const expected = {
          data: { aKey: '1' },
          options: {
            cas: this.subkeysMeta.version,
          },
        };
        assert.true(true, `PATCH request made to ${this.endpoint}`);
        assert.propEqual(
          payload,
          expected,
          `payload: ${JSON.stringify(payload)} matches expected: ${JSON.stringify(payload)}`
        );
        return EXAMPLE_KV_DATA_CREATE_RESPONSE;
      });

      await this.renderComponent();
      await fillIn(FORM.keyInput('new'), 'aKey');
      await fillIn(FORM.valueInput('new'), '1');
      await click(FORM.saveBtn);
      const [route] = this.transitionStub.lastCall.args;
      assert.strictEqual(
        route,
        'vault.cluster.secrets.backend.kv.secret.index',
        `it transitions on save to: ${route}`
      );
    });
  });

  module('it does not submit', function (hooks) {
    hooks.beforeEach(async function () {
      this.endpoint = `${encodePath(this.backend)}/data/${encodePath(this.path)}`;
      this.flashSpy = sinon.spy(this.owner.lookup('service:flash-messages'), 'info');
    });

    test('if no changes from kv editor form', async function (assert) {
      assert.expect(3);
      this.server.patch(this.endpoint, () =>
        overrideResponse(500, `Request made to: ${this.endpoint}. This should not have happened!`)
      );
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
      this.server.patch(this.endpoint, () =>
        overrideResponse(500, `Request made to: ${this.endpoint}. This should not have happened!`)
      );
      await this.renderComponent();
      await click(GENERAL.inputByAttr('JSON'));
      await waitUntil(() => find('.CodeMirror'));
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
      this.endpoint = `${encodePath(this.backend)}/data/${encodePath(this.path)}`;
      this.server.patch(this.endpoint, () => {
        return overrideResponse(403);
      });
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
      await waitUntil(() => find('.CodeMirror'));
      await codemirror().setValue('{ "foo": "foovalue", "bar":null, "number":1 }');
      await click(FORM.saveBtn);
      await click(FORM.saveBtn);
      assert.dom(GENERAL.messageError).hasText('Error permission denied');
      assert.dom(GENERAL.inlineError).hasText('There was an error submitting this form.');
    });
  });
});
