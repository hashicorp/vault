/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import Component from '@glimmer/component';
import { service } from '@ember/service';

import type RouterService from '@ember/routing/router-service';
import type { RoleOrderRouteModel } from 'pki/routes/external/roles/role/order';

interface Args {
  model: RoleOrderRouteModel;
  breadcrumbs: unknown[];
}

export default class ExternalPkiPageRolesRoleOrderComponent extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;

  @action
  refresh() {
    this.router.refresh('vault.cluster.secrets.backend.pki.external.roles.role.order');
  }
}
