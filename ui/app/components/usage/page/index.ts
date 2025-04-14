/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@ember/component';
import { service } from '@ember/service';
import type ApiService from 'vault/services/api';
import type { getUsageDataFunction, UsageDashboardData } from '@hashicorp/vault-reporting/types/index';
import { allEngines } from 'vault/helpers/mountable-secret-engines';
import { allMethods } from 'vault/helpers/mountable-auth-methods';
import type FlagsService from 'vault/services/flags';
export default class ClientsActivityComponent extends Component {
  @service declare readonly api: ApiService;
  @service declare readonly flags: FlagsService;

  handleFetchUsageData: getUsageDataFunction = async () => {
    //TODO: Update client with typed response after the API is updated https://hashicorp.atlassian.net/browse/VAULT-35108
    const response = await this.api.sys.systemReadUtilizationReport();
    const data = response.data as UsageDashboardData;
    // Replace engine names with display names if available
    allEngines().forEach((engine) => {
      if (engine.type in data.secret_engines) {
        data.secret_engines[engine.displayName] = data.secret_engines[engine.type] || 0;
        delete data.secret_engines[engine.type];
      }
    });
    // Replace auth method names with display names if available
    allMethods().forEach((method) => {
      if (method.type in data.auth_methods) {
        data.auth_methods[method.displayName] = data.auth_methods[method.type] || 0;
        delete data.auth_methods[method.type];
      }
    });
    return data as UsageDashboardData;
  };

  get isVaultDedicated(): boolean {
    return this.flags.isHvdManaged;
  }
}
