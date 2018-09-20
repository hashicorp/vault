import { inject as service } from '@ember/service';
import Component from '@ember/component';

export default Component.extend({
  auth: service(),
  router: service(),
  classNames: 'token-expire-warning',

  transitionToRoute: function() {
    this.get('router').transitionTo(...arguments);
  },

  isDismissed: false,

  actions: {
    reauthenticate() {
      this.get('auth').deleteCurrentToken();
      this.transitionToRoute('vault.cluster');
    },

    renewToken() {
      const auth = this.get('auth');
      auth.renew();
      auth.setLastFetch(Date.now());
    },

    dismiss() {
      this.set('isDismissed', true);
    },
  },
});
