import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import sinon from 'sinon';
import hbs from 'htmlbars-inline-precompile';

const SELECTORS = {
  toggle: '[data-test-ttl-toggle]',
  ttlFormGroup: '[data-test-ttl-picker-group]',
  ttlValue: '[data-test-ttl-value]',
  ttlUnit: '[data-test-select="ttl-unit"]',
};

module('Integration | Component | ttl-picker', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.set('onChange', sinon.spy());
  });

  test('it renders time and unit inputs when TTL enabled', async function (assert) {
    await render(hbs`
      <TtlPicker
        @onChange={{this.onChange}}
        @initialEnabled={{true}}
      />
    `);
    assert.dom(SELECTORS.toggle).isChecked('Toggle is checked when initialEnabled is true');
    assert.dom(SELECTORS.ttlFormGroup).exists('TTL Form is rendered');
  });

  test('it does not show time and unit inputs when TTL disabled', async function (assert) {
    await render(hbs`
      <TtlPicker
        @onChange={{this.onChange}}
      />
    `);
    assert.dom(SELECTORS.toggle).isNotChecked('Toggle is unchecked by default');
    assert.dom(SELECTORS.ttlFormGroup).doesNotExist('TTL Form is not rendered');
  });

  test('it passes the appropriate data to onChange when toggled on', async function (assert) {
    const changeSpy = sinon.spy();
    this.set('onChange', changeSpy);
    await render(hbs`
      <TtlPicker
        @label="clicktest"
        @initialValue="10m"
        @onChange={{this.onChange}}
      />
    `);
    await click(SELECTORS.toggle);
    assert.ok(changeSpy.calledOnce, 'it calls the passed onChange');
    assert.ok(
      changeSpy.calledWith({
        enabled: true,
        seconds: 600,
        timeString: '10m',
        goSafeTimeString: '10m',
      }),
      'Passes the values back to onChange'
    );
  });

  test('it sets default value to time and unit passed', async function (assert) {
    await render(hbs`
      <TtlPicker
        @onChange={{this.onChange}}
        @initialValue="2h"
        @initialEnabled={{true}}
      />
    `);
    assert.dom(SELECTORS.ttlValue).hasValue('2', 'time value is 2');
    assert.dom(SELECTORS.ttlUnit).hasValue('h', 'unit is hours');
    assert.ok(this.onChange.notCalled, 'it does not call onChange after render when changeOnInit is not set');
  });

  test('toggle is off on init by default', async function (assert) {
    await render(hbs`
      <TtlPicker
        @label="inittest"
        @onChange={{this.onChange}}
        @initialValue="100m"
      />
    `);
    assert.dom(SELECTORS.toggle).isNotChecked('Toggle is off');
    assert.dom(SELECTORS.ttlFormGroup).doesNotExist('TTL Form not shown on mount');
    await click(SELECTORS.toggle);
    assert.dom(SELECTORS.ttlValue).hasValue('100', 'time after toggle is 100');
    assert.dom(SELECTORS.ttlUnit).hasValue('m', 'Unit is minutes after toggle');
  });

  test('it is enabled on init if initialEnabled is true', async function (assert) {
    await render(hbs`
      <TtlPicker
        @label="inittest"
        @onChange={{this.onChange}}
        @initialValue="100m"
        @initialEnabled={{true}}
      />
    `);
    assert.dom(SELECTORS.toggle).isChecked('Toggle is on');
    assert.dom(SELECTORS.ttlFormGroup).exists();
    assert.dom(SELECTORS.ttlValue).hasValue('100', 'time is shown on mount');
    assert.dom(SELECTORS.ttlUnit).hasValue('m', 'Unit is shown on mount');
    await click(SELECTORS.toggle);
    assert.dom(SELECTORS.toggle).isNotChecked('Toggle is off');
    assert.dom(SELECTORS.ttlFormGroup).doesNotExist('TTL Form no longer shows after toggle');
  });

  test('it is enabled on init if initialEnabled evals to truthy', async function (assert) {
    await render(hbs`
      <TtlPicker
        @label="inittest"
        @onChange={{this.onChange}}
        @initialValue="100m"
        @initialEnabled="100m"
      />
    `);
    assert.dom(SELECTORS.toggle).isChecked('Toggle is enabled');
    assert.dom(SELECTORS.ttlValue).hasValue('100', 'time value is shown on mount');
    assert.dom(SELECTORS.ttlUnit).hasValue('m', 'Unit matches what is passed in');
  });

  test('it calls onChange', async function (assert) {
    await render(hbs`
      <TtlPicker
        @label="clicktest"
        @initialValue="2d"
        @onChange={{this.onChange}}
      />
    `);
    await click(SELECTORS.toggle);
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
      <TtlPicker
        @label="changeOnInitTest"
        @onChange={{this.onChange}}
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
      <TtlPicker
        @label="convertunits"
        @onChange={{this.onChange}}
        @initialValue="60000s"
        @initialEnabled="true"
      />
    `);
    assert.dom(SELECTORS.ttlValue).hasValue('1000', 'time value is converted');
    assert.dom(SELECTORS.ttlUnit).hasValue('m', 'unit value is m (minutes)');
  });

  test('it converts to the largest round unit on init when no unit provided', async function (assert) {
    await render(hbs`
      <TtlPicker
        @label="convertunits"
        @onChange={{this.onChange}}
        @initialValue={{86400}}
        @initialEnabled="true"
      />
    `);
    assert.dom(SELECTORS.ttlValue).hasValue('1', 'time value is converted');
    assert.dom(SELECTORS.ttlUnit).hasValue('d', 'unit value is d (days)');
  });
});
