import { inject as service } from '@ember/service';
import { set } from '@ember/object';
import Route from '@ember/routing/route';
import RSVP from 'rsvp';
import DS from 'ember-data';
import UnloadModelRoute from 'vault/mixins/unload-model-route';
import { getOwner } from '@ember/application';

export default Route.extend(UnloadModelRoute, {
  modelPath: 'model.model',
  pathHelp: service('path-help'),

  modelType(backendType, section) {
    const MODELS = {
      'aws-client': 'auth-config/aws/client',
      'aws-identity-whitelist': 'auth-config/aws/identity-whitelist',
      'aws-roletag-blacklist': 'auth-config/aws/roletag-blacklist',
      'azure-configuration': 'auth-config/azure',
      'github-configuration': 'auth-config/github',
      'gcp-configuration': 'auth-config/gcp',
      'jwt-configuration': 'auth-config/jwt',
      'oidc-configuration': 'auth-config/oidc',
      'kubernetes-configuration': 'auth-config/kubernetes',
      'ldap-configuration': 'auth-config/ldap',
      'okta-configuration': 'auth-config/okta',
      'radius-configuration': 'auth-config/radius',
    };
    return MODELS[`${backendType}-${section}`];
  },

  beforeModel() {
    const { section_name } = this.paramsFor(this.routeName);
    if (section_name === 'options') {
      return;
    }
    const { method } = this.paramsFor('vault.cluster.settings.auth.configure');
    const backend = this.modelFor('vault.cluster.settings.auth.configure');
    const modelType = this.modelType(backend.type, section_name);
    let owner = getOwner(this);
    return this.pathHelp.getNewModel(modelType, owner, method);
  },

  model(params) {
    const backend = this.modelFor('vault.cluster.settings.auth.configure');
    const { section_name: section } = params;
    if (section === 'options') {
      return RSVP.hash({
        model: backend,
        section,
      });
    }
    const modelType = this.modelType(backend.get('type'), section);
    if (!modelType) {
      const error = new DS.AdapterError();
      set(error, 'httpStatus', 404);
      throw error;
    }
    const model = this.store.peekRecord(modelType, backend.id);
    if (model) {
      return RSVP.hash({
        model,
        section,
      });
    }
    return this.store
      .findRecord(modelType, backend.id)
      .then(config => {
        config.set('backend', backend);
        return RSVP.hash({
          model: config,
          section,
        });
      })
      .catch(e => {
        let config;
        // if you haven't saved a config, the API 404s, so create one here to edit and return it
        if (e.httpStatus === 404) {
          config = this.store.createRecord(modelType, {
            id: backend.id,
          });
          config.set('backend', backend);

          return RSVP.hash({
            model: config,
            section,
          });
        }
        throw e;
      });
  },

  actions: {
    willTransition() {
      if (this.currentModel.model.constructor.modelName !== 'auth-method') {
        this.unloadModel();
        return true;
      }
    },
  },
});
