/**
 * @module CalendarWidget
 * CalendarWidget components are used to...
 *
 * @example
 * ```js
 * <CalendarWidget @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

import Component from '@glimmer/component';
import layout from '../templates/components/calendar-widget';
import { setComponentTemplate } from '@ember/component';
import { format, sub, add, eachMonthOfInterval } from 'date-fns';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

class CalendarWidget extends Component {
  startMonthRange = format(this.calculateLastMonth(), 'M/yyyy');
  endMonthRange = format(this.currentDate(), 'M/yyyy');

  @tracked
  isActive = false;

  @tracked isSelected = false;

  @tracked
  // will need to be in API appropriate format, using parseInt here for hack-y functionality
  displayYear = parseInt(format(this.currentDate(), 'yyyy'));

  calculateLastMonth() {
    return sub(this.currentDate(), { months: 1 });
  }

  currentDate() {
    return new Date();
  }

  @action
  subYear() {
    this.displayYear -= 1;
  }

  @action
  addYear() {
    this.displayYear += 1;
  }

  @action
  selectMonth(e) {
    if (e.target.className === 'is-not-selected') {
      e.target.className = 'is-selected';
    } else {
      e.target.className = 'is-not-selected';
    }
  }

  @action
  selectMonths(e) {
    const innerText = e.target.textContent;
    let result;
    let months = [];
    switch (innerText) {
      case 'Last month':
        result = this.calculateLastMonth();
        this.selectJuly = true;
        months.push(format(result, 'MMMM'));
        console.log(months);
        break;
      case 'Last 3 months':
        result = eachMonthOfInterval({
          start: sub(this.currentDate(), { months: 3 }),
          end: this.currentDate(),
        });
        result.forEach(date => months.push(format(date, 'MMMM')));
        console.log(months);
        break;
      case 'Last 6 months':
        result = eachMonthOfInterval({
          start: sub(this.currentDate(), { months: 6 }),
          end: this.currentDate(),
        });
        result.forEach(date => months.push(format(date, 'MMMM')));
        console.log(months);
        break;
      case 'Last 12 months':
        result = eachMonthOfInterval({
          start: sub(this.currentDate(), { months: 12 }),
          end: this.currentDate(),
        });
        result.forEach(date => months.push(format(date, 'MMMM')));
        console.log(months);
        break;
      default:
        console.log('Incorrect input');
    }
  }
}

export default setComponentTemplate(layout, CalendarWidget);
