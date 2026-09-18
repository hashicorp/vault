/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type ApiService from 'vault/services/api';
import type CapabilitiesService from 'vault/services/capabilities';

export default class ConfigRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly capabilities: CapabilitiesService;

  async model() {
    const capabilities = await this.capabilities.for('clientsConfig');
    // sys/internal/counters/config is root-namespace only; explicitly clear the namespace
    // header or this 404s when viewed from a child namespace (VAULT-45847)
    const config = await this.api.sys.internalClientActivityReadConfiguration(
      undefined,
      undefined,
      undefined,
      this.api.buildHeaders({ namespace: '' })
    );
    return { capabilities, config };
  }
}
