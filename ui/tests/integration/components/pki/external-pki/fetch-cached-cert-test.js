/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, fillIn, select, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | pki | external-pki | ExternalPki::FetchCachedCert', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.onSubmit = sinon.spy();
    this.errorMessage = undefined;
    this.renderComponent = () =>
      render(
        hbs`<ExternalPki::FetchCachedCert @onSubmit={{this.onSubmit}} @errorMessage={{this.errorMessage}} />`,
        { owner: this.engine }
      );
  });

  test('it renders the card with form elements', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.overviewCard.container('Retrieve cached certificate')).exists('form is rendered');
    assert.dom(GENERAL.inputByAttr('identifiers')).exists('identifiers input is rendered');
    assert.dom(GENERAL.radioByAttr('duration')).exists('duration radio is rendered').isChecked();
    assert.dom(GENERAL.radioByAttr('percentage')).exists('percentage radio is rendered');
    assert.dom(GENERAL.inputByAttr('minValidityValue')).exists('min validity input is rendered');
    assert.dom(GENERAL.selectByAttr('unit')).exists('unit select is rendered by default (duration mode)');
    assert.dom(GENERAL.submitButton).exists('submit button is rendered');
    assert.dom(GENERAL.messageError).doesNotExist();
  });

  test('it renders an error message when @errorMessage is provided', async function (assert) {
    this.errorMessage = 'Something went wrong';
    await this.renderComponent();
    assert.dom(GENERAL.messageError).exists('error banner is rendered').containsText('Something went wrong');
  });

  test('selecting percentage radio swaps unit select for disabled % input', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.radioByAttr('duration')).exists().isChecked();
    assert.dom(GENERAL.radioByAttr('percentage')).exists().isNotChecked();
    assert.dom(GENERAL.selectByAttr('unit')).exists();
    assert.dom(GENERAL.inputByAttr('percentage')).doesNotExist('% input absent in duration mode');

    await click(GENERAL.radioByAttr('percentage'));
    assert.dom(GENERAL.radioByAttr('percentage')).isChecked();
    assert.dom(GENERAL.radioByAttr('duration')).isNotChecked();
    assert.dom(GENERAL.selectByAttr('unit')).doesNotExist('unit select hidden in percentage mode');
    assert.dom(GENERAL.inputByAttr('percentage')).exists('% input shown in percentage mode');
    assert.dom(GENERAL.inputByAttr('percentage')).isDisabled('% input is disabled (label only)');
  });

  test('switching back to duration restores the unit select', async function (assert) {
    await this.renderComponent();

    await click(GENERAL.radioByAttr('percentage'));
    await click(GENERAL.radioByAttr('duration'));

    assert.dom(GENERAL.selectByAttr('unit')).exists('unit select restored after switching back to duration');
    assert.dom(GENERAL.inputByAttr('percentage')).doesNotExist('% input gone after switching back');
  });

  test('the min validity value is preserved when switching between modes', async function (assert) {
    await this.renderComponent();

    await fillIn(GENERAL.inputByAttr('minValidityValue'), '42');
    await click(GENERAL.radioByAttr('percentage'));

    assert
      .dom(GENERAL.inputByAttr('minValidityValue'))
      .hasValue('42', 'value is preserved across mode switch');
  });

  // ---------------------------------------------------------------------------
  // Validation errors
  // ---------------------------------------------------------------------------

  test('it shows a validation error when identifiers is empty on submit', async function (assert) {
    await this.renderComponent();

    await click(GENERAL.submitButton);

    assert
      .dom(GENERAL.validationErrorByAttr('identifiers'))
      .exists('identifiers error shown')
      .containsText('You must provide at least one identifier.');
  });

  test('it shows a duration validation error when min validity is empty', async function (assert) {
    await this.renderComponent();

    await fillIn(GENERAL.inputByAttr('identifiers'), 'example.com');
    await click(GENERAL.submitButton);

    assert
      .dom(GENERAL.validationErrorByAttr('minValidityValue'))
      .exists('min validity error shown')
      .hasText('A minimum validity is required. Duration must be a positive number.');
  });

  test('it shows a duration validation error for non-positive values', async function (assert) {
    await this.renderComponent();

    await fillIn(GENERAL.inputByAttr('identifiers'), 'example.com');
    await fillIn(GENERAL.inputByAttr('minValidityValue'), '0');
    await click(GENERAL.submitButton);

    assert
      .dom(GENERAL.validationErrorByAttr('minValidityValue'))
      .containsText('Duration must be a positive number.');
  });

  test('it shows a percentage validation error when min validity is empty', async function (assert) {
    await this.renderComponent();

    await fillIn(GENERAL.inputByAttr('identifiers'), 'example.com');
    await click(GENERAL.radioByAttr('percentage'));
    await click(GENERAL.submitButton);

    assert
      .dom(GENERAL.validationErrorByAttr('minValidityValue'))
      .hasText('A minimum validity is required. Percentage must be a number between 1 and 100.');
  });

  test('it shows a percentage validation error when value is out of range', async function (assert) {
    await this.renderComponent();

    await click(GENERAL.radioByAttr('percentage'));
    await fillIn(GENERAL.inputByAttr('identifiers'), 'example.com');
    await fillIn(GENERAL.inputByAttr('minValidityValue'), '150');
    await click(GENERAL.submitButton);

    assert
      .dom(GENERAL.validationErrorByAttr('minValidityValue'))
      .containsText('Percentage must be a number between 1 and 100.');
  });

  test('switching modes clears existing validation errors', async function (assert) {
    await this.renderComponent();

    await fillIn(GENERAL.inputByAttr('identifiers'), 'example.com');
    await click(GENERAL.submitButton);
    assert.dom(GENERAL.validationErrorByAttr('minValidityValue')).exists('error is present before switch');

    await click(GENERAL.radioByAttr('percentage'));
    assert
      .dom(GENERAL.validationErrorByAttr('minValidityValue'))
      .doesNotExist('error cleared after mode switch');
  });

  test('re-submitting updates existing validation errors', async function (assert) {
    await this.renderComponent();

    await fillIn(GENERAL.inputByAttr('identifiers'), 'example.com');
    await click(GENERAL.submitButton);
    assert
      .dom(GENERAL.validationErrorByAttr('minValidityValue'))
      .hasText('A minimum validity is required. Duration must be a positive number.');

    await click(GENERAL.radioByAttr('percentage'));
    await click(GENERAL.submitButton);
    assert
      .dom(GENERAL.validationErrorByAttr('minValidityValue'))
      .hasText('A minimum validity is required. Percentage must be a number between 1 and 100.');
  });

  // ---------------------------------------------------------------------------
  // Payload submitting
  // ---------------------------------------------------------------------------

  test('it submits form payload with min_validity_duration', async function (assert) {
    await this.renderComponent();

    await fillIn(GENERAL.inputByAttr('identifiers'), 'example.com');
    await fillIn(GENERAL.inputByAttr('minValidityValue'), '45');
    await click(GENERAL.submitButton);

    assert.true(this.onSubmit.calledOnce, 'onSubmit was called');
    const [payload] = this.onSubmit.lastCall.args;
    const expected = { identifiers: 'example.com', min_validity_duration: 45 };
    assert.propEqual(payload, expected, `it submits with expected payload: ${JSON.stringify(expected)}`);
  });

  test('it submits form payload with min_validity_percentage', async function (assert) {
    await this.renderComponent();

    await fillIn(GENERAL.inputByAttr('identifiers'), 'example.com');
    await fillIn(GENERAL.inputByAttr('minValidityValue'), '45');
    await click(GENERAL.radioByAttr('percentage'));
    await click(GENERAL.submitButton);

    assert.true(this.onSubmit.calledOnce, 'onSubmit was called');
    const [payload] = this.onSubmit.lastCall.args;
    const expected = { identifiers: 'example.com', min_validity_percentage: 45 };
    assert.propEqual(payload, expected, `it submits with expected payload: ${JSON.stringify(expected)}`);
  });

  const durationCases = [
    { value: '30', unit: 's', expected: 30, label: 'seconds' },
    { value: '5', unit: 'm', expected: 300, label: 'minutes to seconds' },
    { value: '2', unit: 'h', expected: 7200, label: 'hours to seconds' },
    { value: '1', unit: 'd', expected: 86400, label: 'days to seconds' },
  ];

  for (const { value, unit, expected, label } of durationCases) {
    test(`it submits min_validity_duration converted to seconds: ${label}`, async function (assert) {
      assert.expect(1);

      await this.renderComponent();

      await fillIn(GENERAL.inputByAttr('identifiers'), 'example.com');
      await fillIn(GENERAL.inputByAttr('minValidityValue'), value);
      await select(GENERAL.selectByAttr('unit'), unit);
      await click(GENERAL.submitButton);

      const [payload] = this.onSubmit.lastCall.args;
      assert.strictEqual(
        payload.min_validity_duration,
        expected,
        `payload carries min_validity_duration=${expected} for ${value}${unit}`
      );
    });
  }

  test('it does not submit when form is invalid', async function (assert) {
    await this.renderComponent();

    await click(GENERAL.submitButton);
    assert.true(this.onSubmit.notCalled, 'onSubmit is not called when validation fails');
  });
});
