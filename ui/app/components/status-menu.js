import { inject as service } from '@ember/service';
import { alias } from '@ember/object/computed';
import Component from '@ember/component';
import { computed } from '@ember/object';
import { next } from '@ember/runloop';

export default Component.extend({
  currentCluster: service('current-cluster'),
  cluster: alias('currentCluster.cluster'),
  auth: service(),
  media: service(),
  type: 'cluster',
  itemTag: null,
  glyphName: computed('type', function () {
    return {
      cluster: 'circle-dot',
      user: 'user',
    }[this.type];
  }),

  actions: {
    onLinkClick(dropdown) {
      if (dropdown) {
        // strange issue where closing dropdown triggers full transition which redirects to auth screen in production builds
        // closing dropdown in next tick of run loop fixes it
        next(() => dropdown.actions.close());
      }
      this.onLinkClick();
    },
  },
});
