/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import KmipConfigForm from 'vault/forms/secrets/kmip/config';

import type SecretMountPath from 'vault/services/secret-mount-path';
import type ApiService from 'vault/services/api';
import { KmipApplicationModel } from './application';

export default class KmipConfigureRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly api: ApiService;

  async model() {
    const { secretsEngine } = this.modelFor('application') as KmipApplicationModel;
    try {
      const { data } = await this.api.secrets.kmipReadConfiguration(this.secretMountPath.currentPath);
      return { form: new KmipConfigForm(data as object), secretsEngine };
    } catch (error) {
      const { status } = await this.api.parseError(error);
      if (status === 404) {
        return { form: new KmipConfigForm(), secretsEngine };
      }
      throw error;
    }
  }
}
