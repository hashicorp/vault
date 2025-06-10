/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import SecretsEngineResource from 'vault/resources/secrets/engine';

import type ApiService from 'vault/services/api';

export default class SecretsBackends extends Route {
  @service declare readonly api: ApiService;

  async model() {
    const { secret } = await this.api.sys.internalUiListEnabledVisibleMounts();
    return this.api.responseObjectToArray(secret, 'path').map((engine) => new SecretsEngineResource(engine));
  }
}
