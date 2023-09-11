/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

import type Transition from '@ember/routing/transition';
import type AdapterError from 'ember-data/adapter'; // eslint-disable-line ember/use-ember-data-rfc-395-imports
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
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
    ];
    controller.backend = this.modelFor('application') as SecretEngineModel;
  }
}
