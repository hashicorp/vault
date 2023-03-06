import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { click } from '@ember/test-helpers';

module('Integration | Component | alert-popup', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.set('message', 'some very important alert');
    this.set('type', 'warning');
    this.set('close', () => this.set('closed', true));
  });

  test('it renders the alert popup input', async function (assert) {
    await render(hbs`
      <AlertPopup @message={{this.message}} @type={{message-types this.type}} @close={{this.close}} />
    `);

    assert.dom(this.element).hasText('Warning some very important alert');
  });

  test('it invokes the close action', async function (assert) {
    assert.expect(1);

    await render(hbs`
      <AlertPopup @message={{this.message}} @type={{message-types this.type}} @close={{this.close}} />
    `);
    await click('.close-button');

    assert.true(this.closed);
  });

  test('it renders the alert popup with different colors based on types', async function (assert) {
    await render(hbs`
      <AlertPopup @message={{this.message}} @type={{message-types this.type}} @close={{this.close}} />
    `);

    assert.dom('.message').hasClass('is-highlight');

    this.set('type', 'info');

    await render(hbs`
    <AlertPopup @message={{this.message}} @type={{message-types this.type}} @close={{this.close}} />
    `);

    assert.dom('.message').hasClass('is-info');

    this.set('type', 'danger');

    await render(hbs`
    <AlertPopup @message={{this.message}} @type={{message-types this.type}} @close={{this.close}} />
    `);

    assert.dom('.message').hasClass('is-danger');
  });
});
