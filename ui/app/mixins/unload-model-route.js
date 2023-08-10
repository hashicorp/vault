/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Mixin from '@ember/object/mixin';
import removeRecord from 'vault/utils/remove-record';

// removes Ember Data records from the cache when the model
// changes or you move away from the current route
export default Mixin.create({
  modelPath: 'model',
  unloadModel() {
    const { modelPath } = this;
    /* eslint-disable-next-line ember/no-controller-access-in-routes */
    const model = this.controller.get(modelPath);
    // error is thrown when you attempt to unload a record that is inFlight (isSaving)
    if (!model || !model.unloadRecord || model.isSaving) {
      return;
    }
    removeRecord(this.store, model);
    // it's important to unset the model on the controller since controllers are singletons
    this.controller.set(modelPath, null);
  },

  actions: {
    willTransition() {
      this.unloadModel();
      return true;
    },
  },
});
