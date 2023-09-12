/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { withConfig } from 'core/decorators/fetch-secrets-engine-config';

import type Store from '@ember-data/store';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type Transition from '@ember/routing/transition';
import type LdapConfigModel from 'vault/models/ldap/config';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/vault/app-types';

interface LdapConfigureController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
}

@withConfig('ldap/config')
export default class LdapConfigureRoute extends Route {
  @service declare readonly store: Store;
  @service declare readonly secretMountPath: SecretMountPath;

  declare configModel: LdapConfigModel;

  model() {
    const backend = this.secretMountPath.currentPath;
    return this.configModel || this.store.createRecord('ldap/config', { backend });
  }

  setupController(
    controller: LdapConfigureController,
    resolvedModel: LdapConfigModel,
    transition: Transition
  ) {
    super.setupController(controller, resolvedModel, transition);

    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'overview' },
      { label: 'Configure' },
    ];
  }
}
