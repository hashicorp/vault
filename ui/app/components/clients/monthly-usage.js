import Component from '@glimmer/component';
import { mean } from 'd3-array';

/**
 * @module MonthlyUsage
 * MonthlyUsage components show how many total clients use Vault each month. Displaying the average totals to the left of a stacked, vertical bar chart.
 * 
 * @example
 * ```js
 *  <Clients::MonthlyUsage
 *    @chartLegend={{this.chartLegend}}
 *    @verticalBarChartData={{this.totalMonthlyClients}}
 *   />
 * ```

 * @param {array} verticalBarChartData - (passed to child chart) must be an array of flattened objects
 * @param {array} chartLegend - (passed to child) array of objects with key names 'key' and 'label' so data can be stacked
 */

export default class MonthlyUsage extends Component {
  get averageTotalClients() {
    let average = mean(this.args.verticalBarChartData?.map((d) => d.total));
    return Math.round(average) || null;
  }

  get averageNewClients() {
    let average = mean(this.args.verticalBarChartData?.map((d) => d.new_clients));
    return Math.round(average) || null;
  }
}
