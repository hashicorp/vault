/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { computed } from '@ember/object';
import Component from '@ember/component';
import ReplicationActions from 'core/mixins/replication-actions';

export default Component.extend(ReplicationActions, {
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
  showModeSummary: false,
  initialReplicationMode: null,
  cluster: null,

  attrsForCurrentMode: computed('cluster', 'rm.mode', function () {
    return this.cluster[this.rm.mode];
  }),
});
