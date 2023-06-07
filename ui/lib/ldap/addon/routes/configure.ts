/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { withConfig } from 'core/decorators/fetch-secret-config';

import type Store from '@ember-data/store';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type Transition from '@ember/routing/transition';
import type LdapConfigModel from 'vault/models/ldap/config';

export interface LdapConfigureRouteModel {
  backend: LdapConfigModel;
}

@withConfig('ldap/config')
export default class LdapConfigureRoute extends Route {
  @service declare readonly store: Store;
  @service declare readonly secretMountPath: SecretMountPath;

  declare configModel: LdapConfigModel;

  async model() {
    const backend = this.secretMountPath.currentPath;
    return this.configModel || this.store.createRecord('kubernetes/config', { backend });
  }

  setupController(controller: any, resolvedModel: LdapConfigureRouteModel, transition: Transition) {
    super.setupController(controller, resolvedModel, transition);

    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'overview' },
      { label: 'configure' },
    ];
  }
}
