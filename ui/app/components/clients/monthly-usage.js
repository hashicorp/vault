import Component from '@glimmer/component';
import { mean } from 'd3-array';

/**
 * @module MonthlyUsage
 * MonthlyUsage components display the top 10 total client counts for namespaces or auth methods (mounts) during a billing period.
 * If view is filtered for a single month, two graphs display and show a breakdown of new and total client counts by namespace or auth method, respectively 
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
