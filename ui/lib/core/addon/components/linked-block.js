import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { encodePath } from 'vault/utils/path-encoding-helpers';

/**
 * @module LinkedBlock
 * LinkedBlock components are linkable divs that yield any content nested within them. They are often used in list views such as when listing the secret engines.
 *
 * @example
 * ```js
 * <LinkedBlock
 *  @params={{array 'vault.cluster.secrets.backend.show 'my-secret-path'}}
 *  @queryParams={{hash version=1}}
 *  @class="list-item-row"
 *  data-test-list-item-link
 *  >
 * // Use any wrapped content here
 * </LinkedBlock>
 * ```
 *
 * @param {Array} params=null - These are values sent to the router's transitionTo method.  First item is route, second is the optional path.
 * @param {Object} [queryParams=null] - queryParams can be passed via this property. It needs to be an object.
 * @param {String} [linkPrefix=null] - Overwrite the params with custom route.  See KMIP.
 * @param {Boolean} [encode=false] - Encode the path.
 * @param {boolean} [disabled] - disable the link -- prevents on click and removes linked-block hover styling
 */

export default class LinkedBlockComponent extends Component {
  @service router;

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
