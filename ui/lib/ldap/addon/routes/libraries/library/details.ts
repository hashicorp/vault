/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ldapBreadcrumbs, libraryRoutes } from 'ldap/utils/ldap-breadcrumbs';

import type { LdapLibraryRouteModel } from '../library';
import type Controller from '@ember/controller';
import type Transition from '@ember/routing/transition';
import type { Breadcrumb } from 'vault/vault/app-types';
import type SecretMountPath from 'vault/services/secret-mount-path';

interface LdapLibraryDetailsController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapLibraryRouteModel;
}

export default class LdapLibraryDetailsRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;

  setupController(
    controller: LdapLibraryDetailsController,
    resolvedModel: LdapLibraryRouteModel,
    transition: Transition
  ) {
    super.setupController(controller, resolvedModel, transition);
    const { currentPath } = this.secretMountPath;

    const routeParams = (childResource: string) => [currentPath, childResource];

    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: currentPath, route: 'overview' },
      { label: 'Libraries', route: 'libraries' },
      ...ldapBreadcrumbs(resolvedModel.library.completeLibraryName, routeParams, libraryRoutes, true),
    ];
  }
}
