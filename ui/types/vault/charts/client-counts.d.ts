/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// Count and EmptyCount are mutually exclusive
// but that's hard to represent in an interface
// so for now we just have both
interface Count {
  clients?: number;
  entity_clients?: number;
  non_entity_clients?: number;
  secret_syncs?: number;
}
interface EmptyCount {
  count?: null;
}
interface Timestamp {
  month: string; // eg. 12/22
  timestamp: string; // ISO 8601
}

export interface MonthlyChartData extends Count, EmptyCount, Timestamp {
  new_clients?: Count;
}
