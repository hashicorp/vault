import { inject as service } from '@ember/service';
import Component from '@ember/component';
import hbs from 'htmlbars-inline-precompile';
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
 */

let LinkedBlockComponent = Component.extend({
  router: service(),

  layout: hbs`{{yield}}`,

  classNames: 'linked-block',

  queryParams: null,
  linkPrefix: null,

  encode: false,

  click(event) {
    const $target = event.target;
    const isAnchorOrButton =
      $target.tagName === 'A' ||
      $target.tagName === 'BUTTON' ||
      $target.closest('button') ||
      $target.closest('a');
    if (!isAnchorOrButton) {
      let params = this.params;
      if (this.encode) {
        params = params.map((param, index) => {
          if (index === 0 || typeof param !== 'string') {
            return param;
          }
          return encodePath(param);
        });
      }
      const queryParams = this.queryParams;
      if (queryParams) {
        params.push({ queryParams });
      }
      if (this.linkPrefix) {
        let targetRoute = this.params[0];
        targetRoute = `${this.linkPrefix}.${targetRoute}`;
        this.params[0] = targetRoute;
      }
      this.router.transitionTo(...params);
    }
  },
});

LinkedBlockComponent.reopenClass({
  positionalParams: 'params',
});

export default LinkedBlockComponent;
