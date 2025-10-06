/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
import engineDisplayData from 'vault/helpers/engines-display-data';
import type SecretsEngineForm from 'vault/forms/secrets/engine';
import type Router from '@ember/routing/router';

const SUPPORTED_BACKENDS = supportedSecretBackends();

export default class VaultClusterSecretsMountsCreateController extends Controller {
  @service declare router: Router;

  declare model: SecretsEngineForm;

  @action
  onMountSuccess(type: string, path: string, useEngineRoute = false) {
    let transition;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    if (SUPPORTED_BACKENDS.includes(type as any)) {
      const engineInfo = engineDisplayData(type);
      if (engineInfo && useEngineRoute) {
        transition = this.router.transitionTo(
          `vault.cluster.secrets.backend.${engineInfo.engineRoute}`,
          path
        );
      } else if (engineInfo) {
        // For keymgmt, we need to land on provider tab by default using query params
        const queryParams = engineInfo.type === 'keymgmt' ? { tab: 'provider' } : {};
        transition = this.router.transitionTo('vault.cluster.secrets.backend.index', path, { queryParams });
      }
    } else {
      transition = this.router.transitionTo('vault.cluster.secrets.backends');
    }
    return transition?.followRedirects();
  }
}
