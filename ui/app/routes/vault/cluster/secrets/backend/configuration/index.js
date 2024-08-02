/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
import { CONFIGURABLE_SECRET_ENGINES } from 'vault/helpers/mountable-secret-engines';
import { hash } from 'rsvp';

export default class SecretsBackendConfigurationRoute extends Route {
  @service store;

  async model() {
    const backend = this.modelFor('vault.cluster.secrets.backend');

    if (backend.isV2KV) {
      const canRead = await this.store
        .findRecord('capabilities', `${backend.id}/config`)
        .then((response) => response.canRead);
      // only set these config params if they can read the config endpoint.
      if (canRead) {
        // design wants specific default to show that can't be set in the model
        backend.casRequired = backend.casRequired ? backend.casRequired : 'False';
        backend.deleteVersionAfter = backend.deleteVersionAfter ? backend.deleteVersionAfter : 'Never delete';
      } else {
        // remove the default values from the model if they don't have read access otherwise it will display the defaults even if they've been set (because they error on returning config data)
        backend.set('casRequired', null);
        backend.set('deleteVersionAfter', null);
        backend.set('maxVersions', null);
      }
    }
    // If the engine is configurable fetch the config model(s) for the engine and return it alongside the model
    if (CONFIGURABLE_SECRET_ENGINES.includes(backend.type)) {
      const configModel = await this.fetchConfig(backend.type, backend.id);
      return hash({
        backend,
        configModel,
      });
    }
    return backend;
  }

  async fetchConfig(type, backend) {
    // Fetch the config for the engine type.
    switch (type) {
      case 'aws':
        return await this.fetchAwsRootConfig(backend);
      // ARG TODO add fetchAwsLeaseConfig
      case 'ssh':
        return await this.fetchSshCaConfig(backend);
    }
  }

  async fetchAwsRootConfig(backend) {
    try {
      return await this.store.queryRecord('aws/root-config', { backend });
    } catch (e) {
      return e;
    }
  }

  async fetchSshCaConfig(backend) {
    try {
      return await this.store.queryRecord('ssh/ca-config', { backend });
    } catch (e) {
      return e;
    }
  }
}
