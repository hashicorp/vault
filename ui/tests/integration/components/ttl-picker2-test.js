import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn } from '@ember/test-helpers';
import sinon from 'sinon';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | ttl-picker2', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.set('onChange', sinon.spy());
  });

  test('it renders time and unit inputs when TTL enabled', async function (assert) {
    await render(hbs`
      <TtlPicker2
        @onChange={{onChange}}
        @enableTTL={{true}}
      />
    `);
    assert.dom('[data-test-ttl-value]').exists('TTL Picker time input exists');
    assert.dom('[data-test-ttl-unit]').exists('TTL Picker unit select exists');
  });

  test('it does not show time and unit inputs when TTL disabled', async function (assert) {
    await render(hbs`
      <TtlPicker2
        @onChange={{onChange}}
        @enableTTL={{false}}
      />
    `);
    assert.dom('[data-test-ttl-value]').doesNotExist('TTL Picker time input exists');
    assert.dom('[data-test-ttl-unit]').doesNotExist('TTL Picker unit select exists');
  });

  test('it passes the appropriate data to onChange when toggled on', async function (assert) {
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
    assert.ok(this.onChange.calledOnce, 'it calls the passed onChange');
    assert.ok(
      this.onChange.calledWith({
        enabled: true,
        seconds: 600,
        timeString: '10m',
        goSafeTimeString: '10m',
      }),
      'Passes the default values back to onChange'
    );
  });

  test('it keeps seconds value when unit is changed', async function (assert) {
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
    assert.ok(this.onChange.calledOnce, 'it calls the passed onChange');
    assert.ok(
      this.onChange.calledWith({
        enabled: true,
        seconds: 360,
        timeString: '360s',
        goSafeTimeString: '360s',
      }),
      'Changes enabled to true on click'
    );
    await fillIn('[data-test-select="ttl-unit"]', 'm');
    assert.ok(
      this.onChange.calledWith({
        enabled: true,
        seconds: 360,
        timeString: '6m',
        goSafeTimeString: '6m',
      }),
      'Units and time update without changing seconds value'
    );
    assert.dom('[data-test-ttl-value]').hasValue('6', 'time value shows as 6');
    assert.dom('[data-test-select="ttl-unit"]').hasValue('m', 'unit value shows as m (minutes)');
  });

  test('it recalculates seconds when unit is changed and recalculateSeconds is on', async function (assert) {
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
      this.onChange.calledWith({
        enabled: true,
        seconds: 7200,
        timeString: '120m',
        goSafeTimeString: '120m',
      }),
      'Seconds value is recalculated based on time and unit'
    );
  });

  test('it sets default value to time and unit passed', async function (assert) {
    await render(hbs`
      <TtlPicker2
        @onChange={{onChange}}
        @initialValue="2h"
        @enableTTL={{true}}
        @time=4
        @unit="d"
      />
    `);
    assert.dom('[data-test-ttl-value]').hasValue('2', 'time value is 2');
    assert.dom('[data-test-select="ttl-unit"]').hasValue('h', 'unit is hours');
    assert.ok(this.onChange.notCalled, 'it does not call onChange after render when changeOnInit is not set');
  });

  test('it is disabled on init if initialEnabled is false', async function (assert) {
    await render(hbs`
      <TtlPicker2
        @label="inittest"
        @onChange={{onChange}}
        @initialValue="100m"
        @initialEnabled={{false}}
      />
    `);
    assert.dom('[data-test-ttl-value]').doesNotExist('Value is not shown on mount');
    assert.dom('[data-test-ttl-unit]').doesNotExist('Unit is not shown on mount');
    await click('[data-test-toggle-input="inittest"]');
    assert.dom('[data-test-ttl-value]').hasValue('100', 'time after toggle is 100');
    assert.dom('[data-test-select="ttl-unit"]').hasValue('m', 'Unit is minutes after toggle');
  });

  test('it is enabled on init if initialEnabled is true', async function (assert) {
    await render(hbs`
      <TtlPicker2
        @label="inittest"
        @onChange={{onChange}}
        @initialValue="100m"
        @initialEnabled={{true}}
      />
    `);
    assert.dom('[data-test-ttl-value]').hasValue('100', 'time is shown on mount');
    assert.dom('[data-test-select="ttl-unit"]').hasValue('m', 'Unit is shown on mount');
    await click('[data-test-toggle-input="inittest"]');
    assert.dom('[data-test-ttl-value]').doesNotExist('Value no longer shows after toggle');
    assert.dom('[data-test-ttl-unit]').doesNotExist('Unit no longer shows after toggle');
  });

  test('it is enabled on init if initialEnabled evals to truthy', async function (assert) {
    await render(hbs`
      <TtlPicker2
        @label="inittest"
        @onChange={{onChange}}
        @initialValue="100m"
        @initialEnabled="true"
      />
    `);
    assert.dom('[data-test-ttl-value]').hasValue('100', 'time value is shown on mount');
    assert.dom('[data-test-ttl-unit]').exists('Unit is shown on mount');
    assert.dom('[data-test-select="ttl-unit"]').hasValue('m', 'Unit matches what is passed in');
  });

  test('it calls onChange', async function (assert) {
    await render(hbs`
      <TtlPicker2
        @label="clicktest"
        @unit="d"
        @time="2"
        @onChange={{onChange}}
        @enableTTL={{false}}
      />
    `);
    await click('[data-test-toggle-input="clicktest"]');
    assert.ok(this.onChange.calledOnce, 'it calls the passed onChange');
    assert.ok(
      this.onChange.calledWith({
        enabled: true,
        seconds: 172800,
        timeString: '2d',
        goSafeTimeString: '48h',
      }),
      'Converts day unit to go safe time'
    );
  });

  test('it calls onChange on init when rendered if changeOnInit is true', async function (assert) {
    await render(hbs`
      <TtlPicker2
        @label="changeOnInitTest"
        @onChange={{onChange}}
        @initialValue="100m"
        @initialEnabled="true"
        @changeOnInit={{true}}
      />
    `);
    assert.ok(
      this.onChange.calledWith({
        enabled: true,
        seconds: 6000,
        timeString: '100m',
        goSafeTimeString: '100m',
      }),
      'Seconds value is recalculated based on time and unit'
    );
    assert.ok(this.onChange.calledOnce, 'it calls the passed onChange after render');
  });

  test('it converts to the largest round unit on init', async function (assert) {
    await render(hbs`
      <TtlPicker2
        @label="convertunits"
        @onChange={{onChange}}
        @initialValue="60000s"
        @initialEnabled="true"
      />
    `);
    assert.dom('[data-test-ttl-value]').hasValue('1000', 'time value is converted');
    assert.dom('[data-test-select="ttl-unit"]').hasValue('m', 'unit value is m (minutes)');
  });

  test('it converts to the largest round unit on init when no unit provided', async function (assert) {
    await render(hbs`
      <TtlPicker2
        @label="convertunits"
        @onChange={{onChange}}
        @initialValue={{86400}}
        @initialEnabled="true"
      />
    `);
    assert.dom('[data-test-ttl-value]').hasValue('1', 'time value is converted');
    assert.dom('[data-test-select="ttl-unit"]').hasValue('d', 'unit value is d (days)');
  });
});
