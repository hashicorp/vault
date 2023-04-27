/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import AdapterError from '@ember-data/adapter/error';
import { set } from '@ember/object';
import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
const CONFIGURABLE_BACKEND_TYPES = ['aws', 'ssh'];

export default Route.extend({
  store: service(),

  model() {
    const { backend } = this.paramsFor(this.routeName);
    return this.store.query('secret-engine', { path: backend }).then((modelList) => {
      const model = modelList && modelList.get('firstObject');
      if (!model || !CONFIGURABLE_BACKEND_TYPES.includes(model.get('type'))) {
        const error = new AdapterError();
        set(error, 'httpStatus', 404);
        throw error;
      }
      return this.store.findRecord('secret-engine', backend).then(
        () => {
          return model;
        },
        () => {
          return model;
        }
      );
    });
  },

  afterModel(model) {
    const type = model.get('type');

    if (type === 'aws') {
      return this.store
        .queryRecord('secret-engine', {
          backend: model.id,
          type,
        })
        .then(
          () => model,
          () => model
        );
    }
    return model;
  },

  setupController(controller, model) {
    if (model.get('publicKey')) {
      controller.set('configured', true);
    }
    return this._super(...arguments);
  },

  resetController(controller, isExiting) {
    if (isExiting) {
      controller.reset();
    }
  },

  actions: {
    refreshRoute() {
      this.refresh();
    },
  },
});
