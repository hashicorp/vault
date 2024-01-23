/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import CountsComponent from '../counts';
import { calculateAverage } from 'vault/utils/chart-helpers';

import type { MonthlyChartData } from 'vault/vault/charts/client-counts';

export default class SyncComponent extends CountsComponent {
  average = (data: MonthlyChartData[], key: string) => {
    return calculateAverage(data, key);
  };
}
