/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { find, render, waitUntil } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { addMinutes, subMinutes } from 'date-fns';

module('Integration | Component | token-expire-warning', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders a warning when the token is expired', async function (assert) {
    const expirationDate = subMinutes(Date.now(), 30);
    this.set('expirationDate', expirationDate);

    await render(hbs`<TokenExpireWarning @expirationDate={{this.expirationDate}}/>`);
    await waitUntil(() => find('#modal-overlays'));
    assert.dom().includesText('Your auth token expired on');
  });

  test('it does not render a warning when the token is not expired', async function (assert) {
    const expirationDate = addMinutes(Date.now(), 30);
    this.set('expirationDate', expirationDate);

    await render(hbs`
      <TokenExpireWarning @expirationDate={{this.expirationDate}}>
        <p>Do not worry, your token has not expired.</p>
      </TokenExpireWarning>
    `);
    await waitUntil(() => find('#modal-overlays'));
    assert.dom().doesNotIncludeText('Your auth token expired on');
    assert.dom().includesText('Do not worry');
  });
});
