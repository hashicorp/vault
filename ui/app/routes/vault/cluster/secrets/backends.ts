/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import SecretsEngineResource from 'vault/resources/secrets/engine';

import type ApiService from 'vault/services/api';
import type Store from '@ember-data/store';

export default class SecretsBackends extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly store: Store;

  async model() {
    const adapter = this.store.adapterFor('application');
    const { data } = await adapter.ajax('/v1/sys/internal/ui/mounts', 'GET');
    const secret = data.secret;
    return this.api.responseObjectToArray(secret, 'path').map((engine) => new SecretsEngineResource(engine));
  }
}
