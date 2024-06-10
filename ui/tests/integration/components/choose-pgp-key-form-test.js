/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import sinon from 'sinon';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

const CHOOSE_PGP = {
  begin: '[data-test-choose-pgp-key-form="begin"]',
  description: '[data-test-choose-pgp-key-description]',
  toggle: '[data-test-text-toggle]',
  useKeyButton: '[data-test-use-pgp-key-button]',
  pgpTextArea: '[data-test-pgp-file-textarea]',
  confirm: '[data-test-pgp-key-confirm]',
  base64Output: '[data-test-pgp-key-copy]',
  submit: '[data-test-confirm-pgp-key-submit]',
  cancel: '[data-test-use-pgp-key-cancel]',
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
    await click(CHOOSE_PGP.toggle);
    assert.dom(CHOOSE_PGP.useKeyButton).isDisabled('use pgp button is disabled');
    await fillIn(CHOOSE_PGP.pgpTextArea, 'base64-pgp-key');
    assert.dom(CHOOSE_PGP.useKeyButton).isNotDisabled('use pgp button is no longer disabled');
    await click(CHOOSE_PGP.useKeyButton);
    assert
      .dom(CHOOSE_PGP.confirm)
      .hasText(
        'Below is the base-64 encoded PGP Key that will be used. Click the "Do it" button to proceed.',
        'Incorporates button text in confirmation'
      );
    assert.dom(CHOOSE_PGP.base64Output).hasText('base64-pgp-key', 'Shows PGP key contents');
    assert.dom(CHOOSE_PGP.submit).hasText('Do it', 'uses passed buttonText');
    await click(CHOOSE_PGP.submit);
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
    await click(CHOOSE_PGP.toggle);
    assert.dom(CHOOSE_PGP.useKeyButton).isDisabled('use pgp button is disabled');
    await fillIn(CHOOSE_PGP.pgpTextArea, 'base64-pgp-key');
    assert.dom(CHOOSE_PGP.useKeyButton).isNotDisabled('use pgp button is no longer disabled');
    await click(CHOOSE_PGP.useKeyButton);
    assert
      .dom(CHOOSE_PGP.confirm)
      .hasText(
        'Below is the base-64 encoded PGP Key that will be used. Click the "Submit" button to proceed.',
        'Confirmation text has buttonText'
      );
    assert.dom(CHOOSE_PGP.base64Output).hasText('base64-pgp-key', 'Shows PGP key contents');
    assert.dom(CHOOSE_PGP.submit).hasText('Submit', 'uses passed buttonText');
    await click(CHOOSE_PGP.submit);
    assert.ok(submitSpy.calledOnceWith('base64-pgp-key'));
  });

  test('it calls cancel on cancel', async function (assert) {
    const cancelSpy = sinon.spy();
    this.set('onCancel', cancelSpy);
    await render(
      hbs`<ChoosePgpKeyForm @onSubmit={{this.onSubmit}} @onCancel={{this.onCancel}} @buttonText="Submit" />`
    );

    await click(CHOOSE_PGP.toggle);
    await fillIn(CHOOSE_PGP.pgpTextArea, 'base64-pgp-key');
    await click(CHOOSE_PGP.cancel);
    assert.ok(cancelSpy.calledOnce);
  });
});
