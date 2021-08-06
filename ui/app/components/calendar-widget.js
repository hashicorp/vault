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
  currentYear = parseInt(format(this.currentDate(), 'yyyy')); // integer
  currentMonth = parseInt(format(this.currentDate(), 'M')); // integer

  @tracked displayYear = this.currentYear;
  @tracked disablePastYear = this.isObsoleteYear(); // if obsolete year, disable left chevron
  @tracked disableFutureYear = this.isCurrentYear(); // if current year, disable right chevron
  @tracked quickMonthsSelection = null;
  @tracked allMonthsNodeList = [];
  @tracked isClearAllMonths = false;
  @tracked areMonthsSelected = false;
  @tracked shiftClickCount = 0;
  @tracked startMonth;
  @tracked endMonth;
  @tracked shiftClickRange = [];

  // HELPER FUNCTIONS //

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
    // if clearMonths was clicked new dom elements are render and we need to clear any selected months
    this.displayYear = parseInt(this.displayYear) - 1;
    this.selectMonths(this.quickMonthsSelection);
    // call disable months action
    this.disableMonths();
    this.disableFutureYear = this.isCurrentYear();
    this.disablePastYear = this.isObsoleteYear();
    if (this.isClearAllMonths) {
      this.deselectAllMonths();
    }
  }

  @action
  addYear() {
    this.displayYear = parseInt(this.displayYear) + 1;
    this.selectMonths(this.quickMonthsSelection);
    // call disable months action
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
    // if one month is selected, then proceed else return
    // if click + shift again, find range
    e.target.classList.contains('is-selected')
      ? this.removeClass(e.target, 'is-selected')
      : this.addClass(e.target, 'is-selected');

    this.allMonthsNodeList.forEach(e => {
      if (e.classList.contains('is-selected')) {
        this.areMonthsSelected = true;
      }
    });

    if (e.shiftKey) {
      let monthArray = [];
      this.allMonthsNodeList.forEach(e => {
        monthArray.push(e);
      });
      let reverseMonthArray = monthArray.reverse();

      // count shift clicks
      this.shiftClickCount = ++this.shiftClickCount;
      if (this.shiftClickCount > 2) {
        this.deselectMonths();
        this.shiftClickCount = 0;
        return;
      }
      // grab start month
      if (this.shiftClickCount === 1) {
        this.allMonthsNodeList.forEach(e => {
          if (e.classList.contains('is-selected')) {
            this.startMonth = e.id;
            return;
          }
        });
      }
      // grab end month
      let isSelectedArray = [];
      if (this.shiftClickCount === 2) {
        this.endMonth = reverseMonthArray.forEach(e => {
          if (e.classList.contains('is-selected')) {
            isSelectedArray.push(e.id);
            return;
          }
        });
        this.endMonth = isSelectedArray[0];

        // create a range
        // split last months
        this.shiftClickRange = this.createRange(
          parseInt(this.startMonth.split('-')[1]),
          parseInt(this.endMonth.split('-')[1])
        ).map(n => `month-${n}`);

        this.shiftClickRange.forEach(id => {
          this.allMonthsNodeList.forEach(e => {
            if (e.id === id) {
              this.addClass(e, 'is-selected');
            }
          });
        });
      }
    }
  }

  @action
  selectMonths(lastXNumberOfMonths) {
    this.quickMonthsSelection = lastXNumberOfMonths;
    // deselect all elements before reselecting
    this.deselectAllMonths();
    // if the user has not selected anything exit function
    if (lastXNumberOfMonths === null) {
      return;
    }
    this.areMonthsSelected = true;
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
    this.allMonthsNodeList.forEach(monthElement => {
      lastYearSelectedRangeIdsArray.includes(monthElement.id)
        ? previousYearMonthElementsArray.push(monthElement)
        : '';
    });

    // array of ids for months selected from current year
    let selectedRangeIdsArray = selectedRange.filter(n => n > 0).map(n => `month-${n}`);

    let currentYearMonthElementsArray = [];
    this.allMonthsNodeList.forEach(monthElement => {
      selectedRangeIdsArray.includes(monthElement.id) ? currentYearMonthElementsArray.push(monthElement) : '';
    });

    // add selector class to month elements for both last year and current year
    currentYearMonthElementsArray.forEach(element => {
      if (this.currentYear === this.displayYear) {
        this.addClass(element, 'is-selected');
      }
    });

    previousYearMonthElementsArray.forEach(element => {
      if (parseInt(this.currentYear) - 1 === this.displayYear) {
        this.addClass(element, 'is-selected');
      }
    });
  }
}
export default setComponentTemplate(layout, CalendarWidget);
