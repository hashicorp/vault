import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { isSameMonth, isAfter } from 'date-fns';

export default class History extends Component {
  // TODO CMB alphabetize and delete unused vars (particularly @tracked)
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

  chartLegend = [
    { key: 'entity_clients', label: 'entity clients' },
    { key: 'non_entity_clients', label: 'non-entity clients' },
  ];

  // needed for startTime modal picker
  months = Array.from({ length: 12 }, (item, i) => {
    return new Date(0, i).toLocaleString('en-US', { month: 'long' });
  });
  years = Array.from({ length: 5 }, (item, i) => {
    return new Date().getFullYear() - i;
  });

  @service store;

  @tracked queriedActivityResponse = null;
  @tracked barChartSelection = false;
  @tracked isEditStartMonthOpen = false;
  @tracked responseRangeDiffMessage = null;
  @tracked startTimeRequested = null;
  @tracked startTimeFromResponse = this.args.model.startTimeFromLicense; // ex: ['2021', 3] is April 2021 (0 indexed)
  @tracked endTimeFromResponse = this.args.model.endTimeFromResponse;
  @tracked startMonth = null;
  @tracked startYear = null;
  @tracked selectedNamespace = null;
  @tracked noActivityDate = '';
  @tracked namespaceArray = this.getActivityResponse.byNamespace.map((namespace) => {
    return { name: namespace['label'], id: namespace['label'] };
  });
  @tracked firstUpgradeVersion = this.args.model.versionHistory[0].id || null; // return 1.9.0 or earliest upgrade post 1.9.0
  @tracked upgradeDate = this.args.model.versionHistory[0].timestampInstalled || null; // returns RFC3339 timestamp

  // on init API response uses license start_date, getter updates when user queries dates
  get getActivityResponse() {
    return this.queriedActivityResponse || this.args.model.activity;
  }

  get startTimeDisplay() {
    if (!this.startTimeFromResponse) {
      // otherwise will return date of new Date(null)
      return null;
    }
    let month = this.startTimeFromResponse[1];
    let year = this.startTimeFromResponse[0];
    return `${this.arrayOfMonths[month]} ${year}`;
  }

  get endTimeDisplay() {
    if (!this.endTimeFromResponse) {
      // otherwise will return date of new Date(null)
      return null;
    }
    let month = this.endTimeFromResponse[1];
    let year = this.endTimeFromResponse[0];
    return `${this.arrayOfMonths[month]} ${year}`;
  }

  get isDateRange() {
    return !isSameMonth(
      new Date(this.getActivityResponse.startTime),
      new Date(this.getActivityResponse.endTime)
    );
  }

  // top level TOTAL client counts for given date range
  get totalUsageCounts() {
    return this.selectedNamespace
      ? this.filterByNamespace(this.selectedNamespace)
      : this.getActivityResponse.total;
  }

  // total client data for horizontal bar chart in attribution component
  get totalClientsData() {
    if (this.selectedNamespace) {
      let filteredNamespace = this.filterByNamespace(this.selectedNamespace);
      return filteredNamespace.mounts ? this.filterByNamespace(this.selectedNamespace).mounts : null;
    } else {
      return this.getActivityResponse?.byNamespace;
    }
  }

  get responseTimestamp() {
    return this.getActivityResponse.responseTimestamp;
  }

  get countsIncludeOlderData() {
    let firstUpgrade = this.args.model.versionHistory[0];
    if (!firstUpgrade) {
      return false;
    }
    let versionDate = new Date(firstUpgrade.timestampInstalled);
    let startTimeFromResponseAsDateObject = new Date(
      Number(this.startTimeFromResponse[0]),
      this.startTimeFromResponse[1]
    );
    // compare against this startTimeFromResponse to show message or not.
    return isAfter(versionDate, startTimeFromResponseAsDateObject) ? versionDate : false;
  }

  // ACTIONS
  @action
  async handleClientActivityQuery(month, year, dateType) {
    if (dateType === 'cancel') {
      return;
    }
    // clicked "Current Billing period" in the calendar widget
    if (dateType === 'reset') {
      this.startTimeRequested = this.args.model.startTimeFromLicense;
      this.endTimeRequested = null;
    }
    // clicked "Edit" Billing start month in Dashboard which opens a modal.
    if (dateType === 'startTime') {
      let monthIndex = this.arrayOfMonths.indexOf(month);
      this.startTimeRequested = [year.toString(), monthIndex]; // ['2021', 0] (e.g. January 2021) // TODO CHANGE TO ARRAY
      this.endTimeRequested = null;
    }
    // clicked "Custom End Month" from the calendar-widget
    if (dateType === 'endTime') {
      // use the currently selected startTime for your startTimeRequested.
      this.startTimeRequested = this.startTimeFromResponse;
      this.endTimeRequested = [year.toString(), month]; // endTime comes in as a number/index whereas startTime comes in as a month name. Hence the difference between monthIndex and month.
    }

    try {
      let response = await this.store.queryRecord('clients/activity', {
        start_time: this.startTimeRequested,
        end_time: this.endTimeRequested,
      });
      if (response.id === 'no-data') {
        // empty response is the only time we want to update the displayed date with the requested time
        this.startTimeFromResponse = this.startTimeRequested;
        this.noActivityDate = this.startTimeDisplay;
      } else {
        // note: this.startTimeDisplay (getter) is updated by this.startTimeFromResponse
        this.startTimeFromResponse = response.formattedStartTime;
        this.endTimeFromResponse = response.formattedEndTime;
      }
      // compare if the response startTime comes after the requested startTime. If true throw a warning.
      // only display if they selected a startTime
      if (
        dateType === 'startTime' &&
        isAfter(
          new Date(this.getActivityResponse.startTime),
          new Date(this.startTimeRequested[0], this.startTimeRequested[1])
        )
      ) {
        this.responseRangeDiffMessage = `You requested data from ${month} ${year}. We only have data from ${this.startTimeDisplay}, and that is what is being shown here.`;
      } else {
        this.responseRangeDiffMessage = null;
      }
      this.queriedActivityResponse = response;
    } catch (e) {
      // ARG TODO handle error
    }
  }

  @action
  handleCurrentBillingPeriod() {
    this.handleClientActivityQuery(0, 0, 'reset');
  }

  @action
  selectNamespace([value]) {
    // value comes in as [namespace0]
    this.selectedNamespace = value;
  }

  @action
  selectStartMonth(month) {
    this.startMonth = month;
  }

  @action
  selectStartYear(year) {
    this.startYear = year;
  }

  // HELPERS
  filterByNamespace(namespace) {
    return this.getActivityResponse.byNamespace.find((ns) => ns.label === namespace);
  }
}
