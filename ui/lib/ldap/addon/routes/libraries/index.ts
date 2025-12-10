/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import LdapLibrariesRoute from '../libraries';
import { service } from '@ember/service';

import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type Transition from '@ember/routing/transition';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/vault/app-types';
import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type { LdapApplicationModel } from '../application';
import type { LdapLibrary } from 'vault/secrets/ldap';
import type CapabilitiesService from 'vault/services/capabilities';

interface LdapLibrariesRouteModel {
  secretsEngine: SecretsEngineResource;
  promptConfig: boolean;
  libraries: LdapLibrary[];
}
interface LdapLibrariesController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapLibrariesRouteModel;
}

export default class LdapLibrariesIndexRoute extends LdapLibrariesRoute {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly capabilities: CapabilitiesService;

  async model() {
    const { secretsEngine, promptConfig } = this.modelFor('application') as LdapApplicationModel;
    const { libraries, capabilities } = await this.fetchLibrariesAndCapabilities();

    return {
      secretsEngine,
      promptConfig,
      libraries,
      capabilities,
    };
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
