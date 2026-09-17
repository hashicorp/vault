/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import LdapStaticRoleForm from 'vault/forms/secrets/ldap/roles/static';
import LdapDynamicRoleForm from 'vault/forms/secrets/ldap/roles/dynamic';

import { ModelFrom } from 'vault/vault/route';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/vault/app-types';
import type { LdapApplicationModel } from '../application';

export type LdapRolesCreateRouteModel = ModelFrom<LdapRolesCreateRoute>;

interface LdapRolesCreateController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapRolesCreateRouteModel;
}

export default class LdapRolesCreateRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;

  model() {
    const { config } = this.modelFor('application') as LdapApplicationModel;
    const isSelfManaged = config?.self_managed === true;
    return {
      // Defaults to 0s: the customer is choosing not to rotate.
      staticForm: new LdapStaticRoleForm({ rotation_period: '0s' }, { isNew: true, isSelfManaged }),
      dynamicForm: new LdapDynamicRoleForm({ default_ttl: '1h', max_ttl: '24h' }, { isNew: true }),
      isSelfManaged,
    };
  }

  setupController(controller: LdapRolesCreateController, resolvedModel: LdapRolesCreateRouteModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
      { label: 'Roles', route: 'roles' },
      { label: 'Create' },
    ];
  }
}
