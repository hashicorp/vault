import { set } from '@ember/object';
import { hash, all } from 'rsvp';
import Route from '@ember/routing/route';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
import { inject as service } from '@ember/service';
import { normalizePath } from 'vault/utils/path-encoding-helpers';

const SUPPORTED_BACKENDS = supportedSecretBackends();

export default Route.extend({
  templateName: 'vault/cluster/secrets/backend/list',
  pathHelp: service('path-help'),
  queryParams: {
    page: {
      refreshModel: true,
    },
    pageFilter: {
      refreshModel: true,
    },
    tab: {
      refreshModel: true,
    },
  },

  modelTypeForTransform(tab) {
    let modelType;
    switch (tab) {
      case 'role':
        modelType = 'transform/role';
        break;
      case 'template':
        modelType = 'transform/template';
        break;
      case 'alphabet':
        modelType = 'transform/alphabet';
        break;
      default:
        modelType = 'transform'; // CBS TODO: transform/transformation
        break;
    }
    return modelType;
  },

  secretParam() {
    let { secret } = this.paramsFor(this.routeName);
    return secret ? normalizePath(secret) : '';
  },

  enginePathParam() {
    let { backend } = this.paramsFor('vault.cluster.secrets.backend');
    return backend;
  },

  beforeModel() {
    let secret = this.secretParam();
    let backend = this.enginePathParam();
    let { tab } = this.paramsFor('vault.cluster.secrets.backend.list-root');
    let secretEngine = this.store.peekRecord('secret-engine', backend);
    let type = secretEngine && secretEngine.get('engineType');
    if (!type || !SUPPORTED_BACKENDS.includes(type)) {
      return this.transitionTo('vault.cluster.secrets');
    }
    if (this.routeName === 'vault.cluster.secrets.backend.list' && !secret.endsWith('/')) {
      return this.replaceWith('vault.cluster.secrets.backend.list', secret + '/');
    }
    let modelType = this.getModelType(backend, tab);
    return this.pathHelp.getNewModel(modelType, backend).then(() => {
      this.store.unloadAll('capabilities');
    });
  },

  getModelType(backend, tab) {
    let secretEngine = this.store.peekRecord('secret-engine', backend);
    let type = secretEngine.get('engineType');
    let types = {
      // ARG TODO, list both static and dynamic roles
      database: tab === 'role' ? 'database/role' : 'database/connection',
      transit: 'transit-key',
      ssh: 'role-ssh',
      transform: this.modelTypeForTransform(tab),
      aws: 'role-aws',
      pki: tab === 'certs' ? 'pki-certificate' : 'role-pki',
      // secret or secret-v2
      cubbyhole: 'secret',
      kv: secretEngine.get('modelTypeForKV'),
      generic: secretEngine.get('modelTypeForKV'),
    };
    return types[type];
  },

  async model(params) {
    const secret = this.secretParam() || '';
    const backend = this.enginePathParam();
    const backendModel = this.modelFor('vault.cluster.secrets.backend');
    const modelType = this.getModelType(backend, params.tab);

    if (modelType === 'database/role') {
      let queryResponse = await this.store.query('database/role', {});
      let staticRoles = await this.store.query('database/static-role', {});
      // concat the internal models of static roles to the role.content internal models
      let combinedContent = [];
      queryResponse.content.forEach(item => {
        combinedContent.push(item);
      });
      staticRoles.content.forEach(item => {
        combinedContent.push(item);
      });
      queryResponse.content = combinedContent;
      return hash({
        secret,
        secrets: this.store.lazyPaginatedQueryTwoModels(queryResponse, {
          id: secret,
          backend,
          responsePath: 'data.keys',
          page: params.page || 1,
          pageFilter: params.pageFilter,
        }),
      });
    }

    return hash({
      secret,
      secrets: this.store
        .lazyPaginatedQuery(modelType, {
          id: secret,
          backend,
          responsePath: 'data.keys',
          page: params.page || 1,
          pageFilter: params.pageFilter,
        })
        .then(model => {
          this.set('has404', false);
          return model;
        })
        .catch(err => {
          console.log(err, 'error');
          // if we're at the root we don't want to throw
          if (backendModel && err.httpStatus === 404 && secret === '') {
            return [];
          } else {
            // else we're throwing and dealing with this in the error action
            throw err;
          }
        }),
    });
  },

  afterModel(model) {
    const { tab } = this.paramsFor(this.routeName);
    const backend = this.enginePathParam();
    if (!tab || tab !== 'certs') {
      return;
    }
    return all(
      // these ids are treated specially by vault's api, but it's also
      // possible that there is no certificate for them in order to know,
      // we fetch them specifically on the list page, and then unload the
      // records if there is no `certificate` attribute on the resultant model
      ['ca', 'crl', 'ca_chain'].map(id => this.store.queryRecord('pki-certificate', { id, backend }))
    ).then(
      results => {
        results.rejectBy('certificate').forEach(record => record.unloadRecord());
        return model;
      },
      () => {
        return model;
      }
    );
  },

  setupController(controller, resolvedModel) {
    let secretParams = this.paramsFor(this.routeName);
    let secret = resolvedModel.secret;
    let model = resolvedModel.secrets;
    let backend = this.enginePathParam();
    let backendModel = this.store.peekRecord('secret-engine', backend);
    let has404 = this.has404;
    // only clear store cache if this is a new model
    if (secret !== controller.get('baseKey.id')) {
      this.store.clearAllDatasets();
    }

    controller.set('hasModel', true);
    controller.setProperties({
      model,
      has404,
      backend,
      backendModel,
      baseKey: { id: secret },
      backendType: backendModel.get('engineType'),
    });
    if (!has404) {
      const pageFilter = secretParams.pageFilter;
      let filter;
      if (secret) {
        filter = secret + (pageFilter || '');
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

  actions: {
    error(error, transition) {
      let secret = this.secretParam();
      let backend = this.enginePathParam();
      let is404 = error.httpStatus === 404;
      let hasModel = this.controllerFor(this.routeName).get('hasModel');

      // this will occur if we've deleted something,
      // and navigate to its parent and the parent doesn't exist -
      // this if often the case with nested keys in kv-like engines
      if (transition.data.isDeletion && is404) {
        throw error;
      }
      set(error, 'secret', secret);
      set(error, 'isRoot', true);
      set(error, 'backend', backend);
      // only swallow the error if we have a previous model
      if (hasModel && is404) {
        this.set('has404', true);
        transition.abort();
        return false;
      }
      return true;
    },

    willTransition(transition) {
      window.scrollTo(0, 0);
      if (transition.targetName !== this.routeName) {
        this.store.clearAllDatasets();
      }
      return true;
    },
    reload() {
      this.store.clearAllDatasets();
      this.refresh();
    },
  },
});
