import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
export default class MfaRoute extends Route {
  @service router;

  queryParams = {
    page: {
      refreshModel: true,
    },
  };

  model(params) {
    return this.store
      .lazyPaginatedQuery('mfa-method', {
        responsePath: 'data.keys',
        page: params.page || 1,
      })
      .then((model) => {
        return model;
      })
      .catch((err) => {
        if (err.httpStatus === 404) {
          return [];
        } else {
          throw err;
        }
      });
  }
  afterModel(model) {
    if (model.get('length') === 0) {
      this.router.transitionTo('vault.cluster.access.mfa.configure');
    }
  }
  setupController(controller, model) {
    controller.set('model', model);
  }
}
