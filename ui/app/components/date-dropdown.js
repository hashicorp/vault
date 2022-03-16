import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

/**
 * @module DateDropdown
 * DateDropdown components are used to display a dropdown of months and years to handle date selection
 *
 * @example
 * ```js
 * <DateDropdown @handleDateSelection={this.actionFromParent} @name={{"startTime"}} @submitText="Save"/>
 * ```
 * @param {function} handleDateSelection - is the action from the parent that the date picker triggers
 * @param {string} [name] - optional argument passed from date dropdown to parent function
 * @param {string} [submitText] - optional argument to change submit button text
 */
export default class DateDropdown extends Component {
  currentDate = new Date();
  currentYear = this.currentDate.getFullYear(); // integer of year
  currentMonth = this.currentDate.getMonth(); // index of month

  @tracked allowedMonthMax = 12;
  @tracked disabledYear = null;
  @tracked startMonth = null;
  @tracked startYear = null;

  months = Array.from({ length: 12 }, (item, i) => {
    return new Date(0, i).toLocaleString('en-US', { month: 'long' });
  });
  years = Array.from({ length: 5 }, (item, i) => {
    return new Date().getFullYear() - i;
  });

  @action
  selectStartMonth(month, event) {
    this.startMonth = month;
    // disables months if in the future
    this.disabledYear = this.months.indexOf(month) >= this.currentMonth ? this.currentYear : null;
    event.close();
  }

  @action
  selectStartYear(year, event) {
    this.startYear = year;
    this.allowedMonthMax = year === this.currentYear ? this.currentMonth : 12;
    event.close();
  }

  @action
  saveDateSelection() {
    this.args.handleDateSelection(this.startMonth, this.startYear, this.args.name);
  }
}
