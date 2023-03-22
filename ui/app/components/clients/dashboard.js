/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { isAfter, isBefore, isSameMonth, format } from 'date-fns';
import getStorage from 'vault/lib/token-storage';
import { parseAPITimestamp } from 'core/utils/date-formatters';
// my sincere apologies to the next dev who has to refactor/debug this (⇀‸↼‶)
export default class Dashboard extends Component {
  @service store;
  @service version;

  chartLegend = [
    { key: 'entity_clients', label: 'entity clients' },
    { key: 'non_entity_clients', label: 'non-entity clients' },
  ];

  // RESPONSE
  @tracked startMonthTimestamp; // when user queries, updates to first month object of response
  @tracked endMonthTimestamp; // when user queries, updates to last month object of response
  @tracked queriedActivityResponse = null;
  // track params sent to /activity request
  @tracked activityQueryParams = {
    start: {}, // updates when user edits billing start month
    end: {}, // updates when user queries end dates via calendar widget
  };

  // SEARCH SELECT FILTERS
  get namespaceArray() {
    return this.getActivityResponse.byNamespace
      ? this.getActivityResponse.byNamespace.map((namespace) => ({
          name: namespace.label,
          id: namespace.label,
        }))
      : [];
  }
  @tracked selectedNamespace = null;
  @tracked selectedAuthMethod = null;
  @tracked authMethodOptions = [];

  // TEMPLATE VIEW
  @tracked noActivityData;
  @tracked showBillingStartModal = false;
  @tracked isLoadingQuery = false;
  @tracked errorObject = null;

  constructor() {
    super(...arguments);
    this.startMonthTimestamp = this.args.model.licenseStartTimestamp;
    this.endMonthTimestamp = this.args.model.currentDate;
    this.activityQueryParams.start.timestamp = this.args.model.licenseStartTimestamp;
    this.activityQueryParams.end.timestamp = this.args.model.currentDate;
    this.noActivityData = this.args.model.activity.id === 'no-data' ? true : false;
  }

  // returns text for empty state message if noActivityData
  get dateRangeMessage() {
    if (!this.startMonthTimestamp && !this.endMonthTimestamp) return null;
    const endMonth = isSameMonth(
      parseAPITimestamp(this.startMonthTimestamp),
      parseAPITimestamp(this.endMonthTimestamp)
    )
      ? ''
      : ` to ${parseAPITimestamp(this.endMonthTimestamp, 'MMMM yyyy')}`;
    // completes the message 'No data received from { dateRangeMessage }'
    return `from ${parseAPITimestamp(this.startMonthTimestamp, 'MMMM yyyy')}` + endMonth;
  }

  get versionText() {
    return this.version.isEnterprise
      ? {
          label: 'Billing start month',
          description:
            'This date comes from your license, and defines when client counting starts. Without this starting point, the data shown is not reliable.',
          title: 'No billing start date found',
          message:
            'In order to get the most from this data, please enter your billing period start month. This will ensure that the resulting data is accurate.',
        }
      : {
          label: 'Client counting start date',
          description:
            'This date is when client counting starts. Without this starting point, the data shown is not reliable.',
          title: 'No start date found',
          message:
            'In order to get the most from this data, please enter a start month above. Vault will calculate new clients starting from that month.',
        };
  }

  get isDateRange() {
    return !isSameMonth(
      parseAPITimestamp(this.getActivityResponse.startTime),
      parseAPITimestamp(this.getActivityResponse.endTime)
    );
  }

  get isCurrentMonth() {
    return (
      isSameMonth(
        parseAPITimestamp(this.getActivityResponse.startTime),
        parseAPITimestamp(this.args.model.currentDate)
      ) &&
      isSameMonth(
        parseAPITimestamp(this.getActivityResponse.endTime),
        parseAPITimestamp(this.args.model.currentDate)
      )
    );
  }

  get startTimeDiscrepancy() {
    // show banner if startTime returned from activity log (response) is after the queried startTime
    const activityStartDateObject = parseAPITimestamp(this.getActivityResponse.startTime);
    const queryStartDateObject = parseAPITimestamp(this.startMonthTimestamp);
    let message = 'You requested data from';
    if (this.startMonthTimestamp === this.args.model.licenseStartTimestamp && this.version.isEnterprise) {
      // on init, date is automatically pulled from license start date and user hasn't queried anything yet
      message = 'Your license start date is';
    }
    if (
      isAfter(activityStartDateObject, queryStartDateObject) &&
      !isSameMonth(activityStartDateObject, queryStartDateObject)
    ) {
      return `${message} ${parseAPITimestamp(this.startMonthTimestamp, 'MMMM yyyy')}. 
        We only have data from ${parseAPITimestamp(this.getActivityResponse.startTime, 'MMMM yyyy')}, 
        and that is what is being shown here.`;
    } else {
      return null;
    }
  }

  get upgradeDuringActivity() {
    const versionHistory = this.args.model.versionHistory;
    if (!versionHistory || versionHistory.length === 0) {
      return null;
    }

    // filter for upgrade data of noteworthy upgrades (1.9 and/or 1.10)
    const upgradeVersionHistory = versionHistory.filter(
      ({ version }) => version.match('1.9') || version.match('1.10')
    );
    if (!upgradeVersionHistory || upgradeVersionHistory.length === 0) {
      return null;
    }

    const activityStart = parseAPITimestamp(this.getActivityResponse.startTime);
    const activityEnd = parseAPITimestamp(this.getActivityResponse.endTime);
    // filter and return all upgrades that happened within date range of queried activity
    const upgradesWithinData = upgradeVersionHistory.filter(({ timestampInstalled }) => {
      const upgradeDate = parseAPITimestamp(timestampInstalled);
      return isAfter(upgradeDate, activityStart) && isBefore(upgradeDate, activityEnd);
    });
    return upgradesWithinData.length === 0 ? null : upgradesWithinData;
  }

  get upgradeVersionAndDate() {
    if (!this.upgradeDuringActivity) return null;

    if (this.upgradeDuringActivity.length === 2) {
      const [firstUpgrade, secondUpgrade] = this.upgradeDuringActivity;
      const firstDate = parseAPITimestamp(firstUpgrade.timestampInstalled, 'MMM d, yyyy');
      const secondDate = parseAPITimestamp(secondUpgrade.timestampInstalled, 'MMM d, yyyy');
      return `Vault was upgraded to ${firstUpgrade.version} (${firstDate}) and ${secondUpgrade.version} (${secondDate}) during this time range.`;
    } else {
      const [upgrade] = this.upgradeDuringActivity;
      return `Vault was upgraded to ${upgrade.version} on ${parseAPITimestamp(
        upgrade.timestampInstalled,
        'MMM d, yyyy'
      )}.`;
    }
  }

  get upgradeExplanation() {
    if (!this.upgradeDuringActivity) return null;
    if (this.upgradeDuringActivity.length === 1) {
      const version = this.upgradeDuringActivity[0].version;
      if (version.match('1.9')) {
        return ' How we count clients changed in 1.9, so keep that in mind when looking at the data.';
      }
      if (version.match('1.10')) {
        return ' We added monthly breakdowns and mount level attribution starting in 1.10, so keep that in mind when looking at the data.';
      }
    }
    // return combined explanation if spans multiple upgrades
    return ' How we count clients changed in 1.9 and we added monthly breakdowns and mount level attribution starting in 1.10. Keep this in mind when looking at the data.';
  }

  get formattedStartDate() {
    if (!this.startMonthTimestamp) return null;
    return parseAPITimestamp(this.startMonthTimestamp, 'MMMM yyyy');
  }

  // GETTERS FOR RESPONSE & DATA

  // on init API response uses license start_date, getter updates when user queries dates
  get getActivityResponse() {
    return this.queriedActivityResponse || this.args.model.activity;
  }

  get byMonthActivityData() {
    if (this.selectedNamespace) {
      return this.filteredActivityByMonth;
    } else {
      return this.getActivityResponse?.byMonth;
    }
  }

  get hasAttributionData() {
    if (this.selectedAuthMethod) return false;
    if (this.selectedNamespace) {
      return this.authMethodOptions.length > 0;
    }
    return !!this.totalClientAttribution && this.totalUsageCounts && this.totalUsageCounts.clients !== 0;
  }

  // (object) top level TOTAL client counts for given date range
  get totalUsageCounts() {
    return this.selectedNamespace ? this.filteredActivityByNamespace : this.getActivityResponse.total;
  }

  // (object) single month new client data with total counts + array of namespace breakdown
  get newClientCounts() {
    return this.isDateRange ? null : this.byMonthActivityData[0]?.new_clients;
  }

  // total client data for horizontal bar chart in attribution component
  get totalClientAttribution() {
    if (this.selectedNamespace) {
      return this.filteredActivityByNamespace?.mounts || null;
    } else {
      return this.getActivityResponse?.byNamespace || null;
    }
  }

  // new client data for horizontal bar chart
  get newClientAttribution() {
    // new client attribution only available in a single, historical month (not a date range or current month)
    if (this.isDateRange || this.isCurrentMonth) return null;

    if (this.selectedNamespace) {
      return this.newClientCounts?.mounts || null;
    } else {
      return this.newClientCounts?.namespaces || null;
    }
  }

  get responseTimestamp() {
    return this.getActivityResponse.responseTimestamp;
  }

  // FILTERS
  get filteredActivityByNamespace() {
    const namespace = this.selectedNamespace;
    const auth = this.selectedAuthMethod;
    if (!namespace && !auth) {
      return this.getActivityResponse;
    }
    if (!auth) {
      return this.getActivityResponse.byNamespace.find((ns) => ns.label === namespace);
    }
    return this.getActivityResponse.byNamespace
      .find((ns) => ns.label === namespace)
      .mounts?.find((mount) => mount.label === auth);
  }

  get filteredActivityByMonth() {
    const namespace = this.selectedNamespace;
    const auth = this.selectedAuthMethod;
    if (!namespace && !auth) {
      return this.getActivityResponse?.byMonth;
    }
    const namespaceData = this.getActivityResponse?.byMonth
      .map((m) => m.namespaces_by_key[namespace])
      .filter((d) => d !== undefined);
    if (!auth) {
      return namespaceData.length === 0 ? null : namespaceData;
    }
    const mountData = namespaceData
      .map((namespace) => namespace.mounts_by_key[auth])
      .filter((d) => d !== undefined);
    return mountData.length === 0 ? null : mountData;
  }

  @action
  async handleClientActivityQuery({ dateType, monthIdx, year }) {
    this.showBillingStartModal = false;
    switch (dateType) {
      case 'cancel':
        return;
      case 'reset': // clicked 'Current billing period' in calendar widget -> reset to initial start/end dates
        this.activityQueryParams.start.timestamp = this.args.model.licenseStartTimestamp;
        this.activityQueryParams.end.timestamp = this.args.model.currentDate;
        break;
      case 'currentMonth': // clicked 'Current month' from calendar widget
        this.activityQueryParams.start.timestamp = this.args.model.currentDate;
        this.activityQueryParams.end.timestamp = this.args.model.currentDate;
        break;
      case 'startDate': // from "Edit billing start" modal
        this.activityQueryParams.start = { monthIdx, year };
        this.activityQueryParams.end.timestamp = this.args.model.currentDate;
        break;
      case 'endDate': // selected month and year from calendar widget
        this.activityQueryParams.end = { monthIdx, year };
        break;
      default:
        break;
    }
    try {
      this.isLoadingQuery = true;
      const response = await this.store.queryRecord('clients/activity', {
        start_time: this.activityQueryParams.start,
        end_time: this.activityQueryParams.end,
      });
      // preference for byMonth timestamps because those correspond to a user's query
      const { byMonth } = response;
      this.startMonthTimestamp = byMonth[0]?.timestamp || response.startTime;
      this.endMonthTimestamp = byMonth[byMonth.length - 1]?.timestamp || response.endTime;
      if (response.id === 'no-data') {
        this.noActivityData = true;
      } else {
        this.noActivityData = false;
        getStorage().setItem('vault:ui-inputted-start-date', this.startMonthTimestamp);
      }
      this.queriedActivityResponse = response;

      // reset search-select filters
      this.selectedNamespace = null;
      this.selectedAuthMethod = null;
      this.authMethodOptions = [];
    } catch (e) {
      this.errorObject = e;
      return e;
    } finally {
      this.isLoadingQuery = false;
    }
  }

  get hasMultipleMonthsData() {
    return this.byMonthActivityData && this.byMonthActivityData.length > 1;
  }

  @action
  selectNamespace([value]) {
    this.selectedNamespace = value;
    if (!value) {
      this.authMethodOptions = [];
      // on clear, also make sure auth method is cleared
      this.selectedAuthMethod = null;
    } else {
      // Side effect: set auth namespaces
      const mounts = this.filteredActivityByNamespace.mounts?.map((mount) => ({
        id: mount.label,
        name: mount.label,
      }));
      this.authMethodOptions = mounts;
    }
  }

  @action
  setAuthMethod([authMount]) {
    this.selectedAuthMethod = authMount;
  }

  // validation function sent to <DateDropdown> selecting 'endDate'
  @action
  isEndBeforeStart(selection) {
    let { start } = this.activityQueryParams;
    start = start?.timestamp ? parseAPITimestamp(start.timestamp) : new Date(start.year, start.monthIdx);
    return isBefore(selection, start) && !isSameMonth(start, selection)
      ? `End date must be after ${format(start, 'MMMM yyyy')}`
      : false;
  }
}
