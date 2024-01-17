/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { calculateAverage } from 'vault/utils/chart-helpers';

/**
 * @module MonthlyUsage
 * MonthlyUsage components show how many total clients use Vault each month. Displaying the average totals to the left of a stacked, vertical bar chart.
 *
 * @example
 * ```js
  <Clients::MonthlyUsage
    @responseTimestamp={{this.responseTimestamp}}
    @verticalBarChartData={{this.byMonthActivityData}}
  />
 * ```
 * @param {string} timestamp -  ISO timestamp created in serializer to timestamp the response
 * @param {array} verticalBarChartData - array of flattened objects
    sample object =
    {
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
    }
 */
export default class MonthlyUsage extends Component {
  monthlyUsageLegend = [
    { key: 'entity_clients', label: 'entity clients' },
    { key: 'non_entity_clients', label: 'non-entity clients' },
  ];

  get averageTotalClients() {
    return calculateAverage(this.args.verticalBarChartData, 'clients') || '0';
  }

  get averageNewClients() {
    return (
      calculateAverage(
        this.args.verticalBarChartData?.map((d) => d.new_clients),
        'clients'
      ) || '0'
    );
  }
}
