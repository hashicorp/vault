import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
export default Route.extend({
  flashMessages: service(),
  oldModel: null,
  model(params) {
    let { backend } = params;
    return this.store
      .query('secret-engine', {
        path: backend,
      })
      .then(model => {
        if (model) {
          return model.get('firstObject');
        }
      });
  },

  afterModel(model, transition) {
    let path = model && model.get('path');
    if (transition.targetName === this.routeName) {
      return this.replaceWith('vault.cluster.secrets.backend.list-root', path);
    }
  },
});
