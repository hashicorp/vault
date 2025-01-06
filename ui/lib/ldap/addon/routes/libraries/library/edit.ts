/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';

import type LdapLibraryModel from 'vault/models/ldap/library';
import type Controller from '@ember/controller';
import type Transition from '@ember/routing/transition';
import type { Breadcrumb } from 'vault/vault/app-types';

interface LdapLibraryEditController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapLibraryModel;
}

export default class LdapLibraryEditRoute extends Route {
  setupController(
    controller: LdapLibraryEditController,
    resolvedModel: LdapLibraryModel,
    transition: Transition
  ) {
    super.setupController(controller, resolvedModel, transition);

    controller.breadcrumbs = [
      { label: resolvedModel.backend, route: 'overview' },
      { label: 'Libraries', route: 'libraries' },
      { label: resolvedModel.name, route: 'libraries.library.details' },
      { label: 'Edit' },
    ];
  }
}
