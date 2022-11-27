import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import sinon from 'sinon';
import hbs from 'htmlbars-inline-precompile';
import calendarDropdown from 'vault/tests/pages/components/calendar-widget';
import { ARRAY_OF_MONTHS } from 'core/utils/date-formatters';
import { subYears } from 'date-fns';

module('Integration | Component | calendar-widget', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    const CURRENT_YEAR = new Date().getFullYear();
    const PREVIOUS_YEAR = subYears(new Date(), 1).getFullYear();
    this.set('currentYear', CURRENT_YEAR);
    this.set('previousYear', PREVIOUS_YEAR);
    this.set('handleClientActivityQuery', sinon.spy());
    this.set('handleCurrentBillingPeriod', sinon.spy());
    this.set('arrayOfMonths', ARRAY_OF_MONTHS);
    this.set('endTimeFromResponse', [CURRENT_YEAR, 0]);
  });

  test('it renders and can open the calendar view', async function (assert) {
    await render(hbs`
      <CalendarWidget
        @arrayOfMonths={{this.arrayOfMonths}}
        @endTimeDisplay={{concat "January " this.currentYear}}
        @endTimeFromResponse={{this.endTimeFromResponse}}
        @handleClientActivityQuery={{this.handleClientActivityQuery}}
        @handleCurrentBillingPeriod={{this.handleCurrentBillingPeriod}}
        @startTimeDisplay={{concat "February " this.previousYear}}
      />
    `);

    await calendarDropdown.openCalendar();
    assert.ok(calendarDropdown.showsCalendar, 'renders the calendar component');
  });

  test('it does not allow a user to click to a future year but does allow a user to click to previous year', async function (assert) {
    await render(hbs`
      <CalendarWidget
        @arrayOfMonths={{this.arrayOfMonths}}
        @endTimeDisplay={{concat "March " this.currentYear}}
        @endTimeFromResponse={{this.endTimeFromResponse}}
        @handleClientActivityQuery={{this.handleClientActivityQuery}}
        @handleCurrentBillingPeriod={{this.handleCurrentBillingPeriod}}
        @startTimeDisplay={{concat "February " this.previousYear}}
      />
    `);

    await calendarDropdown.openCalendar();
    assert.dom('[data-test-future-year]').isDisabled('Future year is disabled');

    await calendarDropdown.clickPreviousYear();
    assert.dom('[data-test-display-year]').hasText(this.previousYear.toString(), 'shows the previous year');
    assert
      .dom('[data-test-calendar-month="January"]')
      .hasClass(
        'is-readOnly',
        `January ${this.previousYear} is disabled because it comes before startTimeDisplay`
      );
  });

  test('it disables the current month', async function (assert) {
    await render(hbs`
      <CalendarWidget
        @arrayOfMonths={{this.arrayOfMonths}}
        @endTimeDisplay={{concat "January " this.currentYear}}
        @endTimeFromResponse={{this.endTimeFromResponse}}
        @handleClientActivityQuery={{this.handleClientActivityQuery}}
        @handleCurrentBillingPeriod={{this.handleCurrentBillingPeriod}}
        @startTimeDisplay={{concat "February " this.previousYear}}
      />
    `);
    await calendarDropdown.openCalendar();
    const month = this.arrayOfMonths[new Date().getMonth()];
    assert
      .dom(`[data-test-calendar-month="${month}"]`)
      .hasClass('is-readOnly', `${month} ${this.currentYear} is disabled`);
    // The component also disables all months after the current one, but this
    // is tricky to test since it's based on browser time, so the behavior
    // would be different in december than other months
  });

  test('it allows you to reset the billing period', async function (assert) {
    await render(hbs`
    <CalendarWidget
      @arrayOfMonths={{this.arrayOfMonths}}
      @endTimeDisplay={{concat "January " this.currentYear}}
      @endTimeFromResponse={{this.endTimeFromResponse}}
      @handleClientActivityQuery={{this.handleClientActivityQuery}}
      @handleCurrentBillingPeriod={{this.handleCurrentBillingPeriod}}
      @startTimeDisplay={{concat "February " this.previousYear}}
    />
  `);
    await calendarDropdown.menuToggle();
    await calendarDropdown.clickCurrentBillingPeriod();
    assert.ok(this.handleCurrentBillingPeriod.calledOnce, 'it calls the parents handleCurrentBillingPeriod');
  });

  test('it passes the appropriate data to the handleCurrentBillingPeriod when a date is selected', async function (assert) {
    await render(hbs`
    <CalendarWidget
      @arrayOfMonths={{this.arrayOfMonths}}
      @endTimeDisplay={{concat "January " this.currentYear}}
      @endTimeFromResponse={{this.endTimeFromResponse}}
      @handleClientActivityQuery={{this.handleClientActivityQuery}}
      @handleCurrentBillingPeriod={{this.handleCurrentBillingPeriod}}
      @startTimeDisplay={{concat "February " this.previousYear}}
    />
  `);
    await calendarDropdown.openCalendar();
    await calendarDropdown.clickPreviousYear();
    await click('[data-test-calendar-month="October"]'); // select endTime of October 2021
    assert.ok(this.handleClientActivityQuery.calledOnce, 'it calls the parents handleClientActivityQuery');
    assert.ok(
      this.handleClientActivityQuery.calledWith(9, this.previousYear, 'endTime'),
      'Passes the month as an index, year and date type to the parent'
    );
  });

  test('it displays the year from endTimeDisplay when opened', async function (assert) {
    this.set('endTimeFromResponse', [this.previousYear, 11]);
    await render(hbs`
    <CalendarWidget
      @arrayOfMonths={{this.arrayOfMonths}}
      @endTimeDisplay={{concat "December " this.previousYear}}
      @endTimeFromResponse={{this.endTimeFromResponse}}
      @handleClientActivityQuery={{this.handleClientActivityQuery}}
      @handleCurrentBillingPeriod={{this.handleCurrentBillingPeriod}}
      @startTimeDisplay="March 2020"
    />
  `);
    await calendarDropdown.openCalendar();
    assert
      .dom('[data-test-display-year]')
      .hasText(this.previousYear.toString(), 'Shows year from the end response');
  });
});
