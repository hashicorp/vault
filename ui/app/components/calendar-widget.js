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
 * @param {array} arrayOfMonths - An array of all the months that the calendar widget iterates through.
 * @param {string} endTimeDisplay - The formatted display value of the endTime. Ex: January 2022.
 * @param {array} endTimeFromResponse - The value returned on the counters/activity endpoint, which shows the true endTime not the selected one, which can be different. Ex: ['2022', 0]
 * @param {function} handleClientActivityQuery - a function passed from parent. This component sends the month and year to the parent via this method which then calculates the new data.
 * @param {function} handleCurrentBillingPeriod - a function passed from parent. This component makes the parent aware that the user selected Current billing period and it handles resetting the data.
 * @param {string} startTimeDisplay - The formatted display value of the endTime. Ex: January 2022. This component is only responsible for modifying the endTime which is sends to the parent to make the network request.
 * />
 *
 * ```
 */
class CalendarWidget extends Component {
  currentDate = new Date();
  currentYear = this.currentDate.getFullYear(); // integer
  currentMonth = parseInt(this.currentDate.getMonth()); // integer and zero index

  @tracked allMonthsNodeList = [];
  @tracked displayYear = this.currentYear; // init to currentYear and then changes as a user clicks on the chevrons
  @tracked disablePastYear = this.isObsoleteYear(); // if obsolete year, disable left chevron
  @tracked disableFutureYear = this.isCurrentYear(); // if current year, disable right chevron
  @tracked showCalendar = false;
  @tracked tooltipTarget = null;
  @tracked tooltipText = null;

  // HELPER FUNCTIONS (alphabetically) //
  addClass(element, classString) {
    element.classList.add(classString);
  }

  isCurrentYear() {
    return this.currentYear === this.displayYear;
  }

  isObsoleteYear() {
    // do not allow them to choose a year before the this.args.startTimeDisplay
    let startYear = this.args.startTimeDisplay.split(' ')[1];
    return this.displayYear.toString() === startYear; // if on startYear then don't let them click back to the year prior
  }

  removeClass(element, classString) {
    element.classList.remove(classString);
  }

  // ACTIONS (alphabetically) //
  @action
  addTooltip() {
    if (this.isObsoleteYear()) {
      let previousYear = Number(this.displayYear) - 1;
      this.tooltipText = `${previousYear} is unavailable because it is before your billing start month. Change your billing start month to a date in ${previousYear} to see data for this year.`; // set tooltip text
      this.tooltipTarget = '#previous-year';
    }
  }

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
        let startMonthIndex = this.args.arrayOfMonths.indexOf(startMonth);
        // then add readOnly class to any month less than the startMonth index.
        if (startMonthIndex > elementMonthId) {
          e.classList.add('is-readOnly');
        }
      }
      // Compare values so the user cannot select an endTime after the endTime returned from counters/activity response on page load.
      let yearEndTimeFromResponse = Number(this.args.endTimeFromResponse[0]);
      let endMonth = this.args.endTimeFromResponse[1];
      if (this.displayYear === yearEndTimeFromResponse) {
        // add readOnly class to any month that is older (higher) than the endMonth index. (e.g. if nov is the endMonth of the endTimeDisplay, then 11 and 12 should not be displayed 10 < 11 and 10 < 12.)
        if (endMonth < elementMonthId) {
          e.classList.add('is-readOnly');
        }
      }
      // if the year display higher than the endTime e.g. you're looking at 2022 and the returned endTime is 2021, all months should be disabled.
      if (this.displayYear > yearEndTimeFromResponse) {
        // all months should be disabled.
        e.classList.add('is-readOnly');
      }
    });
  }

  @action removeTooltip() {
    this.tooltipTarget = null;
  }

  @action
  selectCurrentBillingPeriod(D) {
    this.args.handleCurrentBillingPeriod(); // resets the billing startTime and endTime to what it is on init via the parent.
    this.showCalendar = false;
    D.actions.close(); // close the dropdown.
  }
  @action
  selectEndMonth(month, year, D) {
    this.toggleShowCalendar();
    this.args.handleClientActivityQuery(month, year, 'endTime');
    D.actions.close(); // close the dropdown.
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
