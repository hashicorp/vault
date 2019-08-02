import { inject as service } from '@ember/service';
import { get } from '@ember/object';
import Mixin from '@ember/object/mixin';
import RSVP from 'rsvp';
const INIT = 'vault.cluster.init';
const UNSEAL = 'vault.cluster.unseal';
const AUTH = 'vault.cluster.auth';
const CLUSTER = 'vault.cluster';
const CLUSTER_INDEX = 'vault.cluster.index';
const OIDC_CALLBACK = 'vault.cluster.oidc-callback';
const DR_REPLICATION_SECONDARY = 'vault.cluster.replication-dr-promote';

export { INIT, UNSEAL, AUTH, CLUSTER, CLUSTER_INDEX, DR_REPLICATION_SECONDARY };

export default Mixin.create({
  auth: service(),
  store: service(),
  router: service(),

  transitionToTargetRoute(transition = {}) {
    const targetRoute = this.targetRouteName(transition);

    if (
      targetRoute &&
      targetRoute !== this.routeName &&
      targetRoute !== transition.targetName &&
      targetRoute !== this.router.currentRouteName
    ) {
      if (
        // only want to redirect if we're going to authenticate
        targetRoute === AUTH &&
        transition.targetName !== CLUSTER_INDEX
      ) {
        return this.transitionTo(targetRoute, { queryParams: { redirect_to: this.router.currentURL } });
      }
      return this.transitionTo(targetRoute);
    }

    return RSVP.resolve();
  },

  beforeModel(transition) {
    return this.transitionToTargetRoute(transition);
  },

  clusterModel() {
    return this.modelFor(CLUSTER) || this.store.peekRecord('cluster', 'vault');
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
