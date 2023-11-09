/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { calculateAverage } from 'vault/utils/chart-helpers';

/**
 * @module RunningTotal
 * RunningTotal components display total and new client counts in a given date range by month.
 * A line chart shows total monthly clients and below a stacked, vertical bar chart shows new clients per month.
 *
 *
 * @example
      <Clients::RunningTotal
        @chartLegend={{this.chartLegend}}
        @selectedNamespace={{this.selectedNamespace}}
        @byMonthActivityData={{this.byMonth}}
        @runningTotals={{this.runningTotals}}
        @upgradeData={{if this.countsIncludeOlderData this.latestUpgradeData}}
      />
 *

 * @param {array} chartLegend - array of objects with key names 'key' and 'label' so data can be stacked
 * @param {string} selectedAuthMethod - string of auth method label for empty state message in bar chart
 * @param {array} byMonthActivityData - array of objects from /activity response, from the 'months' key, includes total and new clients per month
    object structure: {
      month: '1/22',
      entity_clients: 23,
      non_entity_clients: 45,
      clients: 68,
      namespaces: [],
      new_clients: {
        entity_clients: 11,
        non_entity_clients: 36,
        clients: 47,
        namespaces: [],
      },
    };
 * @param {object} runningTotals - top level totals from /activity response { clients: 3517, entity_clients: 1593, non_entity_clients: 1924 }
 * @param {object} upgradeData -  object containing version upgrade data e.g.: {version: '1.9.0', previousVersion: null, timestampInstalled: '2021-11-03T10:23:16Z'}
 * @param {string} timestamp -  ISO timestamp created in serializer to timestamp the response
 *
 */
export default class RunningTotal extends Component {
  get byMonthNewClients() {
    if (this.args.byMonthActivityData) {
      return this.args.byMonthActivityData?.map((m) => m.new_clients);
    }
    return null;
  }

  get entityClientData() {
    return {
      runningTotal: this.args.runningTotals.entity_clients,
      averageNewClients: calculateAverage(this.byMonthNewClients, 'entity_clients'),
    };
  }

  get nonEntityClientData() {
    return {
      runningTotal: this.args.runningTotals.non_entity_clients,
      averageNewClients: calculateAverage(this.byMonthNewClients, 'non_entity_clients'),
    };
  }

  get hasRunningTotalClients() {
    return (
      typeof this.entityClientData.runningTotal === 'number' ||
      typeof this.nonEntityClientData.runningTotal === 'number'
    );
  }

  get hasAverageNewClients() {
    return (
      typeof this.entityClientData.averageNewClients === 'number' ||
      typeof this.nonEntityClientData.averageNewClients === 'number'
    );
  }

  get singleMonthData() {
    return this.args?.byMonthActivityData[0];
  }
}
