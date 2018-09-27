import { inject as service } from '@ember/service';
import { alias } from '@ember/object/computed';
import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  currentCluster: service('current-cluster'),
  cluster: alias('currentCluster.cluster'),
  auth: service(),
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
