import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

/**
 * @module DateDropdown
 * DateDropdown components are used to display a dropdown of months and years to handle date selection
 *
 * @example
 * ```js
 * <DateDropdown @handleDateSelection={this.actionFromParent} @name={{"startTime"}}/>
 * ```
 * @param {function} handleDateSelection - is the action from the parent that the date picker triggers
 * @param {string} [name] - optional argument passed from date dropdown to parent function
 */

export default class DateDropdown extends Component {
  @tracked startMonth = null;
  @tracked startYear = null;

  months = Array.from({ length: 12 }, (item, i) => {
    return new Date(0, i).toLocaleString('en-US', { month: 'long' });
  });
  years = Array.from({ length: 5 }, (item, i) => {
    return new Date().getFullYear() - i;
  });

  @action
  selectStartMonth(month) {
    this.startMonth = month;
  }

  @action
  selectStartYear(year) {
    this.startYear = year;
  }
}
