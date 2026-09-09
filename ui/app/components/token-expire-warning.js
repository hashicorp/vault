/**
 * Copyright IBM Corp. 2016, 2025
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
  @service capabilities;
  @tracked canDismiss = true;
  @tracked canRenewSelf = true;

  constructor(owner, args) {
    super(owner, args);
    this.fetchRenewCapability();
  }

  async fetchRenewCapability() {
    const { canUpdate } = await this.capabilities.fetchPathCapabilities('auth/token/renew-self');
    this.canRenewSelf = canUpdate;
  }

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

  get canShowRenew() {
    return this.auth?.authData?.renewable && this.canRenewSelf;
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
