import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import sinon from 'sinon';
import hbs from 'htmlbars-inline-precompile';
import calendarDropdown from 'vault/tests/pages/components/calendar-widget';
import { ARRAY_OF_MONTHS } from 'core/utils/date-formatters';
import { subMonths, subYears } from 'date-fns';

const CURRENT_DATE = new Date('2020-03-15T14:15:00');

module('Integration | Component | calendar-widget', function (hooks) {
  setupRenderingTest(hooks);

  hooks.before(function () {
    sinon.stub(Date, 'now').returns(CURRENT_DATE);
  });

  hooks.beforeEach(function () {
    this.set('currentDate', CURRENT_DATE);
    this.set('calendarStartDate', subMonths(CURRENT_DATE, 12));
    this.set('calendarEndDate', CURRENT_DATE);
    this.set('startTimestamp', subMonths(CURRENT_DATE, 12).toISOString());
    this.set('endTimestamp', CURRENT_DATE.toISOString());
    this.set('handleClientActivityQuery', sinon.spy());
  });

  test('it renders and disables correct months when start date is 12 months ago', async function (assert) {
    assert.expect(14);
    await render(hbs`
      <CalendarWidget
        @startTimestamp={{this.startTimestamp}}
        @endTimestamp={{this.endTimestamp}}
        @selectMonth={{this.handleClientActivityQuery}}
      />
    `);

    assert
      .dom(calendarDropdown.dateRangeTrigger)
      .hasText(`Mar 2019 - Mar 2020`, 'renders and formats start and end dates');
    await calendarDropdown.openCalendar();
    assert.ok(calendarDropdown.showsCalendar, 'renders the calendar component');
    // assert months in current year are disabled/enabled correctly
    const enabledMonths = ['January', 'February', 'March'];
    ARRAY_OF_MONTHS.forEach(function (month) {
      if (enabledMonths.includes(month)) {
        assert
          .dom(`[data-test-calendar-month="${month}"]`)
          .doesNotHaveClass('is-readOnly', `${month} is enabled`);
      } else {
        assert.dom(`[data-test-calendar-month="${month}"]`).hasClass('is-readOnly', `${month} is read only`);
      }
    });
  });

  test('it renders and disables months before start timestamp', async function (assert) {
    await render(hbs`
      <CalendarWidget
        @startTimestamp={{this.startTimestamp}}
        @endTimestamp={{this.endTimestamp}}
        @selectMonth={{this.handleClientActivityQuery}}
      />
    `);

    await calendarDropdown.openCalendar();
    assert.dom('[data-test-next-year]').isDisabled('Future year is disabled');
    await calendarDropdown.clickPreviousYear();
    assert
      .dom('[data-test-display-year]')
      .hasText(`${subYears(this.currentDate, 1).getFullYear()}`, 'shows the previous year');
    assert.dom('[data-test-previous-year]').isDisabled('disables previous year');

    // assert months in previous year are disabled/enabled correctly
    const disabledMonths = ['January', 'February'];
    ARRAY_OF_MONTHS.forEach(function (month) {
      if (disabledMonths.includes(month)) {
        assert.dom(`[data-test-calendar-month="${month}"]`).hasClass('is-readOnly', `${month} is read only`);
      } else {
        assert
          .dom(`[data-test-calendar-month="${month}"]`)
          .doesNotHaveClass('is-readOnly', `${month} is enabled`);
      }
    });
  });

  test('it calls parent callback with correct arg when clicking "Current billing period"', async function (assert) {
    await render(hbs`
      <CalendarWidget
        @startTimestamp={{this.startTimestamp}}
        @endTimestamp={{this.endTimestamp}}
        @selectMonth={{this.handleClientActivityQuery}}
      />
    `);
    await calendarDropdown.menuToggle();
    await calendarDropdown.clickCurrentBillingPeriod();
    assert.propEqual(
      this.handleClientActivityQuery.args[0][0],
      { dateType: 'reset' },
      'it calls parent function with reset dateType'
    );
  });

  test('it calls parent callback with correct arg when clicking "Current month"', async function (assert) {
    await render(hbs`
      <CalendarWidget
        @startTimestamp={{this.startTimestamp}}
        @endTimestamp={{this.endTimestamp}}
        @selectMonth={{this.handleClientActivityQuery}}
      />
    `);
    await calendarDropdown.menuToggle();
    await calendarDropdown.clickCurrentMonth();
    assert.propEqual(
      this.handleClientActivityQuery.args[0][0],
      { dateType: 'currentMonth' },
      'it calls parent function with currentMoth dateType'
    );
  });

  test('it calls parent callback with correct arg when selecting a month', async function (assert) {
    await render(hbs`
      <CalendarWidget
        @startTimestamp={{this.startTimestamp}}
        @endTimestamp={{this.endTimestamp}}
        @selectMonth={{this.handleClientActivityQuery}}
      />
    `);
    await calendarDropdown.openCalendar();
    await click(`[data-test-calendar-month="March"`);
    assert.propEqual(
      this.handleClientActivityQuery.lastCall.lastArg,
      {
        dateType: 'endDate',
        monthIdx: 2,
        monthName: 'March',
        year: 2020,
      },
      'it calls parent function with end date (current) month/year'
    );

    await calendarDropdown.openCalendar();
    await calendarDropdown.clickPreviousYear();
    await click(`[data-test-calendar-month="March"]`);
    assert.propEqual(
      this.handleClientActivityQuery.lastCall.lastArg,
      {
        dateType: 'endDate',
        monthIdx: 2,
        monthName: 'March',
        year: 2019,
      },
      'it calls parent function with start date month/year'
    );
  });

  test('it disables correct months when start date 6 months ago', async function (assert) {
    this.set('calendarStartDate', subMonths(this.currentDate, 6));
    this.set('startTimestamp', subMonths(this.currentDate, 6).toISOString());
    await render(hbs`
      <CalendarWidget
        @startTimestamp={{this.startTimestamp}}
        @endTimestamp={{this.endTimestamp}}
        @selectMonth={{this.handleClientActivityQuery}}
      />
    `);

    await calendarDropdown.openCalendar();
    assert.dom('[data-test-next-year]').isDisabled('Future year is disabled');

    // Check start year disables correct months
    await calendarDropdown.clickPreviousYear();
    assert.dom('[data-test-previous-year]').isDisabled('previous year is disabled');
    const enabled2019 = ['September', 'October', 'November', 'December'];
    ARRAY_OF_MONTHS.forEach(function (month) {
      if (enabled2019.includes(month)) {
        assert
          .dom(`[data-test-calendar-month="${month}"]`)
          .doesNotHaveClass('is-readOnly', `${month} is enabled`);
      } else {
        assert.dom(`[data-test-calendar-month="${month}"]`).hasClass('is-readOnly', `${month} is read only`);
      }
    });

    // Check end year disables correct months
    await click('[data-test-next-year]');
    const enabled2020 = ['January', 'February', 'March'];
    ARRAY_OF_MONTHS.forEach(function (month) {
      if (enabled2020.includes(month)) {
        assert
          .dom(`[data-test-calendar-month="${month}"]`)
          .doesNotHaveClass('is-readOnly', `${month} is enabled`);
      } else {
        assert.dom(`[data-test-calendar-month="${month}"]`).hasClass('is-readOnly', `${month} is read only`);
      }
    });
  });

  test('it disables correct months when start date 36 months ago', async function (assert) {
    this.set('calendarStartDate', subMonths(this.currentDate, 36));
    this.set('startTimestamp', subMonths(this.currentDate, 36).toISOString());
    await render(hbs`
      <CalendarWidget
        @startTimestamp={{this.startTimestamp}}
        @endTimestamp={{this.endTimestamp}}
        @selectMonth={{this.handleClientActivityQuery}}
      />
    `);

    await calendarDropdown.openCalendar();
    assert.dom('[data-test-next-year]').isDisabled('Future year is disabled');

    await calendarDropdown.clickPreviousYear();
    assert.dom('[data-test-display-year]').hasText('2019');
    await calendarDropdown.clickPreviousYear();
    assert.dom('[data-test-display-year]').hasText('2018');
    await calendarDropdown.clickPreviousYear();
    assert.dom('[data-test-display-year]').hasText('2017');

    assert.dom('[data-test-previous-year]').isDisabled('previous year is disabled');
    assert.dom('[data-test-next-year]').isEnabled('next year is enabled');

    assert.dom('.calendar-widget .is-readOnly').exists('Some months disabled');

    const disabled2017 = ['January', 'February'];
    ARRAY_OF_MONTHS.forEach(function (month) {
      if (disabled2017.includes(month)) {
        assert.dom(`[data-test-calendar-month="${month}"]`).hasClass('is-readOnly', `${month} is read only`);
      } else {
        assert
          .dom(`[data-test-calendar-month="${month}"]`)
          .doesNotHaveClass('is-readOnly', `${month} is enabled`);
      }
    });

    await click('[data-test-next-year]');
    assert.dom('.calendar-widget .is-readOnly').doesNotExist('All months enabled for 2018');
    await click('[data-test-next-year]');
    assert.dom('.calendar-widget .is-readOnly').doesNotExist('All months enabled for 2019');
    await click('[data-test-next-year]');
    assert.dom('.calendar-widget .is-readOnly').exists('Some months disabled for 2020');

    const enabled2020 = ['January', 'February', 'March'];
    ARRAY_OF_MONTHS.forEach(function (month) {
      if (enabled2020.includes(month)) {
        assert
          .dom(`[data-test-calendar-month="${month}"]`)
          .doesNotHaveClass('is-readOnly', `${month} is enabled`);
      } else {
        assert.dom(`[data-test-calendar-month="${month}"]`).hasClass('is-readOnly', `${month} is read only`);
      }
    });
  });
});
