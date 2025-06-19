/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { baseResourceFactory } from 'vault/resources/base-factory';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
import { isAddonEngine } from 'vault/utils/all-engines-metadata';
import engineDisplayData from 'vault/helpers/engines-display-data';

import type { SecretsEngine } from 'vault/secrets/engine';

export default class SecretsEngineResource extends baseResourceFactory<SecretsEngine>() {
  id: string;

  #LIST_EXCLUDED_BACKENDS = ['system', 'identity'];

  constructor(data: SecretsEngine) {
    super(data);
    // strip trailing slash from path for id since it is used in routing
    this.id = data.path.replace(/\/$/, '');
  }

  get version() {
    const { version } = this.options || {};
    return version ? Number(version) : 1;
  }

  get engineType() {
    return (this.type || '').replace(/^ns_/, '');
  }

  get icon() {
    const engineData = engineDisplayData(this.engineType);

    return engineData?.glyph || 'lock';
  }

  get isV2KV() {
    return this.version === 2 && (this.engineType === 'kv' || this.engineType === 'generic');
  }

  get shouldIncludeInList() {
    return !this.#LIST_EXCLUDED_BACKENDS.includes(this.engineType);
  }

  get isSupportedBackend() {
    return supportedSecretBackends().includes(this.engineType);
  }

  get backendLink() {
    if (this.engineType === 'database') {
      return 'vault.cluster.secrets.backend.overview';
    }
    if (isAddonEngine(this.engineType, this.version)) {
      const engine = engineDisplayData(this.engineType);
      if (engine?.engineRoute) {
        return `vault.cluster.secrets.backend.${engine.engineRoute}`;
      }
    }
    if (this.isV2KV) {
      // if it's KV v2 but not registered as an addon, it's type generic
      return 'vault.cluster.secrets.backend.kv.list';
    }
    return `vault.cluster.secrets.backend.list-root`;
  }

  get backendConfigurationLink() {
    if (isAddonEngine(this.engineType, this.version)) {
      return `vault.cluster.secrets.backend.${this.engineType}.configuration`;
    }
    return `vault.cluster.secrets.backend.configuration`;
  }

  get localDisplay() {
    return this.local ? 'local' : 'replicated';
  }
}
