/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click, fillIn, settled, typeIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const {
  searchSelect: { option, options },
  kvSuggestion: { input },
} = GENERAL;

module('Integration | Component | kv-suggestion-input', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.label = 'Input something';
    this.mountPath = 'secret/';
    this.keys = ['foo/', 'my-secret'];
    const response = () => ({ data: { keys: this.keys } });
    this.server.get('/:mount/metadata', response);
    this.server.get('/:mount/metadata/*', response);
    return render(hbs`
      <KvSuggestionInput
        @label={{this.label}}
        @subText="Suggestions will display"
        @value={{this.secretPath}}
        @mountPath={{this.mountPath}}
        @onChange={{fn (mut this.secretPath)}}
      />
    `);
  });

  test('it should render label and sub text', async function (assert) {
    assert.dom('label').hasText('Input something', 'Label renders');
    assert.dom('[data-test-label-subtext]').hasText('Suggestions will display', 'Label subtext renders');
    this.set('label', null);
    await settled();
    assert.dom('label').doesNotExist('Label is hidden when arg is not provided');
  });

  test('it should disable input when mountPath is not provided', async function (assert) {
    this.set('mountPath', null);
    assert.dom(input).isDisabled('Input is disabled when mountPath is not provided');
  });

  test('it should fetch suggestions for initial mount path', async function (assert) {
    assert.dom(options).doesNotExist('Suggestions are hidden initially');
    await click(input);
    ['foo/', 'my-secret'].forEach((secret, index) => {
      assert.dom(option(index)).hasText(secret, 'Suggestion renders for initial mount path');
    });
  });

  test('it should fetch secrets and update suggestions on mountPath change', async function (assert) {
    this.keys = ['test1'];
    this.set('mountPath', 'foo/');
    await settled();
    await click(input);
    assert.dom(option()).hasText('test1', 'Suggestions are fetched and render on mountPath change');
  });

  test('it should filter current result set', async function (assert) {
    await click(input);
    await typeIn(input, 'sec');
    assert.dom(options).exists({ count: 1 }, 'Correct number of options render based on input value');
    assert.dom(option()).hasText('my-secret', 'Result set is filtered');
  });

  test('it should replace filter terms with full path to secret', async function (assert) {
    await fillIn(input, 'sec');
    await click(option());
    assert.dom(input).hasValue('my-secret', 'Partial term replaced with selected secret');

    await fillIn(input, '');
    this.keys = ['secret-nested', 'bar', 'baz'];
    await click(option());
    await fillIn(input, 'nest');
    await click(option());
    assert
      .dom(input)
      .hasValue('foo/secret-nested', 'Partial term in nested path replaced with selected secret');
  });

  test('it should fetch secrets at nested paths', async function (assert) {
    this.keys = ['bar/'];
    await click(input);
    await click(option());
    assert.dom(input).hasValue('foo/', 'Input value updates on select');
    assert.dom(option()).hasText('bar/', 'Suggestions are fetched at new path');

    this.keys = ['baz/'];
    await click(option());
    assert.dom(input).hasValue('foo/bar/', 'Input value updates on select');
    assert.dom(option()).hasText('baz/', 'Suggestions are fetched at new path');

    this.keys = ['nested-secret'];
    await typeIn(input, 'baz/');
    assert.dom(option()).hasText('nested-secret', 'Suggestions are fetched at new path');
  });

  test('it should only render dropdown when suggestions exist', async function (assert) {
    await click(input);
    assert.dom(options).exists({ count: 2 }, 'Suggestions render');

    await fillIn(input, 'bar');
    assert.dom(options).doesNotExist('Drop down is hidden when there are no suggestions');

    await fillIn(input, '');
    assert.dom(options).exists({ count: 2 }, 'Suggestions render');

    this.keys = [];
    await click(option());
    assert.dom(options).doesNotExist('Drop down is hidden when there are no suggestions');
  });
});
