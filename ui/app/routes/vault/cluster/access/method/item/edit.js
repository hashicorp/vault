/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { singularize } from 'ember-inflector';
import { service } from '@ember/service';
import { fetchGeneratedItemForm } from 'vault/routes/vault/cluster/access/method/item/show';

export default Route.extend({
  api: service(),

  model(params) {
    return fetchGeneratedItemForm(this, params.item_id);
  },

  setupController(controller) {
    this._super(...arguments);
    const { itemType, urls } = this.modelFor('vault.cluster.access.method.item');
    const { path: method } = this.paramsFor('vault.cluster.access.method');
    controller.set('itemType', singularize(itemType));
    controller.set('method', method);
    controller.set('urls', urls);
  },
});
