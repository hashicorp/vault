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
}

export default setComponentTemplate(layout, CalendarWidget);
