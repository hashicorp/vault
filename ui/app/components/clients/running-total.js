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
 *   
 */
export default class RunningTotal extends Component {
  get getTotalClients() {
    return (
      this.args.chartLegend?.map((legend) => {
        return {
          label: legend.label,
          total: this.args.runningTotals[legend.key],
        };
      }) || null
    );
  }

  get getAverageNewClients() {
    // maps through legend and creates array of objects
    // e.g. {label: 'unique entities', average: 43}
    return (
      this.args.chartLegend?.map((legend) => {
        return {
          label: legend.label,
          average: Math.round(mean(this.args.barChartData?.map((d) => d[legend.key]))),
        };
      }) || null
    );
  }
}
