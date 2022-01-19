import Component from '@glimmer/component';
import { mean } from 'd3-array';

/**
 * @module RunningTotal
 * RunningTotal components display total and new client counts in a given date range by month. A line chart shows total monthly clients, and below a stacked, vertical bar chart shows new clients per month.
 * 
 *
 * @example
 * ```js
 *  <Clients::RunningTotal
 *    @chartLegend={{this.chartLegend}}
 *    @lineChartData={{this.newAndTotalMonthlyClients}}
 *    @barChartData={{this.newMonthlyClients}}
 *   />
 * ```

 * @param {array} lineChartData - (passed to child chart) must be an array of objects
 * @param {array} barChartData - (passed to child chart) must be an array of flattened objects
 * @param {array} chartLegend - (passed to child) array of objects with key names 'key' and 'label' so data can be stacked
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
