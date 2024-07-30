/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AdapterError from '@ember-data/adapter/error';
import { set } from '@ember/object';
import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';
import { CONFIGURABLE_SECRET_ENGINES } from 'vault/helpers/mountable-secret-engines';

export default Route.extend({
  store: service(),

  async model() {
    const { backend } = this.paramsFor(this.routeName);
    return this.store.query('secret-engine', { path: backend }).then(async (modelList) => {
      const model = modelList && modelList[0];
      const type = model.type;
      if (!model || !CONFIGURABLE_SECRET_ENGINES.includes(type)) {
        const error = new AdapterError();
        set(error, 'httpStatus', 404);
        throw error;
      }
      if (type === 'aws') {
        // For AWS we have two models used for configuration. Similar to KVv2, we want to keep these models separate for permissions reasons.
        let leaseConfig;
        try {
          leaseConfig = await this.store.queryRecord('aws/lease-config', {
            backend: model.id,
            type,
          });
        } catch (e) {
          // if you haven't saved a lease config, the API 404s, so create one here to edit and return it
          if (e.httpStatus === 404) {
            leaseConfig = await this.store.createRecord('aws/lease-config', {
              backend: model.id,
              type,
            });
          } else {
            leaseConfig = e; // assign error to model and handle in ConfigureAwsSecret
          }
        }
        let rootConfig;
        try {
          rootConfig = await this.store.queryRecord('aws/root-config', {
            backend: model.id,
            type,
          });
        } catch (e) {
          rootConfig = e; // assign error to model and handle in ConfigureAwsSecret component
        }
        // reassign model
        return hash({
          leaseConfig,
          rootConfig,
          path: backend,
          type,
        });
      }
      // For all other secret engines we only use the one secret-engine model used for configuration that pulls in configuration data from a belongsTo relationship.
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

  setupController(controller, model) {
    if (model?.publicKey) {
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
