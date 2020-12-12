import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

module('Integration | Component | ttl-form', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.changeSpy = sinon.spy();
    this.set('onChange', this.changeSpy);
  });

  test('it shows no initial time and initial unit of s when not time or unit passed in', async function(assert) {
    await render(hbs`<TtlForm @onChange={{onChange}} />`);
    assert.dom('[data-test-ttlform-value]').hasValue('');
    assert.dom('[data-test-select="ttl-unit"]').hasValue('s');
  });

  test('it calls the change fn with the correct values', async function(assert) {
    await render(hbs`<TtlForm @onChange={{onChange}} @unit="m" />`);

    assert.dom('[data-test-select="ttl-unit"]').hasValue('m', 'unit value initially shows m (minutes)');
    await fillIn('[data-test-ttlform-value]', '10');
    await assert.ok(this.changeSpy.calledOnce, 'it calls the passed onChange');
    assert.ok(
      this.changeSpy.calledWith({
        seconds: 600,
        timeString: '10m',
      }),
      'Passes the default values back to onChange'
    );
  });

  test('it correctly shows initial unit', async function(assert) {
    let changeSpy = sinon.spy();
    this.set('onChange', changeSpy);
    await render(hbs`
      <TtlForm
        @unit="h"
        @time="3"
        @onChange={{onChange}}
      />
    `);

    assert.dom('[data-test-select="ttl-unit"]').hasValue('h', 'unit value initially shows as h (hours)');
  });
});
