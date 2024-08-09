/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { action, set } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { MountClients } from 'core/utils/client-count-utils';

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
    if (activity?.byMonth && ns && mountPath) {
      const filtered = activity.byNamespace
        .find((namespace) => namespace.label === ns)
        ?.mounts.find((mount: MountClients) => mount.label === mountPath);
      return filtered;
    }
    return activity?.total;
  }

  get filteredByMonthActivity() {
    const { activity } = this.model as ClientsCountsRouteModel;
    const { ns, mountPath } = this;

    // only do this if we have a mountPath filter.
    // namespace is filtered on API layer
    if (activity?.byMonth && ns && mountPath) {
      const namespaceData = activity.byMonth
        ?.map((m) => m.namespaces_by_key[ns])
        .filter((d) => d !== undefined);

      const mountData = namespaceData
        ?.map((namespace) => namespace?.mounts_by_key[mountPath])
        .filter((d) => d !== undefined);

      return mountData || [];
    }
    return activity?.byMonth;
  }
}
