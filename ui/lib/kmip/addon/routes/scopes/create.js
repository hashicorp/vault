/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default Route.extend({
  store: service(),
  secretMountPath: service(),
  beforeModel() {
    this.store.unloadAll('kmip/scope');
  },
  model() {
    const model = this.store.createRecord('kmip/scope', {
      backend: this.secretMountPath.currentPath,
    });
    return model;
  },
});
