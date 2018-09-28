import { set } from '@ember/object';
import Route from '@ember/routing/route';
import DS from 'ember-data';

const CONFIGURABLE_BACKEND_TYPES = ['aws', 'ssh', 'pki'];

export default Route.extend({
  model() {
    const { backend } = this.paramsFor(this.routeName);
    return this.store.query('secret-engine', { path: backend }).then(modelList => {
      let model = modelList && modelList.get('firstObject');
      if (!model || !CONFIGURABLE_BACKEND_TYPES.includes(model.get('type'))) {
        const error = new DS.AdapterError();
        set(error, 'httpStatus', 404);
        throw error;
      }
      return this.store.findRecord('secret-engine', backend).then(
        () => {
          return model;
        },
        () => {
          return model;
        }
      );
    });
  },

  afterModel(model, transition) {
    const type = model.get('type');
    if (type === 'pki') {
      if (transition.targetName === this.routeName) {
        return this.transitionTo(`${this.routeName}.section`, 'cert');
      } else {
        return;
      }
    }
    if (type === 'aws') {
      return this.store
        .queryRecord('secret-engine', {
          backend: model.id,
          type,
        })
        .then(() => model, () => model);
    }
    return model;
  },

  setupController(controller, model) {
    if (model.get('publicKey')) {
      controller.set('configured', true);
    }
    return this._super(...arguments);
  },

  resetController(controller, isExiting) {
    if (isExiting) {
      controller.reset();
    }
  },

  actions: {
    refreshRoute() {
      this.refresh();
    },
  },
});
