/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@ember/component';
import { service } from '@ember/service';
import type ApiService from 'vault/services/api';
import type { getUsageDataFunction, UsageDashboardData } from '@hashicorp/vault-reporting/types/index';
import type FlagsService from 'vault/services/flags';
export default class ClientsActivityComponent extends Component {
  @service declare readonly api: ApiService;
  @service declare readonly flags: FlagsService;

  handleFetchUsageData: getUsageDataFunction = async () => {
    //TODO: Update client with typed response after the API is updated https://hashicorp.atlassian.net/browse/VAULT-35108
    const { data = {} } = await this.api.sys.systemReadUtilizationReport();
    return data as UsageDashboardData;
  };

  get isVaultDedicated(): boolean {
    return this.flags.isHvdManaged;
  }
}
