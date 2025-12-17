/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import LdapLibrariesRoute from '../libraries';
import { service } from '@ember/service';
import { hash } from 'rsvp';
import { ldapBreadcrumbs, libraryRoutes } from 'ldap/utils/ldap-breadcrumbs';

import type SecretMountPath from 'vault/services/secret-mount-path';
import type Transition from '@ember/routing/transition';
import type { LdapLibrary } from 'vault/secrets/ldap';
import type LdapLibrariesSubdirectoryController from 'ldap/controllers/libraries/subdirectory';
import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type { LdapApplicationModel } from '../application';

interface RouteModel {
  secretsEngine: SecretsEngineResource;
  path_to_library: string;
  libraries: Array<LdapLibrary>;
}

type RouteController = LdapLibrariesSubdirectoryController;
interface RouteParams {
  path_to_library?: string;
}

export default class LdapLibrariesSubdirectoryRoute extends LdapLibrariesRoute {
  @service declare readonly secretMountPath: SecretMountPath;

  async model(params: RouteParams) {
    const { secretsEngine } = this.modelFor('application') as LdapApplicationModel;
    const { path_to_library } = params;

    // Ensure path_to_library has trailing slash for proper API calls and model construction
    const normalizedPath = path_to_library?.endsWith('/') ? path_to_library : `${path_to_library}/`;
    const { libraries, capabilities } = await this.fetchLibrariesAndCapabilities(normalizedPath);

    return hash({
      secretsEngine,
      path_to_library: normalizedPath,
      libraries,
      capabilities,
    });
  }

  setupController(controller: RouteController, resolvedModel: RouteModel, transition: Transition) {
    super.setupController(controller, resolvedModel, transition);

    const routeParams = (childResource: string) => {
      return [resolvedModel.secretsEngine.id, childResource];
    };

    const currentLevelPath = resolvedModel.path_to_library;

    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.secretsEngine.id, route: 'overview' },
      { label: 'Libraries', route: 'libraries' },
      ...ldapBreadcrumbs(currentLevelPath, routeParams, libraryRoutes, true),
    ];
  }
}
