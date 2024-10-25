/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import type RouterService from '@ember/routing/router-service';
import { LdapRolesRoleRouteModel } from '../role';

export default class LdapRolesRoleIndexRoute extends Route {
  @service('app-router') declare readonly router: RouterService;

  redirect(model: LdapRolesRoleRouteModel) {
    if (model?.roles) {
      this.router.transitionTo('vault.cluster.secrets.backend.ldap.roles.subdirectoy');
    }
    this.router.transitionTo('vault.cluster.secrets.backend.ldap.roles.role.details');
  }
}
