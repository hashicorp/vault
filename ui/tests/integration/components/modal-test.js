import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import sinon from 'sinon';
import { click, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | modal', function (hooks) {
  setupRenderingTest(hooks);
  const closeAction = sinon.spy();

  hooks.beforeEach(function () {
    this.set('onClose', closeAction);
  });

  test('it renders', async function (assert) {
    await render(
      hbs`<Modal @isActive={{true}} @onClose={{this.onClose}}></Modal><div id="modal-wormhole"></div>`
    );

    assert.dom(this.element).hasText('', 'renders without interior content');
    assert.dom('[data-test-modal-div]').hasAttribute('class', 'modal is-active');
    assert.dom('[data-test-modal-close-button]').doesNotExist('does not render close modal button');

    // Template block usage:
    await render(hbs`
      <Modal @isActive={{true}} @showCloseButton={{true}} @onClose={{this.onClose}} >
        template block text
      </Modal>
      <div id="modal-wormhole"></div>
    `);

    assert.dom(this.element).hasText('template block text', 'renders with interior content');
    assert.dom('[data-test-modal-close-button]').exists({ count: 1 }, 'renders close modal button');
    assert.dom('[data-test-modal-glyph]').doesNotExist('Glyph is not rendered by default');
    await click('[data-test-modal-close-button]');
    assert.true(closeAction.called, 'executes passed in onConfirm function');
  });

  test('it adds the correct type class', async function (assert) {
    await render(hbs`
      <Modal @type="warning" @onClose={{this.onClose}}>
        template block text
      </Modal>
      <div id="modal-wormhole"></div>
    `);

    assert.dom('.modal.is-highlight').exists('Modal exists with is-highlight class');
    assert.dom('[data-test-modal-glyph]').exists('Glyph is rendered');
  });
});
