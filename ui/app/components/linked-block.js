import { inject as service } from '@ember/service';
import Component from '@ember/component';
import hbs from 'htmlbars-inline-precompile';

let LinkedBlockComponent = Component.extend({
  router: service(),

  layout: hbs`{{yield}}`,

  classNames: 'linked-block',

  queryParams: null,

  click(event) {
    const $target = this.$(event.target);
    const isAnchorOrButton =
      $target.is('a') ||
      $target.is('button') ||
      $target.closest('button', event.currentTarget).length > 0 ||
      $target.closest('a', event.currentTarget).length > 0;
    if (!isAnchorOrButton) {
      const params = this.get('params');
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
