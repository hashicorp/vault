import Route from '@ember/routing/route';
import ListRoute from 'core/mixins/list-route';
import { inject as service } from '@ember/service';

export default Route.extend(ListRoute, {
  store: service(),
  secretMountPath: service(),
  credParams() {
    let { role_name: role, scope_name: scope } = this.paramsFor('credentials');
    return {
      role,
      scope,
    };
  },
  model(params) {
    let { role, scope } = this.credParams();
    return this.store
      .lazyPaginatedQuery('kmip/credential', {
        role,
        scope,
        backend: this.secretMountPath.currentPath,
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
    let { scope_name: scope, role_name: role } = this.paramsFor('credentials');
    controller.setProperties({ role, scope });
  },
});
