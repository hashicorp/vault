/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { ARRAY_OF_MONTHS, parseAPITimestamp } from 'core/utils/date-formatters';
import { addYears, isSameYear, subYears } from 'date-fns';
import timestamp from 'core/utils/timestamp';
/**
 * @module CalendarWidget
 * CalendarWidget component is used in the client counts dashboard to select a month/year to query the /activity endpoint.
 * The component returns an object with selected date info, example: { dateType: 'endDate', monthIdx: 0, monthName: 'January', year: 2022 }
 *
 * @example
 * ```js
 * <CalendarWidget @startTimestamp={{this.startTime}} @endTimestamp={{this.endTime}} @selectMonth={{this.handleSelection}} />
 *
 *    @param {string} startTimestamp - ISO timestamp string of the calendar widget's start time, displays in dropdown trigger
 *    @param {string} endTimestamp - ISO timestamp string for the calendar widget's end time, displays in dropdown trigger
 *    @param {function} selectMonth - callback function from parent - fires when selecting a month or clicking "Current billing period"
 *  />
 * ```
 */
export default class CalendarWidget extends Component {
  currentDate = timestamp.now();
  @tracked calendarDisplayDate = this.currentDate; // init to current date, updates when user clicks on calendar chevrons
  @tracked showCalendar = false;

  // both date getters return a date object
  get startDate() {
    return parseAPITimestamp(this.args.startTimestamp);
  }
  get endDate() {
    return parseAPITimestamp(this.args.endTimestamp);
  }
  get displayYear() {
    return this.calendarDisplayDate.getFullYear();
  }
  get disableFutureYear() {
    return isSameYear(this.calendarDisplayDate, this.currentDate);
  }
  get disablePastYear() {
    // calendar widget should only go as far back as the passed in start time
    return isSameYear(this.calendarDisplayDate, this.startDate);
  }
  get widgetMonths() {
    const startYear = this.startDate.getFullYear();
    const startMonthIdx = this.startDate.getMonth();
    return ARRAY_OF_MONTHS.map((month, index) => {
      let readonly = false;

      // if widget is showing same year as @startTimestamp year, disable if month is before start month
      if (startYear === this.displayYear && index < startMonthIdx) {
        readonly = true;
      }

      // if widget showing current year, disable if month is later than current month
      if (this.displayYear === this.currentDate.getFullYear() && index > this.currentDate.getMonth()) {
        readonly = true;
      }
      return {
        index,
        year: this.displayYear,
        name: month,
        readonly,
      };
    });
  }

  @action
  addYear() {
    this.calendarDisplayDate = addYears(this.calendarDisplayDate, 1);
  }

  @action
  subYear() {
    this.calendarDisplayDate = subYears(this.calendarDisplayDate, 1);
  }

  @action
  toggleShowCalendar() {
    this.showCalendar = !this.showCalendar;
    this.calendarDisplayDate = this.endDate;
  }

  @action
  handleDateShortcut(dropdown, { target }) {
    this.args.selectMonth({ dateType: target.name }); // send clicked shortcut to parent callback
    this.showCalendar = false;
    dropdown.close();
  }

  @action
  selectMonth(month, dropdown) {
    const { index, year, name } = month;
    this.toggleShowCalendar();
    this.args.selectMonth({ monthIdx: index, monthName: name, year, dateType: 'endDate' });
    dropdown.close();
  }
}
