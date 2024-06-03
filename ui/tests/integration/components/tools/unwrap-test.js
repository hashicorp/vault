/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import codemirror from 'vault/tests/helpers/codemirror';
import { TOOLS_SELECTORS as TS } from 'vault/tests/helpers/tools-selectors';

module('Integration | Component | tools/unwrap', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.onClear = sinon.spy();
    this.renderComponent = async () => {
      await render(hbs`
    <Tools::Unwrap
      @onClear={{this.onClear}}
      @unwrap_data={{this.unwrap_data}}
      @details={{this.details}}
      @errors={{this.errors}}
      @token={{this.token}}
    />`);
    };
  });

  test('it renders defaults', async function (assert) {
    await this.renderComponent();

    assert.dom('h1').hasText('Unwrap Data', 'Title renders');
    assert.dom(TS.submit).hasText('Unwrap data');
    assert.dom(TS.toolsInput('wrapping-token')).hasValue('');
    assert.dom(TS.tab('data')).doesNotExist();
    assert.dom(TS.tab('details')).doesNotExist();
    assert.dom('.CodeMirror').doesNotExist();
    assert.dom(TS.button('Done')).doesNotExist();
  });

  test('it renders errors', async function (assert) {
    this.errors = ['Something is wrong!'];
    await this.renderComponent();
    assert.dom(GENERAL.messageError).hasText('Error Something is wrong!', 'Error renders');
  });

  test('it renders unwrapped data', async function (assert) {
    this.unwrap_data = { foo: 'bar' };
    await this.renderComponent();

    assert.dom('label').hasText('Unwrapped Data');
    assert.strictEqual(codemirror().getValue(' '), '{   "foo": "bar" }', 'it renders unwrapped data');
    assert.dom(TS.tab('data')).hasAttribute('aria-selected', 'true');
    await click(TS.button('Done'));
    assert.true(this.onClear.calledOnce, 'onClear is called');
  });

  test('it renders falsy unwrapped details', async function (assert) {
    this.unwrap_data = { foo: 'bar' };
    this.details = {
      'Request ID': '5810d40e-ce93-3e99-1c72-9dbb58ed3c67',
      'Lease ID': 'None',
      Renewable: 'No',
      'Lease Duration': 'None',
    };
    await this.renderComponent();

    await click(TS.tab('details'));
    assert.dom(TS.tab('details')).hasAttribute('aria-selected', 'true');
    assert.dom(GENERAL.icon('x-circle')).exists({ count: 3 }, 'renders falsy icon for each row');
    for (const property in this.details) {
      assert.dom(GENERAL.infoRowValue(property)).hasText(this.details[property]);
    }
  });

  test('it renders truthy unwrapped details', async function (assert) {
    this.unwrap_data = { foo: 'bar' };
    this.details = {
      'Lease ID': 'Yes',
      Renewable: 'Yes',
      'Lease Duration': '5h',
    };
    await this.renderComponent();

    await click(TS.tab('details'));
    assert.dom(TS.tab('details')).hasAttribute('aria-selected', 'true');
    assert.dom(GENERAL.icon('check-circle')).exists({ count: 2 }, 'renders truthy icon for each row');
    for (const property in this.details) {
      assert.dom(GENERAL.infoRowValue(property)).hasText(this.details[property]);
    }
  });

  test('it calls updates arg when inputs change', async function (assert) {
    await this.renderComponent();
    await fillIn(TS.toolsInput('wrapping-token'), 'my-token');
    assert.strictEqual(this.token, 'my-token', '@token updates when input changes');
  });
});
