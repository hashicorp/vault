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
  @tracked isClearAllMonths = false;
  @tracked areAnyMonthsSelected = false;

  @action
  disableMonths() {
    let getMonths = document.querySelectorAll('.is-month-list');
    this.allMonthsArray = getMonths;
    this.allMonthsArray.forEach(e => {
      // clear all is-readOnly classes and start over.
      e.classList.remove('is-readOnly');
      let elementMonthId = parseInt(e.id.split('-')[1]);
      if (this.currentMonth <= elementMonthId) {
        // only disable months when current year is selected
        if (this.selectedYear === this.currentYear) {
          e.classList.add('is-readOnly');
        }
      }
    });
  }
  @action
  clearSelectedMonths() {
    this.isClearAllMonths = !this.isClearAllMonths;
    this.areAnyMonthsSelected = false;
    this.allMonthsArray.forEach(e => {
      // clear all selected months
      e.classList.remove('is-selected');
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
    // if clearMonths was clicked new dom elements are render and we need to clear any selected months
    console.log(this.isClearAllMonths, 'clearAllMOnths');

    this.selectedYear = parseInt(this.selectedYear) - 1;
    this.selectMonths(this.quickMonthsSelection);
    // call disable months action
    this.disableMonths();
    this.isDisabled = this.currentYear === this.selectedYear;
    if (this.isClearAllMonths) {
      this.allMonthsArray.forEach(e => {
        e.classList.remove('is-selected');
      });
    }
  }

  @action
  addYear() {
    this.selectedYear = parseInt(this.selectedYear) + 1;
    this.selectMonths(this.quickMonthsSelection);
    // call disable months action
    this.disableMonths();
    this.isDisabled = this.currentYear === this.selectedYear;
    // if clearMonths was clicked new dom elements are render and we need to clear any selected months
    if (this.isClearAllMonths) {
      this.allMonthsArray.forEach(e => {
        e.classList.remove('is-selected');
      });
    }
  }

  @action
  selectMonth(e) {
    e.target.classList.contains('is-selected')
      ? e.target.classList.remove('is-selected')
      : e.target.classList.add('is-selected');

    this.allMonthsArray.forEach(e => {
      if (e.classList.contains('is-selected')) {
        this.areAnyMonthsSelected = true;
      }
    });
  }

  createRange(start, end) {
    return Array(end - start + 1)
      .fill()
      .map((_, idx) => start + idx);
  }

  @action
  selectMonths(lastXNumberOfMonths) {
    this.quickMonthsSelection = lastXNumberOfMonths;
    // deselect all elements before reselecting
    this.allMonthsArray.forEach(monthElement => {
      monthElement.classList.remove('is-selected');
    });
    // if the user has not selected anything exit function
    if (lastXNumberOfMonths === null) {
      return;
    }
    this.areAnyMonthsSelected = true;
    // reports are not available for current month so we don't want it in range
    let endMonth = this.currentMonth - 1;
    // start range X months back, subtract one to skip current month
    let startRange = endMonth - lastXNumberOfMonths;
    // creates array of selected months (integers)
    let selectedRange = this.createRange(startRange, endMonth);
    // array of ids for months selected from previous year
    let lastYearSelectedRangeIdsArray = selectedRange.filter(n => n < 0).map(n => `month-${n + 13}`);

    // array of month elements
    let previousYearMonthElementsArray = [];
    this.allMonthsArray.forEach(monthElement => {
      lastYearSelectedRangeIdsArray.includes(monthElement.id)
        ? previousYearMonthElementsArray.push(monthElement)
        : '';
    });

    // array of ids for months selected from current year
    let selectedRangeIdsArray = selectedRange.filter(n => n > 0).map(n => `month-${n}`);

    let currentYearMonthElementsArray = [];
    this.allMonthsArray.forEach(monthElement => {
      selectedRangeIdsArray.includes(monthElement.id) ? currentYearMonthElementsArray.push(monthElement) : '';
    });

    // add selector class to month elements for both last year and current year
    currentYearMonthElementsArray.forEach(element => {
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
