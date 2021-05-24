import { inject as service } from '@ember/service';
import { reads } from '@ember/object/computed';
import Component from '@ember/component';

/**
 * @module ClusterInfo
 *
 * @example
 * ```js
 * <ClusterInfo @cluster={{cluster}} @onLinkClick={{action}} />
 * ```
 *
 * @param {object} cluster - details of the current cluster, passed from the parent.
 * @param {Function} onLinkClick - parent action which determines the behavior on link click
 */
export default Component.extend({
  auth: service(),
  store: service(),
  version: service(),
  cluster: null,

  transitionToRoute: function() {
    this.router.transitionTo(...arguments);
  },

  currentToken: reads('auth.currentToken'),
});
