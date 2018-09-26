import Route from '@ember/routing/route';

export default Route.extend({
  queryParams: {
    page: {
      refreshModel: true,
    },
    pageFilter: {
      refreshModel: true,
    },
  },

  model() {
    return this.store.findAll('auth-method');
  },
});
