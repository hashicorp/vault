import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click } from '@ember/test-helpers';
import sinon from 'sinon';
import hbs from 'htmlbars-inline-precompile';
import calendarDropdown from 'vault/tests/pages/components/calendar-widget';
import { ARRAY_OF_MONTHS } from 'core/utils/date-formatters';

module('Integration | Component | calendar-widget', function (hooks) {
  setupRenderingTest(hooks);

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

  test('it does not allow a user to click to a future year but does allow a user to click to previous year', async function (assert) {
    await render(hbs`
      <CalendarWidget
        @arrayOfMonths={{arrayOfMonths}}
        @endTimeDisplay={{"March 2022"}}
        @endTimeFromResponse={{endTimeFromResponse}}
        @handleClientActivityQuery={{handleClientActivityQuery}}
        @handleCurrentBillingPeriod={{handleCurrentBillingPeriod}}
        @startTimeDisplay={{"February 2021"}}
      />
    `);

    await calendarDropdown.openCalendar();
    assert.dom('[data-test-future-year]').isDisabled('Future year is disabled');

    await calendarDropdown.clickPreviousYear();
    assert.dom('[data-test-display-year]').hasText('2021', 'shows the previous year');
    assert
      .dom('[data-test-calendar-month="January"]')
      .hasClass('is-readOnly', 'January 2021 is disabled because it comes before February 2021');
  });

  test('it enables the current month but disables future months', async function (assert) {
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
    assert
      .dom('[data-test-calendar-month="January"]')
      .doesNotHaveClass('is-readOnly', 'January 2022 is enabled');
    assert.dom('[data-test-calendar-month="February"]').hasClass('is-readOnly', 'February 2022 is enabled');
  });

  test('it allows you to reset the billing period', async function (assert) {
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
    await calendarDropdown.menuToggle();
    await calendarDropdown.clickCurrentBillingPeriod();
    assert.ok(this.handleCurrentBillingPeriod.calledOnce, 'it calls the parents handleCurrentBillingPeriod');
  });

  test('it passes the appropriate data to the handleCurrentBillingPeriod when a date is selected', async function (assert) {
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
    await calendarDropdown.clickPreviousYear();
    await click('[data-test-calendar-month="October"]'); // select endTime of October 2021
    assert.ok(this.handleClientActivityQuery.calledOnce, 'it calls the parents handleClientActivityQuery');
    assert.ok(
      this.handleClientActivityQuery.calledWith(9, 2021, 'endTime'),
      'Passes the month as an index, year and date type to the parent'
    );
  });
});
