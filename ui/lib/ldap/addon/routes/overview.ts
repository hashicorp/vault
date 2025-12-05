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
import type LdapRoleModel from 'vault/models/ldap/role';
import type LdapLibraryModel from 'vault/models/ldap/library';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/app-types';
import type { LdapLibraryAccountStatus } from 'vault/adapters/ldap/library';
import type { LdapApplicationModel } from './application';
import type SecretsEngineResource from 'vault/resources/secrets/engine';

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
}
interface RouteModel {
  secretsEngine: SecretsEngineResource;
  promptConfig: boolean;
  roles: Array<LdapRoleModel>;
  libraries: Array<LdapLibraryModel>;
  librariesStatus: Array<LdapLibraryAccountStatus>;
}

export default class LdapOverviewRoute extends Route {
  @service declare readonly store: Store;
  @service declare readonly secretMountPath: SecretMountPath;

  declare promptConfig: boolean;

  async model() {
    const { promptConfig, secretsEngine } = this.modelFor('application') as LdapApplicationModel;
    const backend = this.secretMountPath.currentPath;
    return hash({
      promptConfig,
      secretsEngine,
      roles: this.store.query('ldap/role', { backend }).catch(() => []),
    });
  }

  setupController(controller: RouteController, resolvedModel: RouteModel, transition: Transition) {
    super.setupController(controller, resolvedModel, transition);

    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.secretsEngine.id },
    ];
  }
}
