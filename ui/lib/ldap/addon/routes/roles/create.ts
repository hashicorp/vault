/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

import type Store from '@ember-data/store';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type LdapRoleModel from 'vault/models/ldap/role';
import type Controller from '@ember/controller';
import type Transition from '@ember/routing/transition';
import type { Breadcrumb } from 'vault/vault/app-types';

interface LdapRolesCreateController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapRoleModel;
}

export default class LdapRolesCreateRoute extends Route {
  @service declare readonly store: Store;
  @service declare readonly secretMountPath: SecretMountPath;

  model() {
    const backend = this.secretMountPath.currentPath;
    return this.store.createRecord('ldap/role', { backend });
  }

  setupController(
    controller: LdapRolesCreateController,
    resolvedModel: LdapRoleModel,
    transition: Transition
  ) {
    super.setupController(controller, resolvedModel, transition);

    controller.breadcrumbs = [
      { label: resolvedModel.backend, route: 'overview' },
      { label: 'roles', route: 'roles' },
      { label: 'create' },
    ];
  }
}
