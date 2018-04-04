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

  capabilities(secret) {
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    const path = backend + '/' + secret;
    return this.store.findRecord('capabilities', path);
  },

  getModelType(backend, tab) {
    const types = {
      transit: 'transit-key',
      ssh: 'role-ssh',
      aws: 'role-aws',
      cubbyhole: 'secret-cubbyhole',
      pki: tab === 'certs' ? 'pki-certificate' : 'role-pki',
    };
    const backendModel = this.store.peekRecord('secret-engine', backend);
    return types[backendModel.get('type')] || 'secret';
  },

  model(params) {
    const secret = params.secret ? params.secret : '';
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    const backends = this.modelFor('vault.cluster.secrets').mapBy('id');
    return Ember.RSVP.hash({
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
          if (backends.includes(backend) && err.httpStatus === 404 && secret === '') {
            return [];
          } else {
            throw err;
          }
        }),
      capabilities: this.capabilities(secret),
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

  setupController(controller, model) {
    const secretParams = this.paramsFor(this.routeName);
    const secret = secretParams.secret ? secretParams.secret : '';
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    const backendModel = this.store.peekRecord('secret-engine', backend);
    const has404 = this.get('has404');
    controller.set('hasModel', true);
    controller.setProperties({
      model: model.secrets,
      capabilities: model.capabilities,
      baseKey: { id: secret },
      has404,
      backend,
      backendModel,
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
        page: model.secrets.get('meta.currentPage') || 1,
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
      const backends = this.modelFor('vault.cluster.secrets').mapBy('id');

      Ember.set(error, 'secret', secret);
      Ember.set(error, 'isRoot', true);
      Ember.set(error, 'hasBackend', backends.includes(backend));
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
  },
});
