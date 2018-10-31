import { inject as service } from '@ember/service';
import Component from '@ember/component';

export default Component.extend({
  auth: service(),
  router: service(),
  classNames: 'token-expire-warning',

  actions: {
    reauthenticate() {
      this.get('auth').deleteCurrentToken();
      this.get('router').transitionTo('vault.cluster');
    },
  },
});
