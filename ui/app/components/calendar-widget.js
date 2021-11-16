import Component from '@glimmer/component';
import layout from '../templates/components/calendar-widget';
import { setComponentTemplate } from '@ember/component';
import { format, sub, isWithinInterval, isBefore } from 'date-fns';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

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

class CalendarWidget extends Component {
  currentDate = new Date();
  currentYear = parseInt(format(this.currentDate, 'yyyy')); // integer
  currentMonth = parseInt(format(this.currentDate, 'M')); // integer

  @tracked displayYear = this.currentYear; // init to currentYear and then changes as a user clicks on the chevrons
  @tracked disablePastYear = this.isObsoleteYear(); // if obsolete year, disable left chevron
  @tracked disableFutureYear = this.isCurrentYear(); // if current year, disable right chevron
  @tracked preselectRangeOfMonths = null; // the number of months selected by the quickSelect options (e.g. last month, last 3 months, etc.)
  @tracked allMonthsNodeList = [];
  @tracked isClearAllMonths = false;
  @tracked areMonthsSelected = false;
  @tracked mouseClickCount = 0;
  @tracked clickRange = []; // the range of months by individually selecting them
  @tracked startDate; // the older time e.g. 10/2020
  @tracked endDate; // the newer time e.g. 11/2021
  @tracked firstClick;
  @tracked secondClick;

  constructor() {
    super(...arguments);
    this.startDate = this.calculateStartDate();
    this.endDate = this.calculateEndDate();
  }

  // HELPER FUNCTIONS (alphabetically)//

  addClass(element, classString) {
    element.classList.add(classString);
  }

  calculateEndDate() {
    let date = new Date(); // need to modify this variable so we're defining it here and not globally
    date.setMonth(date.getMonth() - 1);
    return format(date, 'MM-yyyy');
  }

  calculateLastMonth() {
    return sub(this.currentDate, { months: 1 });
  }

  calculateStartDate(quickSelectNumber) {
    let date = new Date(); // need to modify this variable so we're defining it here and not globally
    // will never query the same month that you are in, always subtract a month to get the N months prior
    date.setMonth(date.getMonth() - (quickSelectNumber ? quickSelectNumber : 11) - 1); // defaults to one year selected (11 months) if no quickSelectNumber
    return format(date, 'MM-yyyy');
  }

  idToDateObject(id) {
    let reverse = id
      .split('-')
      .reverse()
      .join(', ');
    return new Date(reverse);
  }

  dateObjectToHandleQueryParam(dateObject) {
    return format(dateObject, 'MM-yyyy');
  }

  checkIfMonthsSelected(mouseClickCount, id) {
    if (mouseClickCount === 1) {
      this.startDate = id; // this will equal a string number of the month selected and year e.g. '6-2020'
      this.endDate = id;
      this.firstClick = id;
    } else if (mouseClickCount === 2) {
      // reverse this '6-2020' to this '2020-6'
      this.secondClick = id;
      // this.idToDateObject(this.firstClick),
      // this.idToDateObject(this.secondClick)

      this.handleSelectRange(this.idToDateObject(this.firstClick), this.idToDateObject(this.secondClick));
      // select that range

      // if the second click is before the first click
      // reset the startDate and endDate to the first click
    }
    // then set the date on the query mm-yyyy
    this.args.handleQuery(this.startDate, this.endDate);
  }

  createRange(start, end) {
    return Array(end - start + 1)
      .fill()
      .map((_, idx) => start + idx);
  }

  deselectAllMonths() {
    this.allMonthsNodeList.forEach(element => {
      this.removeClass(element, 'is-selected');
    });
  }

  handleSelectRange(firstClickDateObject, secondClickDateObject) {
    let start, end;
    // figure out start date by which click is oldest
    let isFirstClickOlder = isBefore(firstClickDateObject, secondClickDateObject);
    if (isFirstClickOlder) {
      start = firstClickDateObject;
      end = secondClickDateObject;
    } else {
      start = secondClickDateObject;
      end = firstClickDateObject;
    }
    this.startDate = this.dateObjectToHandleQueryParam(start);
    this.endDate = this.dateObjectToHandleQueryParam(end);
    // set interval
    let interval = { start: start, end: end };
    // then use isWithinInterval to check an array our nodeList and their ids and see if they match isToDateObject
    this.allMonthsNodeList.forEach(e => {
      if (isWithinInterval(this.idToDateObject(e.id), interval)) {
        this.addClass(e, 'is-selected');
      }
    });

    // this.clickRange = this.createRange(startAndEndMonths[0], startAndEndMonths[1]).map(n => `month-${n}`);
    // this.clickRange.forEach(id => {
    //   this.allMonthsNodeList.forEach(e => {
    //     if (e.id === id) {
    //       this.addClass(e, 'is-selected');
    //     }
    //   });
    // });
  }

  isCurrentYear() {
    return this.currentYear === this.displayYear;
  }

  isObsoleteYear() {
    return this.displayYear === this.currentYear - 4; // won't display more than 5 years ago
  }

  removeClass(element, classString) {
    element.classList.remove(classString);
  }

  // ACTIONS //
  @action
  disableMonths() {
    this.allMonthsNodeList = document.querySelectorAll('.is-month-list');
    this.allMonthsNodeList.forEach(e => {
      // clear all is-readOnly classes and start over.
      this.removeClass(e, 'is-readOnly');
      let elementMonthId = parseInt(e.id.split('-')[0]);
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

    this.checkIfMonthsSelected(this.mouseClickCount, e.target.id);
  }

  reverseMonthNodeList() {
    let reverseMonthArray = [];
    this.allMonthsNodeList.forEach(e => {
      reverseMonthArray.unshift(e);
    });
    return reverseMonthArray;
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
