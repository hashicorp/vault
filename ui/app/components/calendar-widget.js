// ARG TODO documentation takes start and end for handQuery
/**
 * @module CalendarWidget
 * CalendarWidget components are used in the client counts metrics. It helps users understand the ranges they can select.
 *
 * @example
 * ```js
 * <CalendarWidget
 * @param {function} handleQuery - calls the parent pricing-metrics-dates handleQueryFromCalendar method which sends the data for the network request.
 * />
 *
 * ```
 */

import Component from '@glimmer/component';
import layout from '../templates/components/calendar-widget';
import { setComponentTemplate } from '@ember/component';
import { format, sub } from 'date-fns';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

class CalendarWidget extends Component {
  currentDate = new Date();
  currentYear = parseInt(format(this.currentDate, 'yyyy')); // integer
  currentMonth = parseInt(format(this.currentDate, 'M')); // integer

  @tracked displayYear = this.currentYear;
  @tracked disablePastYear = this.isObsoleteYear(); // if obsolete year, disable left chevron
  @tracked disableFutureYear = this.isCurrentYear(); // if current year, disable right chevron
  @tracked preselectRangeOfMonths = null;
  @tracked allMonthsNodeList = [];
  @tracked isClearAllMonths = false;
  @tracked areMonthsSelected = false;
  @tracked mouseClickCount = 0;
  @tracked clickRange = [];
  @tracked startDate; // the older time e.g. 10/2020
  @tracked endDate; // the newer time e.g. 11/2021

  constructor() {
    super(...arguments);
    this.startDate = this.calculateStartDate();
    this.endDate = this.calculateEndDate();
  }

  // HELPER FUNCTIONS //
  calculateStartDate(quickSelectNumber) {
    let date = new Date(); // need to modify so define here and not globally
    // will never query same month that you are in, always add a month to get the N months prior
    // defaults to one year selected if no quickSelectNumber
    date.setMonth(date.getMonth() - (quickSelectNumber ? quickSelectNumber : 11) - 1);
    console.log(date, 'start');
    return format(date, 'MM-yyyy');
  }

  calculateEndDate() {
    let date = new Date(); // need to modify so define here and not globally
    date.setMonth(date.getMonth() - 1);
    return format(date, 'MM-yyyy');
  }

  checkIfMonthsSelected() {
    // ARG TODO going to have issue with display Year and gather multiple years
    // ARG TODO we should also automatically select the years between if they select two months
    let selectedArray = [];
    this.allMonthsNodeList.forEach(e => {
      if (e.classList.contains('is-selected')) {
        this.areMonthsSelected = true;
        selectedArray.push(e.id);
        // set start date the older time
        let sortedSelected = selectedArray.sort();
        this.startDate = `${sortedSelected[0].split('-')[1]}-${this.displayYear}`;

        // set end date the newer time
        let reverseSelected = selectedArray.reverse();
        this.endDate = `${reverseSelected[0].split('-')[1]}-${this.displayYear}`;
      }
    });
    // then set the date on the query
    this.args.handleQuery(this.startDate, this.endDate);
  }

  calculateLastMonth() {
    return sub(this.currentDate, { months: 1 });
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
    this.quickSelectMonths(this.preselectRangeOfMonths);
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
    this.quickSelectMonths(this.preselectRangeOfMonths);
    this.disableMonths();
    this.disableFutureYear = this.isCurrentYear();
    this.disablePastYear = this.isObsoleteYear();
    // if clearMonths was clicked new dom elements are render and we need to clear any selected months
    if (this.isClearAllMonths) {
      this.deselectAllMonths();
    }
  }

  @action // individually click on months
  selectMonth(e) {
    // if three clicks you want to clear the months and not send the start date.
    this.mouseClickCount = ++this.mouseClickCount;
    if (this.mouseClickCount > 2) {
      this.deselectMonths();
      this.mouseClickCount = 0;
      return;
    }

    e.target.classList.contains('is-selected')
      ? this.removeClass(e.target, 'is-selected')
      : this.addClass(e.target, 'is-selected');

    this.checkIfMonthsSelected();
    this.handleSelectRange(e);
  }

  reverseMonthNodeList() {
    let reverseMonthArray = [];
    this.allMonthsNodeList.forEach(e => {
      reverseMonthArray.unshift(e);
    });
    return reverseMonthArray;
  }

  handleSelectRange() {
    let startAndEndMonths = [];

    this.allMonthsNodeList.forEach(e => {
      if (e.classList.contains('is-selected')) {
        startAndEndMonths.push(parseInt(e.id.split('-')[1]));
        return;
      }
    });

    if (startAndEndMonths.length < 2) {
      // exit because you have one month selected
      return;
    }
    this.clickRange = this.createRange(startAndEndMonths[0], startAndEndMonths[1]).map(n => `month-${n}`);
    this.clickRange.forEach(id => {
      this.allMonthsNodeList.forEach(e => {
        if (e.id === id) {
          this.addClass(e, 'is-selected');
        }
      });
    });
  }

  @action
  quickSelectMonths(quickSelectNumber) {
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
    let previousYearMonthIds = selectedRange.filter(n => n <= 0).map(n => `month-${n + 12}`);

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

    // TODO: fix selection, does not select 12 months back properly

    this.startDate = this.calculateStartDate(quickSelectNumber);
    this.args.handleQuery(this.startDate, this.calculateEndDate());
  }
}
export default setComponentTemplate(layout, CalendarWidget);
