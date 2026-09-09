/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { SupportedSecretBackendsEnum } from 'vault/helpers/supported-secret-backends';

export default class PkiRoute extends Route {
  @service('app-router') router;

  redirect(model) {
    if (model.type === SupportedSecretBackendsEnum.PKI_EXTERNAL) {
      return this.router.transitionTo('vault.cluster.secrets.backend.pki.external.overview');
    }
    return this.router.transitionTo('vault.cluster.secrets.backend.pki.overview');
  }
}
