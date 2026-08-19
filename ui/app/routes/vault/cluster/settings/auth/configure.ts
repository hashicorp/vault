/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import AuthMethodResource from 'vault/resources/auth/method';

import type ApiService from 'vault/services/api';
import type { ModelFrom } from 'vault/route';

export type ClusterSettingsAuthConfigureRouteModel = ModelFrom<ClusterSettingsAuthConfigureRoute>;

export default class ClusterSettingsAuthConfigureRoute extends Route {
  @service declare readonly api: ApiService;

  async model(params: { method: string }) {
    const { method: path } = params;
    const { data } = await this.api.sys.authListEnabledMethods();
    const method = this.api
      .responseObjectToArray(data as Record<string, object>, 'path')
      .map((m) => new AuthMethodResource(m, this))
      .find((m) => m.id === path);
    if (!method) throw { httpStatus: 404, path };
    return {
      methodOptions: method,
      method,
    };
  }
}
