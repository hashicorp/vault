import Ember from 'ember';
import DS from 'ember-data';
import { fragment } from 'ember-data-model-fragments/attributes';
const { hasMany, attr } = DS;
const { computed, get, inject } = Ember;
const { alias, gte, not } = computed;

export default DS.Model.extend({
  version: inject.service(),

  nodes: hasMany('nodes', { async: false }),
  name: attr('string'),
  status: attr('string'),
  standby: attr('boolean'),

  needsInit: computed('nodes', 'nodes.[]', function() {
    // needs init if no nodes are initialized
    return this.get('nodes').isEvery('initialized', false);
  }),

  type: computed(function() {
    return this.constructor.modelName;
  }),

  unsealed: computed('nodes', 'nodes.[]', 'nodes.@each.sealed', function() {
    // unsealed if there's at least one unsealed node
    return !!this.get('nodes').findBy('sealed', false);
  }),

  sealed: not('unsealed'),

  leaderNode: computed('nodes', 'nodes.[]', function() {
    const nodes = this.get('nodes');
    if (nodes.get('length') === 1) {
      return nodes.get('firstObject');
    } else {
      return nodes.findBy('isLeader');
    }
  }),

  sealThreshold: alias('leaderNode.sealThreshold'),
  sealProgress: alias('leaderNode.progress'),
  hasProgress: gte('sealProgress', 1),

  //replication mode - will only ever be 'unsupported'
  //otherwise the particular mode will have the relevant mode attr through replication-attributes
  mode: attr('string'),
  allReplicationDisabled: computed.and('{dr,performance}.replicationDisabled'),

  anyReplicationEnabled: computed.or('{dr,performance}.replicationEnabled'),

  stateDisplay(state) {
    if (!state) {
      return null;
    }
    const defaultDisp = 'Synced';
    const displays = {
      'stream-wals': 'Streaming',
      'merkle-diff': 'Determining sync status',
      'merkle-sync': 'Syncing',
    };

    return displays[state] || defaultDisp;
  },

  drStateDisplay: computed('dr.state', function() {
    return this.stateDisplay(this.get('dr.state'));
  }),

  performanceStateDisplay: computed('performance.state', function() {
    return this.stateDisplay(this.get('performance.state'));
  }),

  stateGlyph(state) {
    const glyph = 'checkmark-circled-outline';

    const glyphs = {
      'stream-wals': 'android-sync',
      'merkle-diff': 'android-sync',
      'merkle-sync': null,
    };

    return glyphs[state] || glyph;
  },

  drStateGlyph: computed('dr.state', function() {
    return this.stateGlyph(this.get('dr.state'));
  }),

  performanceStateGlyph: computed('performance.state', function() {
    return this.stateGlyph(this.get('performance.state'));
  }),

  dr: fragment('replication-attributes'),
  performance: fragment('replication-attributes'),
  // this service exposes what mode the UI is currently viewing
  // replicationAttrs will then return the relevant `replication-attributes` fragment
  rm: Ember.inject.service('replication-mode'),
  replicationMode: computed.alias('rm.mode'),
  replicationAttrs: computed('dr.mode', 'performance.mode', 'replicationMode', function() {
    const replicationMode = this.get('replicationMode');
    return replicationMode ? get(this, replicationMode) : null;
  }),
});
