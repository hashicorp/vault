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
import type LdapLibraryModel from 'vault/models/ldap/library';
import type SecretEngineModel from 'vault/models/secret-engine';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/vault/app-types';

interface LdapLibrariesRouteModel {
  backendModel: SecretEngineModel;
  promptConfig: boolean;
  libraries: Array<LdapLibraryModel>;
}
interface LdapLibrariesController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapLibrariesRouteModel;
}

@withConfig('ldap/config')
export default class LdapLibrariesRoute extends Route {
  @service declare readonly store: Store;
  @service declare readonly secretMountPath: SecretMountPath;

  declare promptConfig: boolean;

  model() {
    const backendModel = this.modelFor('application') as SecretEngineModel;
    return hash({
      backendModel,
      promptConfig: this.promptConfig,
      libraries: this.store.query('ldap/library', { backend: backendModel.id }),
    });
  }

  setupController(
    controller: LdapLibrariesController,
    resolvedModel: LdapLibrariesRouteModel,
    transition: Transition
  ) {
    super.setupController(controller, resolvedModel, transition);

    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backendModel.id },
    ];
  }
}
