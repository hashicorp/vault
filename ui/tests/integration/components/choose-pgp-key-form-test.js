/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import sinon from 'sinon';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const CHOOSE_PGP = {
  begin: '[data-test-dr-token-flow-step="choose-pgp-key"]',
  description: '[data-test-choose-pgp-key-description]',
  confirm: '[data-test-pgp-key-confirm]',
};
module('Integration | Component | choose-pgp-key-form', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.set('onCancel', () => {});
    this.set('onSubmit', () => {});
  });

  test('it renders correctly', async function (assert) {
    await render(
      hbs`<ChoosePgpKeyForm @onSubmit={{this.onSubmit}} @onCancel={{this.onCancel}} @formText="my custom form text" @buttonText="Do it" />`
    );
    assert.dom(CHOOSE_PGP.begin).exists('PGP key selection form exists');
    assert.dom(CHOOSE_PGP.description).hasText('my custom form text', 'uses custom form text');
    await click(GENERAL.textToggle);
    assert.dom(GENERAL.button('use-pgp-key')).isDisabled('use pgp button is disabled');
    await fillIn(GENERAL.textareaByAttr('pgp-key'), 'base64-pgp-key');
    assert.dom(GENERAL.button('use-pgp-key')).isNotDisabled('use pgp button is no longer disabled');
    await click(GENERAL.button('use-pgp-key'));
    assert
      .dom(CHOOSE_PGP.confirm)
      .hasText(
        'Below is the base-64 encoded PGP Key that will be used. Click the "Do it" button to proceed.',
        'Incorporates button text in confirmation'
      );
    assert.dom(GENERAL.copySnippet('pgp-key')).hasText('base64-pgp-key', 'Shows PGP key contents');
    assert.dom(GENERAL.submitButton).hasText('Do it', 'uses passed buttonText');
    await click(GENERAL.submitButton);
  });

  test('it calls onSubmit correctly', async function (assert) {
    const submitSpy = sinon.spy();
    this.set('onSubmit', submitSpy);
    await render(
      hbs`<ChoosePgpKeyForm @onSubmit={{this.onSubmit}} @onCancel={{this.onCancel}} @buttonText="Submit" />`
    );

    assert.dom(CHOOSE_PGP.begin).exists('PGP key selection form exists');
    assert
      .dom(CHOOSE_PGP.description)
      .hasText('Choose a PGP Key from your computer or paste the contents of one in the form below.');
    await click(GENERAL.textToggle);
    assert.dom(GENERAL.button('use-pgp-key')).isDisabled('use pgp button is disabled');
    await fillIn(GENERAL.textareaByAttr('pgp-key'), 'base64-pgp-key');
    assert.dom(GENERAL.button('use-pgp-key')).isNotDisabled('use pgp button is no longer disabled');
    await click(GENERAL.button('use-pgp-key'));
    assert
      .dom(CHOOSE_PGP.confirm)
      .hasText(
        'Below is the base-64 encoded PGP Key that will be used. Click the "Submit" button to proceed.',
        'Confirmation text has buttonText'
      );
    assert.dom(GENERAL.copySnippet('pgp-key')).hasText('base64-pgp-key', 'Shows PGP key contents');
    assert.dom(GENERAL.submitButton).hasText('Submit', 'uses passed buttonText');
    await click(GENERAL.submitButton);
    assert.ok(submitSpy.calledOnceWith('base64-pgp-key'));
  });

  test('it calls cancel on cancel', async function (assert) {
    const cancelSpy = sinon.spy();
    this.set('onCancel', cancelSpy);
    await render(
      hbs`<ChoosePgpKeyForm @onSubmit={{this.onSubmit}} @onCancel={{this.onCancel}} @buttonText="Submit" />`
    );

    await click(GENERAL.textToggle);
    await fillIn(GENERAL.textareaByAttr('pgp-key'), 'base64-pgp-key');
    await click(GENERAL.cancelButton);
    assert.ok(cancelSpy.calledOnce);
  });
});
