/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type Controller from '@ember/controller';
import type Transition from '@ember/routing/transition';
import type { Breadcrumb } from 'vault/vault/app-types';

import type RouterService from '@ember/routing/router-service';
import { LdapRoleRouteModel } from '../role';
interface LdapRoleDetailsController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapRoleRouteModel;
}

export default class LdapRoleIndexRoute extends Route {
  @service('app-router') declare readonly router: RouterService;

  redirect(model: LdapRoleRouteModel) {
    if (!model.roles) {
      this.router.transitionTo('vault.cluster.secrets.backend.ldap.roles.role.details');
    }
  }

  setupController(
    controller: LdapRoleDetailsController,
    resolvedModel: LdapRoleRouteModel,
    transition: Transition
  ) {
    super.setupController(controller, resolvedModel, transition);
    controller.breadcrumbs = [
      { label: resolvedModel.backendModel.id, route: 'overview' },
      { label: 'roles', route: 'roles' },
      { label: resolvedModel.roleName },
    ];
  }
}
