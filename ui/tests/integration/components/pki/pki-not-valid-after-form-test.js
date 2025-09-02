/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { PKI_NOT_VALID_AFTER } from 'vault/tests/helpers/pki/pki-selectors';

module('Integration | Component | pki-not-valid-after-form', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.model = this.store.createRecord('pki/role', { backend: 'pki' });
    this.attr = {
      helpText: '',
      options: {
        helperTextEnabled: 'toggled on and shows text',
      },
    };
  });

  test('it should render the component with ttl selected by default', async function (assert) {
    assert.expect(3);
    await render(
      hbs`
      <div class="has-top-margin-xxl">
        <PkiNotValidAfterForm
          @model={{this.model}}
          @attr={{this.attr}}
        />
       </div>
  `,
      { owner: this.engine }
    );
    assert.dom(PKI_NOT_VALID_AFTER.ttlForm).exists('shows the TTL picker');
    assert.dom(PKI_NOT_VALID_AFTER.ttlTimeInput).hasValue('', 'default TTL is empty');
    assert.dom(PKI_NOT_VALID_AFTER.radioTtl).isChecked('ttl is selected by default');
  });

  test('it clears and resets model properties from cache when changing radio selection', async function (assert) {
    await render(
      hbs`
      <div class="has-top-margin-xxl">
        <PkiNotValidAfterForm
          @model={{this.model}}
          @attr={{this.attr}}
        />
       </div>
  `,
      { owner: this.engine }
    );
    assert.dom(PKI_NOT_VALID_AFTER.radioTtl).isChecked('notBeforeDate radio is selected');
    assert.dom(PKI_NOT_VALID_AFTER.ttlForm).exists({ count: 1 }, 'shows TTL form');
    assert.dom(PKI_NOT_VALID_AFTER.radioDate).isNotChecked('NotAfter selection not checked');
    assert.dom(PKI_NOT_VALID_AFTER.dateInput).doesNotExist('does not show date input field');

    await click(PKI_NOT_VALID_AFTER.radioDateLabel);

    assert.dom(PKI_NOT_VALID_AFTER.radioDate).isChecked('selects NotAfter radio when label clicked');
    assert.dom(PKI_NOT_VALID_AFTER.dateInput).exists({ count: 1 }, 'shows date input field');
    assert.dom(PKI_NOT_VALID_AFTER.radioTtl).isNotChecked('notBeforeDate radio is deselected');
    assert.dom(PKI_NOT_VALID_AFTER.ttlForm).doesNotExist('hides TTL form');

    const utcDate = '1994-11-05';
    const notAfterExpected = '1994-11-05T00:00:00.000Z';
    const ttlDate = 1;
    await fillIn('[data-test-input="not_after"]', utcDate);
    assert.strictEqual(
      this.model.notAfter,
      notAfterExpected,
      'sets the model property notAfter when this value is selected and filled in.'
    );
    await click('[data-test-radio-button="ttl"]');
    assert.strictEqual(
      this.model.notAfter,
      '',
      'The notAfter is cleared on the model because the radio button was selected.'
    );
    await fillIn('[data-test-ttl-value="TTL"]', ttlDate);
    assert.strictEqual(
      this.model.ttl,
      '1s',
      'The ttl is now saved on the model because the radio button was selected.'
    );

    await click('[data-test-radio-button="not_after"]');
    assert.strictEqual(this.model.ttl, '', 'TTL is cleared after radio select.');
    assert.strictEqual(this.model.notAfter, notAfterExpected, 'notAfter gets populated from local cache');
  });
  test('Form renders properly for edit when TTL present', async function (assert) {
    this.model = this.store.createRecord('pki/role', { backend: 'pki', ttl: 6000 });
    await render(
      hbs`
      <div class="has-top-margin-xxl">
        <PkiNotValidAfterForm
          @model={{this.model}}
          @attr={{this.attr}}
        />
       </div>
  `,
      { owner: this.engine }
    );
    assert.dom(PKI_NOT_VALID_AFTER.radioTtl).isChecked('notBeforeDate radio is selected');
    assert.dom(PKI_NOT_VALID_AFTER.ttlForm).exists({ count: 1 }, 'shows TTL form');
    assert.dom(PKI_NOT_VALID_AFTER.radioDate).isNotChecked('NotAfter selection not checked');
    assert.dom(PKI_NOT_VALID_AFTER.dateInput).doesNotExist('does not show date input field');

    assert.dom(PKI_NOT_VALID_AFTER.ttlTimeInput).hasValue('100', 'TTL value is correctly shown');
    assert.dom(PKI_NOT_VALID_AFTER.ttlUnitInput).hasValue('m', 'TTL unit is correctly shown');
  });
  test('Form renders properly for edit when notAfter present', async function (assert) {
    const utcDate = '1994-11-05T00:00:00.000Z';
    this.model = this.store.createRecord('pki/role', { backend: 'pki', notAfter: utcDate });
    await render(
      hbs`
      <div class="has-top-margin-xxl">
        <PkiNotValidAfterForm
          @model={{this.model}}
          @attr={{this.attr}}
        />
       </div>
  `,
      { owner: this.engine }
    );
    assert.dom(PKI_NOT_VALID_AFTER.radioDate).isChecked('notAfter radio is selected');
    assert.dom(PKI_NOT_VALID_AFTER.dateInput).exists({ count: 1 }, 'shows date picker');
    assert.dom(PKI_NOT_VALID_AFTER.radioTtl).isNotChecked('ttl radio not selected');
    assert.dom(PKI_NOT_VALID_AFTER.ttlForm).doesNotExist('does not show date TTL picker');
    // Due to timezones, can't check specific match on input date
    assert.dom(PKI_NOT_VALID_AFTER.dateInput).hasAnyValue('date input shows date');
  });
});
