/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@glimmer/component';

import type RouterService from '@ember/routing/router-service';
import type SecretsEngineResource from 'vault/resources/secrets/engine';

/**
 * @module ExternalTabs
 * The `ExternalTabs` is used to display external PKI tabs.
 */

interface Args {
  backend: SecretsEngineResource;
}

export default class ExternalTabs extends Component<Args> {
  @service('app-router') declare readonly router: RouterService;

  tabs = [
    { label: 'Overview', route: 'external.overview' },
    { label: 'Roles', route: 'external.roles' },
    { label: 'Recent orders', route: 'external.orders' },
    { label: 'DNS providers', route: 'external.dns-providers' },
    { label: 'ACME accounts', route: 'external.acme-accounts' },
  ];

  // Hide tabs for nested child routes
  get renderTabs() {
    switch (this.router.currentRouteName) {
      case 'vault.cluster.secrets.backend.pki.external.error':
        return false;
      case 'vault.cluster.secrets.backend.pki.external.roles.role.details':
        return false;
      case 'vault.cluster.secrets.backend.pki.external.roles.role.active-orders':
        return false;
      case 'vault.cluster.secrets.backend.pki.external.orders.order.details':
        return false;
      default:
        return true;
    }
  }
}
