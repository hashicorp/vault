/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { IUsageDashboardService, UsageDashboardData } from '@hashicorp/vault-reporting/types/index';
import Service from '@ember/service';
import { service } from '@ember/service';

export default class UsageService extends Service implements IUsageDashboardService {
  @service declare readonly auth: { currentToken: string };

  async getUsageData() {
    const token = this.auth.currentToken;
    const res = await fetch('/v1/sys/utilization-report', {
      headers: {
        'X-Vault-Token': token,
      },
    });
    if (!res.ok) {
      throw new Error('Failed to fetch usage data');
    }
    const { data } = await res.json();
    return data as UsageDashboardData;
  }
}
