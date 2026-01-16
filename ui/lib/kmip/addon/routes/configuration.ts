/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';

export default class KmipConfigurationRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  async model() {
    try {
      const { currentPath } = this.secretMountPath;
      const { data } = await this.api.secrets.kmipReadConfiguration(currentPath);
      const config = data as Record<string, unknown>;

      try {
        const { data } = await this.api.secrets.kmipReadCaPem(currentPath);
        const ca = data as Record<string, unknown>;
        return { ...config, ...ca };
      } catch (error) {
        // ignore error if CA PEM is not found
        // component will conditionally render the field if present
      }
      return config;
    } catch (error) {
      const { status } = await this.api.parseError(error);
      if (status !== 404) {
        throw error;
      }
      return undefined;
    }
  }
}
