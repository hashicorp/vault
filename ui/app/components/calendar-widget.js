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
  currentYear = parseInt(format(this.currentDate(), 'yyyy'));
  currentMonth = parseInt(format(this.currentDate(), 'M'));

  @tracked isActive = false;
  @tracked isSelected = false;
  @tracked selectedYear = this.currentYear;
  @tracked quickMonthsSelection = null;
  @tracked isDisabled = this.currentYear === this.selectedYear;
  @tracked allMonthsArray = [];

  @action
  disableMonths() {
    let getMonths = document.querySelectorAll('.month-list');
    this.allMonthsArray = getMonths;
    this.allMonthsArray.forEach(e => {
      let elementMonthId = parseInt(e.id.split('-')[1]);
      if (this.currentMonth <= elementMonthId) {
        e.classList.add('is-readOnly');
      }
    });
  }

  calculateLastMonth() {
    return sub(this.currentDate(), { months: 1 });
  }

  currentDate() {
    return new Date();
  }

  @action
  subYear() {
    this.selectedYear = parseInt(this.selectedYear) - 1;
    this.selectMonths(this.quickMonthsSelection);
    this.isDisabled = this.currentYear === this.selectedYear;
  }

  @action
  addYear() {
    this.selectedYear = parseInt(this.selectedYear) + 1;
    this.selectMonths(this.quickMonthsSelection);
    this.isDisabled = this.currentYear === this.selectedYear;
  }

  @action
  selectMonth(e) {
    e.target.classList.contains('is-selected')
      ? e.target.classList.remove('is-selected')
      : e.target.classList.add('is-selected');
  }

  createRange(start, end) {
    return Array(end - start + 1)
      .fill()
      .map((_, idx) => start + idx);
  }

  @action
  selectMonths(number) {
    this.quickMonthsSelection = number;
    this.allMonthsArray.forEach(monthElement => {
      monthElement.classList.remove('is-selected');
    });
    // define current month
    let lastMonth = parseInt(format(this.currentDate(), 'M')) - 1; // subtract one to skip current month
    let startRange = lastMonth - number;
    let selectedRange = this.createRange(startRange, lastMonth); // returns array of integers
    console.log(selectedRange);
    let previousYearMonthElementsArray = [];
    let lastYearSelectedRangeIdsArray = selectedRange.filter(n => n < 0).map(n => `month-${n + 13}`);
    this.allMonthsArray.forEach(monthElement => {
      lastYearSelectedRangeIdsArray.includes(monthElement.id)
        ? previousYearMonthElementsArray.push(monthElement)
        : '';
    });

    // select current year months
    let selectedRangeIdsArray = selectedRange.filter(n => n > 0).map(n => `month-${n}`);
    let currentYearMonthElementsArray = [];
    this.allMonthsArray.forEach(monthElement => {
      selectedRangeIdsArray.includes(monthElement.id) ? currentYearMonthElementsArray.push(monthElement) : '';
    });
    currentYearMonthElementsArray.forEach(element => {
      console.log(this.currentYear, 'current');
      console.log(this.selectedYear, 'selected');
      if (this.currentYear === this.selectedYear) {
        element.classList.add('is-selected');
      }
    });

    previousYearMonthElementsArray.forEach(element => {
      if (parseInt(this.currentYear) - 1 === this.selectedYear) {
        element.classList.add('is-selected');
      }
    });
  }
}
export default setComponentTemplate(layout, CalendarWidget);
