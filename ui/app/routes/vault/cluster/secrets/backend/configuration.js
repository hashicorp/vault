/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';

export default class SecretsBackendConfigurationRoute extends Route {
  @service store;
  @service secretMountPath;

  async fetchAwsRootConfig(backend) {
    return await this.store.queryRecord('aws/root-config', { backend });
  }

  async fetchSshCaConfig(backend) {
    return await this.store.queryRecord('ssh/ca-config', { backend });
  }

  async model() {
    const backend = this.modelFor('vault.cluster.secrets.backend');
    backend.configModel = null; // reset the config model
    backend.configError = null; // reset the config error
    // Currently two secret engines that return configuration data and that can be configured by the user on the ui: aws, and ssh.
    if (backend.type === 'aws') {
      try {
        backend.configModel = await this.fetchAwsRootConfig(backend.id);
      } catch (e) {
        backend.configError = e;
      }
    }
    if (backend.type === 'ssh') {
      try {
        backend.configModel = await this.fetchSshCaConfig(backend.id);
      } catch (e) {
        backend.configError = e;
      }
    }
    if (backend.isV2KV) {
      const canRead = this.store.findRecord('capabilities', `${backend.id}/config`).canRead;
      // only set these config params if they can read the config endpoint.
      if (canRead) {
        // design wants specific default to show that can't be set in the model
        backend.set('casRequired', backend.casRequired ? backend.casRequired : 'False');
        backend.set(
          'deleteVersionAfter',
          backend.deleteVersionAfter !== '0s' ? backend.deleteVersionAfter : 'Never delete'
        );
      } else {
        // remove the default values from the model if they don't have read access otherwise it will display the defaults even if they've been set (because they error on returning config data)
        backend.set('casRequired', null);
        backend.set('deleteVersionAfter', null);
        backend.set('maxVersions', null);
      }
    }
    return backend;
  }
}
