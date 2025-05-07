/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@ember/component';
import { service } from '@ember/service';
import { allEngines } from 'vault/helpers/mountable-secret-engines';
import { allMethods } from 'vault/helpers/mountable-auth-methods';

import type FlagsService from 'vault/services/flags';
import type ApiService from 'vault/services/api';
import type { getUsageDataFunction, UsageDashboardData } from '@hashicorp/vault-reporting/types/index';

/**
 * @module UsagePage
 * @description This component is responsible for fetching usage data and mounting the vault-reporting dashboard view.
 * It uses the `api` service to make a request to the sys/utilization-report endpoint to get the usage data.
 * The data is then processed to replace engine and auth method names with their display names if available.
 * The component also uses the `flags` service to determine if the cluster is HVD managed or not.
 *
 * The logic is self-contained and this component has no args.
 *
 * @example ```js
 *   <Usage::Page />
 * ```
 */

export default class UsagePage extends Component {
  @service declare readonly api: ApiService;
  @service declare readonly flags: FlagsService;

  handleFetchUsageData: getUsageDataFunction = async () => {
    //TODO: Update client with typed response after the API is updated https://hashicorp.atlassian.net/browse/VAULT-35108
    const response = await this.api.sys.generateUtilizationReport();
    const data = response as UsageDashboardData;
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
}
