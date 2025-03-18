import type { IUsageDashboardService, UsageDashboardData } from 'shared-secure-ui/types/reporting/index';
import Service from '@ember/service';
import { service } from '@ember/service';

export default class UsageService extends Service implements IUsageDashboardService {
  //TODO: Typed?
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
