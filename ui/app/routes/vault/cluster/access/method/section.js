/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AdapterError from '@ember-data/adapter/error';
import { set } from '@ember/object';
import Route from '@ember/routing/route';

export default Route.extend({
  model(params) {
    const { section_name: section } = params;
    if (section !== 'configuration') {
      const error = new AdapterError();
      set(error, 'httpStatus', 404);
      throw error;
    }
    return this.modelFor('vault.cluster.access.method');
  },

  setupController(controller, model) {
    const { section_name: section } = this.paramsFor(this.routeName);
    this._super(...arguments);
    controller.set('section', section);
    controller.set(
      'paths',
      model.paths.paths.filter((path) => path.navigation)
    );
  },
});
