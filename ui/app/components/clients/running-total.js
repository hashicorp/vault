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
      />
 * ```

 * @param {array} lineChartData - array of objects
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
 * @param {array} barChartData - array of objects, object example: { month: '1/22', entity_clients: 11, non_entity_clients: 36, total: 47, namespaces: [] };
 * @param {array} chartLegend - array of objects with key names 'key' and 'label' so data can be stacked
 * @param {object} runningTotals - top level totals from /activity response { clients: 3517, entity_clients: 1593, non_entity_clients: 1924 }
 * @param {string} timestamp -  ISO timestamp created in serializer to timestamp the response
 *   
 */
export default class RunningTotal extends Component {
  get entityClientData() {
    return {
      runningTotal: this.args.runningTotals.entity_clients,
      averageNewClients: Math.round(mean(this.args.barChartData?.map((d) => d.entity_clients))),
    };
  }

  get nonEntityClientData() {
    return {
      runningTotal: this.args.runningTotals.non_entity_clients,
      averageNewClients: Math.round(mean(this.args.barChartData?.map((d) => d.non_entity_clients))),
    };
  }
}
