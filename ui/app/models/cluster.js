import { inject as service } from '@ember/service';
import { not, gte, alias, and, or } from '@ember/object/computed';
import { get, computed } from '@ember/object';
import DS from 'ember-data';
import { fragment } from 'ember-data-model-fragments/attributes';
const { hasMany, attr } = DS;

export default DS.Model.extend({
  version: service(),

  nodes: hasMany('nodes', { async: false }),
  name: attr('string'),
  status: attr('string'),
  standby: attr('boolean'),
  type: attr('string'),

  needsInit: computed('nodes', 'nodes.@each.initialized', function() {
    // needs init if no nodes are initialized
    return this.get('nodes').isEvery('initialized', false);
  }),

  unsealed: computed('nodes', 'nodes.{[],@each.sealed}', function() {
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
  sealType: alias('leaderNode.type'),
  hasProgress: gte('sealProgress', 1),

  //replication mode - will only ever be 'unsupported'
  //otherwise the particular mode will have the relevant mode attr through replication-attributes
  mode: attr('string'),
  allReplicationDisabled: and('{dr,performance}.replicationDisabled'),

  anyReplicationEnabled: or('{dr,performance}.replicationEnabled'),

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
  rm: service('replication-mode'),
  replicationMode: alias('rm.mode'),
  replicationAttrs: computed('dr.mode', 'performance.mode', 'replicationMode', function() {
    const replicationMode = this.get('replicationMode');
    return replicationMode ? get(this, replicationMode) : null;
  }),
});
