/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ldapBreadcrumbs, libraryRoutes } from 'ldap/utils/ldap-breadcrumbs';
import LdapLibraryForm from 'vault/forms/secrets/ldap/library';
import { ModelFrom } from 'vault/route';

import type { LdapLibraryRouteModel } from 'ldap/routes/libraries/library';
import type Controller from '@ember/controller';
import type Transition from '@ember/routing/transition';
import type { Breadcrumb } from 'vault/vault/app-types';
import type SecretMountPath from 'vault/services/secret-mount-path';

export type LdapLibraryEditRouteModel = ModelFrom<LdapLibraryEditRoute>;

interface LdapLibraryEditController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapLibraryEditRouteModel;
}

export default class LdapLibraryEditRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;

  model() {
    const { library } = this.modelFor('libraries.library') as LdapLibraryRouteModel;
    const disable_check_in_enforcement = library.disable_check_in_enforcement ? 'Disabled' : 'Enabled';
    return new LdapLibraryForm({ ...library, disable_check_in_enforcement });
  }

  setupController(
    controller: LdapLibraryEditController,
    resolvedModel: LdapLibraryEditRouteModel,
    transition: Transition
  ) {
    super.setupController(controller, resolvedModel, transition);

    const { currentPath } = this.secretMountPath;
    const { library } = this.modelFor('libraries.library') as LdapLibraryRouteModel;
    const routeParams = (childResource: string) => {
      return [currentPath, childResource];
    };
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: currentPath, route: 'overview' },
      { label: 'Libraries', route: 'libraries' },
      ...ldapBreadcrumbs(library.completeLibraryName, routeParams, libraryRoutes),
      { label: 'Edit' },
    ];
  }
}
