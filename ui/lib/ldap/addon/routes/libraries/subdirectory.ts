/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';

import type Store from '@ember-data/store';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type Transition from '@ember/routing/transition';
import type LdapLibraryModel from 'vault/models/ldap/library';
import type SecretEngineModel from 'vault/models/secret-engine';

import { ldapBreadcrumbs, libraryRoutes } from 'ldap/utils/ldap-breadcrumbs';
import type LdapLibrariesSubdirectoryController from 'ldap/controllers/libraries/subdirectory';

interface RouteModel {
  backendModel: SecretEngineModel;
  path_to_library: string;
  libraries: Array<LdapLibraryModel>;
}

type RouteController = LdapLibrariesSubdirectoryController;
interface RouteParams {
  path_to_library?: string;
}

export default class LdapLibrariesSubdirectoryRoute extends Route {
  @service declare readonly store: Store;
  @service declare readonly secretMountPath: SecretMountPath;

  model(params: RouteParams) {
    const backendModel = this.modelFor('application') as SecretEngineModel;
    const { path_to_library } = params;

    // Ensure path_to_library has trailing slash for proper API calls and model construction
    const normalizedPath = path_to_library?.endsWith('/') ? path_to_library : `${path_to_library}/`;

    return hash({
      backendModel,
      path_to_library: normalizedPath,
      libraries: this.store.query('ldap/library', {
        backend: backendModel.id,
        path_to_library: normalizedPath,
      }),
    });
  }

  setupController(controller: RouteController, resolvedModel: RouteModel, transition: Transition) {
    super.setupController(controller, resolvedModel, transition);

    const routeParams = (childResource: string) => {
      return [resolvedModel.backendModel.id, childResource];
    };

    const currentLevelPath = resolvedModel.path_to_library;

    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backendModel.id, route: 'overview' },
      { label: 'Libraries', route: 'libraries' },
      ...ldapBreadcrumbs(currentLevelPath, routeParams, libraryRoutes, true),
    ];
  }
}
