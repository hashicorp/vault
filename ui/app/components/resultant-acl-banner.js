/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

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
}
