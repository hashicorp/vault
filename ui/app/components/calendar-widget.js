import Component from '@glimmer/component';
import layout from '../templates/components/calendar-widget';
import { setComponentTemplate } from '@ember/component';
import { format, formatRFC3339, parseISO } from 'date-fns';
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
 * @param {object} startDate - The start date is calculated from the parent. This component is only responsible for modifying the end Date. ANd effecting single month.
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
  @tracked endDateDisplay;
  @tracked showCalendar = false;
  @tracked showSingleMonth = false;

  constructor() {
    super(...arguments);
    let endDate = parseISO(this.args.endDate);
    this.endDateDisplay = format(endDate, 'MMMM yyyy');
  }

  // HELPER FUNCTIONS (alphabetically) //
  addClass(element, classString) {
    element.classList.add(classString);
  }

  isCurrentYear() {
    return this.currentYear === this.displayYear;
  }

  isObsoleteYear() {
    // do not allow them to choose a end year before the this.args.startDate
    let startYear = this.args.startDate.split(' ')[1];
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
      if (this.displayYear.toString() === this.args.startDate.split(' ')[1]) {
        // if they are on the view where the start year equals the display year, check which months should not show.
        let startMonth = this.args.startDate.split(' ')[0]; // returns month name e.g. January
        // return the index of the startMonth
        let startMonthIndex = this.arrayOfMonths.indexOf(startMonth);
        // then add readOnly class to any month less than the startMonth index.
        if (startMonthIndex > elementMonthId) {
          e.classList.add('is-readOnly');
        }
      }
    });
  }

  @action
  selectCurrentBillingPeriod() {
    // ARG TOOD send to dashboard the select current billing period. The parent may know this it's just a boolean.
    // Turn the calendars off if they are showing.
    this.showCalendar = false;
    this.showSingleMonth = false;
  }
  @action
  selectEndMonth(month, year, element) {
    // select month
    this.addClass(element.target, 'is-selected');
    // API requires start and end time as EPOCH or RFC3339 timestamp
    let endMonthSelected = formatRFC3339(new Date(year, month));
    this.args.handleEndMonth(endMonthSelected);
    let monthName = this.arrayOfMonths.find((element, index) => {
      if (index === month) {
        return element;
      }
    });
    this.endDateDisplay = `${monthName} ${year}`;
    this.toggleShowCalendar();
  }

  @action
  selectSingleMonth(month, year, element) {
    // select month
    this.addClass(element.target, 'is-selected');
    // API requires start and end time as EPOCH or RFC3339 timestamp
    // let singleMonthSelected = formatRFC3339(new Date(year, month));
    // ARG TODO unsure on what we're going to do yet with date. Depends on API.
    // this.args.handleEndMonth(singleMonthSelected);
    // let monthName = this.arrayOfMonths.find((e, index) => {
    //   if (index === month) {
    //     return e;
    //   }
    // });
    // this.endDateDisplay = `${monthName} ${year}`;
    this.toggleSingleMonth();
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
    this.showSingleMonth = false;
  }

  @action
  toggleSingleMonth() {
    this.showSingleMonth = !this.showSingleMonth;
    this.showCalendar = false;
  }
}
export default setComponentTemplate(layout, CalendarWidget);
