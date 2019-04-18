import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import { hash } from 'rsvp';
import { getOwner } from '@ember/application';
import { singularize } from 'ember-inflector';
import { normalizePath } from 'vault/utils/path-encoding-helpers';
import ListRoute from 'vault/mixins/list-route';

export default Route.extend(ListRoute, {
  wizard: service(),
  pathHelp: service('path-help'),

  beforeModel() {
    const { item_type: itemType } = this.paramsFor(this.routeName);
    const { path: method } = this.paramsFor('vault.cluster.access.method');
    let methodModel = this.modelFor('vault.cluster.access.method');
    let { apiPath, type } = methodModel;
    let modelType = `generated-${singularize(itemType)}-${type}`;
    return this.pathHelp.getNewModel(modelType, getOwner(this), method, `${apiPath}${itemType}/example`);
  },

  model(params) {
    let { item_type: itemType } = this.paramsFor(this.routeName);
    itemType = normalizePath(itemType) || '';
    itemType = singularize(itemType);
    const { path: method } = this.paramsFor('vault.cluster.access.method');
    const methodModel = this.modelFor('vault.cluster.access.method');
    let { type } = methodModel;

    return hash({
      itemType,
      items: this.store
        .lazyPaginatedQuery(`generated-${itemType}-${type}`, {
          id: itemType,
          method,
          responsePath: 'data.keys',
          page: params.page,
          pageFilter: params.pageFilter,
        })
        .then(model => {
          this.set('has404', false);
          return model;
        })
        .catch(err => {
          // if we're at the root we don't want to throw
          if (methodModel && err.httpStatus === 404 && itemType === '') {
            return [];
          } else {
            // else we're throwing and dealing with this in the error action
            throw err;
          }
        }),
    });
  },
  setupController(controller, resolvedModel) {
    let itemParams = this.paramsFor(this.routeName);
    let itemType = resolvedModel.itemType;
    const { path: method } = this.paramsFor('vault.cluster.access.method');
    let modelType = `generated-${singularize(itemType)}-${type}`;
    const methodModel = this.store.peekRecord(modelType, type);
    let { type } = methodModel;
    let model = resolvedModel.items;
    let has404 = this.get('has404');
    // only clear store cache if this is a new model
    if (secret !== controller.get('baseKey.id')) {
      this.store.clearAllDatasets();
    }

    controller.set('hasModel', true);
    controller.setProperties({
      model,
      has404,
      method,
      methodModel,
      baseKey: { id: itemType },
      methodType: type,
    });
    if (!has404) {
      const pageFilter = itemParams.pageFilter;
      let filter;
      if (itemType) {
        filter = itemType + (pageFilter || '');
      } else if (pageFilter) {
        filter = pageFilter;
      }
      controller.setProperties({
        filter: filter || '',
        page: model.get('meta.currentPage') || 1,
      });
    }
  },

  resetController(controller, isExiting) {
    this._super(...arguments);
    if (isExiting) {
      controller.set('pageFilter', null);
      controller.set('filter', null);
    }
  },
});
