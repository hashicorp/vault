import { inject as service } from '@ember/service';
import { reads } from '@ember/object/computed';
import Component from '@ember/component';

export default Component.extend({
  auth: service(),
  store: service(),
  version: service(),

  transitionToRoute: function() {
    this.router.transitionTo(...arguments);
  },

  currentToken: reads('auth.currentToken'),
});
