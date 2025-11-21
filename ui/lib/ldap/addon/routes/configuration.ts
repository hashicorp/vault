/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { withConfig } from 'core/decorators/fetch-secrets-engine-config';

import type Store from '@ember-data/store';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type Transition from '@ember/routing/transition';
import type LdapConfigModel from 'vault/models/ldap/config';
import type SecretEngineModel from 'vault/models/secret-engine';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/vault/app-types';
import type AdapterError from '@ember-data/adapter/error';
import RouterService from '@ember/routing/router-service';

interface RouteModel {
  backendModel: SecretEngineModel;
  configModel: LdapConfigModel;
  configError: AdapterError;
}

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: RouteModel;
}

@withConfig('ldap/config')
export default class LdapConfigurationRoute extends Route {
  @service declare readonly store: Store;
  @service declare readonly secretMountPath: SecretMountPath;
  @service('app-router') declare readonly router: RouterService;

  declare configModel: LdapConfigModel;
  declare configError: AdapterError;
  declare promptConfig: boolean;

  model() {
    const backendModel: SecretEngineModel = this.modelFor('application') as SecretEngineModel;

    return {
      backendModel,
      promptConfig: this.promptConfig,
      configModel: this.configModel,
      configError: this.configError,
    };
  }

  setupController(controller: RouteController, resolvedModel: RouteModel, transition: Transition) {
    super.setupController(controller, resolvedModel, transition);

    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backendModel.id, route: 'overview', model: resolvedModel.backendModel.id },
      { label: 'Configuration' },
    ];
  }

  afterModel(resolvedModel: RouteModel) {
    if (!resolvedModel.configModel) {
      this.router.transitionTo('vault.cluster.secrets.backend.ldap.configure');
    }
  }
}
