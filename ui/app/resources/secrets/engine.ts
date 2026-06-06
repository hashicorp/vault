/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import engineDisplayData from 'vault/helpers/engines-display-data';
import {
  supportedSecretBackends,
  SupportedSecretBackendsEnum,
} from 'vault/helpers/supported-secret-backends';
import { baseResourceFactory } from 'vault/resources/base-factory';
import { INTERNAL_ENGINE_TYPES, isAddonEngine } from 'vault/utils/all-engines-metadata';
import { getEffectiveEngineType } from 'vault/utils/external-plugin-helpers';

import type { Mount } from 'vault/mount';

export const SUPPORTS_RECOVERY = [
  SupportedSecretBackendsEnum.CUBBYHOLE,
  SupportedSecretBackendsEnum.KV, // only kv v1
  SupportedSecretBackendsEnum.DATABASE,
] as const;

export type RecoverySupportedEngines = (typeof SUPPORTS_RECOVERY)[number];

export default class SecretsEngineResource extends baseResourceFactory<Mount>() {
  id: string;

  #LIST_EXCLUDED_BACKENDS = INTERNAL_ENGINE_TYPES;

  constructor(data: Mount) {
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

  get effectiveEngineType() {
    return getEffectiveEngineType(this.engineType);
  }

  get isV2KV() {
    return (
      this.version === 2 &&
      (this.effectiveEngineType === SupportedSecretBackendsEnum.KV || this.effectiveEngineType === 'generic')
    );
  }

  get shouldIncludeInList() {
    return !this.#LIST_EXCLUDED_BACKENDS.includes(this.engineType);
  }

  get isSupportedBackend() {
    return supportedSecretBackends().includes(this.effectiveEngineType as SupportedSecretBackendsEnum);
  }

  get backendLink() {
    if (this.effectiveEngineType === 'database') {
      return 'vault.cluster.secrets.backend.overview';
    }
    if (isAddonEngine(this.effectiveEngineType, this.version)) {
      const engine = engineDisplayData(this.effectiveEngineType); // Use effective type to get proper metadata
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
    const { isConfigurable, configRoute } = engineDisplayData(this.effectiveEngineType);
    if (isConfigurable) {
      const route = configRoute || 'configuration.plugin-settings';
      return `vault.cluster.secrets.backend.${route}`;
    }
    return `vault.cluster.secrets.backend.configuration.general-settings`;
  }

  get localDisplay() {
    return this.local ? 'local' : 'replicated';
  }

  get supportsRecovery() {
    if (!SUPPORTS_RECOVERY.includes(this.effectiveEngineType as RecoverySupportedEngines)) {
      return false;
    }

    if (this.effectiveEngineType === SupportedSecretBackendsEnum.KV) {
      return !this.isV2KV;
    }

    return true;
  }
}
