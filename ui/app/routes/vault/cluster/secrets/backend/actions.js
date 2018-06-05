import EditBase from './secret-edit';
import utils from 'vault/lib/key-utils';

export default EditBase.extend({
  queryParams: {
    selectedAction: {
      replace: true,
    },
  },

  templateName: 'vault/cluster/secrets/backend/transitActionsLayout',

  beforeModel() {
    const { secret } = this.paramsFor(this.routeName);
    const parentKey = utils.parentKeyForKey(secret);
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    if (this.backendType(backend) !== 'transit') {
      if (parentKey) {
        return this.transitionTo('vault.cluster.secrets.backend.show', parentKey);
      } else {
        return this.transitionTo('vault.cluster.secrets.backend.show-root');
      }
    }
  },
  setupController(controller, model) {
    this._super(...arguments);
    const { selectedAction } = this.paramsFor(this.routeName);
    controller.set('selectedAction', selectedAction || model.secret.get('supportedActions.firstObject'));
  },
});
