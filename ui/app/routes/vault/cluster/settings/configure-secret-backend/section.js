/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import AdapterError from '@ember-data/adapter/error';
import { set } from '@ember/object';
import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

// ARG TODO glimmerize
const SECTIONS_FOR_TYPE = {
  pki: ['cert', 'urls', 'crl', 'tidy'],
};

export default Route.extend({
  store: service(),

  fetchModel() {
    const { section_name: sectionName } = this.paramsFor(this.routeName);
    const backendModel = this.modelFor('vault.cluster.settings.configure-secret-backend');
    const type = backendModel.get('type');
    let modelType;
    if (type === 'pki') {
      // pki models are in models/pki
      modelType = `${type}/${type}-config`;
    } else {
      modelType = `${type}-config`;
    }
    return this.store
      .queryRecord(modelType, {
        backend: backendModel.id,
        section: sectionName,
      })
      .then((model) => {
        model.set('backendType', type);
        model.set('section', sectionName);
        return model;
      });
  },

  model(params) {
    const { section_name: sectionName } = params;
    const backendModel = this.modelFor('vault.cluster.settings.configure-secret-backend');
    const sections = SECTIONS_FOR_TYPE[backendModel.get('type')];
    const hasSection = sections.includes(sectionName);
    if (!backendModel || !hasSection) {
      const error = new AdapterError();
      set(error, 'httpStatus', 404);
      throw error;
    }
    return this.fetchModel();
  },

  setupController(controller) {
    this._super(...arguments);
    controller.set('onRefresh', () => this.fetchModel());
  },
});
