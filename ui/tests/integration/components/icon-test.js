/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import waitForError from 'vault/tests/helpers/wait-for-error';

module('Integration | Component | icon', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    await render(hbs`<Icon @name="vault-color" class="i-con" />`);
    assert.dom('.i-con').exists('renders');

    // non-flight icon
    await render(hbs`<Icon @name="vault-logo" />`);
    assert.dom('.vault-logo').exists('inlines the SVG');
    assert.dom('.hs-icon').hasClass('hs-icon-l', 'Default hs class applied');

    await render(hbs`<Icon @name="vault-color" class="ah" aria-hidden="true" />`);
    assert.dom('.ah').hasAttribute('aria-hidden', 'true', 'renders aria-hidden');

    await render(hbs`<Icon @name="vault-color" class="al" aria-label="Testing" />`);
    assert.dom('.al').hasAttribute('aria-label', 'Testing', 'renders aria-label');

    // non-flight icon
    await render(hbs`<Icon @name="vault-logo" @size="24"/>`);
    assert.dom('.hs-icon').hasClass('hs-icon-xlm', 'adds the larger size class');
  });

  test('it throws error when size attribute is invalid', async function (assert) {
    const promise = waitForError();
    render(hbs`<Icon @name="vault-color" @size="12"/>`);
    const err = await promise;
    assert.strictEqual(
      err.message,
      'Assertion Failed: Icon component size argument must be either "16" or "24"',
      'Error is thrown when supported size is not provided'
    );
  });

  test('it throws error when name attribute is not passed', async function (assert) {
    const promise = waitForError();
    render(hbs`<Icon />`);
    const err = await promise;
    assert.strictEqual(
      err.message,
      'Assertion Failed: Icon name argument must be provided',
      'Error is thrown when name not provided'
    );
  });

  test('it should render Hds::Icon', async function (assert) {
    assert.expect(3);

    await render(hbs`<Icon @name="x" />`);
    assert.dom('.hds-icon').exists('Hds::Icon renders when provided name of icon in set');
    assert.dom('.hds-icon').hasAttribute('width', '16', 'Default size applied svg');

    await render(hbs`<Icon @name="x" @size="24" />`);
    assert.dom('.hds-icon').hasAttribute('width', '24', 'Size applied to svg');
  });
});
