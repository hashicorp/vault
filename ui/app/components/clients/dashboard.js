import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { format, formatRFC3339, isSameMonth, parseISO } from 'date-fns';
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
    { key: 'entity_clients', label: 'unique entities' },
    { key: 'non_entity_clients', label: 'non-entity tokens' },
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
  @tracked requestedStartTime = null;
  @tracked startTime = this.args.model.startTime;
  @tracked endTime = this.args.model.endTime;
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
    if (!this.startTime) {
      // otherwise will return date of new Date(null)
      return null;
    }
    // unable to use date-fns format here because of the local timestamp attached to the date when the user selects a new start time from the modal
    let formattedAsDate = new Date(this.startTime);
    let test = this.utcDate(formattedAsDate).toISOString();
    let year = test.substring(0, 4);
    let month = Number(test.substring(5, 7)) - 1;
    return `${this.arrayOfMonths[month]} ${year}`;
  }

  get endTimeDisplay() {
    if (!this.endTime) {
      // otherwise will return date of new Date(null)
      return null;
    }
    let formattedAsDate = new Date(this.endTime);
    return format(formattedAsDate, 'MMMM yyyy');
  }

  get isDateRange() {
    return !isSameMonth(new Date(this.startTime), new Date(this.endTime));
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

  // ACTIONS

  @action
  async handleClientActivityQuery(month, year, dateType) {
    if (dateType === 'cancel') {
      return;
    }
    // dateType is either startTime or endTime
    let monthIndex = this.arrayOfMonths.indexOf(month);
    if (dateType === 'startTime') {
      let utcDateObject = this.utcDate(new Date(year, monthIndex));
      this.requestedStartTime = utcDateObject.toISOString();
      this.endTime = null;
    }
    if (dateType === 'endTime') {
      // this month comes in as an index
      // ARG TODO remove timezone!!!! user helper
      this.endTime = formatRFC3339(new Date(year, month));
    }
    try {
      let response = await this.adapter.queryClientActivity(this.requestedStartTime, this.endTime);
      if (!response) {
        // this.endTime will be null and use this to show EmptyState message on the template.
        return;
      }
      // compare year and month of the RFC33393 times from the start time set to what is returned on the response
      if (this.requestedStartTime.substring(0, 7) !== response.data.start_time.substring(0, 7)) {
        // reset the starTime
        // show warning that the response is different than what was set
        this.responseRangeDiffMessage = `You requested data from ${month} ${year}. We only have data from ${this.startTimeDisplay}, and that is what is being shown here.`;
      } else {
        this.responseRangeDiffMessage = null;
      }
      // change the startTime to the time returned on the response
      this.startTime = response.data.start_time;

      // resets the endTime to what is returned on the response
      this.endTime = response.data.end_time;
      this.updateData = response;
      return response;
      // ARG TODO this is the response you need to use to repopulate the chart data
    } catch (e) {
      console.log('ERROR', e);
      // ARG TODO handle error
    }
  }

  @action
  handleCurrentBillingPeriod() {
    let parsed = format(parseISO(this.args.model.startTime), 'MMMM yyyy'); // this.startTime is a tracked property and changes, thus use the startTime first passed into the component from the model.
    let month = parsed.split(' ')[0];
    let year = parsed.split(' ')[1];
    this.handleClientActivityQuery(month, year, 'startTime');
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
