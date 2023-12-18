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
import type LdapRoleModel from 'vault/models/ldap/role';
import type LdapLibraryModel from 'vault/models/ldap/library';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/vault/app-types';
import { LdapLibraryAccountStatus } from 'vault/vault/adapters/ldap/library';

interface LdapOverviewController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
}
interface LdapOverviewRouteModel {
  backendModel: SecretEngineModel;
  promptConfig: boolean;
  roles: Array<LdapRoleModel>;
  libraries: Array<LdapLibraryModel>;
  librariesStatus: Array<LdapLibraryAccountStatus>;
}

@withConfig('ldap/config')
export default class LdapOverviewRoute extends Route {
  @service declare readonly store: Store;
  @service declare readonly secretMountPath: SecretMountPath;

  declare promptConfig: boolean;

  async fetchLibrariesStatus(libraries: Array<LdapLibraryModel>): Promise<Array<LdapLibraryAccountStatus>> {
    const allStatuses: Array<LdapLibraryAccountStatus> = [];

    for (const library of libraries) {
      try {
        const statuses = await library.fetchStatus();
        allStatuses.push(...statuses);
      } catch (error) {
        // suppressing error
      }
    }
    return allStatuses;
  }

  async fetchLibraries(backend: string) {
    return this.store.query('ldap/library', { backend }).catch(() => []);
  }

  async model() {
    const backend = this.secretMountPath.currentPath;
    const libraries = await this.fetchLibraries(backend);
    return hash({
      promptConfig: this.promptConfig,
      backendModel: this.modelFor('application'),
      roles: this.store.query('ldap/role', { backend }).catch(() => []),
      libraries,
      librariesStatus: this.fetchLibrariesStatus(libraries as Array<LdapLibraryModel>),
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
