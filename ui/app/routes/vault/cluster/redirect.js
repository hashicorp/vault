/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { AUTH, CLUSTER } from 'vault/lib/route-paths';

export default class VaultClusterRedirectRoute extends Route {
  @service auth;
  @service router;

  beforeModel({ to: { queryParams } }) {
    let transition;
    const isAuthed = this.auth.currentToken;
    // eslint-disable-next-line ember/no-controller-access-in-routes
    const controller = this.controllerFor('vault');
    const { redirect_to, ...otherParams } = queryParams;

    if (isAuthed && redirect_to) {
      // if authenticated and redirect exists, redirect to that place and strip other params
      transition = this.router.replaceWith(redirect_to);
    } else if (isAuthed) {
      // if authed no redirect, go to cluster
      transition = this.router.replaceWith(CLUSTER, { queryParams: otherParams });
    } else {
      // default go to Auth
      transition = this.router.replaceWith(AUTH, { queryParams: otherParams });
    }
    transition.followRedirects().then(() => {
      controller.set('redirectTo', '');
    });
  }
}
