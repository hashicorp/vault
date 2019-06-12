import Route from '@ember/routing/route';
import ListRoute from 'core/mixins/list-route';
import { inject as service } from '@ember/service';

export default Route.extend(ListRoute, {
  store: service(),
  model(params) {
    let model = [{ id: 'ARole' }];
    model.set('meta', { total: 1 });
    return model;
    return this.store
      .lazyPaginatedQuery('kmip/role', {
        responsePath: 'data.keys',
        page: params.page,
        pageFilter: params.pageFilter,
      })
      .catch(err => {
        if (err.httpStatus === 404) {
          return [];
        } else {
          throw err;
        }
      });
  },

  setupController(controller) {
    this._super(...arguments);
    controller.set('scope', this.paramsFor('scope').scope_name);
  },
});
