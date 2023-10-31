/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';

const SEAL_WHEN_STANDBY_MSG = 'vault cannot seal when in standby mode; please restart instead';

module('Integration | Component | seal-action', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.sealSuccess = sinon.spy(() => new Promise((resolve) => resolve({})));
    this.sealError = sinon.stub().throws({ message: SEAL_WHEN_STANDBY_MSG });
  });

  test('it handles success', async function (assert) {
    this.set('handleSeal', this.sealSuccess);
    await render(hbs`<SealAction @onSeal={{action this.handleSeal}} />`);

    // attempt seal
    await click('[data-test-seal] button');
    await click('[data-test-confirm-button]');

    assert.ok(this.sealSuccess.calledOnce, 'called onSeal action');
    assert.dom('[data-test-seal-error]').doesNotExist('Does not show error when successful');
  });

  test('it handles error', async function (assert) {
    this.set('handleSeal', this.sealError);
    await render(hbs`<SealAction @onSeal={{action this.handleSeal}} />`);

    // attempt seal
    await click('[data-test-seal] button');
    await click('[data-test-confirm-button]');

    assert.ok(this.sealError.calledOnce, 'called onSeal action');
    assert.dom('[data-test-seal-error]').includesText(SEAL_WHEN_STANDBY_MSG, 'Shows error returned from API');
  });
});
