import Component from '@glimmer/component';
import layout from '../templates/components/calendar-widget';
import { setComponentTemplate } from '@ember/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

/**
 * @module CalendarWidget
 * CalendarWidget components are used in the client counts metrics. It helps users understand the ranges they can select.
 *
 * @example
 * ```js
 * <CalendarWidget
 * @param {string} endTimeDisplay - The display value of the endTime. Ex: January 2022.
 * @param {string} endTimeAtInit - Because endTimeDisplay is a tracked property of the parent it changes. But we need to keep a value unchanged that records the endTime returned by the counters/activity endpoint. This is that value.
 * @param {string} startTimeDisplay - The display value of startTime that the parent manages. This component is only responsible for modifying the endTime which is sends to the parent to make the network request.
 * @param {function} handleClientActivityQuery - a function passed from parent. This component sends the month and year to the parent via this method which then calculates the new data.
 * @param {function} handleCurrentBillingPeriod - a function passed from parent. This component makes the parent aware that the user selected Current billing period and it handles resetting the data.
 * />
 *
 * ```
 */
class CalendarWidget extends Component {
  arrayOfMonths = [
    'January',
    'February',
    'March',
    'April',
    'May',
    'June',
    'July',
    'August',
    'September',
    'October',
    'November',
    'December',
  ];
  currentDate = new Date();
  currentYear = this.currentDate.getFullYear(); // integer
  currentMonth = parseInt(this.currentDate.getMonth()); // integer and zero index

  @tracked allMonthsNodeList = [];
  @tracked displayYear = this.currentYear; // init to currentYear and then changes as a user clicks on the chevrons
  @tracked disablePastYear = this.isObsoleteYear(); // if obsolete year, disable left chevron
  @tracked disableFutureYear = this.isCurrentYear(); // if current year, disable right chevron
  @tracked showCalendar = false;

  // HELPER FUNCTIONS (alphabetically) //
  addClass(element, classString) {
    element.classList.add(classString);
  }

  isCurrentYear() {
    return this.currentYear === this.displayYear;
  }

  isObsoleteYear() {
    // do not allow them to choose a end year before the this.args.startTimeDisplay
    let startYear = this.args.startTimeDisplay.split(' ')[1];
    return this.displayYear.toString() === startYear; // if on startYear then don't let them click back to the year prior
  }

  removeClass(element, classString) {
    element.classList.remove(classString);
  }

  // ACTIONS (alphabetically) //
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
    this.allMonthsNodeList.forEach((e) => {
      // clear all is-readOnly classes and start over.
      this.removeClass(e, 'is-readOnly');

      let elementMonthId = parseInt(e.id.split('-')[0]); // dependent on the shape of the element id
      // for current year
      if (this.currentMonth <= elementMonthId) {
        // only disable months when current year is selected
        if (this.isCurrentYear()) {
          e.classList.add('is-readOnly');
        }
      }
      // compare for startYear view
      if (this.displayYear.toString() === this.args.startTimeDisplay.split(' ')[1]) {
        // if they are on the view where the start year equals the display year, check which months should not show.
        let startMonth = this.args.startTimeDisplay.split(' ')[0]; // returns month name e.g. January
        // return the index of the startMonth
        let startMonthIndex = this.arrayOfMonths.indexOf(startMonth);
        // then add readOnly class to any month less than the startMonth index.
        if (startMonthIndex > elementMonthId) {
          e.classList.add('is-readOnly');
        }
      }
      // Compare values so the user cannot select an endTime after the endTime returned from counters/activity response on page load.
      // ARG TODO will need to test if no data is returned on page load.
      if (this.displayYear.toString() === this.args.endTimeAtInit.split(' ')[1]) {
        let endMonth = this.args.endTimeAtInit.split(' ')[0];
        let endMonthIndex = this.arrayOfMonths.indexOf(endMonth);
        // add readOnly class to any month that is older (higher) than the endMonth index. (e.g. if nov is the endMonth of the endTimeDisplay, then 11 and 12 should not be displayed 10 < 11 and 10 < 12.)
        if (endMonthIndex < elementMonthId) {
          e.classList.add('is-readOnly');
        }
      }
    });
  }

  @action
  selectCurrentBillingPeriod() {
    this.args.handleCurrentBillingPeriod(); // resets the billing startTime and endTime to what it is on init via the parent.
    this.showCalendar = false;
  }
  @action
  selectEndMonth(month, year, element) {
    this.addClass(element.target, 'is-selected');
    this.toggleShowCalendar();
    this.args.handleClientActivityQuery(month, year, 'endTime');
  }

  @action
  selectSingleMonth(month, year, element) {
    // select month
    this.addClass(element.target, 'is-selected');
    this.toggleSingleMonth();
    // ARG TODO similar to selectEndMonth
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
