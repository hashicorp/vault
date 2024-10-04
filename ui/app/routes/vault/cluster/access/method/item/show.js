/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { singularize } from 'ember-inflector';
import { service } from '@ember/service';
import Route from '@ember/routing/route';

export default Route.extend({
  store: service(),
  pathHelp: service('path-help'),

  model(params) {
    const id = params.item_id;
    const { item_type: itemType } = this.paramsFor('vault.cluster.access.method.item');
    const methodModel = this.modelFor('vault.cluster.access.method');
    const modelType = `generated-${singularize(itemType)}-${methodModel.type}`;
    return this.store.queryRecord(modelType, { id, authMethodPath: methodModel.id });
  },

  setupController(controller) {
    this._super(...arguments);
    const { item_type: itemType } = this.paramsFor('vault.cluster.access.method.item');
    const { path: method } = this.paramsFor('vault.cluster.access.method');
    controller.set('itemType', singularize(itemType));
    controller.set('method', method);
  },
});
