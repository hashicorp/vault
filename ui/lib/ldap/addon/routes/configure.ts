/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import LdapConfigForm from 'vault/forms/secrets/ldap/config';

import type SecretMountPath from 'vault/services/secret-mount-path';
import type Transition from '@ember/routing/transition';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/vault/app-types';
import type { LdapApplicationModel } from './application';
import type { ModelFrom } from 'vault/route';

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
}

export type LdapConfigureModel = ModelFrom<LdapConfigureRoute>;

export default class LdapConfigureRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;

  model() {
    const { secretsEngine, promptConfig, config } = this.modelFor('application') as LdapApplicationModel;
    const form = new LdapConfigForm(config, { isNew: !config });
    return {
      secretsEngine,
      promptConfig,
      form,
    };
  }

  setupController(controller: RouteController, resolvedModel: LdapConfigureModel, transition: Transition) {
    super.setupController(controller, resolvedModel, transition);

    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: resolvedModel.secretsEngine.id, route: 'overview' },
      ...(resolvedModel.promptConfig ? [] : [{ label: 'Configuration', route: 'configuration' }]),
      { label: 'Configure' },
    ];
  }
}
