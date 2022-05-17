import Component from '@glimmer/component';
import { calculateAverageClients } from 'vault/utils/chart-helpers';

/**
 * @module RunningTotal
 * RunningTotal components display total and new client counts in a given date range by month.
 * A line chart shows total monthly clients and below a stacked, vertical bar chart shows new clients per month.
 *
 *
 * @example
 * ```js
      <Clients::RunningTotal
        @chartLegend={{this.chartLegend}}
        @selectedNamespace={{this.selectedNamespace}}
        @barChartData={{this.byMonthNewClients}}
        @lineChartData={{this.byMonth}}
        @runningTotals={{this.runningTotals}}
        @upgradeData={{if this.countsIncludeOlderData this.latestUpgradeData}}
      />
 * ```

 * @param {array} chartLegend - array of objects with key names 'key' and 'label' so data can be stacked
 * @param {string} selectedAuthMethod - string of auth method label for empty state message in bar chart
 * @param {array} barChartData - array of objects from /activity response, from the 'months' key
    object example: {
      month: '1/22',
      entity_clients: 23,
      non_entity_clients: 45,
      total: 68,
      namespaces: [],
      new_clients: {
        entity_clients: 11,
        non_entity_clients: 36,
        total: 47,
        namespaces: [],
      },
    };
 * @param {array} lineChartData - array of objects from /activity response, from the 'months' key
 * @param {object} runningTotals - top level totals from /activity response { clients: 3517, entity_clients: 1593, non_entity_clients: 1924 }
 * @param {object} upgradeData -  object containing version upgrade data e.g.: {id: '1.9.0', previousVersion: null, timestampInstalled: '2021-11-03T10:23:16Z'}
 * @param {string} timestamp -  ISO timestamp created in serializer to timestamp the response
 *
 */
export default class RunningTotal extends Component {
  get entityClientData() {
    return {
      runningTotal: this.args.runningTotals.entity_clients,
      averageNewClients: calculateAverageClients(this.args.barChartData, 'entity_clients') || '0',
    };
  }

  get nonEntityClientData() {
    return {
      runningTotal: this.args.runningTotals.non_entity_clients,
      averageNewClients: calculateAverageClients(this.args.barChartData, 'non_entity_clients') || '0',
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

  get showSingleMonth() {
    if (this.args.barChartData.length === 1) {
      const monthData = this.args.lineChartData[0];
      return {
        total: {
          total: monthData.clients,
          entityClients: monthData.entity_clients,
          nonEntityClients: monthData.non_entity_clients,
        },
        new: {
          total: monthData.new_clients.clients,
          entityClients: monthData.new_clients.entity_clients,
          nonEntityClients: monthData.new_clients.non_entity_clients,
        },
      };
    }
    return null;
  }
}
