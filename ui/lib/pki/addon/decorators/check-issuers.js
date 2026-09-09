/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { SecretsApiPkiListIssuersListEnum } from '@hashicorp/vault-client-typescript';

/**
 * the overview, roles, issuers, certificates, and key routes all need to be aware of the whether there is a config for the engine
 * if the user has not configured they are prompted to do so in each of the routes
 * decorate the necessary routes to perform the check in the beforeModel hook since that may change what is returned for the model
 */

export function withConfig() {
  return function decorator(SuperClass) {
    if (!Object.prototype.isPrototypeOf.call(Route, SuperClass)) {
      // eslint-disable-next-line
      console.error(
        'withConfig decorator must be used on an instance of ember Route class. Decorator not applied to returned class'
      );
      return SuperClass;
    }
    class CheckConfig extends SuperClass {
      @service secretMountPath;
      @service api;

      pkiMountHasConfig = false;

      async beforeModel() {
        super.beforeModel(...arguments);

        // When the engine is configured, it creates a default issuer.
        // If the issuers list is empty, we know it hasn't been configured
        try {
          await this.api.secrets.pkiListIssuers(
            this.secretMountPath.currentPath,
            SecretsApiPkiListIssuersListEnum.TRUE
          );
          this.pkiMountHasConfig = true;
        } catch (e) {
          this.pkiMountHasConfig = false;
        }
      }
    }
    return CheckConfig;
  };
}
