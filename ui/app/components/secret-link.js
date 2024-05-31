/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { debug } from '@ember/debug';

/**
 * @module SecretLink
 *
 * @example
 * ```js
 * <SecretLink
 *   @mode="edit"
 *   @backend="backend-path"
 *   @secret="my/secret path"
 *   @queryParams={{hash tab="version"}}
 *   @disabled={{true}}
 * />
 *
 * @param {string} mode - *required* controls the route link. added to the base route vault.cluster.secrets.backend
 * @param {string} backend - *required* backend path. Is encoded in the component
 * @param {string} secret - secret path. Is encoded in the component
 * @param {object} queryParams - params passed to the link
 * @param {boolean} disabled - passed to LinkTo to disable link
 * @param {CallableFunction} onLinkClick - side effect when link is clicked
 */
export default class SecretLink extends Component {
  get link() {
    const { mode, secret, backend } = this.args;
    if (!backend) {
      debug(`Arg "backend" missing from secret-link with mode: ${mode} secret: ${secret}`);
    }
    const route = `vault.cluster.secrets.backend.${mode}`;
    const models = backend ? [encodePath(backend)] : [];
    if ((mode !== 'versions' && !secret) || secret === ' ') {
      return { route: `${route}-root`, models };
    } else {
      models.push(encodePath(secret));
      return { route, models };
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
