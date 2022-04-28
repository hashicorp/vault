import Component from '@glimmer/component';
import { mean } from 'd3-array';

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
        @barChartData={{this.byMonthNewClients}}
        @lineChartData={{this.byMonth}}
        @runningTotals={{this.runningTotals}}
        @upgradeData={{if this.countsIncludeOlderData this.latestUpgradeData}}
      />
 * ```

 * @param {array} chartData - array of objects from /activity response
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
 * @param {array} chartLegend - array of objects with key names 'key' and 'label' so data can be stacked
 * @param {object} runningTotals - top level totals from /activity response { clients: 3517, entity_clients: 1593, non_entity_clients: 1924 }
 * @param {string} timestamp -  ISO timestamp created in serializer to timestamp the response
 * @param {object} upgradeData -  object containing version upgrade data e.g.: {id: '1.9.0', previousVersion: null, timestampInstalled: '2021-11-03T10:23:16Z'}
 *   
 */
export default class RunningTotal extends Component {
  get lineChartData() {
    if (!this.args.selectedNamespace) {
      return this.args.chartData;
    }
    return this.filterBySelectedNamespace(this.args.chartData, this.args.selectedNamespace);
  }

  get barChartData() {
    if (!this.args.selectedNamespace) {
      return this.args.chartData.map((m) => m.new_clients);
    }

    let test = this.filterBySelectedNamespace(this.args.chartData, this.args.selectedNamespace);
    return test.map((m) => {
      let { month, new_clients } = m;
      return { month, ...new_clients };
    });
  }

  get entityClientData() {
    return {
      runningTotal: this.args.runningTotals.entity_clients,
      averageNewClients: Math.round(mean(this.barChartData?.map((d) => d.entity_clients))),
    };
  }

  get nonEntityClientData() {
    return {
      runningTotal: this.args.runningTotals.non_entity_clients,
      averageNewClients: Math.round(mean(this.barChartData?.map((d) => d.non_entity_clients))),
    };
  }

  filterBySelectedNamespace(chartData, namespace) {
    return chartData.map((m) => {
      let { by_namespace_key, month } = m;
      if (by_namespace_key && by_namespace_key[namespace]) {
        return { month, ...by_namespace_key[namespace] };
      } else {
        return m;
      }
    });
  }
}
