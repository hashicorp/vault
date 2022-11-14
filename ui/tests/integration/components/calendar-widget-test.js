import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, findAll, find } from '@ember/test-helpers';
import sinon from 'sinon';
import hbs from 'htmlbars-inline-precompile';
import calendarDropdown from 'vault/tests/pages/components/calendar-widget';
import { ARRAY_OF_MONTHS } from 'core/utils/date-formatters';
import { subMonths, subYears } from 'date-fns';
import format from 'date-fns/format';

module('Integration | Component | calendar-widget', function (hooks) {
  setupRenderingTest(hooks);

  const isDisplayingSameYear = (comparisonDate, calendarYear) => {
    return comparisonDate.getFullYear() === parseInt(calendarYear);
  };

  hooks.beforeEach(function () {
    const CURRENT_DATE = new Date();
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

    assert.dom(calendarDropdown.dateRangeTrigger).hasText(
      `${format(this.calendarStartDate, 'MMM yyyy')} - 
      ${format(this.calendarEndDate, 'MMM yyyy')}`,
      'renders and formats start and end dates'
    );
    await calendarDropdown.openCalendar();
    assert.ok(calendarDropdown.showsCalendar, 'renders the calendar component');

    // assert months in current year are disabled/enabled correctly
    const monthButtons = findAll('[data-test-calendar-month]');
    const enabledMonths = [],
      disabledMonths = [];
    for (let monthIdx = 0; monthIdx < 12; monthIdx++) {
      if (monthIdx > this.calendarEndDate.getMonth()) {
        disabledMonths.push(monthButtons[monthIdx]);
      } else {
        enabledMonths.push(monthButtons[monthIdx]);
      }
    }
    enabledMonths.forEach((btn) => {
      assert
        .dom(btn)
        .doesNotHaveClass(
          'is-readOnly',
          `${ARRAY_OF_MONTHS[btn.id] + this.calendarEndDate.getFullYear()} is enabled`
        );
    });
    disabledMonths.forEach((btn) => {
      assert
        .dom(btn)
        .hasClass(
          'is-readOnly',
          `${ARRAY_OF_MONTHS[btn.id] + this.calendarEndDate.getFullYear()} is read only`
        );
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
    const monthButtons = findAll('[data-test-calendar-month]');
    const enabledMonths = [],
      disabledMonths = [];
    for (let monthIdx = 0; monthIdx < 12; monthIdx++) {
      if (monthIdx < this.calendarStartDate.getMonth()) {
        disabledMonths.push(monthButtons[monthIdx]);
      } else {
        enabledMonths.push(monthButtons[monthIdx]);
      }
    }
    disabledMonths.forEach((btn) => {
      assert
        .dom(btn)
        .hasClass(
          'is-readOnly',
          `${ARRAY_OF_MONTHS[btn.id] + this.calendarEndDate.getFullYear()} is read only`
        );
    });
    enabledMonths.forEach((btn) => {
      assert
        .dom(btn)
        .doesNotHaveClass(
          'is-readOnly',
          `${ARRAY_OF_MONTHS[btn.id] + this.calendarEndDate.getFullYear()} is enabled`
        );
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
    await click(`[data-test-calendar-month="${ARRAY_OF_MONTHS[this.calendarEndDate.getMonth()]}"]`);
    assert.propEqual(
      this.handleClientActivityQuery.lastCall.lastArg,
      {
        dateType: 'endDate',
        monthIdx: this.currentDate.getMonth(),
        monthName: ARRAY_OF_MONTHS[this.currentDate.getMonth()],
        year: this.currentDate.getFullYear(),
      },
      'it calls parent function with end date (current) month/year'
    );

    await calendarDropdown.openCalendar();
    await calendarDropdown.clickPreviousYear();
    await click(`[data-test-calendar-month="${ARRAY_OF_MONTHS[this.calendarStartDate.getMonth()]}"]`);
    assert.propEqual(
      this.handleClientActivityQuery.lastCall.lastArg,
      {
        dateType: 'endDate',
        monthIdx: this.currentDate.getMonth(),
        monthName: ARRAY_OF_MONTHS[this.currentDate.getMonth()],
        year: this.currentDate.getFullYear() - 1,
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

    const displayYear = find('[data-test-display-year]').innerText;
    const isRangeSameYear = isDisplayingSameYear(this.calendarStartDate, displayYear);

    // only click previous year if 6 months ago was last year
    if (!isRangeSameYear) {
      await calendarDropdown.clickPreviousYear();
    }
    assert.dom('[data-test-previous-year]').isDisabled('previous year is disabled');

    // DOM calendar is viewing start date year
    findAll('[data-test-calendar-month]').forEach((m) => {
      // months before start month should always be disabled
      if (m.id < this.calendarStartDate.getMonth()) {
        assert.dom(m).hasClass('is-readOnly', `${ARRAY_OF_MONTHS[m.id] + displayYear} is read only`);
      }
      // if start/end dates are in the same year, DOM is also showing end date
      if (isRangeSameYear) {
        // months after end date should be disabled
        if (m.id > this.calendarEndDate.getMonth()) {
          assert.dom(m).hasClass('is-readOnly', `${ARRAY_OF_MONTHS[m.id] + displayYear} is read only`);
        }
        // months between including start/end month should be enabled
        if (m.id >= this.calendarStartDate.getMonth() && m.id <= this.calendarEndDate.getMonth()) {
          assert.dom(m).doesNotHaveClass('is-readOnly', `${ARRAY_OF_MONTHS[m.id] + displayYear} is enabled`);
        }
      }
    });

    // click back to current year if duration spans multiple years
    if (!isRangeSameYear) {
      await click('[data-test-next-year]');
      findAll('[data-test-calendar-month]').forEach((m) => {
        // DOM is no longer showing start month, all months before current date should be enabled
        if (m.id <= this.currentDate.getMonth()) {
          assert.dom(m).doesNotHaveClass('is-readOnly', `${ARRAY_OF_MONTHS[m.id] + displayYear} is enabled`);
        }
        // future months should be disabled
        if (m.id > this.currentDate.getMonth()) {
          assert.dom(m).hasClass('is-readOnly', `${ARRAY_OF_MONTHS[m.id] + displayYear} is read only`);
        }
      });
    }
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

    let displayYear = find('[data-test-display-year]').innerText;

    while (!isDisplayingSameYear(this.calendarStartDate, displayYear)) {
      await calendarDropdown.clickPreviousYear();
      displayYear = find('[data-test-display-year]').innerText;
    }

    assert.dom('[data-test-previous-year]').isDisabled('previous year is disabled');
    assert.dom('[data-test-next-year]').isEnabled('next year is enabled');

    // DOM calendar is viewing start date year (3 years ago)
    findAll('[data-test-calendar-month]').forEach((m) => {
      // months before start month should always be disabled
      if (m.id < this.calendarStartDate.getMonth()) {
        assert.dom(m).hasClass('is-readOnly', `${ARRAY_OF_MONTHS[m.id] + displayYear} is read only`);
      }
      if (m.id >= this.calendarStartDate.getMonth()) {
        assert.dom(m).doesNotHaveClass('is-readOnly', `${ARRAY_OF_MONTHS[m.id] + displayYear} is enabled`);
      }
    });

    await click('[data-test-next-year]');
    displayYear = await find('[data-test-display-year]').innerText;

    if (!isDisplayingSameYear(this.currentDate, displayYear)) {
      await findAll('[data-test-calendar-month]').forEach((m) => {
        // years between should have all months enabled
        assert.dom(m).doesNotHaveClass('is-readOnly', `${ARRAY_OF_MONTHS[m.id] + displayYear} is enabled`);
      });
    }

    await click('[data-test-next-year]');
    displayYear = await find('[data-test-display-year]').innerText;

    if (!isDisplayingSameYear(this.currentDate, displayYear)) {
      await findAll('[data-test-calendar-month]').forEach((m) => {
        // years between should have all months enabled
        assert.dom(m).doesNotHaveClass('is-readOnly', `${ARRAY_OF_MONTHS[m.id] + displayYear} is enabled`);
      });
    }

    await click('[data-test-next-year]');
    displayYear = await find('[data-test-display-year]').innerText;
    // now DOM is showing current year
    assert.dom('[data-test-next-year]').isDisabled('Future year is disabled');
    if (isDisplayingSameYear(this.currentDate, displayYear)) {
      findAll('[data-test-calendar-month]').forEach((m) => {
        //  all months before current month should be enabled
        if (m.id <= this.currentDate.getMonth()) {
          assert.dom(m).doesNotHaveClass('is-readOnly', `${ARRAY_OF_MONTHS[m.id] + displayYear} is enabled`);
        }
        // future months should be disabled
        if (m.id > this.currentDate.getMonth()) {
          assert.dom(m).hasClass('is-readOnly', `${ARRAY_OF_MONTHS[m.id] + displayYear} is read only`);
        }
      });
    }
  });
});
