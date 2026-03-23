/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import RouterService from '@ember/routing/router-service';
import { service } from '@ember/service';

import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import { KmipApplicationModel } from './application';

export default class KmipConfigurationRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service('app-router') declare readonly router: RouterService;

  async model() {
    const { secretsEngine } = this.modelFor('application') as KmipApplicationModel;

    try {
      const { currentPath } = this.secretMountPath;
      // the spec changed and now the operation ids are the same for reading both the config and ca pem
      const { data } = await this.api.secrets.kmipReadConfiguration_1(currentPath);
      const config = data as Record<string, unknown>;
      const { secretsEngine } = this.modelFor('application') as KmipApplicationModel;
      try {
        // this method now calls the same endpoint as the former kmipReadCaPem
        const { data } = await this.api.secrets.kmipReadConfiguration(currentPath);
        const ca = data as Record<string, unknown>;
        return { config: { ...config, ...ca }, secretsEngine };
      } catch (error) {
        // ignore error if CA PEM is not found
        // component will conditionally render the field if present
      }
      return { config, secretsEngine };
    } catch (error) {
      const { status } = await this.api.parseError(error);
      if (status !== 404) {
        throw error;
      }
      return { config: null, secretsEngine };
    }
  }

  afterModel(resolvedModel: KmipApplicationModel) {
    if (!resolvedModel.config) {
      this.router.transitionTo('vault.cluster.secrets.backend.kmip.configure');
    }
  }
}
