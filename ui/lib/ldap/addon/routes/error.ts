/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type Transition from '@ember/routing/transition';
import type AdapterError from '@ember-data/adapter/error';
import type SecretEngineModel from 'vault/models/secret-engine';
import type { Breadcrumb } from 'vault/vault/app-types';
import type Controller from '@ember/controller';
import type SecretMountPath from 'vault/services/secret-mount-path';

interface LdapErrorController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  backend: SecretEngineModel;
}

export default class LdapErrorRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;

  setupController(controller: LdapErrorController, resolvedModel: AdapterError, transition: Transition) {
    super.setupController(controller, resolvedModel, transition);
    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
    ];
    controller.backend = this.modelFor('application') as SecretEngineModel;
  }
}
