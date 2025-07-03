/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import AuthMethodResource from 'vault/resources/auth/method';

import type ApiService from 'vault/services/api';

export default class VaultClusterAccessMethodsRoute extends Route {
  @service declare readonly api: ApiService;

  queryParams = {
    page: {
      refreshModel: true,
    },
    pageFilter: {
      refreshModel: true,
    },
  };

  async model() {
    const { auth } = await this.api.sys.internalUiListEnabledVisibleMounts();
    return this.api.responseObjectToArray(auth, 'path').map((method) => new AuthMethodResource(method, this));
  }
}
