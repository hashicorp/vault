/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { later } from '@ember/runloop';
import { task } from 'ember-concurrency';

export default class TokenExpireWarning extends Component {
  @service auth;
  @service router;
  @tracked canDismiss = true;

  handleRenew() {
    return new Promise((resolve) => {
      later(() => {
        this.auth
          .renew()
          .then(() => {
            // This renewal was triggered by an explicit user action,
            // so this will reset the time inactive calculation
            this.auth.setLastFetch(Date.now());
          })
          .finally(() => {
            resolve();
          });
      }, 200);
    });
  }

  @task
  *renewToken() {
    yield this.handleRenew();
  }

  get queryParams() {
    // Bring user back to current page after login
    return { redirect_to: this.router.currentURL };
  }

  get showWarning() {
    const currentRoute = this.router.currentRouteName;
    if ('vault.cluster.oidc-provider' === currentRoute) {
      return false;
    }

    return !!this.args.expirationDate;
  }
}
