/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { singularize } from 'ember-inflector';
import GeneratedItemForm from 'vault/forms/generated-item';

export default Route.extend({
  model() {
    const { formMetadata } = this.modelFor('vault.cluster.access.method.item');
    return new GeneratedItemForm(formMetadata, {}, { isNew: true });
  },

  setupController(controller) {
    this._super(...arguments);
    const { itemType, urls } = this.modelFor('vault.cluster.access.method.item');
    const { path: method } = this.paramsFor('vault.cluster.access.method');
    controller.set('itemType', singularize(itemType));
    controller.set('mode', 'create');
    controller.set('method', method);
    controller.set('urls', urls);
  },
});
