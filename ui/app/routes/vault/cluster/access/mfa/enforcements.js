import Route from '@ember/routing/route';

export default class MfaEnforcementsRoute extends Route {
  queryParams = {
    page: {
      refreshModel: true,
    },
  };

  model(params) {
    return this.store
      .lazyPaginatedQuery('mfa-login-enforcement', {
        responsePath: 'data.keys',
        page: params.page || 1,
      })
      .catch((err) => {
        if (err.httpStatus === 404) {
          return [];
        } else {
          throw err;
        }
      });
  }
  setupController(controller, model) {
    controller.set('model', model);
  }
}
