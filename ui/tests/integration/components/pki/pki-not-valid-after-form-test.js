import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, typeIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { SELECTORS } from 'vault/tests/helpers/pki/pki-not-valid-after-form';

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
    assert.dom(SELECTORS.ttlForm).exists('shows the TTL picker');
    assert.dom(SELECTORS.ttlTimeInput).hasValue('', 'default TTL is empty');
    assert.dom(SELECTORS.radioTtl).isChecked('ttl is selected by default');
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
    assert.dom(SELECTORS.radioTtl).isChecked('notBeforeDate radio is selected');
    assert.dom(SELECTORS.ttlForm).exists({ count: 1 }, 'shows TTL form');
    assert.dom(SELECTORS.radioDate).isNotChecked('NotAfter selection not checked');
    assert.dom(SELECTORS.dateInput).doesNotExist('does not show date input field');

    await click(SELECTORS.radioDateLabel);

    assert.dom(SELECTORS.radioDate).isChecked('selects NotAfter radio when label clicked');
    assert.dom(SELECTORS.dateInput).exists({ count: 1 }, 'shows date input field');
    assert.dom(SELECTORS.radioTtl).isNotChecked('notBeforeDate radio is deselected');
    assert.dom(SELECTORS.ttlForm).doesNotExist('hides TTL form');

    const utcDate = '1994-11-05';
    const notAfterExpected = '1994-11-05T00:00:00.000Z';
    const ttlDate = 1;
    await typeIn('[data-test-input="not_after"]', utcDate);
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
    await typeIn('[data-test-ttl-value="TTL"]', ttlDate);
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
    assert.dom(SELECTORS.radioTtl).isChecked('notBeforeDate radio is selected');
    assert.dom(SELECTORS.ttlForm).exists({ count: 1 }, 'shows TTL form');
    assert.dom(SELECTORS.radioDate).isNotChecked('NotAfter selection not checked');
    assert.dom(SELECTORS.dateInput).doesNotExist('does not show date input field');

    assert.dom(SELECTORS.ttlTimeInput).hasValue('100', 'TTL value is correctly shown');
    assert.dom(SELECTORS.ttlUnitInput).hasValue('m', 'TTL unit is correctly shown');
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
    assert.dom(SELECTORS.radioDate).isChecked('notAfter radio is selected');
    assert.dom(SELECTORS.dateInput).exists({ count: 1 }, 'shows date picker');
    assert.dom(SELECTORS.radioTtl).isNotChecked('ttl radio not selected');
    assert.dom(SELECTORS.ttlForm).doesNotExist('does not show date TTL picker');
    // Due to timezones, can't check specific match on input date
    assert.dom(SELECTORS.dateInput).hasAnyValue('date input shows date');
  });
});
