/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { action, set } from '@ember/object';
import { ClientFilters } from 'core/utils/client-count-utils';

import type { ClientsCountsRouteParams } from 'vault/routes/vault/cluster/clients/counts';

// these params refire the request to /sys/internal/counters/activity
const ACTIVITY_QUERY_PARAMS = ['start_time', 'end_time', 'ns'];
// these params client-side filter table data
const DROPDOWN_FILTERS = Object.values(ClientFilters);
const queryParamKeys = [...ACTIVITY_QUERY_PARAMS, ...DROPDOWN_FILTERS];
export default class ClientsCountsController extends Controller {
  queryParams = queryParamKeys;

  start_time: string | number | undefined = undefined;
  end_time: string | number | undefined = undefined;
  ns: string | undefined = undefined; // TODO delete when filter toolbar is removed
  mountPath: string | undefined = undefined; // TODO delete when filter toolbar is removed
  namespace_path = '';
  mount_path = '';
  mount_type = '';

  // using router.transitionTo to update the query params results in the model hook firing each time
  // this happens when the queryParams object is not added to the route or refreshModel is explicitly set to false
  // updating the bound properties does however respect the refreshModel settings and functions expectedly
  @action
  updateQueryParams(updatedParams: ClientsCountsRouteParams) {
    if (!updatedParams) {
      this.queryParams.forEach((key) => (this[key as keyof ClientsCountsRouteParams] = ''));
    } else {
      Object.keys(updatedParams).forEach((key) => {
        if (queryParamKeys.includes(key)) {
          const value = updatedParams[key as keyof ClientsCountsRouteParams];
          set(this, key as keyof ClientsCountsRouteParams, value);
        }
      });
    }
  }
}
