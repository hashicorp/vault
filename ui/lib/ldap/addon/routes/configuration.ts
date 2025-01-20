/**
 * Copyright (c) HashiCorp, Inc.
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
import type AdapterError from 'ember-data/adapter'; // eslint-disable-line ember/use-ember-data-rfc-395-imports

interface LdapConfigurationRouteModel {
  backendModel: SecretEngineModel;
  configModel: LdapConfigModel;
  configError: AdapterError;
}
interface LdapConfigurationController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapConfigurationRouteModel;
}

@withConfig('ldap/config')
export default class LdapConfigurationRoute extends Route {
  @service declare readonly store: Store;
  @service declare readonly secretMountPath: SecretMountPath;

  declare configModel: LdapConfigModel;
  declare configError: AdapterError;

  model() {
    return {
      backendModel: this.modelFor('application'),
      configModel: this.configModel,
      configError: this.configError,
    };
  }

  setupController(
    controller: LdapConfigurationController,
    resolvedModel: LdapConfigurationRouteModel,
    transition: Transition
  ) {
    super.setupController(controller, resolvedModel, transition);

    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backendModel.id, route: 'overview', model: resolvedModel.backend },
      { label: 'Configuration' },
    ];
  }
}
