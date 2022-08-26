import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, settled, find, waitUntil } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | alert-inline', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.set('message', 'some very important alert');
  });

  test('it renders alert message with correct class args', async function (assert) {
    await render(hbs`
    <AlertInline
      @paddingTop={{true}}
      @isMarginless={{true}}
      @sizeSmall={{true}}
      @message={{this.message}}
    />
    `);
    assert.dom('[data-test-inline-error-message]').hasText('some very important alert');
    assert
      .dom('[data-test-inline-alert]')
      .hasAttribute('class', 'message-inline padding-top is-marginless size-small');
  });

  test('it yields to block text', async function (assert) {
    await render(hbs`
    <AlertInline @message={{this.message}}> 
      A much more important alert
    </AlertInline>
    `);
    assert.dom('[data-test-inline-error-message]').hasText('A much more important alert');
  });

  test('it renders correctly for type=danger', async function (assert) {
    this.set('type', 'danger');
    await render(hbs`
    <AlertInline 
      @type={{this.type}}
      @message={{this.message}}
    />
    `);
    assert
      .dom('[data-test-inline-error-message]')
      .hasAttribute('class', 'has-text-danger', 'has danger text');
    assert.dom('[data-test-icon="x-square-fill"]').exists('danger icon exists');
  });

  test('it renders correctly for type=warning', async function (assert) {
    this.set('type', 'warning');
    await render(hbs`
    <AlertInline 
      @type={{this.type}}
      @message={{this.message}}
    />
    `);
    assert.dom('[data-test-inline-error-message]').doesNotHaveAttribute('class', 'does not have styled text');
    assert.dom('[data-test-icon="alert-triangle-fill"]').exists('warning icon exists');
  });

  test('it mimics loading when message changes', async function (assert) {
    await render(hbs`
    <AlertInline 
      @message={{this.message}}
      @mimicRefresh={{true}} 
    />
    `);
    assert
      .dom('[data-test-inline-error-message]')
      .hasText('some very important alert', 'it renders original message');

    this.set('message', 'some changed alert!!!');
    await waitUntil(() => find('[data-test-icon="loading"]'));
    assert.ok(find('[data-test-icon="loading"]'), 'it shows loading icon when message changes');
    await settled();
    assert
      .dom('[data-test-inline-error-message]')
      .hasText('some changed alert!!!', 'it shows updated message');
  });
});
