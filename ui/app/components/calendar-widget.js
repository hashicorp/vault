import Component from '@glimmer/component';
import layout from '../templates/components/calendar-widget';
import { setComponentTemplate } from '@ember/component';
import { format } from 'date-fns';
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
  currentYear = this.currentDate.getFullYear(); // integer
  currentMonth = parseInt(format(this.currentDate, 'M')); // integer

  @tracked showCalendar = false;

  @tracked displayYear = this.currentYear; // init to currentYear and then changes as a user clicks on the chevrons
  @tracked disablePastYear = this.isObsoleteYear(); // if obsolete year, disable left chevron
  @tracked disableFutureYear = this.isCurrentYear(); // if current year, disable right chevron
  @tracked preselectRangeOfMonths = null; // the number of months selected by the quickSelect options (e.g. last month, last 3 months, etc.)
  @tracked allMonthsNodeList = [];
  @tracked areMonthsSelected = false;
  @tracked mouseClickCount = 0;
  @tracked clickRange = []; // the range of months by individually selecting them
  @tracked startDate = this.currentDate;
  @tracked endDate; // ARG TODO: For now, until you return the data from the parent
  @tracked firstClick;
  @tracked secondClick;

  // HELPER FUNCTIONS (alphabetically) //

  addClass(element, classString) {
    element.classList.add(classString);
  }

  isCurrentYear() {
    return this.currentYear === this.displayYear;
  }

  isObsoleteYear() {
    return this.displayYear === this.currentYear - 4; // won't display more than 5 years ago
  }

  // calculateEndDate() {
  //   let date = new Date(); // need to modify this variable so we're defining it here and not globally
  //   date.setMonth(date.getMonth() - 1);
  //   return format(date, 'MM-yyyy');
  // }

  // calculateLastMonth() {
  //   return sub(this.currentDate, { months: 1 });
  // }

  // calculateStartDate(quickSelectNumber) {
  //   let date = new Date(); // need to modify this variable so we're defining it here and not globally
  //   // will never query the same month that you are in, always subtract a month to get the N months prior
  //   date.setMonth(date.getMonth() - (quickSelectNumber ? quickSelectNumber : 11) - 1); // defaults to one year selected (11 months) if no quickSelectNumber
  //   return format(date, 'MM-yyyy');
  // }

  // idToDateObject(id) {
  //   let reverse = id
  //     .split('-')
  //     .reverse()
  //     .join(', ');
  //   return new Date(reverse);
  // }

  // dateObjectToHandleQueryParam(dateObject) {
  //   return format(dateObject, 'MM-yyyy');
  // }

  removeClass(element, classString) {
    element.classList.remove(classString);
  }

  // ACTIONS //

  @action
  addYear() {
    this.displayYear = this.displayYear + 1;
    this.disableMonths();
    this.disableFutureYear = this.isCurrentYear();
    this.disablePastYear = this.isObsoleteYear();
  }

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

  // action to parent Dashboard
  @action
  selectEndMonth(month, year, e) {
    // select month
    this.addClass(e.target, 'is-selected');
    // when ready send to handleQuery
    let endMonthSelected = e.target.id;
    this.args.handleQuery(endMonthSelected); // ARG TODO might need to change format, you have options
    this.endDate = endMonthSelected; // ARG TODO will likely have to modify?
    this.toggleShowCalendar();
  }

  @action
  subYear() {
    this.displayYear = this.displayYear - 1;
    this.disableMonths();
    this.disableFutureYear = this.isCurrentYear();
    this.disablePastYear = this.isObsoleteYear();
  }

  @action
  toggleShowCalendar() {
    this.showCalendar = !this.showCalendar;
  }
}
export default setComponentTemplate(layout, CalendarWidget);
