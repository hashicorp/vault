/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';

import type Store from '@ember-data/store';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type Transition from '@ember/routing/transition';
import type Model from '@ember-data/model';
import type AdapterError from 'ember-data/adapter'; // eslint-disable-line ember/use-ember-data-rfc-395-imports

/**
 * for use in routes that need to be aware of the config for a secrets engine
 * if the user has not configured they are prompted to do so in each of the routes
 * decorate the necessary routes to perform the check in the beforeModel hook since that may change what is returned for the model
 */

interface BaseRoute extends Route {
  store: Store;
  secretMountPath: SecretMountPath;
}

export function withConfig(modelName: string) {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  return function <RouteClass extends new (...args: any[]) => BaseRoute>(SuperClass: RouteClass) {
    if (!Object.prototype.isPrototypeOf.call(Route, SuperClass)) {
      // eslint-disable-next-line
      console.error(
        'withConfig decorator must be used on an instance of Ember Route class. Decorator not applied to returned class'
      );
      return SuperClass;
    }

    return class FetchSecretsEngineConfig extends SuperClass {
      configModel: Model | null = null;
      configError: AdapterError | null = null;
      promptConfig = false;

      async beforeModel(transition: Transition) {
        super.beforeModel(transition);

        const backend = this.secretMountPath.currentPath;
        // check the store for record first
        this.configModel = this.store.peekRecord(modelName, backend);
        if (!this.configModel) {
          return this.store
            .queryRecord(modelName, { backend })
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
