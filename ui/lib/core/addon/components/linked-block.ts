/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import routerLookup from 'core/utils/router-lookup';

import type RouterService from '@ember/routing/router-service';

type TransitionArgs = [string, ...unknown[]];

interface Args {
  params: TransitionArgs;
  queryParams?: Record<string, unknown>;
  linkPrefix?: string;
  encode?: boolean;
  disabled?: boolean;
}

/**
 * @module LinkedBlock
 * LinkedBlock components are linkable divs that yield any content nested within them. They are often used in list views such as when listing the secret engines.
 *
 * @example
 * <LinkedBlock @params={{array "vault.cluster.secrets.backend.show" "my-secret-path"}} class="list-item-row" >
 * My wrapped content
 * </LinkedBlock>
 *
 *
 * @param {Array} params=null - These are values sent to the router's transitionTo method.  First item is route, second is the optional path.
 * @param {Object} [queryParams=null] - queryParams can be passed via this property. It needs to be an object.
 * @param {String} [linkPrefix=null] - Overwrite the params with custom route.  Needed for use in engines (KMIP and PKI). ex: vault.cluster.secrets.backend.kmip
 * @param {Boolean} [encode=false] - Encode the path.
 * @param {boolean} [disabled] - disable the link -- prevents on click and removes linked-block hover styling
 */

export default class LinkedBlockComponent extends Component<Args> {
  get router(): RouterService {
    return routerLookup(this);
  }

  @action
  onClick(event: MouseEvent): void {
    if (this.args.disabled) {
      return;
    }

    const target = event.target;
    if (!(target instanceof Element)) {
      return;
    }

    const isAnchorOrButton =
      target.tagName === 'A' ||
      target.tagName === 'BUTTON' ||
      !!target.closest('button') ||
      !!target.closest('a');
    if (isAnchorOrButton) {
      return;
    }

    let params = [...this.args.params] as TransitionArgs;
    if (this.args.encode) {
      params = params.map((param, index) => {
        if (index === 0 || typeof param !== 'string') {
          return param;
        }
        return encodePath(param);
      }) as TransitionArgs;
    }

    const queryParams = this.args.queryParams;
    if (queryParams) {
      params.push({ queryParams });
    }

    if (this.args.linkPrefix) {
      const [targetRoute, ...rest] = params;
      params = [`${this.args.linkPrefix}.${targetRoute}`, ...rest] as TransitionArgs;
    }

    this.router.transitionTo(...params);
  }
}
