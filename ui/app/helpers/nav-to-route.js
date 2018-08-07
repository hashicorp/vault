import Ember from 'ember';

const { Helper, inject } = Ember;

export default Helper.extend({
  router: inject.service(),

  compute([routeName, ...models], { replace = false }) {
    return () => {
      const router = this.get('router');
      const method = replace ? router.replaceWith : router.transitionTo;
      return method.call(router, routeName, ...models);
    };
  },
});
