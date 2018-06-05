import Ember from 'ember';

const { get } = Ember;
const INIT = 'vault.cluster.init';
const UNSEAL = 'vault.cluster.unseal';
const AUTH = 'vault.cluster.auth';
const CLUSTER = 'vault.cluster';
const DR_REPLICATION_SECONDARY = 'vault.cluster.replication-dr-promote';

export { INIT, UNSEAL, AUTH, CLUSTER, DR_REPLICATION_SECONDARY };

export default Ember.Mixin.create({
  auth: Ember.inject.service(),

  transitionToTargetRoute() {
    const targetRoute = this.targetRouteName();
    if (targetRoute && targetRoute !== this.routeName) {
      return this.transitionTo(targetRoute);
    }
    return Ember.RSVP.resolve();
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

  targetRouteName() {
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
