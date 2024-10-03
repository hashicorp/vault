/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AdapterError from '@ember-data/adapter/error';
import { next } from '@ember/runloop';
import { hash } from 'rsvp';
import { set } from '@ember/object';
import Route from '@ember/routing/route';
import { TABS } from 'vault/helpers/tabs-for-identity-show';
import { service } from '@ember/service';

export default Route.extend({
  router: service(),
  store: service(),

  model(params) {
    const { section } = params;
    const itemType = this.modelFor('vault.cluster.access.identity');
    const tabs = TABS[itemType];
    const modelType = `identity/${itemType}`;
    if (!tabs.includes(section)) {
      const error = new AdapterError();
      set(error, 'httpStatus', 404);
      throw error;
    }

    // if the record is in the store use that
    let model = this.store.peekRecord(modelType, params.item_id);

    // if we don't have creationTime, we only have a partial model so reload
    if (model && !model?.creationTime) {
      model = model.reload();
    }

    // if there's no model, we need to fetch it
    if (!model) {
      model = this.store.findRecord(modelType, params.item_id);
    }

    return hash({
      model,
      section,
    });
  },

  activate() {
    // if we're just entering the route, and it's not a hard reload
    // reload to make sure we have the newest info
    if (this.currentModel) {
      next(() => {
        /* eslint-disable-next-line ember/no-controller-access-in-routes */
        this.controller.model.reload();
      });
    }
  },

  afterModel(resolvedModel) {
    const { section, model } = resolvedModel;
    if (model?.identityType === 'group' && model?.type === 'internal' && section === 'aliases') {
      return this.router.transitionTo('vault.cluster.access.identity.show', model.id, 'details');
    }
  },

  setupController(controller, resolvedModel) {
    const { model, section } = resolvedModel;
    controller.setProperties({
      model,
      section,
    });
  },
});
