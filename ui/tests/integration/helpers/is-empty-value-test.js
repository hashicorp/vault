/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

const template = hbs`
{{#if (is-empty-value this.inputValue hasDefault=this.defaultValue)}}
Empty
{{else}}
Full
{{/if}}
`;

const emptyObject = {};

const nonEmptyObject = { thing: 0 };

module('Integration | Helper | is-empty-value', function (hooks) {
  setupRenderingTest(hooks);

  test('it is truthy if the value evaluated is undefined and no default', async function (assert) {
    this.set('inputValue', undefined);
    this.set('defaultValue', false);

    await render(template);

    assert.dom(this.element).hasText('Empty');
  });

  test('it is truthy if the value evaluated is an empty string and no default', async function (assert) {
    this.set('inputValue', '');
    this.set('defaultValue', false);

    await render(template);

    assert.dom(this.element).hasText('Empty');
  });

  test('it is truthy if the value evaluated is an empty object and no default', async function (assert) {
    this.set('inputValue', emptyObject);
    this.set('defaultValue', false);

    await render(template);

    assert.dom(this.element).hasText('Empty');
  });

  test('it is falsy if the value evaluated is not an empty object and no default', async function (assert) {
    this.set('inputValue', nonEmptyObject);
    this.set('defaultValue', false);

    await render(template);

    assert.dom(this.element).hasText('Full');
  });

  test('it is falsy if the value evaluated is empty but a default exists', async function (assert) {
    this.set('defaultValue', 'Some default');
    this.set('inputValue', emptyObject);

    await render(template);
    assert.dom(this.element).hasText('Full', 'shows default when value is empty object');

    this.set('inputValue', '');
    await render(template);
    assert.dom(this.element).hasText('Full', 'shows default when value is empty string');

    this.set('inputValue', undefined);
    await render(template);
    assert.dom(this.element).hasText('Full', 'shows default when value is undefined');
  });
});
