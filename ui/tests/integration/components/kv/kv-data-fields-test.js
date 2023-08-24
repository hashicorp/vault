/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { hbs } from 'ember-cli-htmlbars';
import { fillIn, render } from '@ember/test-helpers';
import codemirror from 'vault/tests/helpers/codemirror';
import { FORM } from 'vault/tests/helpers/kv/kv-selectors';

module('Integration | Component | kv-v2 | KvDataFields', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.backend = 'my-kv-engine';
    this.path = 'my-secret';
    this.secret = this.store.createRecord('kv/data', { backend: this.backend });
  });

  test('it updates the secret model', async function (assert) {
    assert.expect(2);

    await render(hbs`<KvDataFields @showJson={{false}} @secret={{this.secret}} />`, { owner: this.engine });

    await fillIn(FORM.inputByAttr('path'), this.path);
    await fillIn(FORM.keyInput(), 'foo');
    await fillIn(FORM.maskedValueInput(), 'bar');
    assert.strictEqual(this.secret.path, this.path);
    assert.propEqual(this.secret.secretData, { foo: 'bar' });
  });

  test('it JSON editor initializes with empty object and modifies secretData', async function (assert) {
    assert.expect(3);

    await render(hbs`<KvDataFields @showJson={{true}} @secret={{this.secret}} />`, { owner: this.engine });

    assert.strictEqual(
      codemirror().getValue(' '),
      `{   \"\": \"\" }`, // eslint-disable-line no-useless-escape
      'json editor initializes with empty object'
    );
    await fillIn(`${FORM.jsonEditor} textarea`, 'blah');
    assert.strictEqual(codemirror().state.lint.marked.length, 1, 'codemirror lints input');
    codemirror().setValue(`{ "hello": "there"}`);
    assert.propEqual(this.secret.secretData, { hello: 'there' }, 'json editor updates secret data');
  });

  test('it disables path and prefills secret data when creating a new secret version', async function (assert) {
    assert.expect(5);
    this.secret.secretData = { foo: 'bar' };
    this.secret.path = this.path;

    this.newVersion = this.store.createRecord('kv/data', {
      backend: this.backend,
      path: this.path,
      secretData: this.secret.secretData,
    });

    await render(hbs`<KvDataFields @showJson={{false}} @isEdit={{true}} @secret={{this.secret}} />`, {
      owner: this.engine,
    });

    assert.dom(FORM.inputByAttr('path')).isDisabled();
    assert.dom(FORM.inputByAttr('path')).hasValue(this.path);
    assert.dom(FORM.keyInput()).hasValue('foo');
    assert.dom(FORM.maskedValueInput()).hasValue('bar');
    assert.dom(FORM.dataInputLabel({ isJson: false })).hasText('Version data');
  });
});
