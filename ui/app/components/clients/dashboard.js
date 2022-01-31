import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { zonedTimeToUtc } from 'date-fns-tz'; // https://github.com/marnusw/date-fns-tz#zonedtimetoutc

export default class Dashboard extends Component {
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
  maxNamespaces = 10;
  chartLegend = [
    { key: 'entity_clients', label: 'entity clients' },
    { key: 'non_entity_clients', label: 'non-entity clients' },
  ];
  adapter = this.store.adapterFor('clients/activity');

  // needed for startTime modal picker
  months = Array.from({ length: 12 }, (item, i) => {
    return new Date(0, i).toLocaleString('en-US', { month: 'long' });
  });
  years = Array.from({ length: 5 }, (item, i) => {
    return new Date().getFullYear() - i;
  });

  @service store;

  @tracked barChartSelection = false;
  @tracked isEditStartMonthOpen = false;
  @tracked responseRangeDiffMessage = null;
  @tracked startTimeRequested = null;
  @tracked startTimeFromResponse = this.args.model.startTimeFromLicense; // ex: "3,2021"
  @tracked endTimeFromResponse = this.args.model.endTimeFromLicense;
  @tracked startMonth = null;
  @tracked startYear = null;
  @tracked selectedNamespace = null;
  // @tracked selectedNamespace = 'namespacelonglonglong4/'; // for testing namespace selection view

  // HELPER
  utcDate(dateObject) {
    // To remove the timezone of the local user (API returns and expects Zulu time/UTC) we need to use a method provided by date-fns-tz to return the UTC date
    let timeZone = Intl.DateTimeFormat().resolvedOptions().timeZone; // browser API method
    return zonedTimeToUtc(dateObject, timeZone);
  }

  get startTimeDisplay() {
    if (!this.startTimeFromResponse) {
      // otherwise will return date of new Date(null)
      return null;
    }
    let month = Number(this.startTimeFromResponse.split(',')[0]) - 1;
    let year = this.startTimeFromResponse.split(',')[1];
    return `${this.arrayOfMonths[month]} ${year}`;
  }

  get endTimeDisplay() {
    if (!this.endTimeFromResponse) {
      // otherwise will return date of new Date(null)
      return null;
    }
    let month = Number(this.endTimeFromResponse.split(',')[0]) - 1;
    let year = this.endTimeFromResponse.split(',')[1];
    return `${this.arrayOfMonths[month]} ${year}`;
  }

  get isDateRange() {
    if (!this.startTimeFromResponse || !this.endTimeFromResponse) {
      // need to check because the moment one of these tracked properties changes, this getter is fired of, even if the other tracked property hasn't been set yet
      return false;
    }
    return this.startTimeFromResponse.split(',')[0] !== this.endTimeFromResponse.split(',')[0];
  }

  // Determine if we have client count data based on the current tab
  get hasClientData() {
    if (this.args.tab === 'current') {
      // Show the current numbers as long as config is on
      return this.args.model.config?.enabled !== 'Off';
    }
    return this.args.model.activity && this.args.model.activity.total;
  }

  // top level TOTAL client counts from response for given date range
  get runningTotals() {
    if (!this.args.model.activity || !this.args.model.activity.total) {
      return null;
    }
    return this.args.model.activity.total;
  }

  // for horizontal bar chart in Attribution component
  get topTenNamespaces() {
    if (!this.args.model.activity || !this.args.model.activity.byNamespace) {
      return null;
    }
    return this.args.model.activity.byNamespace;
  }

  get responseTimestamp() {
    if (!this.args.model.activity || !this.args.model.activity.responseTimestamp) {
      return null;
    }
    return this.args.model.activity.responseTimestamp;
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
      this.startTimeRequested = `${monthIndex + 1},${year}`; // "1, 2021"
      this.endTimeRequested = null;
    }
    // clicked "Custom End Month" from the calendar-widget
    if (dateType === 'endTime') {
      // use the currently selected startTime for your startTimeRequested.
      this.startTimeRequested = this.startTimeFromResponse;
      this.endTimeRequested = `${month},${year}`; // "1, 2021"
    }

    try {
      let response = await this.store.queryRecord('clients/activity', {
        start_time: this.startTimeRequested,
        end_time: this.endTimeRequested,
      });
      if (!response) {
        // this.endTime will be null and use this to show EmptyState message on the template.
        return;
      }
      // serializer transforms response from rfc3339 to "3,2021"
      // ARG TODO move to serializer
      this.startTimeFromResponse = `${response.data.start_time.split('-')[1].replace(/^0+/, '')},${
        response.data.start_time.split('-')[0]
      }`;
      this.endTimeFromResponse = `${response.data.end_time.split('-')[1].replace(/^0+/, '')},${
        response.data.end_time.split('-')[0]
      }`;

      if (this.startTimeRequested !== this.startTimeFromResponse) {
        this.responseRangeDiffMessage = `You requested data from ${month} ${year}. We only have data from ${this.startTimeDisplay}, and that is what is being shown here.`;
      } else {
        this.responseRangeDiffMessage = null;
      }
      return response;
    } catch (e) {
      // ARG TODO handle error
    }
  }

  @action
  handleCurrentBillingPeriod() {
    this.handleClientActivityQuery(0, 0, 'reset');
  }

  @action
  selectNamespace(value) {
    // In case of search select component, value returned is an array
    if (Array.isArray(value)) {
      this.selectedNamespace = this.getNamespace(value[0]);
      this.barChartSelection = false;
    } else if (typeof value === 'object') {
      // While D3 bar selection returns an object
      this.selectedNamespace = this.getNamespace(value.label);
      this.barChartSelection = true;
    }
  }

  @action
  selectStartMonth(month) {
    this.startMonth = month;
  }

  @action
  selectStartYear(year) {
    this.startYear = year;
  }
}
