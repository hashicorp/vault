import { inject as service } from '@ember/service';
import Component from '@ember/component';
import hbs from 'htmlbars-inline-precompile';
import { encodePath } from 'vault/utils/path-encoding-helpers';

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
      let params = this.get('params');
      if (this.encode) {
        params = params.map((param, index) => {
          if (index === 0 || typeof param !== 'string') {
            return param;
          }
          return encodePath(param);
        });
      }
      const queryParams = this.get('queryParams');
      if (queryParams) {
        params.push({ queryParams });
      }
      if (this.linkPrefix) {
        let targetRoute = this.params[0];
        targetRoute = `${this.linkPrefix}.${targetRoute}`;
        this.params[0] = targetRoute;
      }
      this.get('router').transitionTo(...params);
    }
  },
});

LinkedBlockComponent.reopenClass({
  positionalParams: 'params',
});

export default LinkedBlockComponent;
