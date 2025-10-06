/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import AuthMethodResource from 'vault/resources/auth/method';

import type ApiService from 'vault/services/api';
import type { ModelFrom } from 'vault/route';
import type { Mount } from 'vault/vault/mount';

export type ClusterSettingsAuthConfigureRouteModel = ModelFrom<ClusterSettingsAuthConfigureRoute>;

export default class ClusterSettingsAuthConfigureRoute extends Route {
  @service declare readonly api: ApiService;

  async model(params: { method: string }) {
    const path = params.method;
    const methodOptions = (await this.api.sys.authReadConfiguration(path)) as Mount;
    const method = new AuthMethodResource({ ...methodOptions, path }, this);
    return {
      methodOptions,
      method,
    };
  }
}
