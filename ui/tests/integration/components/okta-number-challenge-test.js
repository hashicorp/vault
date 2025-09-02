/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | auth | okta-number-challenge', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.oktaNumberChallengeAnswer = null;
    this.hasError = false;
    this.onCancel = sinon.spy();
    this.renderComponent = async () => {
      return render(hbs`
  <OktaNumberChallenge
    @correctAnswer={{this.oktaNumberChallengeAnswer}}
    @hasError={{this.hasError}}
    @onReturnToLogin={{this.onCancel}}
  />
  `);
    };
  });

  test('it should render correct descriptions', async function (assert) {
    await this.renderComponent();
    assert
      .dom('[data-test-okta-number-challenge-description]')
      .includesText(
        'To finish signing in, you will need to complete an additional MFA step.',
        'Correct description renders'
      );
    assert.dom('[data-test-loading]').includesText('Please wait...', 'Correct loading description renders');
  });

  test('it should show correct number for okta number challenge', async function (assert) {
    this.oktaNumberChallengeAnswer = 1;
    await this.renderComponent();
    assert
      .dom('[data-test-okta-number-challenge-description]')
      .includesText(
        'To finish signing in, you will need to complete an additional MFA step.',
        'Correct description renders'
      );
    assert
      .dom('[data-test-verification-type]')
      .includesText('Okta verification', 'Correct verification type renders');

    assert
      .dom('[data-test-description]')
      .includesText(
        'Select the following number to complete verification:',
        'Correct verification description renders'
      );
    assert.dom('[data-test-answer]').includesText('1', 'Correct okta number challenge answer renders');
  });

  test('it should show error screen', async function (assert) {
    this.hasError = 'Authentication failed: multi-factor authentication denied';
    await this.renderComponent();

    assert
      .dom('[data-test-okta-number-challenge-description]')
      .hasTextContaining(
        'To finish signing in, you will need to complete an additional MFA step.',
        'Correct description renders'
      );
    assert.dom('[data-test-message-error]').hasText(`Error ${this.hasError}`);
    await click(GENERAL.backButton);
    assert.true(this.onCancel.calledOnce, 'onCancel is called');
  });
});
