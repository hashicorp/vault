/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import LdapLibraryForm from 'vault/forms/secrets/ldap/library';
import { ModelFrom } from 'vault/route';

import type SecretMountPath from 'vault/services/secret-mount-path';
import type Controller from '@ember/controller';
import type Transition from '@ember/routing/transition';
import type { Breadcrumb } from 'vault/vault/app-types';

export type LdapLibrariesCreateRouteModel = ModelFrom<LdapLibrariesCreateRoute>;

interface LdapLibrariesCreateController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapLibrariesCreateRouteModel;
}

export default class LdapLibrariesCreateRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;

  model() {
    const defaults = {
      ttl: '24h',
      max_ttl: '24h',
      disable_check_in_enforcement: 'Enabled',
    };
    return new LdapLibraryForm(defaults, { isNew: true });
  }

  setupController(
    controller: LdapLibrariesCreateController,
    resolvedModel: LdapLibrariesCreateRouteModel,
    transition: Transition
  ) {
    super.setupController(controller, resolvedModel, transition);

    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
      { label: 'Libraries', route: 'libraries' },
      { label: 'Create' },
    ];
  }
}
