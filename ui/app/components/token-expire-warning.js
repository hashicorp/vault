import Ember from 'ember';

export default Ember.Component.extend({
  classNames: 'token-expire-warning',
  auth: Ember.inject.service(),

  routing: Ember.inject.service('-routing'),

  transitionToRoute: function() {
    var router = this.get('routing.router');
    router.transitionTo.apply(router, arguments);
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
