/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn } from '@ember/test-helpers';
import sinon from 'sinon';
import hbs from 'htmlbars-inline-precompile';
import selectors from 'vault/tests/helpers/components/ttl-picker';
import { setRunOptions } from 'ember-a11y-testing/test-support';

module('Integration | Component | ttl-picker', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.set('onChange', sinon.spy());
    this.set('label', 'Foobar');
  });

  module('without toggle', function (hooks) {
    hooks.beforeEach(function () {
      this.set('hideToggle', true);
    });

    test('it shows correct time and value when no initialValue set', async function (assert) {
      await render(hbs`<TtlPicker
        @label={{this.label}}
        @hideToggle={{this.hideToggle}}
        @onChange={{this.onChange}} />`);
      assert.dom(selectors.ttlFormGroup).exists('TTL Form fields exist');
      assert.dom(selectors.ttlValue).hasValue('');
      assert.dom(selectors.ttlUnit).hasValue('s');
    });

    test('it calls the change fn with the correct values', async function (assert) {
      const changeSpy = sinon.spy();
      this.set('onChange', changeSpy);
      await render(hbs`
      <TtlPicker
        @label={{this.label}}
        @hideToggle={{this.hideToggle}}
        @onChange={{this.onChange}}
        @initialValue="30m" />
      `);
      assert.dom(selectors.ttlUnit).hasValue('m', 'unit value shows m (minutes)');
      await fillIn(selectors.ttlValue, '10');
      await assert.ok(changeSpy.calledOnce, 'it calls the passed onChange');
      assert.ok(
        changeSpy.calledWithExactly({
          enabled: true,
          seconds: 600,
          timeString: '10m',
          goSafeTimeString: '10m',
        }),
        'Passes the values back to onChange'
      );
    });

    test('it correctly shows initial time and unit', async function (assert) {
      await render(hbs`
        <TtlPicker
          @label={{this.label}}
          @hideToggle={{this.hideToggle}}
          @initialValue="3h"
          @onChange={{this.onChange}}
        />
      `);

      assert.dom(selectors.ttlUnit).hasValue('h', 'unit value initially shows as h (hours)');
      assert.dom(selectors.ttlValue).hasValue('3', 'time value initially shows as 3');
    });

    test('it fails gracefully when initialValue is not parseable', async function (assert) {
      const changeSpy = sinon.spy();
      this.set('onChange', changeSpy);
      await render(hbs`
        <TtlPicker
          @label={{this.label}}
          @hideToggle={{this.hideToggle}}
          @initialValue="foobar"
          @onChange={{this.onChange}}
          @changeOnInit={{true}}
        />
      `);

      assert.dom(selectors.ttlValue).hasValue('', 'time value initially shows as empty');
      assert.dom(selectors.ttlUnit).hasValue('s', 'unit value initially shows as s (seconds)');
      assert.ok(changeSpy.notCalled, 'onChange is not called on init');
    });

    test('it recalculates time when unit is changed', async function (assert) {
      const changeSpy = sinon.spy();
      this.set('onChange', changeSpy);
      await render(hbs`
        <TtlPicker
          @label={{this.label}}
          @hideToggle={{this.hideToggle}}
          @initialValue="1h"
          @onChange={{this.onChange}}
        />
      `);

      assert.dom(selectors.ttlUnit).hasValue('h', 'unit value initially shows as h (hours)');
      assert.dom(selectors.ttlValue).hasValue('1', 'time value initially shows as 1');
      await fillIn(selectors.ttlUnit, 'm');
      assert.dom(selectors.ttlUnit).hasValue('m', 'unit value changed to m (minutes)');
      assert.dom(selectors.ttlValue).hasValue('60', 'time value recalculates to fit unit');
      assert.ok(
        changeSpy.calledWithExactly({
          enabled: true,
          seconds: 3600,
          timeString: '60m',
          goSafeTimeString: '60m',
        }),
        'Passes the values back to onChange'
      );
    });

    test('it skips recalculating time when unit is changed if time is not whole number', async function (assert) {
      const changeSpy = sinon.spy();
      this.set('onChange', changeSpy);
      await render(hbs`
        <TtlPicker
          @label={{this.label}}
          @hideToggle={{this.hideToggle}}
          @initialValue="30s"
          @onChange={{this.onChange}}
        />
      `);

      assert.dom(selectors.ttlUnit).hasValue('s', 'unit value starts as s (seconds)');
      assert.dom(selectors.ttlValue).hasValue('30', 'time value starts as 30');
      await fillIn(selectors.ttlUnit, 'm');
      assert.dom(selectors.ttlUnit).hasValue('m', 'unit value changed to m (minutes)');
      assert.dom(selectors.ttlValue).hasValue('30', 'time value is still 30');
      assert.ok(
        changeSpy.calledWithExactly({
          enabled: true,
          seconds: 1800,
          timeString: '30m',
          goSafeTimeString: '30m',
        }),
        'Passes the values back to onChange'
      );
    });

    test('it calls onChange on init when changeOnInit is true', async function (assert) {
      const changeSpy = sinon.spy();
      this.set('onChange', changeSpy);
      await render(hbs`
        <TtlPicker
          @label={{this.label}}
          @hideToggle={{this.hideToggle}}
          @initialValue="10m"
          @changeOnInit={{true}}
          @onChange={{this.onChange}}
        />
      `);

      assert.ok(changeSpy.calledOnce, 'it calls the passed onChange when rendered');
      assert.ok(
        changeSpy.calledWithExactly({
          enabled: true,
          seconds: 600,
          timeString: '10m',
          goSafeTimeString: '10m',
        }),
        'Passes the values back to onChange'
      );
    });

    test('it shows a label when passed', async function (assert) {
      this.set('label', 'My Label');
      await render(hbs`
        <TtlPicker
          @label={{this.label}}
          @hideToggle={{this.hideToggle}}
          @onChange={{this.onChange}}
        />
      `);

      assert.dom('[data-test-ttl-form-label]').hasText('My Label', 'Renders label correctly');
      assert.dom('[data-test-ttl-form-subtext]').doesNotExist('Subtext not rendered');
      assert.dom('[data-test-tooltip-trigger]').doesNotExist('Description tooltip not rendered');
    });

    test('it shows subtext and description when passed', async function (assert) {
      setRunOptions({
        rules: {
          // TODO: remove Tooltip
          'aria-command-name': { enabled: false },
        },
      });
      this.set('label', 'My Label');
      await render(hbs`
        <TtlPicker
          @label={{this.label}}
          @hideToggle={{this.hideToggle}}
          @helperTextEnabled="Subtext"
          @description="Description"
          @onChange={{this.onChange}}
        />
      `);

      assert.dom('[data-test-ttl-form-label]').hasText('My Label', 'Renders label correctly');
      assert.dom('[data-test-ttl-form-subtext]').hasText('Subtext', 'Renders subtext when present');
      assert
        .dom('[data-test-tooltip-trigger]')
        .exists({ count: 1 }, 'Description tooltip icon shows when description present');
    });

    test('it yields in place of label if block is present', async function (assert) {
      this.set('label', 'My Label');
      await render(hbs`
        <TtlPicker
          @label={{this.label}}
          @hideToggle={{this.hideToggle}}
          @helperTextEnabled="Subtext"
          @description="Description"
          @onChange={{this.onChange}}
        >
          <legend data-test-custom>Different Label</legend>
        </TtlPicker>
      `);

      assert.dom('[data-test-custom]').hasText('Different Label', 'custom block is rendered');
      assert.dom('[data-test-ttl-form-label]').doesNotExist('Label not rendered');
    });
  });

  module('with toggle', function () {
    test('it has toggle off by default', async function (assert) {
      await render(hbs`
        <TtlPicker
          @label={{this.label}}
          @onChange={{this.onChange}}
        />
      `);
      assert.dom(selectors.toggle).isNotChecked('Toggle is unchecked by default');
      assert.dom(selectors.ttlFormGroup).doesNotExist('TTL Form is not rendered');
    });

    test('it shows time and unit inputs when initialEnabled', async function (assert) {
      const changeSpy = sinon.spy();
      this.set('onChange', changeSpy);
      await render(hbs`
        <TtlPicker
          @label={{this.label}}
          @onChange={{this.onChange}}
          @initialEnabled={{true}}
          @changeOnInit={{true}}
        />
      `);
      assert.dom(selectors.toggle).isChecked('Toggle is checked when initialEnabled is true');
      assert.dom(selectors.ttlFormGroup).exists('TTL Form is rendered');
      assert.ok(changeSpy.notCalled, 'onChange not called because initialValue not parsed');
    });

    test('it sets initial value to initialValue', async function (assert) {
      await render(hbs`
        <TtlPicker
          @label={{this.label}}
          @onChange={{this.onChange}}
          @initialValue="2h"
          @initialEnabled={{true}}
        />
      `);
      assert.dom(selectors.ttlValue).hasValue('2', 'time value is 2');
      assert.dom(selectors.ttlUnit).hasValue('h', 'unit is hours');
      assert.ok(
        this.onChange.notCalled,
        'it does not call onChange after render when changeOnInit is not set'
      );
    });

    test('it passes the appropriate data to onChange when toggled on', async function (assert) {
      const changeSpy = sinon.spy();
      this.set('onChange', changeSpy);
      await render(hbs`
        <TtlPicker
          @label={{this.label}} @initialValue="10m"
          @onChange={{this.onChange}}
        />
      `);
      await click(selectors.toggle);
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

    test('inputs reflect initial value when toggled on', async function (assert) {
      await render(hbs`
        <TtlPicker
          @label={{this.label}} @onChange={{this.onChange}}
          @initialValue="100m"
        />
      `);
      assert.dom(selectors.toggle).isNotChecked('Toggle is off');
      assert.dom(selectors.ttlFormGroup).doesNotExist('TTL Form not shown on mount');
      await click(selectors.toggle);
      assert.dom(selectors.ttlValue).hasValue('100', 'time after toggle is 100');
      assert.dom(selectors.ttlUnit).hasValue('m', 'Unit is minutes after toggle');
    });

    test('it is enabled on init if initialEnabled is true', async function (assert) {
      await render(hbs`
        <TtlPicker
          @label={{this.label}} @onChange={{this.onChange}}
          @initialValue="100m"
          @initialEnabled={{true}}
        />
      `);
      assert.dom(selectors.toggle).isChecked('Toggle is on');
      assert.dom(selectors.ttlFormGroup).exists();
      assert.dom(selectors.ttlValue).hasValue('100', 'time is shown on mount');
      assert.dom(selectors.ttlUnit).hasValue('m', 'Unit is shown on mount');
      await click(selectors.toggle);
      assert.dom(selectors.toggle).isNotChecked('Toggle is off');
      assert.dom(selectors.ttlFormGroup).doesNotExist('TTL Form no longer shows after toggle');
    });

    test('it is enabled on init if initialEnabled evals to truthy', async function (assert) {
      await render(hbs`
        <TtlPicker
          @label={{this.label}} @onChange={{this.onChange}}
          @initialValue="100m"
          @initialEnabled="100m"
        />
      `);
      assert.dom(selectors.toggle).isChecked('Toggle is enabled');
      assert.dom(selectors.ttlValue).hasValue('100', 'time value is shown on mount');
      assert.dom(selectors.ttlUnit).hasValue('m', 'Unit matches what is passed in');
    });

    test('it converts days to go safe time', async function (assert) {
      await render(hbs`
        <TtlPicker
          @label={{this.label}} @initialValue="2d"
          @onChange={{this.onChange}}
        />
      `);
      await click(selectors.toggle);
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

    test('it converts to the largest round unit on init', async function (assert) {
      await render(hbs`
        <TtlPicker
          @label={{this.label}} @onChange={{this.onChange}}
          @initialValue="60000s"
          @initialEnabled="true"
        />
      `);
      assert.dom(selectors.ttlValue).hasValue('1000', 'time value is converted');
      assert.dom(selectors.ttlUnit).hasValue('m', 'unit value is m (minutes)');
    });

    test('it converts to the largest round unit on init when no unit provided', async function (assert) {
      await render(hbs`
        <TtlPicker
          @label={{this.label}} @onChange={{this.onChange}}
          @initialValue={{86400}}
          @initialEnabled="true"
        />
      `);
      assert.dom(selectors.ttlValue).hasValue('1', 'time value is converted');
      assert.dom(selectors.ttlUnit).hasValue('d', 'unit value is d (days)');
    });
  });
});
