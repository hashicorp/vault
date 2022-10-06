import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

/**
 * @module DateDropdown
 * DateDropdown components are used to display a dropdown of months and years to handle date selection. Future dates are disabled (current month and year are selectable).
 *
 * @example
 * ```js
 * <DateDropdown @handleSubmit={{this.actionFromParent}} @name="startTime" @submitText="Save" @handleCancel={{this.onCancel}}/>
 * ```
 * @param {function} handleSubmit - callback function from parent that the date picker triggers on submit
 * @param {function} [handleCancel] - optional callback for cancel action, if passed in buttons appear modal style with a light gray background
 * @param {string} [name] - optional argument passed from date dropdown to parent function, could be used to identify dropdown if there are multiple
 * @param {string} [submitText] - optional argument to change submit button text
 */
export default class DateDropdown extends Component {
  currentDate = new Date();
  currentYear = this.currentDate.getFullYear(); // integer of year
  currentMonth = this.currentDate.getMonth(); // index of month

  @tracked maxMonthIdx = 11; // disables months with index greater than this number, initially all months are selectable
  @tracked disabledYear = null; // year as integer if current year should be disabled
  @tracked selectedMonth = null;
  @tracked selectedYear = null;

  months = Array.from({ length: 12 }, (item, i) => {
    return new Date(0, i).toLocaleString('en-US', { month: 'long' });
  });
  years = Array.from({ length: 5 }, (item, i) => {
    return new Date().getFullYear() - i;
  });

  @action
  selectMonth(month, dropdown) {
    this.selectedMonth = month;
    // disable current year if selected month is later than current month
    this.disabledYear = this.months.indexOf(month) > this.currentMonth ? this.currentYear : null;
    dropdown.close();
  }

  @action
  selectYear(year, dropdown) {
    this.selectedYear = year;
    // disable months after current month if selected year is current year
    this.maxMonthIdx = year === this.currentYear ? this.currentMonth : 11;
    dropdown.close();
  }

  @action
  handleSubmit() {
    this.args.handleSubmit(this.selectedMonth, this.selectedYear, this.args.name);
  }

  @action
  handleCancel() {
    this.selectedMonth = null;
    this.selectedYear = null;
    this.args.handleCancel();
  }
}
