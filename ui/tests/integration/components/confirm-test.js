import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

module('Integration | Component | Confirm', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.set('id', 'foo');
    this.set('title', 'Are you sure?');
    this.set('message', 'You will not be able to recover this item later.');
    this.set('triggerText', 'Click me!');
    this.set('onConfirm', sinon.spy());
  });

  test('it renders', async function (assert) {
    await render(hbs`
      <Confirm as |c|>
        <c.Message
          @id={{this.id}}
          @title={{this.title}}
          @triggerText={{this.triggerText}}
          @message={{this.message}}
          @onConfirm={{this.onConfirm}}
        />
      </Confirm>
    `);

    assert.dom('.confirm-wrapper').exists();
    assert.dom('.confirm').containsText(this.triggerText);
  });

  test('does not show the confirmation message until it is triggered', async function (assert) {
    await render(hbs`
      <Confirm as |c|>
        <c.Message
          @id={{this.id}}
          @title={{this.title}}
          @triggerText={{this.triggerText}}
          @message={{this.message}}
          @onConfirm={{this.onConfirm}}
        />
      </Confirm>
    `);
    assert.dom('.confirm-overlay').doesNotContainText(this.message);

    await click('[data-test-confirm-action-trigger]');

    assert.dom('.confirm-overlay').containsText(this.title);
    assert.dom('.confirm-overlay').containsText(this.message);
  });

  test('it calls onConfirm when the confirm button is clicked', async function (assert) {
    await render(hbs`
      <Confirm as |c|>
        <c.Message
          @id={{this.id}}
          @title={{this.title}}
          @triggerText={{this.triggerText}}
          @message={{this.message}}
          @onConfirm={{this.onConfirm}}
        />
      </Confirm>
    `);
    await click('[data-test-confirm-action-trigger]');
    await click('[data-test-confirm-button=true]');

    assert.ok(this.onConfirm.calledOnce);
  });

  test('it shows only the active triggers message', async function (assert) {
    await render(hbs`
      <Confirm as |c|>
        <c.Message
          @id={{this.id}}
          @title={{this.title}}
          @triggerText={{this.triggerText}}
          @message={{this.message}}
          @onConfirm={{this.onConfirm}}
        />
        <c.Message
          @id='bar'
          @title='Wow'
          @message='Bazinga!'
          @onConfirm={{this.onConfirm}}
        />
      </Confirm>
    `);

    await click(`[data-test-confirm-action-trigger=${this.id}]`);
    assert.dom('.confirm-overlay').containsText(this.title);
    assert.dom('.confirm-overlay').containsText(this.message);

    await click('[data-test-confirm-cancel-button]');

    await click("[data-test-confirm-action-trigger='bar']");
    assert.dom('.confirm-overlay').containsText('Wow');
    assert.dom('.confirm-overlay').containsText('Bazinga!');
  });
});
