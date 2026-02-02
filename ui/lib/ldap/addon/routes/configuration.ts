/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type SecretMountPath from 'vault/services/secret-mount-path';
import type Transition from '@ember/routing/transition';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/vault/app-types';
import type RouterService from '@ember/routing/router-service';
import type ApiService from 'vault/services/api';
import type { LdapApplicationModel } from './application';

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapApplicationModel;
}

export default class LdapConfigurationRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service('app-router') declare readonly router: RouterService;

  setupController(controller: RouteController, resolvedModel: LdapApplicationModel, transition: Transition) {
    super.setupController(controller, resolvedModel, transition);

    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.secretsEngine.id, route: 'overview', model: resolvedModel.secretsEngine.id },
      { label: 'Configuration' },
    ];
  }

  afterModel(resolvedModel: LdapApplicationModel) {
    if (!resolvedModel.config) {
      this.router.transitionTo('vault.cluster.secrets.backend.ldap.configure');
    }
  }
}
