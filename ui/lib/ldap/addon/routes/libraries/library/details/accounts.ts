/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type { LdapLibraryRouteModel } from 'ldap/routes/libraries/library';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';

export default class LdapLibraryRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  async model() {
    const { library, capabilities } = this.modelFor('libraries.library') as LdapLibraryRouteModel;
    const response = await this.api.secrets.ldapLibraryCheckStatus(
      library.completeLibraryName,
      this.secretMountPath.currentPath
    );
    const status = response.data as Record<
      string,
      { available: boolean; borrower_client_token?: string; borrower_entity_id?: string }
    >;

    const statuses = [];
    for (const key in status) {
      statuses.push({
        ...status[key],
        account: key,
        library: library.completeLibraryName,
      });
    }

    return {
      library,
      capabilities,
      statuses,
    };
  }
}
