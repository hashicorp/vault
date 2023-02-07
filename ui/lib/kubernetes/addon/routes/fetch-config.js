/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

/**
 * the overview, configure, configuration and roles routes all need to be aware of the config for the engine
 * if the user has not configured they are prompted to do so in each of the routes
 * this route can be extended so the check happens in the beforeModel hook since that may change what is returned from the model hook
 */

export default class KubernetesFetchConfigRoute extends Route {
  @service store;
  @service secretMountPath;

  configModel = null;

  async beforeModel() {
    const backend = this.secretMountPath.get();
    // check the store for record first
    this.configModel = this.store.peekRecord('kubernetes/config', backend);
    if (!this.configModel) {
      return this.store
        .queryRecord('kubernetes/config', { backend })
        .then((record) => {
          this.configModel = record;
        })
        .catch(() => {
          // it's ok! we don't need to transition to the error route
        });
    }
  }
}
