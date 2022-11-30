import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn, select } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

module('Integration | Component | ttl-form', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.changeSpy = sinon.spy();
    this.set('onChange', this.changeSpy);
  });

  test('it shows no initial time and initial unit of s when not time or unit passed in', async function (assert) {
    await render(hbs`<TtlForm @onChange={{this.onChange}} />`);
    assert.dom('[data-test-ttlform-value]').hasValue('');
    assert.dom('[data-test-select="ttl-unit"]').hasValue('s');
  });

  test('it calls the change fn with the correct values', async function (assert) {
    await render(hbs`<TtlForm @onChange={{this.onChange}} @initialValue="30m" />`);
    assert.dom('[data-test-select="ttl-unit"]').hasValue('m', 'unit value shows m (minutes)');
    await fillIn('[data-test-ttlform-value]', '10');
    await assert.ok(this.changeSpy.calledOnce, 'it calls the passed onChange');
    assert.ok(
      this.changeSpy.calledWithExactly({
        seconds: 600,
        timeString: '10m',
        goSafeTimeString: '10m',
      }),
      'Passes the values back to onChange'
    );
  });

  test('it correctly shows initial time and unit', async function (assert) {
    const changeSpy = sinon.spy();
    this.set('onChange', changeSpy);
    await render(hbs`
      <TtlForm
        @initialValue="3h"
        @onChange={{this.onChange}}
      />
    `);

    assert.dom('[data-test-select="ttl-unit"]').hasValue('h', 'unit value initially shows as h (hours)');
    assert.dom('[data-test-ttlform-value]').hasValue('3', 'time value initially shows as 3');
  });

  test('it shows blank initial time when initialValue is not parseable', async function (assert) {
    await render(hbs`
      <TtlForm
        @initialValue="foobar"
        @onChange={{this.onChange}}
      />
    `);

    assert.dom('[data-test-ttlform-value]').hasValue('', 'time value initially shows as empty');
    assert.dom('[data-test-select="ttl-unit"]').hasValue('s', 'unit value initially shows as s (seconds)');
  });

  test('it recalculates time when unit is changed', async function (assert) {
    const changeSpy = sinon.spy();
    this.set('onChange', changeSpy);
    await render(hbs`
      <TtlForm
        @initialValue="1h"
        @onChange={{this.onChange}}
      />
    `);

    assert.dom('[data-test-select="ttl-unit"]').hasValue('h', 'unit value initially shows as h (hours)');
    assert.dom('[data-test-ttlform-value]').hasValue('1', 'time value initially shows as 1');
    await select('[data-test-select="ttl-unit"]', 'm');
    assert.dom('[data-test-select="ttl-unit"]').hasValue('m', 'unit value changed to m (minutes)');
    assert.dom('[data-test-ttlform-value]').hasValue('60', 'time value recalculates to fit unit');
    assert.ok(
      changeSpy.calledWithExactly({ seconds: 3600, timeString: '60m', goSafeTimeString: '60m' }),
      'Passes the values back to onChange'
    );
  });

  test('it calls onChange on init when changeOnInit is true', async function (assert) {
    const changeSpy = sinon.spy();
    this.set('onChange', changeSpy);
    await render(hbs`
      <TtlForm
        @initialValue="10m"
        @changeOnInit={{true}}
        @onChange={{this.onChange}}
      />
    `);

    assert.ok(changeSpy.calledOnce, 'it calls the passed onChange when rendered');
    assert.ok(
      changeSpy.calledWithExactly({
        seconds: 600,
        timeString: '10m',
        goSafeTimeString: '10m',
      }),
      'Passes the values back to onChange'
    );
  });

  test('it shows a label when passed', async function (assert) {
    await render(hbs`
      <TtlForm
        @label="My Label"
        @onChange={{this.onChange}}
      />
    `);

    assert.dom('[data-test-ttl-form-label]').hasText('My Label', 'Renders label correctly');
    assert.dom('[data-test-ttl-form-subtext]').doesNotExist('Subtext not rendered');
    assert.dom('[data-test-tooltip-trigger]').doesNotExist('Description tooltip not rendered');
  });

  test('it shows subtext and description when passed', async function (assert) {
    await render(hbs`
      <TtlForm
        @label="My Label"
        @subText="Subtext"
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
    await render(hbs`
      <TtlForm
        @label="My Label"
        @subText="Subtext"
        @description="Description"
        @onChange={{this.onChange}}
      >
        <legend data-test-custom>Different Label</legend>
      </TtlForm>
    `);

    assert.dom('[data-test-custom]').hasText('Different Label', 'custom block is rendered');
    assert.dom('[data-test-ttl-form-label]').doesNotExist('Label not rendered');
  });
});
