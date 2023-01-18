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
  @tracked showCalendar = false;
  @tracked tooltipTarget = null;
  @tracked tooltipText = null;

  get selectedMonthId() {
    if (!this.args.endTimeFromResponse) return '';
    const [year, monthIndex] = this.args.endTimeFromResponse;
    return `${monthIndex}-${year}`;
  }
  get disableFutureYear() {
    return this.displayYear === this.currentYear;
  }
  get disablePastYear() {
    const startYear = parseInt(this.args.startTimeDisplay.split(' ')[1]);
    return this.displayYear === startYear; // if on startYear then don't let them click back to the year prior
  }
  get widgetMonths() {
    const displayYear = this.displayYear;
    const currentYear = this.currentYear;
    const currentMonthIdx = this.currentMonth;
    const [startMonth, startYear] = this.args.startTimeDisplay.split(' ');
    const startMonthIdx = this.args.arrayOfMonths.indexOf(startMonth);
    return this.args.arrayOfMonths.map((month, idx) => {
      const monthId = `${idx}-${displayYear}`;
      let readonly = false;

      // if widget is showing billing start year, disable if month is before start month
      if (parseInt(startYear) === displayYear && idx < startMonthIdx) {
        readonly = true;
      }

      // if widget showing current year, disable if month is current or later
      if (displayYear === currentYear && idx >= currentMonthIdx) {
        readonly = true;
      }
      return {
        id: monthId,
        month,
        readonly,
        current: monthId === `${currentMonthIdx}-${currentYear}`,
      };
    });
  }

  // HELPER FUNCTIONS (alphabetically) //
  addClass(element, classString) {
    element.classList.add(classString);
  }

  removeClass(element, classString) {
    element.classList.remove(classString);
  }

  resetDisplayYear() {
    let setYear = this.currentYear;
    if (this.args.endTimeDisplay) {
      try {
        const year = this.args.endTimeDisplay.split(' ')[1];
        setYear = parseInt(year);
      } catch (e) {
        console.debug('Error resetting display year', e); // eslint-disable-line
      }
    }
    this.displayYear = setYear;
  }

  // ACTIONS (alphabetically) //
  @action
  addTooltip() {
    if (this.disablePastYear) {
      const previousYear = Number(this.displayYear) - 1;
      this.tooltipText = `${previousYear} is unavailable because it is before your billing start month. Change your billing start month to a date in ${previousYear} to see data for this year.`; // set tooltip text
      this.tooltipTarget = '#previous-year';
    }
  }

  @action
  addYear() {
    this.displayYear = this.displayYear + 1;
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
  selectEndMonth(monthId, D) {
    const [monthIdx, year] = monthId.split('-');
    this.toggleShowCalendar();
    this.args.handleClientActivityQuery(parseInt(monthIdx), parseInt(year), 'endTime');
    D.actions.close(); // close the dropdown.
  }

  @action
  subYear() {
    this.displayYear = this.displayYear - 1;
  }

  @action
  toggleShowCalendar() {
    this.showCalendar = !this.showCalendar;
    this.resetDisplayYear();
  }
}
export default setComponentTemplate(layout, CalendarWidget);
