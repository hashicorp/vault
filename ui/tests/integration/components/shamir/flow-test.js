/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import sinon from 'sinon';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render, settled } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';

// Checks that the correct data were passed around happens in the integration test
// this one is checking that things happen at the right time
module('Integration | Component | shamir/flow', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.server.put('/sys/unseal', () => ({
      sealed: false,
      t: this.threshold,
      n: this.threshold,
      progress: 1,
    }));

    this.keyPart = 'some-key-partition';
    this.progress = 0;
    this.threshold = 3;
    this.updateProgress = sinon.spy();
    this.checkComplete = sinon.stub().returns(false);
    this.onSuccess = sinon.spy();

    this.renderComponent = () =>
      render(hbs`
        <Shamir::Flow
          @action="unseal"
          @threshold={{this.threshold}}
          @progress={{this.progress}}
          @updateProgress={{this.updateProgress}}
          @checkComplete={{this.checkComplete}}
          @onShamirSuccess={{this.onSuccess}}
        />`);
  });

  test('it sends data to the passed action and calls updateProgress', async function (assert) {
    await this.renderComponent();

    await fillIn(GENERAL.inputByAttr('shamir-key'), this.keyPart);
    await click(GENERAL.submitButton);

    assert.true(this.onSuccess.notCalled, 'onShamirSuccess was not called');
    assert.true(this.updateProgress.calledOnce, 'updateProgress was called');
    // Default shamir flow expects the updated values to be passed
    // in from parent model, so this approximates the update happening
    // from a side effect of the updateProgress call
    this.set('progress', 2);
    // Pretend the next call will mean completion
    this.checkComplete.returns(true);
    await settled();

    await fillIn(GENERAL.inputByAttr('shamir-key'), this.keyPart);
    await click(GENERAL.submitButton);

    assert.true(this.onSuccess.calledOnce, 'onShamirSuccess was called');
    assert.true(this.updateProgress.calledTwice, 'updateProgress was called again');
  });

  test('it shows the error when request fails with 400 status', async function (assert) {
    assert.expect(3);

    this.progress = 2;
    this.server.put('/sys/unseal', () => {
      return new Response(
        400,
        { 'Content-Type': 'application/json' },
        JSON.stringify({ errors: ['something is wrong', 'seriously wrong'] })
      );
    });

    await this.renderComponent();
    await fillIn(GENERAL.inputByAttr('shamir-key'), this.keyPart);
    await click(GENERAL.submitButton);
    assert.dom(GENERAL.messageError).exists({ count: 2 }, 'renders errors');
    assert.true(this.checkComplete.notCalled, 'checkComplete was not called');
    assert.true(this.updateProgress.notCalled, 'updateProgress was not called');
  });
});
