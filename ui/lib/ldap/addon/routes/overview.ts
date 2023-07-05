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
import type SecretEngineModel from 'vault/models/secret-engine';
// import type LdapRoleModel from 'vault/models/ldap/role';
// import type LdapLibraryModel from 'vault/models/ldap/library';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/vault/app-types';

interface LdapOverviewController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
}
interface LdapOverviewRouteModel {
  backendModel: SecretEngineModel;
  promptConfig: boolean;
  // roles: Array<LdapRoleModel>;
  // libraries: Array<LdapLibraryModel>;
}

@withConfig('ldap/config')
export default class LdapConfigureRoute extends Route {
  @service declare readonly store: Store;
  @service declare readonly secretMountPath: SecretMountPath;

  declare promptConfig: boolean;

  async model() {
    // roles and libraries will be needed to pass into card components
    // add to hash once models have been created

    // const backend = this.secretMountPath.currentPath;
    return hash({
      promptConfig: this.promptConfig,
      backendModel: this.modelFor('application'),
      // roles: this.store.query('ldap/role', { backend }).catch(() => []),
      // libraries: this.store.query('ldap/libraries', { backend }).catch(() => []),
    });
  }

  setupController(
    controller: LdapOverviewController,
    resolvedModel: LdapOverviewRouteModel,
    transition: Transition
  ) {
    super.setupController(controller, resolvedModel, transition);

    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backendModel.id },
    ];
  }
}
