/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
import { didCancel } from 'ember-concurrency';
import SecretsEngineResource from 'vault/resources/secrets/engine';

export default Route.extend({
  flashMessages: service(),
  router: service(),
  secretMountPath: service(),
  api: service(),
  kvMountRetry: service(),

  oldModel: null,

  async model(params) {
    const { backend } = params;
    this.secretMountPath.update(backend);

    try {
      // A brand-new mount can be transiently unroutable for a short window right after
      // POST /sys/mounts/<path> succeeds; retry instead of surfacing the raw error.
      const secretsEngine = await this.kvMountRetry.loadMountInfo.perform(backend);
      return new SecretsEngineResource({ ...secretsEngine, path: `${backend}/` });
    } catch (e) {
      if (didCancel(e)) throw e;

      // the backend.error template is expecting additional data so for now we will catch and rethrow
      const error = await this.api.parseError(e);
      throw {
        backend,
        httpStatus: error.status,
        ...error,
      };
    }
  },

  afterModel(model, transition) {
    const path = model && model.path;
    if (transition.targetName === this.routeName) {
      return this.router.replaceWith('vault.cluster.secrets.backend.list-root', path);
    }
  },
});
