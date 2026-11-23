/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';

export default Route.extend({
  model(params) {
    const { section_name: section } = params;
    if (section !== 'configuration') {
      throw { httpStatus: 404 };
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
