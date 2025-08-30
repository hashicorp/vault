/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type ApiService from 'vault/services/api';
import type { ModelFrom } from 'vault/route';

export type ClusterSettingsAuthConfigureRouteModel = ModelFrom<ClusterSettingsAuthConfigureRoute>;

export default class ClusterSettingsAuthConfigureRoute extends Route {
  @service declare readonly api: ApiService;

  async model(params: { method: string }) {
    const path = params.method;
    const methodOptions = await this.api.sys.authReadConfiguration(path);

    return {
      methodOptions,
      type: methodOptions.type as string,
      id: path,
    };
  }
}
