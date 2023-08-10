/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';

/**
 * the overview, configure, configuration and roles routes all need to be aware of the config for the engine
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
    return class FetchConfig extends SuperClass {
      configModel = null;
      configError = null;
      promptConfig = false;

      async beforeModel() {
        super.beforeModel(...arguments);

        const backend = this.secretMountPath.get();
        // check the store for record first
        this.configModel = this.store.peekRecord('kubernetes/config', backend);
        if (!this.configModel) {
          return this.store
            .queryRecord('kubernetes/config', { backend })
            .then((record) => {
              this.configModel = record;
              this.promptConfig = false;
            })
            .catch((error) => {
              // we need to ignore if the user does not have permission or other failures so as to not block the other operations
              if (error.httpStatus === 404) {
                this.promptConfig = true;
              } else {
                // not considering 404 an error since it triggers the cta
                // this error is thrown in the configuration route so we can display the error in the view
                this.configError = error;
              }
            });
        } else {
          this.promptConfig = false;
        }
      }
    };
  };
}
