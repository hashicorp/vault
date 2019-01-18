import { inject as service } from '@ember/service';
import { alias } from '@ember/object/computed';
import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  currentCluster: service('current-cluster'),
  cluster: alias('currentCluster.cluster'),
  auth: service(),
  store: service(),
  media: service(),
  version: service(),
  type: 'cluster',
  itemTag: null,
  partialName: computed('type', function() {
    return `partials/status/${this.type}`;
  }),
  glyphName: computed('type', function() {
    const glyphs = {
      cluster: 'status-indicator',
      user: 'person',
    };
    return glyphs[this.type];
  }),
  activeCluster: computed('auth.activeCluster', function() {
    return this.get('store').peekRecord('cluster', this.get('auth.activeCluster'));
  }),
  currentToken: computed('auth.currentToken', function() {
    return this.get('auth.currentToken');
  }),
});
