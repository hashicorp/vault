import Ember from 'ember';

export default Ember.Component.extend({
  auth: Ember.inject.service(),

  routing: Ember.inject.service('-routing'),

  transitionToRoute: function() {
    var router = this.get('routing.router');
    router.transitionTo.apply(router, arguments);
  },

  classNames: 'user-menu auth-info',

  isRenewing: Ember.computed.or('fakeRenew', 'auth.isRenewing'),

  actions: {
    renewToken() {
      this.set('fakeRenew', true);
      Ember.run.later(() => {
        this.set('fakeRenew', false);
        this.get('auth').renew();
      }, 200);
    },

    revokeToken() {
      this.get('auth').revokeCurrentToken().then(() => {
        this.transitionToRoute('vault.cluster.logout');
      });
    },
  },
});
