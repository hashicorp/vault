import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn } from '@ember/test-helpers';
import sinon from 'sinon';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | ttl-picker2', function(hooks) {
  setupRenderingTest(hooks);

  test('it renders time and unit inputs when TTL enabled', async function(assert) {
    let changeSpy = sinon.spy();
    this.set('onChange', changeSpy);
    await render(hbs`
      <TtlPicker2
        @onChange={{onChange}}
        @enableTTL={{true}}
      />
    `);
    assert.dom('[data-test-ttl-value]').exists('TTL Picker time input exists');
    assert.dom('[data-test-ttl-unit]').exists('TTL Picker unit select exists');
  });

  test('it does not show time and unit inputs when TTL disabled', async function(assert) {
    let changeSpy = sinon.spy();
    this.set('onChange', changeSpy);
    await render(hbs`
      <TtlPicker2
        @onChange={{onChange}}
        @enableTTL={{false}}
      />
    `);
    assert.dom('[data-test-ttl-value]').doesNotExist('TTL Picker time input exists');
    assert.dom('[data-test-ttl-unit]').doesNotExist('TTL Picker unit select exists');
  });

  test('it passes the appropriate data to onChange when toggled on', async function(assert) {
    let changeSpy = sinon.spy();
    this.set('onChange', changeSpy);
    await render(hbs`
      <TtlPicker2
        @label="clicktest"
        @unit="m"
        @time="10"
        @onChange={{onChange}}
        @enableTTL={{false}}
      />
    `);
    await click('[data-test-toggle-input="clicktest"]');
    assert.ok(changeSpy.calledOnce, 'it calls the passed onChange');
    assert.ok(
      changeSpy.calledWith({
        enabled: true,
        seconds: 600,
        timeString: '10m',
      }),
      'Passes the default values back to onChange'
    );
  });

  test('it keeps seconds value when unit is changed', async function(assert) {
    let changeSpy = sinon.spy();
    this.set('onChange', changeSpy);
    await render(hbs`
      <TtlPicker2
        @label="clicktest"
        @unit="s"
        @time="360"
        @onChange={{onChange}}
        @enableTTL={{false}}
      />
    `);
    await click('[data-test-toggle-input="clicktest"]');
    assert.ok(changeSpy.calledOnce, 'it calls the passed onChange');
    assert.ok(
      changeSpy.calledWith({
        enabled: true,
        seconds: 360,
        timeString: '360s',
      }),
      'Changes enabled to true on click'
    );
    await fillIn('[data-test-select="ttl-unit"]', 'm');
    assert.ok(
      changeSpy.calledWith({
        enabled: true,
        seconds: 360,
        timeString: '6m',
      }),
      'Units and time update without changing seconds value'
    );
    assert.dom('[data-test-ttl-value]').hasValue('6', 'time value shows as 6');
    assert.dom('[data-test-select="ttl-unit"]').hasValue('m', 'unit value shows as m (minutes)');
  });

  test('it recalculates seconds when unit is changed and recalculateSeconds is on', async function(assert) {
    let changeSpy = sinon.spy();
    this.set('onChange', changeSpy);
    await render(hbs`
      <TtlPicker2
        @label="clicktest"
        @unit="s"
        @time="120"
        @onChange={{onChange}}
        @enableTTL={{true}}
        @recalculateSeconds={{true}}
      />
    `);
    await fillIn('[data-test-select="ttl-unit"]', 'm');
    assert.ok(
      changeSpy.calledWith({
        enabled: true,
        seconds: 7200,
        timeString: '120m',
      }),
      'Seconds value is recalculated based on time and unit'
    );
  });

  test('it sets default value to seconds of parsed value when set', async function(assert) {
    let changeSpy = sinon.spy();
    this.set('onChange', changeSpy);
    await render(hbs`
      <TtlPicker2
        @onChange={{onChange}}
        @initialValue="2h"
        @enableTTL={{true}}
        @time=4
        @unit="d"
      />
    `);
    assert.dom('[data-test-ttl-value]').hasValue('7200', 'time value is initialValue as seconds');
    assert.dom('[data-test-select="ttl-unit"]').hasValue('s', 'unit is seconds');
  });

  test('it is enabled on init if initialEnabled is true', async function(assert) {
    let changeSpy = sinon.spy();
    this.set('onChange', changeSpy);
    await render(hbs`
      <TtlPicker2
        @label="inittest"
        @onChange={{onChange}}
        @initialValue="100m"
        @initialEnabled={{true}}
      />
    `);
    assert.dom('[data-test-ttl-value]').hasValue('6000', 'time value is initialValue as seconds');
    assert.dom('[data-test-ttl-unit]').exists('Unit is shown on mount');
    await click('[data-test-toggle-input="inittest"]');
    assert.dom('[data-test-ttl-value]').doesNotExist('Value no longer shows after toggle');
    assert.dom('[data-test-ttl-unit]').doesNotExist('Unit no longer shows after toggle');
  });

  test('it is enabled on init if initialEnabled evals to truthy', async function(assert) {
    let changeSpy = sinon.spy();
    this.set('onChange', changeSpy);
    await render(hbs`
      <TtlPicker2
        @label="inittest"
        @onChange={{onChange}}
        @initialValue="100m"
        @initialEnabled="true"
      />
    `);
    assert.dom('[data-test-ttl-value]').hasValue('6000', 'time value is initialValue as seconds');
    assert.dom('[data-test-ttl-unit]').exists('Unit is shown on mount');
  });
});
