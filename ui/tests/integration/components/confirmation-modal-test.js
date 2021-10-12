import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import sinon from 'sinon';
import { fillIn, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | confirmation-modal', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders with disabled confirmation button until input matches', async function(assert) {
    let spy = sinon.spy();
    this.set('onConfirm', spy);

    await render(hbs`
      <div id="modal-wormhole"></div>
      <ConfirmationModal
        @isActive={true}
        @onConfirm={this.onConfirm}
        @buttonText="Plz Continue"
        @confirmText="Destructive Thing"
        @testSelector="demote"
      />
    `);

    assert.dom('[data-test-confirm-button]').isDisabled();
    assert.dom('[data-test-confirm-button]').hasText('Plz Continue', 'Confirm button has specified value');

    await fillIn('[data-test-confirmation-modal-input="demote"]', 'Destructive Thing');
    assert.dom('[data-test-confirm-button="demote"]').isNotDisabled();
  });
});
