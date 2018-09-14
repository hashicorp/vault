import Ember from 'ember';

const { Component, inject, computed, run } = Ember;
export default Component.extend({
  auth: inject.service(),
  wizard: inject.service(),
  routing: inject.service('-routing'),

  transitionToRoute: function() {
    var router = this.get('routing.router');
    router.transitionTo.apply(router, arguments);
  },

  classNames: 'user-menu auth-info',

  isRenewing: computed.or('fakeRenew', 'auth.isRenewing'),

  actions: {
    restartGuide() {
      this.get('wizard').restartGuide();
    },
    renewToken() {
      this.set('fakeRenew', true);
      run.later(() => {
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
