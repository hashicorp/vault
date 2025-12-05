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
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/vault/app-types';
import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type { LdapApplicationModel } from '../application';

interface LdapLibrariesRouteModel {
  secretsEngine: SecretsEngineResource;
  promptConfig: boolean;
  libraries: Array<LdapLibraryModel>;
}
interface LdapLibrariesController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapLibrariesRouteModel;
}

export default class LdapLibrariesRoute extends Route {
  @service declare readonly store: Store;
  @service declare readonly secretMountPath: SecretMountPath;

  model() {
    const { secretsEngine, promptConfig } = this.modelFor('application') as LdapApplicationModel;
    return hash({
      secretsEngine,
      promptConfig,
      libraries: this.store.query('ldap/library', { backend: secretsEngine.id }),
    });
  }

  setupController(
    controller: LdapLibrariesController,
    resolvedModel: LdapLibrariesRouteModel,
    transition: Transition
  ) {
    super.setupController(controller, resolvedModel, transition);

    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.secretsEngine.id, route: 'overview', model: resolvedModel.secretsEngine.id },
      { label: 'Libraries' },
    ];
  }
}
