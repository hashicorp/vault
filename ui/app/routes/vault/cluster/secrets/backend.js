import Ember from 'ember';
const { inject } = Ember;
export default Ember.Route.extend({
  flashMessages: inject.service(),
  beforeModel(transition) {
    const target = transition.targetName;
    const { backend } = this.paramsFor(this.routeName);
    const backendModel = this.store.peekRecord('secret-engine', backend);
    const type = backendModel && backendModel.get('type');
    if (type === 'kv' && backendModel.get('options.version') === 2) {
      this.get('flashMessages').stickyInfo(
        `"${backend}" is a newer version of the KV backend. The Vault UI does not currently support the additional versioning features. All actions taken through the UI in this engine will operate on the most recent version of a secret.`
      );
    }

    if (target === this.routeName) {
      return this.replaceWith('vault.cluster.secrets.backend.list-root', backend);
    }
  },
});
