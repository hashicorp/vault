/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { hbs } from 'ember-cli-htmlbars';
import { fillIn, render, click } from '@ember/test-helpers';
import codemirror from 'vault/tests/helpers/codemirror';
import { PAGE, FORM } from 'vault/tests/helpers/kv/kv-selectors';
import { setRunOptions } from 'ember-a11y-testing/test-support';

module('Integration | Component | kv-v2 | KvDataFields', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.backend = 'my-kv-engine';
    this.path = 'my-secret';
    this.secret = this.store.createRecord('kv/data', { backend: this.backend });
    // TODO: Fix JSONEditor/CodeMirror
    setRunOptions({
      rules: {
        label: { enabled: false },
      },
    });
  });

  test('it updates the secret model', async function (assert) {
    assert.expect(2);

    await render(hbs`<KvDataFields @showJson={{false}} @secret={{this.secret}} @type="create" />`, {
      owner: this.engine,
    });
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
      `{ \"\": \"\" }`, // eslint-disable-line no-useless-escape
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

    await render(hbs`<KvDataFields @showJson={{false}} @secret={{this.secret}} @type="edit" />`, {
      owner: this.engine,
    });

    assert.dom(FORM.inputByAttr('path')).isDisabled();
    assert.dom(FORM.inputByAttr('path')).hasValue(this.path);
    assert.dom(FORM.keyInput()).hasValue('foo');
    assert.dom(FORM.maskedValueInput()).hasValue('bar');
    assert.dom(FORM.dataInputLabel({ isJson: false })).hasText('Version data');
  });

  test('it shows readonly info rows when viewing secret details of simple secret', async function (assert) {
    assert.expect(3);
    this.secret.secretData = { foo: 'bar' };
    this.secret.path = this.path;

    await render(hbs`<KvDataFields @showJson={{false}} @secret={{this.secret}} @type="details" />`, {
      owner: this.engine,
    });
    assert.dom(PAGE.infoRow).exists({ count: 1 }, '1 row of data shows');
    assert.dom(PAGE.infoRowValue('foo')).hasText('***********');
    await click(PAGE.infoRowToggleMasked('foo'));
    assert.dom(PAGE.infoRowValue('foo')).hasText('bar', 'secret value shows after toggle');
  });

  test('it shows readonly json editor when viewing secret details of complex secret', async function (assert) {
    assert.expect(3);
    this.secret.secretData = {
      foo: {
        bar: 'baz',
      },
    };
    this.secret.path = this.path;

    await render(hbs`<KvDataFields @showJson={{true}} @secret={{this.secret}} @type="details" />`, {
      owner: this.engine,
    });
    assert.dom(PAGE.infoRowValue('foo')).doesNotExist('does not render rows of secret data');
    assert.dom('[data-test-component="code-mirror-modifier"]').hasClass('readonly-codemirror');
    assert.dom('[data-test-component="code-mirror-modifier"]').includesText(`{ "foo": { "bar": "baz" }}`);
  });
});
