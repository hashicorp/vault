/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@ember/component';
import { service } from '@ember/service';

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
    /**
     * We get a partially typed response from the API client, but only 1 level deep.
     * Casting the nested types here and falling back to defaults in the mappings.
     * We should get typescript errors if top level interfaces in the API client or
     * the vault-reporting addon change.
     */
    const response = await this.api.sys.generateUtilizationReport();
    const leaseCountQuotas = response.leaseCountQuotas as UsageDashboardData['leaseCountQuotas'];
    const replicationStatus = response.replicationStatus as UsageDashboardData['replicationStatus'];
    const pki = response.pki as UsageDashboardData['pki'];
    const secretSync = response.secretSync as UsageDashboardData['secretSync'];

    const data: UsageDashboardData = {
      authMethods: (response.authMethods as Record<string, number>) || {},
      secretEngines: (response.secretEngines as Record<string, number>) || {},
      leasesByAuthMethod: (response.leasesByAuthMethod as Record<string, number>) || {},
      leaseCountQuotas: {
        globalLeaseCountQuota: {
          capacity: leaseCountQuotas?.globalLeaseCountQuota?.capacity || 0,
          count: leaseCountQuotas?.globalLeaseCountQuota?.count || 0,
          name: leaseCountQuotas?.globalLeaseCountQuota?.name || '',
        },
        totalLeaseCountQuotas: leaseCountQuotas?.totalLeaseCountQuotas || 0,
      },
      replicationStatus: {
        drState: replicationStatus?.drState || 'disabled',
        prState: replicationStatus?.prState || 'disabled',
        drPrimary: replicationStatus?.drPrimary ?? false,
        prPrimary: replicationStatus?.prPrimary ?? false,
      },
      kvv1Secrets: response.kvv1Secrets || 0,
      kvv2Secrets: response.kvv2Secrets || 0,
      namespaces: response.namespaces || 0,
      pki: {
        totalIssuers: pki?.totalIssuers || 0,
        totalRoles: pki?.totalRoles || 0,
      },
      secretSync: {
        totalDestinations: secretSync?.totalDestinations || 0,
      },
    };
    return data as UsageDashboardData;
  };
}
