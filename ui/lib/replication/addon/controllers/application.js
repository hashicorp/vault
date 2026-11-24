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
  rm: service('replication-mode'),
  replicationMode: alias('rm.mode'),
  secondaryToRevoke: null,

  init() {
    this._super(...arguments);
    this.set('filterConfig', new PathFilterConfigForm({}, { isNew: true }));
  },

  async submitError(e) {
    const { status, message, response } = await this.api.parseError(e);
    if (response?.errors) {
      this.set('errors', response.errors);
    } else if (status) {
      // errors returned as { data: { error } } or without a body are summarized in the message
      this.set('errors', [message]);
    } else {
      throw e;
    }
  },

  replicationRequest(action, replicationMode, clusterMode, data) {
    const { sys } = this.api;
    if (clusterMode === 'primary') {
      if (action === 'secondary-token') {
        return replicationMode === 'dr'
          ? sys.systemWriteReplicationDrPrimarySecondaryToken(data)
          : sys.systemWriteReplicationPerformancePrimarySecondaryToken(data);
      }
      if (action === 'revoke-secondary') {
        return replicationMode === 'dr'
          ? sys.systemWriteReplicationDrPrimaryRevokeSecondary(data)
          : sys.systemWriteReplicationPerformancePrimaryRevokeSecondary(data);
      }
    }
    throw new Error(`Unsupported replication action: ${replicationMode}/${clusterMode}/${action}`);
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

  async submitHandler(action, clusterMode, data, event) {
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

    try {
      let resp;
      try {
        resp = await this.replicationRequest(action, replicationMode, clusterMode, data);
      } catch (e) {
        return await this.submitError(e);
      }
      await this.saveFilterConfig();
      return this.submitSuccess(resp, action, clusterMode);
    } finally {
      this.set('secondaryToRevoke', null);
      waiter.endAsync(waiterToken);
    }
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
