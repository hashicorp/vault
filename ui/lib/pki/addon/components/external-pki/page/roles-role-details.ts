/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import Component from '@glimmer/component';

import type { HTTPQuery, PkiExternalCaReadRoleCachedResponse } from '@hashicorp/vault-client-typescript';
import type { RoleDetailsRouteModel } from 'pki/routes/external/roles/role/details';
import type ApiService from 'vault/services/api';

interface Args {
  model: RoleDetailsRouteModel;
}

export default class ExternalPkiPageRolesRoleDetailsComponent extends Component<Args> {
  @service declare readonly api: ApiService;

  @tracked errorMessage = '';
  @tracked fetchedCert: PkiExternalCaReadRoleCachedResponse | undefined;

  @action
  async fetchCertificate(payload: HTTPQuery) {
    const { role, engine } = this.args.model;
    try {
      this.fetchedCert = await this.api.secrets.pkiExternalCaReadRoleCached(
        role?.name as string,
        engine.id,
        (context) => this.api.addQueryParams(context, payload)
      );
    } catch (error) {
      const { message } = await this.api.parseError(error);
      this.errorMessage = message;
    }
  }
}
