import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { format, formatRFC3339 } from 'date-fns';

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

  months = Array.from({ length: 12 }, (item, i) => {
    return new Date(0, i).toLocaleString('en-US', { month: 'long' });
  });
  years = Array.from({ length: 5 }, (item, i) => {
    return new Date().getFullYear() - i;
  });
  adapter = this.store.adapterFor('clients/new-init-activity');

  @service store;

  @tracked barChartSelection = false;
  @tracked isEditStartMonthOpen = false;
  @tracked startTime = this.args.model.startTime;
  @tracked endTime = this.args.model.endTime;
  @tracked startMonth = null;
  @tracked startYear = null;
  @tracked selectedNamespace = null;

  // filtering conditionals for HBS file
  // set initially based on response received from API & update based on filters
  @tracked isDateRange = true; // rename to isSingleMonth ?
  @tracked isAllNamespaces; // false when filtered down to auth methods (mounts)

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
    if (!this.args.model.newInitActivity || !this.args.model.newInitActivity.total) {
      return null;
    }
    return this.args.model.newInitActivity.total;
  }

  // for horizontal bar chart in Attribution component
  get topTenNamespaces() {
    if (!this.args.model.newInitActivity || !this.args.model.newInitActivity.byNamespace) {
      return null;
    }
    return this.args.model.newInitActivity.byNamespace;
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
      // TODO CMB set tracked property isDateRange depending on dates from response
      return response;
      // ARG TODO this is the response you need to use to repopulate the chart data
    } catch (e) {
      // ARG TODO handle error
    }
  }

  // ARG TODO this might be a carry over from history, will need to confirm
  @action
  resetData() {
    this.barChartSelection = false;
    this.selectedNamespace = null;
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
