import Ember from 'ember';
const { inject } = Ember;
export default Ember.Route.extend({
  flashMessages: inject.service(),
  beforeModel(transition) {
    const target = transition.targetName;
    const { backend } = this.paramsFor(this.routeName);
    const backendModel = this.store.peekRecord('secret-engine', backend);
    const type = backendModel && backendModel.get('type');
    if (type === 'kv' && backendModel.get('isVersioned')) {
      this.get('flashMessages').stickyInfo(
        `"${backend}" is a versioned kv secrets engine. The Vault UI does not currently support the additional versioning features. All actions taken through the UI in this engine will operate on the most recent version of a secret.`
      );
    }

    if (target === this.routeName) {
      return this.replaceWith('vault.cluster.secrets.backend.list-root', backend);
    }
  },
});
