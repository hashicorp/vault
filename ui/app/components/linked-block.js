import { inject as service } from '@ember/service';
import Component from '@ember/component';
import hbs from 'htmlbars-inline-precompile';
import { encodePath } from 'vault/utils/path-encoding-helpers';

let LinkedBlockComponent = Component.extend({
  router: service(),

  layout: hbs`{{yield}}`,

  classNames: 'linked-block',

  queryParams: null,

  encode: false,

  click(event) {
    const $target = this.$(event.target);
    const isAnchorOrButton =
      $target.is('a') ||
      $target.is('button') ||
      $target.closest('button', event.currentTarget).length > 0 ||
      $target.closest('a', event.currentTarget).length > 0;
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
      this.get('router').transitionTo(...params);
    }
  },
});

LinkedBlockComponent.reopenClass({
  positionalParams: 'params',
});

export default LinkedBlockComponent;
