import Ember from 'ember';

const { inject, computed } = Ember;

export default Ember.Component.extend({
  currentCluster: inject.service('current-cluster'),
  cluster: computed.alias('currentCluster.cluster'),
  auth: inject.service(),
  type: 'cluster',
  partialName: computed('type', function() {
    return `partials/status/${this.get('type')}`;
  }),
  glyphName: computed('type', function() {
    const glyphs = {
      cluster: 'unlocked',
      user: 'android-person',
      replication: 'replication',
    };
    return glyphs[this.get('type')];
  }),
});
