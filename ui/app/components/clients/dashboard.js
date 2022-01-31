import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { format, formatRFC3339, isSameMonth } from 'date-fns';

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

  @service store;

  @tracked barChartSelection = false;
  @tracked isEditStartMonthOpen = false;
  @tracked startTime = this.args.model.startTime;
  @tracked endTime = this.args.model.endTime;
  @tracked startMonth = null;
  @tracked startYear = null;
  @tracked selectedNamespace = null;
  // @tracked selectedNamespace = 'namespace18anotherlong/'; // for testing namespace selection view with mirage

  get startTimeDisplay() {
    if (!this.startTime) {
      // otherwise will return date of new Date(null)
      return null;
    }
    let formattedAsDate = new Date(this.startTime); // on init it's formatted as a Date object, but when modified by modal it's formatted as RFC3339
    return format(formattedAsDate, 'MMMM yyyy');
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

  // top level TOTAL client counts from response for given date range
  get totalUsageCounts() {
    return this.selectedNamespace
      ? this.filterByNamespace(this.selectedNamespace)
      : this.args.model.activity?.total;
  }

  // by namespace client count data for date range
  get byNamespaceActivity() {
    return this.args.model.activity?.byNamespace || null;
  }

  // for horizontal bar chart in attribution component
  get topTenChartData() {
    if (this.selectedNamespace) {
      let filteredNamespace = this.filterByNamespace(this.selectedNamespace);
      return filteredNamespace.mounts
        ? this.filterByNamespace(this.selectedNamespace).mounts.slice(0, 10)
        : null;
    } else {
      return this.byNamespaceActivity.slice(0, 10);
    }
  }

  get responseTimestamp() {
    return this.args.model.activity?.responseTimestamp;
  }

  @action
  async handleClientActivityQuery(month, year, dateType) {
    if (dateType === 'cancel') {
      return;
    }
    // dateType is either startTime or endTime
    let monthIndex = this.arrayOfMonths.indexOf(month);
    if (dateType === 'startTime') {
      this.startTime = formatRFC3339(new Date(year, monthIndex));
      this.endTime = null;
    }
    if (dateType === 'endTime') {
      // this month comes in as an index
      this.endTime = formatRFC3339(new Date(year, month));
    }
    try {
      let response = await this.adapter.queryClientActivity(this.startTime, this.endTime);
      // resets the endTime to what is returned on the response
      this.endTime = response.data.end_time;

      return response;
      // ARG TODO this is the response you need to use to repopulate the chart data
    } catch (e) {
      // ARG TODO handle error
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

  // HELPERS
  filterByNamespace(namespace) {
    return this.byNamespaceActivity.find((ns) => ns.label === namespace);
  }
}
