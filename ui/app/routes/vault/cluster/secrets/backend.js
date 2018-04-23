import Ember from 'ember';
const { inject } = Ember;
export default Ember.Route.extend({
  flashMessages: inject.service(),
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
    let target = transition.targetName;
    let path = model && model.get('path');
    let type = model && model.get('type');
    if (type === 'kv' && model.get('options.version') === 2) {
      this.get('flashMessages').stickyInfo(
        `"${path}" is a newer version of the KV backend. The Vault UI does not currently support the additional versioning features. All actions taken through the UI in this engine will operate on the most recent version of a secret.`
      );
    }

    if (target === this.routeName) {
      return this.replaceWith('vault.cluster.secrets.backend.list-root', path);
    }
  },
});
