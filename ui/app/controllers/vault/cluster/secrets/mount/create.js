/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Controller from '@ember/controller';
import { action } from '@ember/object';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
import engineDisplayData from 'vault/helpers/engines-display-data';

const SUPPORTED_BACKENDS = supportedSecretBackends();

export default class MountSecretBackendController extends Controller {
  @service router;

  @action
  onMountSuccess(type, path, useEngineRoute = false) {
    let transition;
    if (SUPPORTED_BACKENDS.includes(type)) {
      const engineInfo = engineDisplayData(type);
      if (useEngineRoute) {
        transition = this.router.transitionTo(
          `vault.cluster.secrets.backend.${engineInfo.engineRoute}`,
          path
        );
      } else {
        // For keymgmt, we need to land on provider tab by default using query params
        const queryParams = engineInfo.type === 'keymgmt' ? { tab: 'provider' } : {};
        transition = this.router.transitionTo('vault.cluster.secrets.backend.index', path, { queryParams });
      }
    } else {
      transition = this.router.transitionTo('vault.cluster.secrets.backends');
    }
    return transition.followRedirects();
  }
}
