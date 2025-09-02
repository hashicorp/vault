/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AdapterError from '@ember-data/adapter/error';
import { hash } from 'rsvp';
import { set } from '@ember/object';
import Route from '@ember/routing/route';
import { TABS } from 'vault/helpers/tabs-for-identity-show';
import { service } from '@ember/service';

export default Route.extend({
  store: service(),

  model(params) {
    const { section } = params;
    const itemType = this.modelFor('vault.cluster.access.identity') + '-alias';
    const tabs = TABS[itemType];
    const modelType = `identity/${itemType}`;
    if (!tabs.includes(section)) {
      const error = new AdapterError();
      set(error, 'httpStatus', 404);
      throw error;
    }
    // TODO peekRecord here to see if we have the record already
    return hash({
      model: this.store.findRecord(modelType, params.item_alias_id),
      section,
    });
  },

  setupController(controller, resolvedModel) {
    const { model, section } = resolvedModel;
    controller.setProperties({
      model,
      section,
    });
  },
});
