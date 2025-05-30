/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
import SecretsEngineResource from 'vault/resources/secrets/engine';

export default Route.extend({
  flashMessages: service(),
  router: service(),
  secretMountPath: service(),
  api: service(),

  oldModel: null,

  async model(params) {
    const { backend } = params;
    this.secretMountPath.update(backend);

    const secretsEngine = await this.api.sys.mountsReadConfiguration(backend);
    return new SecretsEngineResource({ ...secretsEngine, path: `${backend}/` });
  },

  afterModel(model, transition) {
    const path = model && model.path;
    if (transition.targetName === this.routeName) {
      return this.router.replaceWith('vault.cluster.secrets.backend.list-root', path);
    }
  },
});
