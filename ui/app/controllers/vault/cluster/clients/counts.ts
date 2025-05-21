/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { action, set } from '@ember/object';

import type { ClientsCountsRouteParams } from 'vault/routes/vault/cluster/clients/counts';

const queryParamKeys = ['start_time', 'end_time', 'ns', 'mountPath'];
export default class ClientsCountsController extends Controller {
  queryParams = queryParamKeys;

  start_time: string | number | undefined = undefined;
  end_time: string | number | undefined = undefined;
  ns: string | undefined = undefined;
  mountPath: string | undefined = undefined;

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
}
