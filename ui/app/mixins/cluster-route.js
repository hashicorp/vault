import { inject as service } from '@ember/service';
import { get } from '@ember/object';
import Mixin from '@ember/object/mixin';
import RSVP from 'rsvp';
const INIT = 'vault.cluster.init';
const UNSEAL = 'vault.cluster.unseal';
const AUTH = 'vault.cluster.auth';
const CLUSTER = 'vault.cluster';
const OIDC_CALLBACK = 'vault.cluster.oidc-callback';
const DR_REPLICATION_SECONDARY = 'vault.cluster.replication-dr-promote';

export { INIT, UNSEAL, AUTH, CLUSTER, DR_REPLICATION_SECONDARY };

export default Mixin.create({
  auth: service(),

  transitionToTargetRoute(transition) {
    const targetRoute = this.targetRouteName(transition);
    if (targetRoute && targetRoute !== this.routeName) {
      return this.transitionTo(targetRoute);
    }

    return RSVP.resolve();
  },

  beforeModel() {
    return this.transitionToTargetRoute();
  },

  clusterModel() {
    return this.modelFor(CLUSTER);
  },

  authToken() {
    return get(this, 'auth.currentToken');
  },

  hasKeyData() {
    return !!get(this.controllerFor(INIT), 'keyData');
  },

  targetRouteName(transition) {
    const cluster = this.clusterModel();
    const isAuthed = this.authToken();
    if (get(cluster, 'needsInit')) {
      return INIT;
    }
    if (this.hasKeyData() && this.routeName !== UNSEAL && this.routeName !== AUTH) {
      return INIT;
    }
    if (get(cluster, 'sealed')) {
      return UNSEAL;
    }
    if (get(cluster, 'dr.isSecondary')) {
      return DR_REPLICATION_SECONDARY;
    }
    if (!isAuthed) {
      if ((transition && transition.targetName === OIDC_CALLBACK) || this.routeName === OIDC_CALLBACK) {
        return OIDC_CALLBACK;
      }
      return AUTH;
    }
    if (
      (!get(cluster, 'needsInit') && this.routeName === INIT) ||
      (!get(cluster, 'sealed') && this.routeName === UNSEAL) ||
      (!get(cluster, 'dr.isSecondary') && this.routeName === DR_REPLICATION_SECONDARY) ||
      (isAuthed && this.routeName === AUTH)
    ) {
      return CLUSTER;
    }
    return null;
  },
});
