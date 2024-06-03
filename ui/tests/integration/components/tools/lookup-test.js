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
import { TOOLS_SELECTORS as TS } from 'vault/tests/helpers/tools-selectors';
import { format } from 'date-fns';

module('Integration | Component | tools/lookup', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.onClear = sinon.spy();
    this.renderComponent = async () => {
      await render(hbs`
    <Tools::Lookup
      @onClear={{this.onClear}}
      @creation_time={{this.creation_time}}
      @creation_ttl={{this.creation_ttl}}
      @creation_path={{this.creation_path}}
      @token={{this.token}}
      @errors={{this.errors}}
    />`);
    };
  });

  test('it renders defaults', async function (assert) {
    await this.renderComponent();

    assert.dom('h1').hasText('Lookup Token', 'Title renders');
    assert.dom('label').hasText('Wrapped token');
    assert.dom(TS.submit).hasText('Lookup token');
    assert.dom(TS.toolsInput('wrapping-token')).hasValue('');
    assert.dom(TS.button('Done')).doesNotExist();
  });

  test('it renders errors', async function (assert) {
    this.errors = ['Something is wrong!'];
    await this.renderComponent();
    assert.dom(GENERAL.messageError).hasText('Error Something is wrong!', 'Error renders');
  });

  test('it renders lookup details', async function (assert) {
    // not stubbing the timestamp util here because this component uses the date-fns formatDistanceToNow method
    // so we need an actual now date for testing (which is why we don't assert the timestamp below, just the day of the month)
    const now = new Date();
    this.creation_path = 'sys/wrapping/wrap';
    this.creation_time = now.toISOString();
    this.creation_ttl = 1800;
    format;
    await this.renderComponent();
    assert.dom(GENERAL.infoRowValue('Creation path')).hasText(this.creation_path);
    assert.dom(GENERAL.infoRowValue('Creation time')).hasText(this.creation_time);
    assert.dom(GENERAL.infoRowValue('Creation TTL')).hasText(`${this.creation_ttl}`);
    assert.dom(GENERAL.infoRowValue('Expiration date')).hasTextContaining(format(now, 'MMM dd yyyy')); // intentionally exclude time to avoid race conditions
    // remove below assertion if flaky (but unlikely this test would take longer than a minute..)
    assert.dom(GENERAL.infoRowValue('Expires in')).hasText('30 minutes'); // from 1800s ttl
    await click(TS.button('Done'));
    assert.true(this.onClear.calledOnce, 'onClear is called');
  });

  test('it calls updates arg when inputs change', async function (assert) {
    await this.renderComponent();
    await fillIn(TS.toolsInput('wrapping-token'), 'my-token');
    assert.strictEqual(this.token, 'my-token', '@token updates when input changes');
  });
});
