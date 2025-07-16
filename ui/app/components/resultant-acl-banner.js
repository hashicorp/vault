/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { PERMISSIONS_BANNER_STATES } from 'vault/services/permissions';

export default class ResultantAclBannerComponent extends Component {
  @service namespace;
  @service router;
  @tracked hideBanner = false;

  get ns() {
    return this.namespace.path || 'root';
  }

  get queryParams() {
    // Bring user back to current page after login
    return { redirect_to: this.router.currentURL };
  }

  get title() {
    return this.args.failType === PERMISSIONS_BANNER_STATES.noAccess
      ? 'You do not have access to this namespace'
      : 'Resultant ACL check failed';
  }

  get message() {
    return this.args.failType === PERMISSIONS_BANNER_STATES.noAccess
      ? 'Log into the namespace directly, or contact your administrator if you think you should have access.'
      : "Links might be shown that you don't have access to. Contact your administrator to update your policy.";
  }
}
