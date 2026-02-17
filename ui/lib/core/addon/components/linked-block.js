/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import routerLookup from 'core/utils/router-lookup';

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

export default class LinkedBlockComponent extends Component {
  get router() {
    return routerLookup(this);
  }

  @action
  onClick(event) {
    if (!this.args.disabled) {
      const $target = event.target;
      const isAnchorOrButton =
        $target.tagName === 'A' ||
        $target.tagName === 'BUTTON' ||
        $target.closest('button') ||
        $target.closest('a');
      if (!isAnchorOrButton) {
        let params = this.args.params;
        if (this.args.encode) {
          params = params.map((param, index) => {
            if (index === 0 || typeof param !== 'string') {
              return param;
            }
            return encodePath(param);
          });
        }
        const queryParams = this.args.queryParams;
        if (queryParams) {
          params.push({ queryParams });
        }
        if (this.args.linkPrefix) {
          let targetRoute = this.args.params[0];
          targetRoute = `${this.args.linkPrefix}.${targetRoute}`;
          this.args.params[0] = targetRoute;
        }
        this.router.transitionTo(...params);
      }
    }
  }
}
