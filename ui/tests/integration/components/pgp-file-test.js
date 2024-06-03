/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn, triggerEvent, waitUntil, find } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

let file;
const fileEvent = () => {
  const data = { some: 'content' };
  file = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
  file.name = 'file.json';
  return ['change', { files: [file] }];
};

module('Integration | Component | pgp file', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    file = null;
    this.lastOnChangeCall = null;
    this.set('change', (index, key) => {
      this.lastOnChangeCall = [index, key];
      this.set('key', key);
    });
  });

  test('it renders', async function (assert) {
    this.set('key', { value: '' });
    this.set('index', 0);

    await render(hbs`
      <PgpFile
        @index={{this.index}}
        @key={{this.key}}
        @onChange={{action this.change}}
      />
    `);

    assert.dom('[data-test-pgp-label]').hasText('PGP KEY 1');
    assert.dom('[data-test-pgp-file-input-label]').hasText('Choose a file…');
  });

  test('it accepts files', async function (assert) {
    const key = { value: '' };
    const event = fileEvent();
    this.set('key', key);
    this.set('index', 0);

    await render(hbs`
      <PgpFile
        @index={{this.index}}
        @key={{this.key}}
        @onChange={{action this.change}}
      />
    `);
    triggerEvent('[data-test-pgp-file-input]', ...event);

    // FileReader is async, but then we need extra run loop wait to re-render
    await waitUntil(() => {
      return !!this.lastOnChangeCall;
    });
    assert.dom('[data-test-pgp-file-input-label]').hasText(file.name, 'the file input shows the file name');
    assert.notDeepEqual(this.lastOnChangeCall[1].value, key.value, 'onChange was called with the new key');
    assert.strictEqual(this.lastOnChangeCall[0], 0, 'onChange is called with the index value');
    await click('[data-test-pgp-clear]');
    assert.strictEqual(
      this.lastOnChangeCall[1].value,
      key.value,
      'the key gets reset when the input is cleared'
    );
  });

  test('it allows for text entry', async function (assert) {
    const key = { value: '' };
    const text = 'a really long pgp key';
    this.set('key', key);
    this.set('index', 0);

    await render(hbs`
      <PgpFile
        @index={{this.index}}
        @key={{this.key}}
        @onChange={{action this.change}}
      />
    `);
    await click('[data-test-text-toggle]');
    assert.dom('[data-test-pgp-file-textarea]').exists({ count: 1 }, 'renders the textarea on toggle');

    fillIn('[data-test-pgp-file-textarea]', text);
    await waitUntil(() => {
      return !!this.lastOnChangeCall;
    });
    assert.strictEqual(this.lastOnChangeCall[1].value, text, 'the key value is passed to onChange');
  });

  test('toggling back and forth', async function (assert) {
    const key = { value: '' };
    const event = fileEvent();
    this.set('key', key);
    this.set('index', 0);

    await render(hbs`
      <PgpFile
        @index={{this.index}}
        @key={{this.key}}
        @onChange={{action this.change}}
      />
    `);
    await triggerEvent('[data-test-pgp-file-input]', ...event);
    await waitUntil(() => find('[data-test-pgp-file-input-label]').innerText === 'file.json');
    await click('[data-test-text-toggle]');
    assert.dom('[data-test-pgp-file-textarea]').exists({ count: 1 }, 'renders the textarea on toggle');
    assert
      .dom('[data-test-pgp-file-textarea]')
      .hasText(this.lastOnChangeCall[1].value, 'textarea shows the value of the base64d key');
  });
});
