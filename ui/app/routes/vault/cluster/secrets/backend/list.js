import Ember from 'ember';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';

const SUPPORTED_BACKENDS = supportedSecretBackends();

export default Ember.Route.extend({
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

  templateName: 'vault/cluster/secrets/backend/list',

  beforeModel() {
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    const backendModel = this.store.peekRecord('secret-engine', backend);
    const type = backendModel && backendModel.get('type');
    if (!type || !SUPPORTED_BACKENDS.includes(type)) {
      return this.transitionTo('vault.cluster.secrets');
    }
  },

  getModelType(backend, tab) {
    let backendModel = this.store.peekRecord('secret-engine', backend);
    let type = backendModel.get('type');
    let types = {
      transit: 'transit-key',
      ssh: 'role-ssh',
      aws: 'role-aws',
      pki: tab === 'certs' ? 'pki-certificate' : 'role-pki',
      // secret or secret-v2
      kv: backendModel.get('modelTypeForKV'),
      generic: backendModel.get('modelTypeForKV'),
    };
    return types[type];
  },

  model(params) {
    const secret = params.secret ? params.secret : '';
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    const backendModel = this.modelFor('vault.cluster.secrets.backend');
    return Ember.RSVP.hash({
      secret,
      secrets: this.store
        .lazyPaginatedQuery(this.getModelType(backend, params.tab), {
          id: secret,
          backend,
          responsePath: 'data.keys',
          page: params.page,
          pageFilter: params.pageFilter,
          size: 100,
        })
        .then(model => {
          this.set('has404', false);
          return model;
        })
        .catch(err => {
          if (backendModel && err.httpStatus === 404 && secret === '') {
            return [];
          } else {
            throw err;
          }
        }),
    });
  },

  afterModel(model) {
    const { tab } = this.paramsFor(this.routeName);
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    if (!tab || tab !== 'certs') {
      return;
    }
    return Ember.RSVP
      .all(
        // these ids are treated specially by vault's api, but it's also
        // possible that there is no certificate for them in order to know,
        // we fetch them specifically on the list page, and then unload the
        // records if there is no `certificate` attribute on the resultant model
        ['ca', 'crl', 'ca_chain'].map(id => this.store.queryRecord('pki-certificate', { id, backend }))
      )
      .then(
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
    let { backend } = this.paramsFor('vault.cluster.secrets.backend');
    let backendModel = this.store.peekRecord('secret-engine', backend);
    let has404 = this.get('has404');
    controller.set('hasModel', true);
    controller.setProperties({
      model,
      has404,
      backend,
      backendModel,
      baseKey: { id: secret },
      backendType: backendModel.get('type'),
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
      controller.set('filter', '');
    }
  },

  actions: {
    error(error, transition) {
      const { secret } = this.paramsFor(this.routeName);
      const { backend } = this.paramsFor('vault.cluster.secrets.backend');

      Ember.set(error, 'secret', secret);
      Ember.set(error, 'isRoot', true);
      Ember.set(error, 'backend', backend);
      const hasModel = this.controllerFor(this.routeName).get('hasModel');
      // only swallow the error if we have a previous model
      if (hasModel && error.httpStatus === 404) {
        this.set('has404', true);
        transition.abort();
      } else {
        return true;
      }
    },

    willTransition(transition) {
      window.scrollTo(0, 0);
      if (transition.targetName !== this.routeName) {
        this.store.clearAllDatasets();
      }
      return true;
    },
    reload() {
      this.refresh();
      this.store.clearAllDatasets();
    },
  },
});
