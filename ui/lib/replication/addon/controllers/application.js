/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { isPresent } from '@ember/utils';
import { alias } from '@ember/object/computed';
import { inject as service } from '@ember/service';
import Controller from '@ember/controller';
import { copy } from 'ember-copy';
import { resolve } from 'rsvp';
import decodeConfigFromJWT from 'replication/utils/decode-config-from-jwt';

const DEFAULTS = {
  token: null,
  id: null,
  loading: false,
  errors: [],
  primary_api_addr: null,
  primary_cluster_addr: null,
  filterConfig: {
    mode: null,
    paths: [],
  },
};

export default Controller.extend(copy(DEFAULTS, true), {
  isModalActive: false,
  expirationDate: null,
  store: service(),
  rm: service('replication-mode'),
  replicationMode: alias('rm.mode'),
  flashMessages: service(),

  submitError(e) {
    if (e.errors) {
      this.set('errors', e.errors);
    } else {
      throw e;
    }
  },

  saveFilterConfig() {
    const config = this.filterConfig;
    const id = this.id;
    config.id = id;
    // if there is no mode, then they don't want to filter, so we don't save a filter config
    if (!config.mode) {
      return resolve();
    }
    const configRecord = this.store.createRecord('path-filter-config', config);
    return configRecord.save().catch((e) => this.submitError(e));
  },

  reset() {
    this.setProperties(copy(DEFAULTS, true));
  },

  submitSuccess(resp, action) {
    const cluster = this.model;
    if (!cluster) {
      return;
    }

    if (resp && resp.wrap_info) {
      this.set('token', resp.wrap_info.token);
    }
    if (action === 'secondary-token') {
      this.setProperties({
        loading: false,
        primary_api_addr: null,
        primary_cluster_addr: null,
      });

      // decode token and return epoch expiration, convert to timestamp
      const expirationDate = new Date(decodeConfigFromJWT(this.token).exp * 1000);
      this.set('expirationDate', expirationDate);

      // open modal
      this.toggleProperty('isModalActive');
      return cluster.reload();
    }
    this.reset();
    this.send('refresh');
    return;
  },

  submitHandler(action, clusterMode, data, event) {
    const replicationMode = this.replicationMode;
    if (event && event.preventDefault) {
      event.preventDefault();
    }
    this.setProperties({
      loading: true,
      errors: [],
    });
    if (data) {
      data = Object.keys(data).reduce((newData, key) => {
        var val = data[key];
        if (isPresent(val)) {
          newData[key] = val;
        }
        return newData;
      }, {});
    }

    return this.store
      .adapterFor('cluster')
      .replicationAction(action, replicationMode, clusterMode, data)
      .then(
        (resp) => {
          return this.saveFilterConfig().then(() => {
            return this.submitSuccess(resp, action, clusterMode);
          });
        },
        (...args) => this.submitError(...args)
      );
  },

  actions: {
    onSubmit(/*action, mode, data, event*/) {
      return this.submitHandler(...arguments);
    },
    copyClose(successMessage) {
      // separate action for copy & close button so it does not try and use execCommand to copy token to clipboard
      if (!!successMessage && typeof successMessage === 'string') {
        this.flashMessages.success(successMessage);
      }
      this.toggleProperty('isModalActive');
      this.transitionToRoute('mode.secondaries');
    },
    toggleModal(successMessage) {
      if (!!successMessage && typeof successMessage === 'string') {
        this.flashMessages.success(successMessage);
      }
      // use copy browser extension to copy token if you close the modal by clicking outside of it.
      const htmlSelectedToken = document.querySelector('#token-textarea');
      htmlSelectedToken.select();
      document.execCommand('copy');

      this.toggleProperty('isModalActive');
      this.transitionToRoute('mode.secondaries');
    },
    clear() {
      this.reset();
      this.setProperties({
        token: null,
        id: null,
      });
    },
    refresh() {
      // bubble to the route
      return true;
    },
  },
});
