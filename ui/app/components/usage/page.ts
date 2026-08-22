/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@ember/component';
import { service } from '@ember/service';
import { getRelativePath, sanitizePath } from 'core/utils/sanitize-path';

import type FlagsService from 'vault/services/flags';
import type ApiService from 'vault/services/api';
import type NamespaceService from 'vault/services/namespace';
import type AuthService from 'vault/services/auth';
import type {
  getUsageDataFunction,
  getNamespaceDataFunction,
  UsageDashboardData,
} from 'vault/types/usage-reporting';
import type { UtilizationReport } from 'vault/usage';

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
  @service declare readonly namespace: NamespaceService;
  @service declare readonly auth: AuthService;

  handleFetchUsageData: getUsageDataFunction = async (namespace?: string) => {
    /**
     * We get a partially typed response from the API client, but only 1 level deep.
     * Casting the nested types here and falling back to defaults in the mappings.
     * We should get typescript errors if top level interfaces in the API client or
     * the vault-reporting addon change.
     */

    // Fetch utilization data from the user's root namespace context and pass a relative namespace query.
    const userRootNamespace = sanitizePath(this.namespace.userRootNamespace);
    const selectedNamespace = namespace === 'root' ? '' : sanitizePath(namespace || '');
    const relativeNamespace = getRelativePath(selectedNamespace, userRootNamespace);

    const initOverrides = userRootNamespace
      ? this.api.buildHeaders({ namespace: userRootNamespace })
      : undefined;

    const response = (await this.api.sys.generateUtilizationReport(
      relativeNamespace || undefined,
      initOverrides
    )) as UtilizationReport;

    const { lease_count_quotas, replication_status, pki, secret_sync } = response;

    const data: UsageDashboardData = {
      authMethods: (response.auth_methods as Record<string, number>) || {},
      secretEngines: (response.secret_engines as Record<string, number>) || {},
      leasesByAuthMethod: (response.leases_by_auth_method as Record<string, number>) || {},
      leaseCountQuotas: {
        globalLeaseCountQuota: {
          capacity: lease_count_quotas?.global_lease_count_quota?.capacity || 0,
          count: lease_count_quotas?.global_lease_count_quota?.count || 0,
          name: lease_count_quotas?.global_lease_count_quota?.name || '',
        },
        totalLeaseCountQuotas: lease_count_quotas?.total_lease_count_quotas || 0,
      },
      replicationStatus: {
        drState: replication_status?.dr_state || 'disabled',
        prState: replication_status?.pr_state || 'disabled',
        drPrimary: replication_status?.dr_primary ?? false,
        prPrimary: replication_status?.pr_primary ?? false,
      },
      kvv1Secrets: response.kvv1_secrets || 0,
      kvv2Secrets: response.kvv2_secrets || 0,
      namespaces: response.namespaces || 0,
      pki: {
        totalIssuers: pki?.total_issuers || 0,
        totalRoles: pki?.total_roles || 0,
      },
      secretSync: {
        destinations: secret_sync?.destinations || {},
        totalDestinations: secret_sync?.total_destinations || 0,
      },
    };

    return data;
  };

  handleFetchNamespaceData: getNamespaceDataFunction = async () => {
    await this.namespace?.findNamespacesForUser?.perform();
    const options = this.namespace.getOptions();
    const data = {
      keys: options.map((option) => option.label),
    };

    return data;
  };
}
