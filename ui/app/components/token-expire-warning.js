/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';

export default class TokenExpireWarning extends Component {
  @service router;

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
