/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { inject as service } from '@ember/service';
import Mixin from '@ember/object/mixin';
import RSVP from 'rsvp';
import {
  INIT,
  UNSEAL,
  AUTH,
  CLUSTER,
  CLUSTER_INDEX,
  OIDC_CALLBACK,
  OIDC_PROVIDER,
  NS_OIDC_PROVIDER,
  DR_REPLICATION_SECONDARY,
  DR_REPLICATION_SECONDARY_DETAILS,
  EXCLUDED_REDIRECT_URLS,
  REDIRECT,
  LOGOUT,
} from 'vault/lib/route-paths';

export default Mixin.create({
  auth: service(),
  store: service(),
  router: service(),

  _shouldInterceptTransition(targetRoute, transitionTarget, currentRouteName) {
    const attemptingLogout = transitionTarget === LOGOUT || this.routeName === LOGOUT;
    if (
      // not attempting to logout
      !attemptingLogout &&
      // has target route to intercept to
      targetRoute &&
      // target route different from current route
      targetRoute !== currentRouteName &&
      // target route different from transition target
      targetRoute !== transitionTarget
    ) {
      return true;
    }
    return false;
  },

  transitionToTargetRoute(transition = {}) {
    const targetRoute = this.targetRouteName(transition);
    const intercept = this._shouldInterceptTransition(
      targetRoute,
      transition.targetName,
      this.routeName || this.router.currentRouteName
    );
    if (intercept) {
      // there may be query params so check for inclusion rather than exact match
      const isExcluded = EXCLUDED_REDIRECT_URLS.find((url) => this.router.currentURL?.includes(url));
      if (
        // only want to add redirect_to if we're going to authenticate
        targetRoute === AUTH &&
        // and if we aren't redirecting from logout
        !isExcluded &&
        // and the target isn't the default
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
    return this.auth.currentToken;
  },

  hasKeyData() {
    /* eslint-disable-next-line ember/no-controller-access-in-routes */
    return !!this.controllerFor(INIT).keyData;
  },

  targetRouteName(transition) {
    const cluster = this.clusterModel();
    const isAuthed = this.authToken();
    if (cluster.needsInit) {
      return INIT;
    }
    if (this.hasKeyData() && this.routeName !== UNSEAL && this.routeName !== AUTH) {
      return INIT;
    }
    if (cluster.sealed) {
      return UNSEAL;
    }
    if (cluster?.dr?.isSecondary) {
      if (transition && transition.targetName === DR_REPLICATION_SECONDARY_DETAILS) {
        return DR_REPLICATION_SECONDARY_DETAILS;
      }
      if (this.router.currentRouteName === DR_REPLICATION_SECONDARY_DETAILS) {
        return DR_REPLICATION_SECONDARY_DETAILS;
      }

      return DR_REPLICATION_SECONDARY;
    }
    if (!isAuthed) {
      if ((transition && transition.targetName === OIDC_PROVIDER) || this.routeName === OIDC_PROVIDER) {
        return OIDC_PROVIDER;
      }
      if ((transition && transition.targetName === NS_OIDC_PROVIDER) || this.routeName === NS_OIDC_PROVIDER) {
        return NS_OIDC_PROVIDER;
      }
      if ((transition && transition.targetName === OIDC_CALLBACK) || this.routeName === OIDC_CALLBACK) {
        return OIDC_CALLBACK;
      }
      return AUTH;
    }
    if (
      (!cluster.needsInit && this.routeName === INIT) ||
      (!cluster.sealed && this.routeName === UNSEAL) ||
      (!cluster?.dr?.isSecondary && this.routeName === DR_REPLICATION_SECONDARY)
    ) {
      return CLUSTER;
    }
    if (isAuthed && this.routeName === AUTH) {
      // if you're already authed and you wanna go to auth, you probably want to redirect
      return REDIRECT;
    }
    return null;
  },
});
