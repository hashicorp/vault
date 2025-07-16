/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import UnsavedModelRoute from 'vault/mixins/unsaved-model-route';
import { singularize } from 'ember-inflector';
import { service } from '@ember/service';

export default Route.extend(UnsavedModelRoute, {
  store: service(),

  model() {
    const { item_type: itemType } = this.paramsFor('vault.cluster.access.method.item');
    const methodModel = this.modelFor('vault.cluster.access.method');
    const { type } = methodModel;
    const { path: method } = this.paramsFor('vault.cluster.access.method');
    const modelType = `generated-${singularize(itemType)}-${type}`;
    return this.store.createRecord(modelType, {
      itemType,
      method,
      adapterOptions: { path: `${method}/${itemType}` },
    });
  },

  setupController(controller) {
    this._super(...arguments);
    const { item_type: itemType } = this.paramsFor('vault.cluster.access.method.item');
    const { path: method } = this.paramsFor('vault.cluster.access.method');
    controller.set('itemType', singularize(itemType));
    controller.set('mode', 'create');
    controller.set('method', method);
  },
});
