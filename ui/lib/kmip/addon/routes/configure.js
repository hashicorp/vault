/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default Route.extend({
  store: service(),
  secretMountPath: service(),
  pathHelp: service(),
  beforeModel() {
    return this.pathHelp.getNewModel('kmip/config', this.secretMountPath.currentPath);
  },
  model() {
    return this.store.findRecord('kmip/config', this.secretMountPath.currentPath).catch((err) => {
      if (err.httpStatus === 404) {
        const model = this.store.createRecord('kmip/config');
        model.set('id', this.secretMountPath.currentPath);
        return model;
      } else {
        throw err;
      }
    });
  },
});
