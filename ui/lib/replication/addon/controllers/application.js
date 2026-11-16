/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { isPresent } from '@ember/utils';
import { alias } from '@ember/object/computed';
import { service } from '@ember/service';
import Controller from '@ember/controller';
import PathFilterConfigForm from 'vault/forms/path-filter-config';
import decodeConfigFromJWT from 'replication/utils/decode-config-from-jwt';
import { buildWaiter } from '@ember/test-waiters';

const DEFAULTS = {
  token: null,
  id: null,
  loading: false,
  errors: [],
  primary_api_addr: null,
  primary_cluster_addr: null,
};
const waiter = buildWaiter('replication-actions');

/**
 * @type Class
 */
export default Controller.extend(structuredClone(DEFAULTS), {
  api: service(),
  filterConfig: null,
  isModalActive: false,
  isTokenCopied: false,
  expirationDate: null,
  router: service('app-router'),
  store: service(),
  rm: service('replication-mode'),
  replicationMode: alias('rm.mode'),
  secondaryToRevoke: null,

  init() {
    this._super(...arguments);
    this.set('filterConfig', new PathFilterConfigForm({}, { isNew: true }));
  },

  submitError(e) {
    if (e.errors) {
      this.set('errors', e.errors);
    } else {
      throw e;
    }
  },

  async saveFilterConfig() {
    const { data } = this.filterConfig.toJSON();
    if (!data.mode) {
      return;
    }
    try {
      await this.api.sys.systemWriteReplicationPerformancePrimaryPathsFilterId(this.id, data);
    } catch (error) {
      const { message } = await this.api.parseError(error);
      // Keep filter errors on the form without discarding an already-created secondary token.
      this.set('errors', [message]);
    }
  },

  reset() {
    this.setProperties({
      ...structuredClone(DEFAULTS),
      filterConfig: new PathFilterConfigForm({}, { isNew: true }),
    });
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
    const waiterToken = waiter.beginAsync();
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
      )
      .finally(() => {
        this.set('secondaryToRevoke', null);
        waiter.endAsync(waiterToken);
      });
  },

  actions: {
    onSubmit(/*action, mode, data, event*/) {
      return this.submitHandler(...arguments);
    },
    closeTokenModal() {
      this.toggleProperty('isModalActive');
      this.router.transitionTo('vault.cluster.replication.mode.secondaries');
      this.set('isTokenCopied', false);
    },
    onCopy() {
      this.set('isTokenCopied', true);
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
