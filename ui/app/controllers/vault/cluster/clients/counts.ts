/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { action, set } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { ByMonthClients, emptyCounts, MountByKey, MountClients } from 'core/utils/client-count-utils';

import type {
  ClientsCountsRouteModel,
  ClientsCountsRouteParams,
} from 'vault/routes/vault/cluster/clients/counts';
import type NamespaceService from 'vault/services/namespace';

const queryParamKeys = ['start_time', 'end_time', 'ns', 'mountPath'];
export default class ClientsCountsController extends Controller {
  queryParams = queryParamKeys;
  @service('namespace') declare readonly namespaceSvc: NamespaceService;

  @tracked start_time: string | number | undefined = undefined;
  @tracked end_time: string | number | undefined = undefined;
  @tracked ns: string | undefined = undefined;
  @tracked mountPath: string | undefined = undefined;

  // using router.transitionTo to update the query params results in the model hook firing each time
  // this happens when the queryParams object is not added to the route or refreshModel is explicitly set to false
  // updating the bound properties does however respect the refreshModel settings and functions expectedly
  @action
  updateQueryParams(updatedParams: ClientsCountsRouteParams) {
    if (!updatedParams) {
      this.queryParams.forEach((key) => (this[key as keyof ClientsCountsRouteParams] = undefined));
    } else {
      Object.keys(updatedParams).forEach((key) => {
        if (queryParamKeys.includes(key)) {
          const value = updatedParams[key as keyof ClientsCountsRouteParams];
          set(this, key as keyof ClientsCountsRouteParams, value as keyof ClientsCountsRouteParams);
        }
      });
    }
  }

  get filteredActivityTotals() {
    const { activity } = this.model as ClientsCountsRouteModel;
    const { ns, mountPath } = this;

    // only do this if we have a mountPath filter.
    // namespace is filtered on API layer
    if (activity?.byNamespace && ns && mountPath) {
      const filtered = activity.byNamespace
        .find((namespace) => namespace.label === ns)
        ?.mounts.find((mount: MountClients) => mount.label === mountPath);
      return filtered;
    }
    return activity?.total;
  }

  get filteredByMonthActivity(): (ByMonthClients | MountByKey)[] {
    const { activity } = this.model as ClientsCountsRouteModel;
    const { ns, mountPath } = this;

    // only do this if we have a mountPath filter.
    // namespace is filtered on API layer
    if (activity?.byMonth && ns && mountPath) {
      const mountData = activity.byMonth
        ?.map((m) => {
          if (!m?.clients) {
            // if the month doesn't have data (null) or it's zero, we can just return the block
            return m;
          }
          const namespace = m.namespaces_by_key[ns];
          const mount = namespace?.mounts_by_key[mountPath];
          if (mount) return mount;

          // if the month has data but none for this mount, return mocked zeros
          return {
            label: mountPath,
            timestamp: m.timestamp,
            month: m.month,
            ...emptyCounts(),
            new_clients: {
              timestamp: m.timestamp,
              month: m.month,
              label: mountPath,
              ...emptyCounts(),
            },
          } as MountByKey;
        })
        .filter((d) => d !== undefined);

      return mountData || [];
    }
    return activity?.byMonth;
  }
}
