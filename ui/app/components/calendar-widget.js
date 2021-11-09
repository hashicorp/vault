/**
 * @module CalendarWidget
 * CalendarWidget components are used in the client counts metrics. It helps user understand the ranges they can select.
 *
 * @example
 * ```js
 * <CalendarWidget/>
 * ```
 */

import Component from '@glimmer/component';
import layout from '../templates/components/calendar-widget';
import { setComponentTemplate } from '@ember/component';
import { format, sub } from 'date-fns';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

class CalendarWidget extends Component {
  currentYear = parseInt(format(this.currentDate(), 'yyyy')); // integer
  currentMonth = parseInt(format(this.currentDate(), 'M')); // integer

  @tracked displayYear = this.currentYear;
  @tracked disablePastYear = this.isObsoleteYear(); // if obsolete year, disable left chevron
  @tracked disableFutureYear = this.isCurrentYear(); // if current year, disable right chevron
  @tracked preselectRangeOfMonths = null;
  @tracked allMonthsNodeList = [];
  @tracked isClearAllMonths = false;
  @tracked areMonthsSelected = false;
  @tracked shiftClickCount = 0;
  @tracked shiftClickRange = [];

  // HELPER FUNCTIONS //

  checkIfMonthsSelected() {
    this.allMonthsNodeList.forEach(e => {
      if (e.classList.contains('is-selected')) {
        this.areMonthsSelected = true;
      }
    });
  }
  calculateLastMonth() {
    return sub(this.currentDate(), { months: 1 });
  }

  currentDate() {
    return new Date();
  }

  isCurrentYear() {
    return this.currentYear === this.displayYear;
  }

  isObsoleteYear() {
    return this.displayYear === this.currentYear - 4; // won't display more than 5 years ago
  }

  deselectAllMonths() {
    this.allMonthsNodeList.forEach(element => {
      this.removeClass(element, 'is-selected');
    });
  }

  removeClass(element, classString) {
    element.classList.remove(classString);
  }

  addClass(element, classString) {
    element.classList.add(classString);
  }

  createRange(start, end) {
    return Array(end - start + 1)
      .fill()
      .map((_, idx) => start + idx);
  }

  // ACTIONS //

  @action
  disableMonths() {
    this.allMonthsNodeList = document.querySelectorAll('.is-month-list');
    this.allMonthsNodeList.forEach(e => {
      // clear all is-readOnly classes and start over.
      this.removeClass(e, 'is-readOnly');
      let elementMonthId = parseInt(e.id.split('-')[1]);
      if (this.currentMonth <= elementMonthId) {
        // only disable months when current year is selected
        if (this.isCurrentYear()) {
          e.classList.add('is-readOnly');
        }
      }
    });
  }

  @action
  deselectMonths() {
    this.isClearAllMonths = !this.isClearAllMonths;
    this.areMonthsSelected = false;
    this.deselectAllMonths();
  }

  @action
  subYear() {
    this.displayYear = this.displayYear - 1;
    this.selectMonths(this.preselectRangeOfMonths);
    // call disable months action
    this.disableMonths();
    this.disableFutureYear = this.isCurrentYear();
    this.disablePastYear = this.isObsoleteYear();
    // if clearMonths was clicked new dom elements are rendered and we need to clear any selected months
    if (this.isClearAllMonths) {
      this.deselectAllMonths();
    }
  }

  @action
  addYear() {
    this.displayYear = this.displayYear + 1;
    this.selectMonths(this.preselectRangeOfMonths);
    this.disableMonths();
    this.disableFutureYear = this.isCurrentYear();
    this.disablePastYear = this.isObsoleteYear();
    // if clearMonths was clicked new dom elements are render and we need to clear any selected months
    if (this.isClearAllMonths) {
      this.deselectAllMonths();
    }
  }

  @action
  selectMonth(e) {
    e.target.classList.contains('is-selected')
      ? this.removeClass(e.target, 'is-selected')
      : this.addClass(e.target, 'is-selected');

    this.checkIfMonthsSelected();

    if (e.shiftKey) {
      this.handleShift();
    }
  }

  reverseMonthNodeList() {
    let reverseMonthArray = [];
    this.allMonthsNodeList.forEach(e => {
      reverseMonthArray.unshift(e);
    });
    return reverseMonthArray;
  }

  handleShift() {
    // count shift clicks
    this.shiftClickCount = ++this.shiftClickCount;

    // if going wild with shift clicks, reset count and deselect all months
    if (this.shiftClickCount > 2) {
      this.deselectMonths();
      this.shiftClickCount = 0;
      return;
    }

    let startAndEndMonths = [];
    if (this.shiftClickCount === 2) {
      this.allMonthsNodeList.forEach(e => {
        if (e.classList.contains('is-selected')) {
          startAndEndMonths.push(parseInt(e.id.split('-')[1]));
          return;
        }
      });
    }

    this.shiftClickRange = this.createRange(startAndEndMonths[0], startAndEndMonths[1]).map(
      n => `month-${n}`
    );

    this.shiftClickRange.forEach(id => {
      this.allMonthsNodeList.forEach(e => {
        if (e.id === id) {
          this.addClass(e, 'is-selected');
        }
      });
    });
  }

  @action
  selectMonths(quickSelectNumber) {
    this.preselectRangeOfMonths = quickSelectNumber;
    this.deselectAllMonths();
    // if the user has not selected anything exit function
    if (quickSelectNumber === null) {
      return;
    }
    this.areMonthsSelected = true;
    // exclude current month in range
    let endMonth = this.currentMonth - 1;
    // start range quickSelectNumber of months back
    let startRange = endMonth - quickSelectNumber;
    // creates array of integers correlating to selected months
    let selectedRange = this.createRange(startRange, endMonth);
    // array of month-ids for months selected from previous year
    let previousYearMonthIds = selectedRange.filter(n => n < 0).map(n => `month-${n + 13}`);

    // array of month-ids for months selected from current year
    let currentYearMonthIds = selectedRange.filter(n => n > 0).map(n => `month-${n}`);

    let previousYearMonthElements = [];
    this.allMonthsNodeList.forEach(monthElement => {
      previousYearMonthIds.includes(monthElement.id) ? previousYearMonthElements.push(monthElement) : '';
    });

    let currentYearMonthElements = [];
    this.allMonthsNodeList.forEach(monthElement => {
      currentYearMonthIds.includes(monthElement.id) ? currentYearMonthElements.push(monthElement) : '';
    });

    // iterate array of current year month elements and select
    currentYearMonthElements.forEach(element => {
      if (this.currentYear === this.displayYear) {
        this.addClass(element, 'is-selected');
      }
    });

    // iterate array of previous year month elements and select
    previousYearMonthElements.forEach(element => {
      if (this.currentYear - 1 === this.displayYear) {
        this.addClass(element, 'is-selected');
      }
    });

    // get quick select number
    // use current date to determine range in MM/yyyy format
    // currentMonth

    // TODO: fix selection, does not select 12 months back properly

    let date = new Date();
    // will never query same month that you are in, always add a month to get the N months prior
    date.setMonth(date.getMonth() - quickSelectNumber - 1);

    let dateAgain = new Date();
    dateAgain.setMonth(dateAgain.getMonth() - 1);

    let startDate = format(dateAgain, 'MM-yyyy');
    let endDate = format(date, 'MM-yyyy');
    this.args.handleQuery(endDate, startDate);
  }
}
export default setComponentTemplate(layout, CalendarWidget);
