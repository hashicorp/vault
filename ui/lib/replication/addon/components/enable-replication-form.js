/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import decodeConfigFromJwt from 'replication/utils/decode-config-from-jwt';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import errorMessage from 'vault/utils/error-message';
import { isPresent } from '@ember/utils';
import { waitFor } from '@ember/test-waiters';

/**
 * @module EnableReplicationFormComponent
 * EnableReplicationForm component is used in the replication engine to enable replication. It must be passed the replicationMode,
 * but otherwise it handles the rest of the form inputs. On success it will clear the form and call the onSuccess callback.
 *
 * @example
 * ```js
 * <EnableReplicationForm @replicationMode="dr" @canEnablePrimary={{true}} @canEnableSecondary={{false}} @performanceReplicationDisabled={{false}} @onSuccess={{this.reloadCluster}} />
 *    @param {string} replicationMode - should be one of "dr" or "performance"
 *    @param {boolean} canEnablePrimary - if the capabilities allow the user to enable a primary cluster
 *    @param {boolean} canEnableSecondary - if the capabilities allow the user to enable a secondary cluster
 *    @param {boolean} performanceMode - should be "primary", "secondary", or "disabled". If enabled, form will show a warning when attempting to enable DR secondary
 *    @param {Promise} onSuccess - (optional) callback called after successful replication enablement. Must be a promise.
 *    @param {boolean} doTransition - (optional) if provided, passed to onSuccess callback to determine if a transition should be done
 *  />
 * ```
 */
export default class EnableReplicationFormComponent extends Component {
  @service version;
  @service store;

  @tracked error = '';
  @tracked showExplanation = false;
  data = new EnablePayload();

  get performanceReplicationEnabled() {
    return this.args.performanceMode !== 'disabled';
  }

  get tokenIncludesAPIAddr() {
    const config = decodeConfigFromJwt(this.token);
    return config && config.addr ? true : false;
  }

  get disallowEnable() {
    if (this.args.replicationMode === 'performance' && this.version.hasPerfReplication === false) {
      return true;
    }
    const { mode, tokenIncludesAPIAddr, primary_api_addr } = this.data;
    if (mode !== 'secondary' || tokenIncludesAPIAddr || (!tokenIncludesAPIAddr && primary_api_addr)) {
      return false;
    }
    return true;
  }

  async onSuccess(resp, clusterMode) {
    // clear form
    this.data.reset();
    // call callback
    if (this.args.onSuccess) {
      await this.args.onSuccess(resp, this.args.replicationMode, clusterMode, this.args.doTransition);
    }
  }

  @action inputChange(evt) {
    const name = evt.target.name;
    const val = evt.target.value;
    this.data[name] = val;
  }

  @task
  @waitFor
  *enableReplication(replicationMode, clusterMode, data) {
    const payload = data.allKeys.reduce((newData, key) => {
      var val = data[key];
      if (isPresent(val)) {
        newData[key] = val;
      }
      return newData;
    }, {});
    delete payload.mode;
    try {
      const resp = yield this.store
        .adapterFor('cluster')
        .replicationAction('enable', replicationMode, clusterMode, payload);
      yield this.onSuccess(resp, clusterMode);
    } catch (e) {
      this.error = errorMessage(e, 'Enable replication failed. Check Vault logs for details.');
    }
  }

  @action onSubmit(payload, evt) {
    evt.preventDefault();
    this.error = '';
    this.enableReplication.perform(this.args.replicationMode, this.data.mode, payload);
  }
}

class EnablePayload {
  @tracked mode = 'primary';
  @tracked token = '';
  @tracked primary_api_addr = '';
  @tracked primary_cluster_addr = '';
  @tracked ca_file = '';
  @tracked ca_path = '';
  get tokenIncludesAPIAddr() {
    const config = decodeConfigFromJwt(this.token);
    return config && config.addr ? true : false;
  }
  get allKeys() {
    return ['mode', 'token', 'primary_api_addr', 'primary_cluster_addr', 'ca_file', 'ca_path'];
  }
  reset() {
    // reset all but mode
    this.token = '';
    this.primary_api_addr = '';
    this.primary_cluster_addr = '';
    this.ca_file = '';
    this.ca_path = '';
  }
}
