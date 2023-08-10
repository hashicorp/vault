/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class SecretLink extends Component {
  get link() {
    const { mode, secret } = this.args;
    const route = `vault.cluster.secrets.backend.${mode}`;
    if ((mode !== 'versions' && !secret) || secret === ' ') {
      return { route: `${route}-root`, models: [] };
    } else {
      return { route, models: [encodePath(secret)] };
    }
  }
  get query() {
    const qp = this.args.queryParams || {};
    return qp.isQueryParams ? qp.values : qp;
  }

  @action
  onLinkClick() {
    if (this.args.onLinkClick) {
      this.args.onLinkClick(...arguments);
    }
  }
}
