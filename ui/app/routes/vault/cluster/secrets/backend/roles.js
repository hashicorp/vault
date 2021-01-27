import { hash } from 'rsvp';
import Route from '@ember/routing/route';
import { normalizePath } from 'vault/utils/path-encoding-helpers';

import ListRoute from './list';

export default ListRoute.extend({
  templateName: 'vault/cluster/secrets/backend/roles',
  // queryParams: {
  //   page: {
  //     refreshModel: true,
  //   },
  //   pageFilter: {
  //     refreshModel: true,
  //   },
  // },

  // secretParam() {
  //   let { secret } = this.paramsFor(this.routeName);
  //   return secret ? normalizePath(secret) : '';
  // },

  // enginePathParam() {
  //   let { backend } = this.paramsFor('vault.cluster.secrets.backend');
  //   return backend;
  // },

  getModelType(backend, tab) {
    return 'database/role';
  },
  // async model(params) {
  //   const secret = this.secretParam() || '';
  //   const backend = this.enginePathParam();
  //   const backendModel = this.modelFor('vault.cluster.secrets.backend');
  //   // This route is currently only used to display database roles
  //   const combinedModels = ['database/role', 'database/static-role'];
  //     return hash({
  //       secret,
  //       secrets: this.store
  //         .lazyPaginatedQueryTwoModels(combinedModels, {
  //           id: secret,
  //           backend,
  //           responsePath: 'data.keys',
  //           page: params.page || 1,
  //           pageFilter: params.pageFilter,
  //         })
  //         .then(model => {
  //           this.set('has404', false);
  //           // extra filtering here?
  //           console.log({ model });
  //           return model;
  //         })
  //         .catch(err => {
  //           console.error('error getting combined list');
  //           // if we're at the root we don't want to throw
  //           if (backendModel && err.httpStatus === 404 && secret === '') {
  //             return [];
  //           } else {
  //             // else we're throwing and dealing with this in the error action
  //             throw err;
  //           }
  //         }),
  //     });
  // },

  // setupController(controller, resolvedModel) {
  //   let secretParams = this.paramsFor(this.routeName);
  //   let secret = resolvedModel.secret;
  //   let model = resolvedModel.secrets;
  //   let backend = this.enginePathParam();
  //   let backendModel = this.store.peekRecord('secret-engine', backend);
  //   let has404 = this.has404;
  //   let root = {
  //     label: backend,
  //     model: backend,
  //     path: 'vault.cluster.secrets.backend.list-root',
  //     text: backend,
  //   }
  //   controller.set('hasModel', true);
  //   controller.setProperties({
  //     root,
  //     model,
  //     has404,
  //     backend,
  //     backendModel,
  //     baseKey: { id: secret },
  //     backendType: backendModel.get('engineType'),
  //   });
  //   if (!has404) {
  //     const pageFilter = secretParams.pageFilter;
  //     let filter;
  //     if (secret) {
  //       filter = secret + (pageFilter || '');
  //     } else if (pageFilter) {
  //       filter = pageFilter;
  //     }
  //     controller.setProperties({
  //       filter: filter || '',
  //       page: model.meta.currentPage || 1,
  //     });
  //   }
  // },
});
