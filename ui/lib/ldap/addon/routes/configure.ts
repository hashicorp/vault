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
import SecretsEngineResource from 'vault/resources/secrets/engine';

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
}

interface RouteModel {
  backendModel: SecretEngineModel;
  promptConfig: boolean;
  configModel: LdapConfigModel;
  engineDisplayData: SecretsEngineResource;
}

@withConfig('ldap/config')
export default class LdapConfigureRoute extends Route {
  @service declare readonly store: Store;
  @service declare readonly secretMountPath: SecretMountPath;

  declare configModel: LdapConfigModel;
  declare promptConfig: boolean;

  model() {
    const backend = this.secretMountPath.currentPath;
    const backendModel: SecretEngineModel = this.modelFor('application') as SecretEngineModel;

    return {
      backendModel,
      promptConfig: this.promptConfig,
      configModel: this.configModel || this.store.createRecord('ldap/config', { backend }),
    };
  }

  setupController(controller: RouteController, resolvedModel: RouteModel, transition: Transition) {
    super.setupController(controller, resolvedModel, transition);

    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backendModel.id, route: 'overview' },
      ...(this.promptConfig ? [] : [{ label: 'Configuration', route: 'configuration' }]),
      { label: 'Configure' },
    ];
  }
}
