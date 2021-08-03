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
import { format, sub, add } from 'date-fns';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

class CalendarWidget extends Component {
  startMonthRange = format(this.calculateLastMonth(), 'M/yyyy');
  endMonthRange = format(this.currentDate(), 'M/yyyy');

  @tracked
  isActive = false;

  @tracked
  // will need to be in API appropriate format, using parseInt here for hack-y functionality
  displayYear = parseInt(format(this.currentDate(), 'yyyy'));

  calculateLastMonth() {
    return sub(this.currentDate(), { months: 1 });
  }

  currentDate() {
    return Date.now();
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
  selectMonths(e) {
    const innerText = e.target.textContent;
    switch (innerText) {
      case 'Last month':
        console.log('okay just go back a month');
        break;
      case 'Last 3 months':
        console.log('nice 3 months');
        break;
      case 'Last 6 months':
        console.log('woah 6 months');
        break;
      case 'Last 12 months':
        console.log('a whole dang year?!');
        break;
      default:
        console.log('Incorrect input');
    }
  }
}

export default setComponentTemplate(layout, CalendarWidget);
