import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import sinon from 'sinon';
import { fillIn, find, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | confirmation-modal', function(hooks) {
  setupRenderingTest(hooks);

  test('it works as expected', async function(assert) {
    let spy = sinon.spy();
    this.set('onConfirm', spy);

    await render(hbs`
      <div id="modal-wormhole"></div>
      <ConfirmationModal
        @isActive={true}
        @onConfirm={this.onConfirm}
        @buttonText="Plz Continue"
        @confirmText="Destructive Thing"
      />
    `);

    assert.dom('[data-test-confirm-button]').isDisabled();
    assert.equal(
      find('[data-test-confirm-button]').textContent.trim(),
      'Plz Continue',
      'Confirm button has specified value'
    );

    await fillIn('[data-test-confirmation-modal-input="confirmationInput"]', 'Destructive Thing');
    assert.dom('[data-test-confirm-button]').isNotDisabled();
  });
});
