/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { ARRAY_OF_MONTHS } from 'core/utils/date-formatters';
import timestamp from 'core/utils/timestamp';
/**
 * @module DateDropdown
 * DateDropdown components are used to display a dropdown of months and years to handle date selection. Future dates are disabled (current month and year are selectable).
 * The component returns an object with selected date info, example: { dateType: 'start', monthIdx: 0, monthName: 'January', year: 2022 }
 *
 * @example
 * ```js
 * <DateDropdown @handleSubmit={{this.actionFromParent}} @name="startTime" @submitText="Save" @handleCancel={{this.onCancel}}/>
 * ```
 * @param {function} handleSubmit - callback function from parent that the date picker triggers on submit
 * @param {function} [handleCancel] - optional callback for cancel action, if exists then buttons appear modal style with a light gray background
 * @param {string} [dateType] - optional argument to give the selected month/year a type
 * @param {string} [submitText] - optional argument to change submit button text
 * @param {function} [validateDate] - parent function to validate date selection, receives date object and returns an error message that's passed to the inline alert
 */
export default class DateDropdown extends Component {
  currentDate = timestamp.now();
  currentYear = this.currentDate.getFullYear(); // integer of year
  currentMonthIdx = this.currentDate.getMonth(); // integer of month, 0 indexed
  dropdownMonths = ARRAY_OF_MONTHS.map((m, i) => ({ name: m, index: i }));
  dropdownYears = Array.from({ length: 5 }, (item, i) => this.currentYear - i);

  @tracked maxMonthIdx = 11; // disables months with index greater than this number, initially all months are selectable
  @tracked disabledYear = null; // year as integer if current year should be disabled
  @tracked selectedMonth = null;
  @tracked selectedYear = null;
  @tracked invalidDate = null;

  @action
  selectMonth(month, dropdown) {
    this.selectedMonth = month;
    // disable current year if selected month is later than current month
    this.disabledYear = month.index > this.currentMonthIdx ? this.currentYear : null;
    dropdown.close();
  }

  @action
  selectYear(year, dropdown) {
    this.selectedYear = year;
    // disable months after current month if selected year is current year
    this.maxMonthIdx = year === this.currentYear ? this.currentMonthIdx : 11;
    dropdown.close();
  }

  @action
  handleSubmit() {
    if (this.args.validateDate) {
      this.invalidDate = null;
      this.invalidDate = this.args.validateDate(new Date(this.selectedYear, this.selectedMonth.index));
      if (this.invalidDate) return;
    }
    const { index, name } = this.selectedMonth;
    this.args.handleSubmit({
      monthIdx: index,
      monthName: name,
      year: this.selectedYear,
      dateType: this.args.dateType,
    });
    this.resetDropdown();
  }

  @action
  handleCancel() {
    this.args.handleCancel();
    this.resetDropdown();
  }

  resetDropdown() {
    this.maxMonthIdx = 11;
    this.disabledYear = null;
    this.selectedMonth = null;
    this.selectedYear = null;
    this.invalidDate = null;
  }
}
