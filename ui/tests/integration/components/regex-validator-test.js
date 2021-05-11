import EmberObject from '@ember/object';
import sinon from 'sinon';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | regex-validator', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders input and validation messages', async function(assert) {
    let attr = EmberObject.create({
      name: 'example',
    });
    let spy = sinon.spy();
    this.set('onChange', spy);
    this.set('attr', attr);
    this.set('value', '(\\d{4})');
    this.set('labelString', 'Regex Example');

    await render(
      hbs`<RegexValidator
        @onChange={{onChange}}
        @attr={{attr}}
        @value={{value}}
        @labelString={{labelString}}
      />`
    );
    assert.dom('.regex-label label').hasText('Regex Example', 'Label is correct');
    assert.dom('[data-test-toggle-input="example-validation-toggle"]').exists('Validation toggle exists');
    assert.dom('[data-test-regex-validator-test-string]').doesNotExist('Test string input does not show');

    await click('[data-test-toggle-input="example-validation-toggle"]');
    assert.dom('[data-test-regex-validator-test-string]').exists('Test string input shows after toggle');
    assert
      .dom('[data-test-regex-validation-message]')
      .doesNotExist('Validation message does not show if test string is empty');

    await fillIn('[data-test-input="example-testval"]', '123a');
    assert.dom('[data-test-regex-validation-message]').exists('Validation message shows after input filled');
    assert
      .dom('[data-test-inline-error-message]')
      .hasText("Your regex doesn't match the subject string", 'Shows error when regex does not match string');

    await fillIn('[data-test-input="example-testval"]', '1234');
    assert
      .dom('[data-test-inline-success-message]')
      .hasText('Your regex matches the subject string', 'Shows success when regex matches');

    await fillIn('[data-test-input="example-testval"]', '12345');
    assert
      .dom('[data-test-inline-error-message]')
      .hasText(
        "Your regex doesn't match the subject string",
        "Shows error if regex doesn't match complete string"
      );
    await fillIn('[data-test-input="example"]', '(\\d{5})');
    assert.ok(spy.calledOnce, 'Calls the passed onChange function when main input is changed');
  });
});
