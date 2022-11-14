import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';

module('Integration | Component | radio-select-ttl-or-string', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.model = this.store.createRecord('pki/role');
    this.model.backend = 'pki';
    this.attr = {
      helpText: '',
      options: {
        helperTextEnabled: 'toggled on and shows text',
      },
    };
  });

  test('it should render the component and init with ttl selected', async function (assert) {
    assert.expect(3);
    await render(
      hbs`
      <div class="has-top-margin-xxl">
        <RadioSelectTtlOrString
          @model={{this.model}}
          @attr={{this.attr}}
        />
       </div>
  `,
      { owner: this.engine }
    );
    assert.dom('[data-test-input="ttl"]').exists('shows the TTL component');
    const inputValue = document.querySelector('[data-test-ttl-value="TTL"]').value;
    assert.strictEqual(inputValue, '', 'default TTL is empty');
    assert.dom('[data-test-radio-button="ttl"]').isChecked('ttl is selected by default');
  });

  test('it should set the model properties ttl or notAfter based on the radio button selections', async function (assert) {
    assert.expect(8);
    await render(
      hbs`
      <div class="has-top-margin-xxl">
        <RadioSelectTtlOrString
          @model={{this.model}}
          @attr={{this.attr}}
        />
       </div>
  `,
      { owner: this.engine }
    );
    assert.dom('[data-test-input="not_after"]').doesNotExist('does not show input field on initial render');

    await click('[data-test-radio-button="not_after"]');
    assert
      .dom('[data-test-input="not_after"]')
      .exists('does show input field after clicking the radio button');

    const utcDate = '1994-11-05T08:15:30-05:0';
    const ttlDate = 1;
    await fillIn('[data-test-input="not_after"]', utcDate);
    assert.strictEqual(
      this.model.notAfter,
      utcDate,
      'sets the model property notAfter when this value is selected and filled in.'
    );

    await fillIn('[data-test-ttl-value="TTL"]', ttlDate);
    assert.strictEqual(this.model.ttl, '', 'No ttl is set because the radio button was not selected.');

    await click('[data-test-radio-button="ttl"]');
    assert.strictEqual(
      this.model.notAfter,
      '',
      'The notAfter is cleared on the model because the radio button was selected.'
    );
    assert.strictEqual(
      this.model.ttl,
      ttlDate,
      'The ttl is now saved on the model because the radio button was selected.'
    );

    await click('[data-test-radio-button="not_after"]');
    assert.strictEqual(this.model.ttl, '', 'Both ttl and notAfter are cleared.');
    assert.strictEqual(this.model.notAfter, '', 'Both ttl and notAfter are cleared.');
  });
});
