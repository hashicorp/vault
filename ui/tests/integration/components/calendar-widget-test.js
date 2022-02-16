import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn } from '@ember/test-helpers';
import sinon from 'sinon';
import hbs from 'htmlbars-inline-precompile';
import calendarDropdown from 'vault/tests/pages/components/calendar-widget';

module('Integration | Component | calendar-widget', function (hooks) {
  setupRenderingTest(hooks);

  const ARRAY_OF_MONTHS = [
    'January',
    'February',
    'March',
    'April',
    'May',
    'June',
    'July',
    'August',
    'September',
    'October',
    'November',
    'December',
  ];

  hooks.beforeEach(function () {
    this.set('handleClientActivityQuery', sinon.spy());
    this.set('handleCurrentBillingPeriod', sinon.spy());
    this.set('arrayOfMonths', ARRAY_OF_MONTHS);
    this.set('endTimeFromResponse', ['2022', 0]);
  });

  test('it renders and can open the calendar view', async function (assert) {
    await render(hbs`
      <CalendarWidget
        @arrayOfMonths={{arrayOfMonths}}
        @endTimeDisplay={{"January 2022"}}
        @endTimeFromResponse={{endTimeFromResponse}}
        @handleClientActivityQuery={{handleClientActivityQuery}}
        @handleCurrentBillingPeriod={{handleCurrentBillingPeriod}}
        @startTimeDisplay={{"February 2021"}}
      />
    `);

    await calendarDropdown.openCalendar();
    assert.ok(calendarDropdown.showsCalendar, 'renders the calendar component');
  });

  test('it does not allow you to click to a future year but does allow you to click to previous years', async function (assert) {
    await render(hbs`
      <CalendarWidget
        @arrayOfMonths={{arrayOfMonths}}
        @endTimeDisplay={{"January 2022"}}
        @endTimeFromResponse={{endTimeFromResponse}}
        @handleClientActivityQuery={{handleClientActivityQuery}}
        @handleCurrentBillingPeriod={{handleCurrentBillingPeriod}}
        @startTimeDisplay={{"February 2021"}}
      />
    `);

    await calendarDropdown.openCalendar();
    assert.dom('[data-test-ttl-value]').doesNotExist('TTL Picker time input exists');
    assert.dom('[data-test-ttl-unit]').doesNotExist('TTL Picker unit select exists');
  });

  test('it shows a tooltip if the click to a previous year is disabled', async function (assert) {
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

  test('it disables the current month and future months', async function (assert) {
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

  test('it allows you to reset the billing period', async function (assert) {
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
});
