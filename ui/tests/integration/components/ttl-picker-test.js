import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

module('Integration | Component | ttl picker', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.changeSpy = sinon.spy();
    this.set('onChange', this.changeSpy);
  });

  test('it renders error on non-number input', async function(assert) {
    await render(hbs`<TtlPicker />`);

    await fillIn('[data-test-ttl-value]', 'foo');
    assert.equal(this.changeSpy.callCount, 0, "it did't call onChange");
    assert.dom('[data-test-ttl-error]').includesText('Error', 'renders the error box');

    await fillIn('[data-test-ttl-value]', '33');
    assert.dom('[data-test-ttl-error]').doesNotIncludeText('Error', 'removes the error box');
  });
});
