/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { withConfig } from 'core/decorators/fetch-secrets-engine-config';
import { hash } from 'rsvp';

import type Store from '@ember-data/store';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type Transition from '@ember/routing/transition';
import type LdapRoleModel from 'vault/models/ldap/role';
import type SecretEngineModel from 'vault/models/secret-engine';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/vault/app-types';

interface LdapRolesRouteModel {
  backendModel: SecretEngineModel;
  promptConfig: boolean;
  roles: Array<LdapRoleModel>;
}
interface LdapRolesController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapRolesRouteModel;
}

@withConfig('ldap/config')
export default class LdapRolesRoute extends Route {
  @service declare readonly store: Store;
  @service declare readonly secretMountPath: SecretMountPath;

  declare promptConfig: boolean;

  model() {
    const backendModel = this.modelFor('application') as SecretEngineModel;
    return hash({
      backendModel,
      promptConfig: this.promptConfig,
      roles: this.store.query(
        'ldap/role',
        { backend: backendModel.id },
        { adapterOptions: { showPartialError: true } }
      ),
    });
  }

  setupController(
    controller: LdapRolesController,
    resolvedModel: LdapRolesRouteModel,
    transition: Transition
  ) {
    super.setupController(controller, resolvedModel, transition);

    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backendModel.id },
    ];
  }
}
