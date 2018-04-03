import Ember from 'ember';

const { Helper, inject } = Ember;

export default Helper.extend({
  routing: inject.service('-routing'),

  compute([routeName, ...models], { replace = false }) {
    return () => {
      const router = this.get('routing.router');
      const method = replace ? router.replaceWith : router.transitionTo;
      return method.call(router, routeName, ...models);
    };
  },
});
