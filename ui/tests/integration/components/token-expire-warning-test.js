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
    assert.dom('[data-test-token-expired-banner]').includesText('Your auth token expired on');
  });

  test('it does not render a warning when the token is not expired', async function (assert) {
    const expirationDate = addMinutes(Date.now(), 30);
    this.set('expirationDate', expirationDate);

    await render(hbs`
      <TokenExpireWarning @expirationDate={{this.expirationDate}}>
        <p>Do not worry, your token has not expired.</p>
      </TokenExpireWarning>
    `);
    assert.dom().doesNotIncludeText('Your auth token expired on');
    assert.dom().includesText('Do not worry');
  });

  test('it renders a warning when the token is no longer renewing', async function (assert) {
    const expirationDate = addMinutes(Date.now(), 3);
    this.set('expirationDate', expirationDate);
    this.set('allowingExpiration', false);

    await render(
      hbs`
      <TokenExpireWarning @expirationDate={{this.expirationDate}} @allowingExpiration={{this.allowingExpiration}}>
        <p data-test-content>This is the content</p>
      </TokenExpireWarning>
    `
    );
    assert.dom('[data-test-token-expired-banner]').doesNotExist('Does not show token expired banner');
    assert.dom('[data-test-token-expiring-banner]').doesNotExist('Does not show token expiring banner');
    assert.dom('[data-test-content]').hasText('This is the content');

    await this.set('allowingExpiration', true);
    assert.dom('[data-test-token-expired-banner]').doesNotExist('Does not show token expired banner');
    assert.dom('[data-test-token-expiring-banner]').exists('Shows token expiring banner');
    assert.dom('[data-test-content]').hasText('This is the content');
  });

  test('Does not render a warning if no expiration date', async function (assert) {
    this.set('expirationDate', null);
    this.set('allowingExpiration', true);

    await render(
      hbs`
      <TokenExpireWarning @expirationDate={{this.expirationDate}} @allowingExpiration={{this.allowingExpiration}}>
        <p data-test-content>This is the content</p>
      </TokenExpireWarning>
    `
    );
    assert.dom('[data-test-token-expired-banner]').doesNotExist('Does not show token expired banner');
    assert.dom('[data-test-token-expiring-banner]').doesNotExist('Does not show token expiring banner');
    assert.dom('[data-test-content]').hasText('This is the content');
  });
});
