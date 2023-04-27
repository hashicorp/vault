/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model, { attr, belongsTo, hasMany } from '@ember-data/model';
import { inject as service } from '@ember/service';
import { alias, and, equal, gte, not, or } from '@ember/object/computed';
import { get, computed } from '@ember/object';

export default Model.extend({
  version: service(),

  nodes: hasMany('nodes', { async: false }),
  name: attr('string'),
  status: attr('string'),
  standby: attr('boolean'),
  type: attr('string'),
  license: attr('object'),

  /* Licensing concerns */
  licenseExpiry: alias('license.expiry_time'),
  licenseState: alias('license.state'),

  needsInit: computed('nodes', 'nodes.@each.initialized', function () {
    // needs init if no nodes are initialized
    return this.nodes.isEvery('initialized', false);
  }),

  unsealed: computed('nodes', 'nodes.{[],@each.sealed}', function () {
    // unsealed if there's at least one unsealed node
    return !!this.nodes.findBy('sealed', false);
  }),

  sealed: not('unsealed'),

  leaderNode: computed('nodes', 'nodes.[]', function () {
    const nodes = this.nodes;
    if (nodes.get('length') === 1) {
      return nodes.get('firstObject');
    } else {
      return nodes.findBy('isLeader');
    }
  }),

  sealThreshold: alias('leaderNode.sealThreshold'),
  sealProgress: alias('leaderNode.progress'),
  sealType: alias('leaderNode.type'),
  storageType: alias('leaderNode.storageType'),
  hcpLinkStatus: alias('leaderNode.hcpLinkStatus'),
  hasProgress: gte('sealProgress', 1),
  usingRaft: equal('storageType', 'raft'),

  //replication mode - will only ever be 'unsupported'
  //otherwise the particular mode will have the relevant mode attr through replication-attributes
  mode: attr('string'),
  allReplicationDisabled: and('{dr,performance}.replicationDisabled'),
  anyReplicationEnabled: or('{dr,performance}.replicationEnabled'),

  dr: belongsTo('replication-attributes', { async: false, inverse: null }),
  performance: belongsTo('replication-attributes', { async: false, inverse: null }),
  // this service exposes what mode the UI is currently viewing
  // replicationAttrs will then return the relevant `replication-attributes` model
  rm: service('replication-mode'),
  drMode: alias('dr.mode'),
  replicationMode: alias('rm.mode'),
  replicationModeForDisplay: computed('replicationMode', function () {
    return this.replicationMode === 'dr' ? 'Disaster Recovery' : 'Performance';
  }),
  replicationIsInitializing: computed('dr.mode', 'performance.mode', function () {
    // a mode of null only happens when a cluster is being initialized
    // otherwise the mode will be 'disabled', 'primary', 'secondary'
    return !this.dr.mode || !this.performance.mode;
  }),
  replicationAttrs: computed('dr.mode', 'performance.mode', 'replicationMode', function () {
    const replicationMode = this.replicationMode;
    return replicationMode ? get(this, replicationMode) : null;
  }),
});
