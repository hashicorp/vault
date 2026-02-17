/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { KmipListClientCertificatesListEnum } from '@hashicorp/vault-client-typescript';
import { paginate } from 'core/utils/paginate-list';

import type Controller from '@ember/controller';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type CapabilitiesService from 'vault/services/capabilities';

interface KmipCredentialsController extends Controller {
  pageFilter: string | undefined;
  page: number | undefined;
}

export default class KmipCredentialsRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly capabilities: CapabilitiesService;

  async model(params: { page: number; pageFilter: string }) {
    const { page, pageFilter } = params;
    const { role_name, scope_name } = this.paramsFor('credentials');
    const { currentPath } = this.secretMountPath;
    const model = {
      roleName: role_name,
      scopeName: scope_name,
      credentials: [],
      capabilities: {},
      filterValue: pageFilter,
    };

    try {
      const { keys } = await this.api.secrets.kmipListClientCertificates(
        role_name as string,
        scope_name as string,
        currentPath,
        KmipListClientCertificatesListEnum.TRUE
      );
      const credentials = keys ? paginate(keys, { page: Number(page) || 1, filter: pageFilter }) : [];
      // capabilities exist at root path, not for individual credentials
      const capabilities = await this.capabilities.for('kmipCredentialsRevoke', {
        backend: currentPath,
        role: role_name,
        scope: scope_name,
      });

      return { ...model, credentials, capabilities };
    } catch (error) {
      const { status } = await this.api.parseError(error);
      if (status === 404) {
        return model;
      }
      throw error;
    }
  }

  resetController(controller: KmipCredentialsController, isExiting: boolean) {
    if (isExiting) {
      controller.set('pageFilter', undefined);
      controller.set('page', undefined);
    }
  }
}
