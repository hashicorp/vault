/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { Model } from 'vault/app-types';
import type {
  ByMonthClients,
  ByNamespaceClients,
  TotalClients,
} from 'vault/vault/client-counts/activity-api';

export default interface ClientsActivityModel extends Model {
  byMonth: ByMonthClients[];
  byNamespace: ByNamespaceClients[];
  total: TotalClients;
  startTime: string;
  endTime: string;
  responseTimestamp: string;
}
