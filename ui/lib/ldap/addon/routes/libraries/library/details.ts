/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';

import type LdapLibraryModel from 'vault/models/ldap/library';
import type Controller from '@ember/controller';
import type Transition from '@ember/routing/transition';
import type { Breadcrumb } from 'vault/vault/app-types';
import { ldapBreadcrumbs, libraryRoutes } from 'ldap/utils/ldap-breadcrumbs';

interface LdapLibraryDetailsController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapLibraryModel;
}

export default class LdapLibraryDetailsRoute extends Route {
  setupController(
    controller: LdapLibraryDetailsController,
    resolvedModel: LdapLibraryModel,
    transition: Transition
  ) {
    super.setupController(controller, resolvedModel, transition);

    const routeParams = (childResource: string) => {
      return [resolvedModel.backend, childResource];
    };

    controller.breadcrumbs = [
      { label: resolvedModel.backend, route: 'overview' },
      { label: 'Libraries', route: 'libraries' },
      ...ldapBreadcrumbs(resolvedModel.name, routeParams, libraryRoutes, true),
    ];
  }
}
