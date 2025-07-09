/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Controller from '@ember/controller';
import { set } from '@ember/object';

export default Controller.extend({
  store: service(),
  loading: false,
  emptyData: '{\n}',
  actions: {
    sign() {
      this.set('loading', true);
      this.model.save().finally(() => {
        this.set('loading', false);
      });
    },

    editorUpdated(attr, val) {
      // wont set invalid JSON to the model
      try {
        set(this.model, attr, JSON.parse(val));
      } catch {
        // linting is handled by the component
      }
    },

    updateTtl(path, val) {
      const model = this.model;
      const valueToSet = val.enabled === true ? `${val.seconds}s` : undefined;
      set(model, path, valueToSet);
    },

    newModel() {
      const model = this.model;
      const roleModel = model.role;
      model.unloadRecord();
      const newModel = this.store.createRecord('ssh-sign', {
        role: roleModel,
        id: `${roleModel.backend}-${roleModel.name}`,
      });
      this.set('model', newModel);
    },
  },
});
