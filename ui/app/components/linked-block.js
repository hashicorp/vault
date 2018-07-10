import Ember from 'ember';
import hbs from 'htmlbars-inline-precompile';

let LinkedBlockComponent = Ember.Component.extend({
  layout: hbs`{{yield}}`,

  classNames: 'linked-block',

  routing: Ember.inject.service('-routing'),
  queryParams: null,

  click(event) {
    const $target = this.$(event.target);
    const isAnchorOrButton =
      $target.is('a') ||
      $target.is('button') ||
      $target.closest('button', event.currentTarget).length > 0 ||
      $target.closest('a', event.currentTarget).length > 0;
    if (!isAnchorOrButton) {
      const router = this.get('routing.router');
      const params = this.get('params');
      const queryParams = this.get('queryParams');
      if (queryParams) {
        params.push({ queryParams });
      }
      router.transitionTo.apply(router, params);
    }
  },
});

LinkedBlockComponent.reopenClass({
  positionalParams: 'params',
});

export default LinkedBlockComponent;
