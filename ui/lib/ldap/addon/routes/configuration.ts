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
import type { LdapConfigureRequest } from '@hashicorp/vault-client-typescript';
import type { ModelFrom } from 'vault/vault/route';

export type LdapConfigurationModel = ModelFrom<LdapConfigurationRoute>;

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapConfigurationModel;
}

export default class LdapConfigurationRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service('app-router') declare readonly router: RouterService;

  async model() {
    const { secretsEngine } = this.modelFor('application') as LdapApplicationModel;
    let config: LdapConfigureRequest | undefined;
    let promptConfig = false;
    let configError: unknown;
    // check if engine is configured
    // child routes will handle prompting for configuration if needed
    try {
      const { data } = await this.api.secrets.ldapReadConfiguration(secretsEngine.id);
      config = data as LdapConfigureRequest;
    } catch (error) {
      const { response, status } = await this.api.parseError(error);
      // not considering 404 an error since it triggers the cta
      if (status === 404) {
        promptConfig = true;
      } else {
        // ignore if the user does not have permission or other failures so as to not block the other operations
        // this error is thrown in the configuration route so we can display the error in the view
        configError = response;
      }
    }
    return {
      secretsEngine,
      config,
      configError,
      promptConfig,
    };
  }

  setupController(
    controller: RouteController,
    resolvedModel: LdapConfigurationModel,
    transition: Transition
  ) {
    super.setupController(controller, resolvedModel, transition);

    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: resolvedModel.secretsEngine.id, route: 'overview', model: resolvedModel.secretsEngine.id },
      { label: 'Configuration' },
    ];
  }

  afterModel(resolvedModel: LdapConfigurationModel) {
    if (!resolvedModel.config) {
      this.router.transitionTo('vault.cluster.secrets.backend.ldap.configure');
    }
  }
}
