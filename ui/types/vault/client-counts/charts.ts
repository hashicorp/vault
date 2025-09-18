/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { TotalClients } from './activity-api';

// TotalClients and EmptyCount are mutually exclusive
// but that's hard to represent in an interface
// so for now we just have both
interface EmptyCount {
  count?: null;
}
interface Timestamp {
  month: string; // eg. 12/22
  timestamp: string; // ISO 8601
}

export interface MonthlyChartData extends TotalClients, EmptyCount, Timestamp {
  new_clients?: TotalClients;
}
