/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { service } from '@ember/service';
import { action } from '@ember/object';
import {
  supportedSecretBackends,
  SupportedSecretBackendsEnum,
} from 'vault/helpers/supported-secret-backends';
import engineDisplayData from 'vault/helpers/engines-display-data';
import { getEffectiveEngineType } from 'vault/utils/external-plugin-helpers';
import type SecretsEngineForm from 'vault/forms/secrets/engine';
import type Router from '@ember/routing/router';
import type { EngineVersionInfo } from 'vault/utils/plugin-catalog-helpers';

const SUPPORTED_BACKENDS = supportedSecretBackends();

export default class VaultClusterSecretsEnableCreateController extends Controller {
  @service declare router: Router;

  declare model: {
    form: SecretsEngineForm;
    availableVersions: EngineVersionInfo[];
  };

  @action
  onMountSuccess(type: string, path: string, useEngineRoute = false) {
    let transition;
    const engineInfo = engineDisplayData(type);
    const effectiveType = getEffectiveEngineType(type);

    if (engineInfo && SUPPORTED_BACKENDS.includes(effectiveType as SupportedSecretBackendsEnum)) {
      if (useEngineRoute && engineInfo.engineRoute) {
        transition = this.router.transitionTo(
          `vault.cluster.secrets.backend.${engineInfo.engineRoute}`,
          path
        );
      } else {
        // For keymgmt, we need to land on provider tab by default using query params
        const queryParams = effectiveType === 'keymgmt' ? { tab: 'provider' } : {};
        transition = this.router.transitionTo('vault.cluster.secrets.backend.index', path, { queryParams });
      }
    } else {
      transition = this.router.transitionTo('vault.cluster.secrets.backends');
    }
    return transition?.followRedirects();
  }
}
