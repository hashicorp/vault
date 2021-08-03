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
import { format, sub } from 'date-fns';

class CalendarWidget extends Component {
  lastMonth = format(this.calculateLastMonth(), 'M/yyyy');
  presentMonth = format(Date.now(), 'M/yyyy');

  calculateLastMonth() {
    return sub(Date.now(), { months: 1 });
  }
}

export default setComponentTemplate(layout, CalendarWidget);
