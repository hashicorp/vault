/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

import type Store from '@ember-data/store';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type LdapLibraryModel from 'vault/models/ldap/library';
import type Controller from '@ember/controller';
import type Transition from '@ember/routing/transition';
import type { Breadcrumb } from 'vault/vault/app-types';

interface LdapLibrariesCreateController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapLibraryModel;
}

export default class LdapLibrariesCreateRoute extends Route {
  @service declare readonly store: Store;
  @service declare readonly secretMountPath: SecretMountPath;

  model() {
    const backend = this.secretMountPath.currentPath;
    return this.store.createRecord('ldap/library', { backend });
  }

  setupController(
    controller: LdapLibrariesCreateController,
    resolvedModel: LdapLibraryModel,
    transition: Transition
  ) {
    super.setupController(controller, resolvedModel, transition);

    controller.breadcrumbs = [
      { label: resolvedModel.backend, route: 'overview' },
      { label: 'libraries', route: 'libraries' },
      { label: 'create' },
    ];
  }
}
