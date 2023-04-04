/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';

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
    return class CheckConfig extends SuperClass {
      shouldPromptConfig = false;

      async beforeModel() {
        super.beforeModel(...arguments);

        // When the engine is configured, it creates a default issuer.
        // If the issuers list is empty, we know it hasn't been configured
        return (
          this.store
            .query('pki/issuer', { backend: this.secretMountPath.currentPath })
            .then(() => (this.shouldPromptConfig = true))
            // this endpoint is unauthenticated, so we're not worried about permissions errors
            .catch(() => (this.shouldPromptConfig = false))
        );
      }
    };
  };
}
