import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, findAll } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | modal', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function(assert) {
    // Set any properties with this.set('myProperty', 'value');
    // Handle any actions with this.set('myAction', function(val) { ... });

    await render(hbs`<Modal></Modal><div id="modal-wormhole"></div>`);

    assert.equal(this.element.textContent.trim(), '', 'renders without interior content');
    assert.equal(findAll('[data-test-modal-close-button]').length, 0, 'does not render close modal button');

    // Template block usage:
    await render(hbs`
      <Modal @showCloseButton={{true}}>
        template block text
      </Modal>
      <div id="modal-wormhole"></div>
    `);

    assert.equal(this.element.textContent.trim(), 'template block text', 'renders with interior content');
    assert.equal(findAll('[data-test-modal-close-button]').length, 1, 'renders close modal button');
    assert.dom('[data-test-modal-glyph]').doesNotExist('Glyph is not rendered by default');
  });

  test('it adds the correct type class', async function(assert) {
    await render(hbs`
      <Modal @type="warning">
        template block text
      </Modal>
      <div id="modal-wormhole"></div>
    `);

    assert.dom('.modal.is-highlight').exists('Modal exists with is-highlight class');
    assert.dom('[data-test-modal-glyph]').exists('Glyph is rendered');
  });
});
