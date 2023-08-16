/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import sinon from 'sinon';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

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

    assert.dom('[data-test-choose-pgp-key-form="begin"]').exists('PGP key selection form exists');
    assert
      .dom('[data-test-choose-pgp-key-description]')
      .hasText('my custom form text', 'uses custom form text');
    await click('[data-test-text-toggle]');
    assert.dom('[data-test-use-pgp-key-button]').isDisabled('use pgp button is disabled');
    await fillIn('[data-test-pgp-file-textarea]', 'base64-pgp-key');
    assert.dom('[data-test-use-pgp-key-button]').isNotDisabled('use pgp button is no longer disabled');
    await click('[data-test-use-pgp-key-button]');
    assert
      .dom('[data-test-pgp-key-confirm]')
      .hasText(
        'Below is the base-64 encoded PGP Key that will be used. Click the "Do it" button to proceed.',
        'Incorporates button text in confirmation'
      );
    assert.dom('[data-test-pgp-key-copy]').hasText('base64-pgp-key', 'Shows PGP key contents');
    assert.dom('[data-test-confirm-pgp-key-submit]').hasText('Do it', 'uses passed buttonText');
    await click('[data-test-confirm-pgp-key-submit]');
  });
  test('it calls onSubmit correctly', async function (assert) {
    const submitSpy = sinon.spy();
    this.set('onSubmit', submitSpy);
    await render(
      hbs`<ChoosePgpKeyForm @onSubmit={{this.onSubmit}} @onCancel={{this.onCancel}} @buttonText="Submit" />`
    );

    assert.dom('[data-test-choose-pgp-key-form="begin"]').exists('PGP key selection form exists');
    assert
      .dom('[data-test-choose-pgp-key-description]')
      .hasText('Choose a PGP Key from your computer or paste the contents of one in the form below.');
    await click('[data-test-text-toggle]');
    assert.dom('[data-test-use-pgp-key-button]').isDisabled('use pgp button is disabled');
    await fillIn('[data-test-pgp-file-textarea]', 'base64-pgp-key');
    assert.dom('[data-test-use-pgp-key-button]').isNotDisabled('use pgp button is no longer disabled');
    await click('[data-test-use-pgp-key-button]');
    assert
      .dom('[data-test-pgp-key-confirm]')
      .hasText(
        'Below is the base-64 encoded PGP Key that will be used. Click the "Submit" button to proceed.',
        'Confirmation text has buttonText'
      );
    assert.dom('[data-test-pgp-key-copy]').hasText('base64-pgp-key', 'Shows PGP key contents');
    assert.dom('[data-test-confirm-pgp-key-submit]').hasText('Submit', 'uses passed buttonText');
    await click('[data-test-confirm-pgp-key-submit]');
    assert.ok(submitSpy.calledOnceWith('base64-pgp-key'));
  });

  test('it calls cancel on cancel', async function (assert) {
    const cancelSpy = sinon.spy();
    this.set('onCancel', cancelSpy);
    await render(
      hbs`<ChoosePgpKeyForm @onSubmit={{this.onSubmit}} @onCancel={{this.onCancel}} @buttonText="Submit" />`
    );

    await click('[data-test-text-toggle]');
    await fillIn('[data-test-pgp-file-textarea]', 'base64-pgp-key');
    await click('[data-test-use-pgp-key-cancel]');
    assert.ok(cancelSpy.calledOnce);
  });
});
