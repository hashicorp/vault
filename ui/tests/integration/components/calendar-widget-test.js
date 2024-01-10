/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import sinon from 'sinon';
import hbs from 'htmlbars-inline-precompile';
import calendarDropdown from 'vault/tests/pages/components/calendar-widget';
import { ARRAY_OF_MONTHS } from 'core/utils/date-formatters';
import { subMonths, subYears } from 'date-fns';
import timestamp from 'core/utils/timestamp';

module('Integration | Component | calendar-widget', function (hooks) {
  setupRenderingTest(hooks);

  hooks.before(function () {
    sinon.stub(timestamp, 'now').callsFake(() => new Date('2018-04-03T14:15:30'));
  });
  hooks.beforeEach(function () {
    const CURRENT_DATE = timestamp.now();
    this.set('currentDate', CURRENT_DATE);
    this.set('calendarStartDate', subMonths(CURRENT_DATE, 12));
    this.set('calendarEndDate', CURRENT_DATE);
    this.set('startTimestamp', subMonths(CURRENT_DATE, 12).toISOString());
    this.set('endTimestamp', CURRENT_DATE.toISOString());
    this.set('handleClientActivityQuery', sinon.spy());
  });
  hooks.after(function () {
    timestamp.now.restore();
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
      .dom('[data-test-calendar-widget-trigger]')
      .hasText(`Apr 2017 - Apr 2018`, 'renders and formats start and end dates');
    await calendarDropdown.openCalendar();
    assert.ok(calendarDropdown.showsCalendar, 'renders the calendar component');
    // assert months in current year are disabled/enabled correctly
    const enabledMonths = ['January', 'February', 'March', 'April'];
    ARRAY_OF_MONTHS.forEach(function (month) {
      if (enabledMonths.includes(month)) {
        assert.dom(`[data-test-calendar-month="${month}"]`).isNotDisabled(`${month} is enabled`);
      } else {
        assert.dom(`[data-test-calendar-month="${month}"]`).isDisabled(`${month} is disabled`);
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
    const disabledMonths = ['January', 'February', 'March'];
    ARRAY_OF_MONTHS.forEach(function (month) {
      if (disabledMonths.includes(month)) {
        assert.dom(`[data-test-calendar-month="${month}"]`).isDisabled(`${month} is disabled`);
      } else {
        assert.dom(`[data-test-calendar-month="${month}"]`).isNotDisabled(`${month} is enabled`);
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
    await click(`[data-test-calendar-month="April"`);
    assert.propEqual(
      this.handleClientActivityQuery.lastCall.lastArg,
      {
        dateType: 'endDate',
        monthIdx: 3,
        monthName: 'April',
        year: 2018,
      },
      'it calls parent function with end date (current) month/year'
    );

    await calendarDropdown.openCalendar();
    await calendarDropdown.clickPreviousYear();
    await click(`[data-test-calendar-month="May"]`);
    assert.propEqual(
      this.handleClientActivityQuery.lastCall.lastArg,
      {
        dateType: 'endDate',
        monthIdx: 4,
        monthName: 'May',
        year: 2017,
      },
      'it calls parent function with selected start date month/year'
    );
  });

  test('it disables correct months when start date 6 months ago', async function (assert) {
    this.set('calendarStartDate', subMonths(this.currentDate, 6)); // Nov 3, 2017
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
    const prevYearEnabled = ['October', 'November', 'December'];
    ARRAY_OF_MONTHS.forEach(function (month) {
      if (prevYearEnabled.includes(month)) {
        assert.dom(`[data-test-calendar-month="${month}"]`).isNotDisabled(`${month} is enabled`);
      } else {
        assert.dom(`[data-test-calendar-month="${month}"]`).isDisabled(`${month} is read only`);
      }
    });

    // Check end year disables correct months
    await click('[data-test-next-year]');
    const currYearEnabled = ['January', 'February', 'March', 'April'];
    ARRAY_OF_MONTHS.forEach(function (month) {
      if (currYearEnabled.includes(month)) {
        assert.dom(`[data-test-calendar-month="${month}"]`).isNotDisabled(`${month} is enabled`);
      } else {
        assert.dom(`[data-test-calendar-month="${month}"]`).isDisabled(`${month} is disabled`);
      }
    });
  });

  test('it disables correct months when start date 36 months ago', async function (assert) {
    this.set('calendarStartDate', subMonths(this.currentDate, 36)); // April 3 2015
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
    for (const year of [2017, 2016, 2015]) {
      await calendarDropdown.clickPreviousYear();
      assert.dom('[data-test-display-year]').hasText(year.toString());
    }

    assert.dom('[data-test-previous-year]').isDisabled('previous year is disabled');
    assert.dom('[data-test-next-year]').isEnabled('next year is enabled');

    const disabledMonths = ['January', 'February', 'March'];
    ARRAY_OF_MONTHS.forEach(function (month) {
      if (disabledMonths.includes(month)) {
        assert.dom(`[data-test-calendar-month="${month}"]`).isDisabled(`${month} is disabled`);
      } else {
        assert.dom(`[data-test-calendar-month="${month}"]`).isNotDisabled(`${month} is enabled`);
      }
    });

    await click('[data-test-next-year]');
    ARRAY_OF_MONTHS.forEach(function (month) {
      assert.dom(`[data-test-calendar-month="${month}"]`).isNotDisabled(`${month} is enabled for 2016`);
    });
    await click('[data-test-next-year]');
    ARRAY_OF_MONTHS.forEach(function (month) {
      assert.dom(`[data-test-calendar-month="${month}"]`).isNotDisabled(`${month} is enabled for 2017`);
    });
    await click('[data-test-next-year]');

    const enabledMonths = ['January', 'February', 'March', 'April'];
    ARRAY_OF_MONTHS.forEach(function (month) {
      if (enabledMonths.includes(month)) {
        assert.dom(`[data-test-calendar-month="${month}"]`).isNotDisabled(`${month} is enabled`);
      } else {
        assert.dom(`[data-test-calendar-month="${month}"]`).isDisabled(`${month} is disabled`);
      }
    });
  });
});
