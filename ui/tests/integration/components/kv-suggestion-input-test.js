/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn, typeIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

const {
  searchSelect: { option, options },
  kvSuggestion: { input },
} = GENERAL;

module('Integration | Component | kv-suggestion-input', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.label = 'Input something';
    this.mountPath = 'secret/';
    this.apiStub = sinon
      .stub(this.owner.lookup('service:api').secrets, 'kvV2List')
      .resolves({ keys: ['foo/', 'my-secret'] });

    this.renderComponent = () =>
      render(hbs`
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
    await this.renderComponent();
    assert.dom('label').hasText('Input something', 'Label renders');
    assert.dom('[data-test-label-subtext]').hasText('Suggestions will display', 'Label subtext renders');

    this.label = null;
    await this.renderComponent();
    assert.dom('label').doesNotExist('Label is hidden when arg is not provided');
  });

  test('it should disable input when mountPath is not provided', async function (assert) {
    this.mountPath = null;
    await this.renderComponent();
    assert.dom(input).isDisabled('Input is disabled when mountPath is not provided');
  });

  test('it should fetch suggestions for initial mount path', async function (assert) {
    await this.renderComponent();
    assert.dom(options).doesNotExist('Suggestions are hidden initially');
    await click(input);
    ['foo/', 'my-secret'].forEach((secret, index) => {
      assert.dom(option(index)).hasText(secret, 'Suggestion renders for initial mount path');
    });
  });

  test('it should fetch secrets and update suggestions on mountPath change', async function (assert) {
    this.apiStub.resolves({ keys: ['test1'] });
    this.mountPath = 'foo/';
    await this.renderComponent();
    await click(input);
    assert.dom(option()).hasText('test1', 'Suggestions are fetched and render on mountPath change');
  });

  test('it should filter current result set', async function (assert) {
    await this.renderComponent();
    await click(input);
    await typeIn(input, 'sec');
    assert.dom(options).exists({ count: 1 }, 'Correct number of options render based on input value');
    assert.dom(option()).hasText('my-secret', 'Result set is filtered');
  });

  test('it should replace filter terms with full path to secret', async function (assert) {
    await this.renderComponent();
    await fillIn(input, 'sec');
    await click(option());
    assert.dom(input).hasValue('my-secret', 'Partial term replaced with selected secret');

    await fillIn(input, '');
    this.apiStub.resolves({ keys: ['secret-nested', 'bar', 'baz'] });
    await click(option());
    await fillIn(input, 'nest');
    await click(option());
    assert
      .dom(input)
      .hasValue('foo/secret-nested', 'Partial term in nested path replaced with selected secret');
  });

  test('it should fetch secrets at nested paths', async function (assert) {
    await this.renderComponent();
    this.apiStub.resolves({ keys: ['bar/'] });
    await click(input);
    await click(option());
    assert.dom(input).hasValue('foo/', 'Input value updates on select');
    assert.dom(option()).hasText('bar/', 'Suggestions are fetched at new path');

    this.apiStub.resolves({ keys: ['baz/'] });
    await click(option());
    assert.dom(input).hasValue('foo/bar/', 'Input value updates on select');
    assert.dom(option()).hasText('baz/', 'Suggestions are fetched at new path');

    this.apiStub.resolves({ keys: ['nested-secret'] });
    await typeIn(input, 'baz/');
    assert.dom(option()).hasText('nested-secret', 'Suggestions are fetched at new path');
  });

  test('it should only render dropdown when suggestions exist', async function (assert) {
    await this.renderComponent();
    await click(input);
    assert.dom(options).exists({ count: 2 }, 'Suggestions render');

    await fillIn(input, 'bar');
    assert.dom(options).doesNotExist('Drop down is hidden when there are no suggestions');

    await fillIn(input, '');
    assert.dom(options).exists({ count: 2 }, 'Suggestions render');

    this.apiStub.resolves({ keys: [] });
    await click(option());
    assert.dom(options).doesNotExist('Drop down is hidden when there are no suggestions');
  });
});
