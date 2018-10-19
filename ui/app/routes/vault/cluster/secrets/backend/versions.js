import Route from '@ember/routing/route';
import utils from 'vault/lib/key-utils';
import UnloadModelRoute from 'vault/mixins/unload-model-route';

export default Route.extend(UnloadModelRoute, {
  templateName: 'vault/cluster/secrets/backend/versions',

  beforeModel() {
    let backendModel = this.modelFor('vault.cluster.secrets.backend');
    const { secret } = this.paramsFor(this.routeName);
    const parentKey = utils.parentKeyForKey(secret);
    if (backendModel.get('isV2KV')) {
      return;
    }
    if (parentKey) {
      return this.transitionTo('vault.cluster.secrets.backend.list', parentKey);
    } else {
      return this.transitionTo('vault.cluster.secrets.backend.list-root');
    }
  },

  model(params) {
    let { secret } = params;
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    return this.store.queryRecord('secret-v2', { id: secret, backend });
  },
});
