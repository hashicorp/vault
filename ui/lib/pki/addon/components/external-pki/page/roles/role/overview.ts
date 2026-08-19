/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import Component from '@glimmer/component';

import type { HTTPQuery, PkiExternalCaReadRoleCachedResponse } from '@hashicorp/vault-client-typescript';
import type { RoleRouteModel } from 'pki/routes/external/roles/role';
import type ApiService from 'vault/services/api';
import type RouterService from '@ember/routing/router-service';

interface Args {
  model: RoleRouteModel;
}

export default class ExternalPkiPageRolesRoleOverviewComponent extends Component<Args> {
  @service declare readonly api: ApiService;
  @service('app-router') declare readonly router: RouterService;

  @tracked errorMessage = '';
  @tracked fetchedCert: PkiExternalCaReadRoleCachedResponse | undefined;
  @tracked roleOrderId = '';

  @action
  async getCertificate(payload: HTTPQuery) {
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

  @action
  lookupOrder() {
    this.router.transitionTo(
      'vault.cluster.secrets.backend.pki.external.roles.role.order',
      this.args.model.engine.id,
      this.args.model.role.name,
      this.roleOrderId
    );
  }

  get shouldRenderActiveOrdersCard() {
    const { activeOrders } = this.args.model;
    // If the length is 0 (falsy) but the error status is 404 then the user has permission,
    // there just aren't any orders.
    if (activeOrders?.list?.length || activeOrders?.error?.status === 404) {
      return true;
    }
    return false;
  }
}
