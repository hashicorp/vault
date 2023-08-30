/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { inject as service } from '@ember/service';
import Controller from '@ember/controller';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
import { allEngines } from 'vault/helpers/mountable-secret-engines';
import { action } from '@ember/object';

const SUPPORTED_BACKENDS = supportedSecretBackends();

export default class MountSecretBackendController extends Controller {
  @service router;

  @action
  onMountSuccess(type, path, useEngineRoute = false) {
    let transition;
    if (SUPPORTED_BACKENDS.includes(type)) {
      const engineInfo = allEngines().findBy('type', type);
      if (useEngineRoute) {
        transition = this.router.transitionTo(
          `vault.cluster.secrets.backend.${engineInfo.engineRoute}`,
          path
        );
      } else {
        const queryParams = engineInfo?.routeQueryParams || {};
        transition = this.router.transitionTo('vault.cluster.secrets.backend.index', path, { queryParams });
      }
    } else {
      transition = this.router.transitionTo('vault.cluster.secrets.backends');
    }
    return transition.followRedirects();
  }
}
