/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { computed } from '@ember/object';
import Component from '@ember/component';
import ReplicationActions from 'core/mixins/replication-actions';

/**
 * @module ReplicationSummary
 * ReplicationSummary component is a component to show the mode-specific summary for replication
 *
 * @param {ClusterModel} cluster - the cluster ember-data model
 * @param {string} initialReplicationMode - mode for replication details we want to see, either "dr" or "performance"
 */
export default Component.extend(ReplicationActions, {
  'data-test-replication-summary': true,
  attributeBindings: ['data-test-replication-summary'],
  replicationMode: 'dr',
  mode: 'primary',
  version: service(),
  rm: service('replication-mode'),
  didReceiveAttrs() {
    this._super(...arguments);
    const initialReplicationMode = this.initialReplicationMode;
    if (initialReplicationMode) {
      this.set('replicationMode', initialReplicationMode);
    }
  },
  initialReplicationMode: null,
  cluster: null,

  attrsForCurrentMode: computed('cluster', 'rm.mode', function () {
    return this.cluster[this.rm.mode];
  }),
});
