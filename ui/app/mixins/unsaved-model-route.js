/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Mixin from '@ember/object/mixin';
import Ember from 'ember';

// this mixin relies on `unload-model-route` also being used
export default Mixin.create({
  actions: {
    willTransition(transition) {
      const model = this.controller.get('model');
      if (!model) {
        return true;
      }
      if (model.hasDirtyAttributes) {
        if (
          Ember.testing ||
          window.confirm(
            'You have unsaved changes. Navigating away will discard these changes. Are you sure you want to discard your changes?'
          )
        ) {
          this.unloadModel();
          return true;
        } else {
          transition.abort();
          return false;
        }
      }
      return true;
    },
  },
});
