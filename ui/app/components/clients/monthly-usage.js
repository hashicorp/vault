import Component from '@glimmer/component';
import { mean } from 'd3-array';

/**
 * @module MonthlyUsage
 * MonthlyUsage components show how many total clients use Vault each month. Displaying the average totals to the left of a stacked, vertical bar chart.
 * 
 * @example
 * ```js
  <Clients::MonthlyUsage 
    @chartLegend={{this.chartLegend}} 
    @timestamp={{this.responseTimestamp}}
    @verticalBarChartData={{this.byMonthTotalClients}} 
  />
 * ```
 * @param {array} chartLegend - array of objects with key names 'key' and 'label' so data can be stacked
 * @param {string} timestamp -  ISO timestamp created in serializer to timestamp the response
 * @param {array} verticalBarChartData - array of flattened objects
    sample object = 
    {
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
    }
 * @param {array} chartLegend - array of objects with key names 'key' and 'label' so data can be stacked
 */
export default class MonthlyUsage extends Component {
  barChartData = this.args.verticalBarChartData;

  get averageTotalClients() {
    let average = mean(this.barChartData?.map((d) => d.total));
    return Math.round(average) || null;
  }

  get averageNewClients() {
    let average = mean(this.barChartData?.map((d) => d.new_clients.total));
    return Math.round(average) || null;
  }
}
