/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { later } from '@ember/runloop';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class SidebarUserMenuComponent extends Component {
  @service auth;
  @service session;
  @service currentCluster;
  @service router;
  @service namespace;

  @tracked fakeRenew = false;

  get authData() {
    return this.session.data.authenticated;
  }
  get hasEntityId() {
    // root users will not have an entity_id because they are not associated with an entity.
    // in order to use the MFA end user setup they need an entity_id
    return !!this.authData?.entity_id;
  }
  get isUserpass() {
    return this.authData?.backend?.type === 'userpass';
  }

  get isRenewing() {
    return this.fakeRenew || this.auth.isRenewing;
  }

  transitionToRoute() {
    this.router.transitionTo(...arguments);
  }

  @action
  renewToken() {
    this.fakeRenew = true;
    const { authenticator, backend, token, userRootNamespace } = this.authData;
    later(() => {
      this.session.authenticate(
        authenticator,
        { token },
        {
          renew: true,
          backend: backend.mountPath,
          namespace: userRootNamespace,
        }
      );
      this.auth.renew().then(() => {
        this.fakeRenew = this.auth.isRenewing;
      });
    }, 200);
  }

  @action
  revokeToken() {
    this.session.invalidate({ revoke: true }).then(() => {
      this.transitionToRoute('vault.cluster.logout');
    });
  }
}
