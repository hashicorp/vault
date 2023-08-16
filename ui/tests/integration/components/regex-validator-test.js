/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import EmberObject from '@ember/object';
import sinon from 'sinon';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn, settled } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | regex-validator', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders input and validation messages', async function (assert) {
    const attr = EmberObject.create({
      name: 'example',
    });
    const spy = sinon.spy();
    this.set('onChange', spy);
    this.set('attr', attr);
    this.set('value', '(\\d{4})');
    this.set('labelString', 'Regex Example');

    await render(
      hbs`<RegexValidator
        @onChange={{this.onChange}}
        @attr={{this.attr}}
        @value={{this.value}}
        @labelString={{this.labelString}}
      />`
    );
    assert.dom('.regex-label label').hasText('Regex Example', 'Label is correct');
    assert.dom('[data-test-toggle-input="example-validation-toggle"]').exists('Validation toggle exists');
    assert.dom('[data-test-regex-validator-test-string]').doesNotExist('Test string input does not show');

    await click('[data-test-toggle-input="example-validation-toggle"]');
    assert.dom('[data-test-regex-validator-test-string]').exists('Test string input shows after toggle');
    assert
      .dom('[data-test-regex-validator-test-string] label')
      .hasText('Test string', 'Test input label renders');
    assert
      .dom('[data-test-regex-validator-test-string] .sub-text')
      .doesNotExist('Test input sub text is hidden when not provided');
    assert
      .dom('[data-test-regex-validation-message]')
      .doesNotExist('Validation message does not show if test string is empty');

    await fillIn('[data-test-input="regex-test-val"]', '123a');
    assert.dom('[data-test-regex-validation-message]').exists('Validation message shows after input filled');
    assert
      .dom('[data-test-inline-error-message]')
      .hasText(
        'This test string does not match the pattern regex.',
        'Shows error when regex does not match string'
      );

    await fillIn('[data-test-input="regex-test-val"]', '1234');
    assert
      .dom('[data-test-inline-success-message]')
      .hasText('This test string matches the pattern regex.', 'Shows success when regex matches');

    await fillIn('[data-test-input="regex-test-val"]', '12345');
    assert
      .dom('[data-test-inline-error-message]')
      .hasText(
        'This test string does not match the pattern regex.',
        "Shows error if regex doesn't match complete string"
      );
    await fillIn('[data-test-input="example"]', '(\\d{5})');
    assert.ok(spy.calledOnce, 'Calls the passed onChange function when main input is changed');
  });

  test('it renders test input only when attr is not provided', async function (assert) {
    this.setProperties({
      value: null,
      label: 'Sample input',
      subText: 'Some text to further describe the input',
    });

    await render(hbs`
      <RegexValidator
        @value={{this.value}}
        @testInputLabel={{this.label}}
        @testInputSubText={{this.subText}}
      />`);

    assert
      .dom('[data-test-regex-validator-pattern]')
      .doesNotExist('Pattern input is hidden when attr is not provided');
    assert
      .dom('[data-test-regex-validator-test-string] label')
      .hasText(this.label, 'Test input label renders');
    assert
      .dom('[data-test-regex-validator-test-string] .sub-text')
      .hasText(this.subText, 'Test input sub text renders');

    await fillIn('[data-test-input="regex-test-val"]', 'test');
    assert
      .dom('[data-test-inline-error-message]')
      .hasText(
        'A pattern has not been entered. Enter a pattern to check this sample input against it.',
        'Warning renders when test input has value but not regex exists'
      );

    this.set('value', 'test');
    assert
      .dom('[data-test-inline-success-message]')
      .hasText('This test string matches the pattern regex.', 'Shows success when regex matches');

    this.set('value', 'foo');
    await settled();
    assert
      .dom('[data-test-inline-error-message]')
      .hasText(
        'This test string does not match the pattern regex.',
        'Pattern is validated on external value change'
      );
  });

  test('it renders capture groups', async function (assert) {
    this.set('value', '(test)(?<last>\\d?)');

    await render(hbs`
      <RegexValidator
        @value={{this.value}}
        @showGroups={{true}}
      />`);
    await fillIn('[data-test-input="regex-test-val"]', 'foobar');
    assert
      .dom('[data-test-regex-validator-groups-placeholder]')
      .exists('Placeholder is shown when regex does not match test input value');
    await fillIn('[data-test-input="regex-test-val"]', 'test8');
    assert.dom('[data-test-regex-group-position="$1"]').hasText('$1', 'First capture group position renders');
    assert.dom('[data-test-regex-group-value="$1"]').hasText('test', 'First capture group value renders');
    assert
      .dom('[data-test-regex-group-position="$2"]')
      .hasText('$2', 'Second capture group position renders');
    assert.dom('[data-test-regex-group-value="$2"]').hasText('8', 'Second capture group value renders');
    assert.dom('[data-test-regex-group-position="$last"]').hasText('$last', 'Named capture group renders');
    assert.dom('[data-test-regex-group-value="$last"]').hasText('8', 'Named capture group value renders');
  });
});
