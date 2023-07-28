/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import { SELECTORS as PAGE } from 'vault/tests/helpers/kv/kv-page-selectors';
import { SELECTORS } from 'vault/tests/helpers/kv/kv-general-selectors';
import codemirror from 'vault/tests/helpers/codemirror';

module('Integration | Component | kv | KvSecretForm', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    this.store = this.owner.lookup('service:store');
    this.backend = 'my-kv-engine';
    this.path = 'my-secret';
    this.secret = this.store.createRecord('kv/data', { backend: this.backend });
    this.onSave = () => {};
    this.onCancel = () => {};
  });

  test('it makes post request on save', async function (assert) {
    assert.expect(3);
    this.onSave = () => assert.ok(true, 'onSave callback fires on save success');
    this.server.post(`${this.backend}/data/${this.path}`, (schema, req) => {
      assert.ok(true, 'Request made to save secret');
      const payload = JSON.parse(req.requestBody);
      assert.propEqual(payload, {
        data: { foo: 'bar' },
        options: { cas: 0 },
      });
      return {
        request_id: 'bd76db73-605d-fcbc-0dad-d44a008f9b95',
        data: {
          created_time: '2023-07-28T18:47:32.924809Z',
          custom_metadata: null,
          deletion_time: '',
          destroyed: false,
          version: 1,
        },
      };
    });

    await render(
      hbs`
        <KvSecretForm
          @secret={{this.secret}}
          @onSave={{this.onSave}}
        />`,
      { owner: this.engine }
    );

    await fillIn(PAGE.form.inputByAttr('path'), this.path);
    await fillIn(PAGE.form.keyInput(), 'foo');
    await fillIn(PAGE.form.maskedValueInput(), 'bar');
    await click(PAGE.form.secretSave);
  });

  test('it rolls back model attrs on cancel and JSON editor modifies secretData', async function (assert) {
    assert.expect(5);
    this.onCancel = () => assert.ok(true, 'onCancel callback fires on save success');

    await render(
      hbs`
        <KvSecretForm
          @secret={{this.secret}}
          @onCancel={{this.onCancel}}
        />`,
      { owner: this.engine }
    );

    await click(SELECTORS.toggleJson);
    assert.strictEqual(
      codemirror().getValue(' '),
      `{   \"\": \"\" }`, // eslint-disable-line no-useless-escape
      'json editor initializes with empty object'
    );
    await fillIn(SELECTORS.jsonEditor, 'blah');
    assert.strictEqual(codemirror().state.lint.marked.length, 1, 'codemirror lints input');
    codemirror().setValue(`{ "hello": "there"}`);
    assert.propEqual(this.secret.secretData, { hello: 'there' }, 'json editor updates secret data');
    await click(PAGE.form.secretCancel);
    assert.propEqual(this.secret.secretData, {}, 'it rolls back model attribute');
  });
});
