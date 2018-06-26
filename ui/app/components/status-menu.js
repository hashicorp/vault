import Ember from 'ember';

const { inject, computed } = Ember;

export default Ember.Component.extend({
  currentCluster: inject.service('current-cluster'),
  cluster: computed.alias('currentCluster.cluster'),
  auth: inject.service(),
  type: 'cluster',
  itemTag: null,
  partialName: computed('type', function() {
    let type = this.get('type');
    let partial = type === 'replication-status' ? 'replication' : type;
    return `partials/status/${partial}`;
  }),
  glyphName: computed('type', function() {
    const glyphs = {
      cluster: 'unlocked',
      user: 'android-person',
      'replication-status': 'replication',
    };
    return glyphs[this.get('type')];
  }),
});
