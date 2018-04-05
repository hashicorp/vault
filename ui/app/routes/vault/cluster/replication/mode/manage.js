import Ember from 'ember';
import { replicationActionForMode } from 'vault/helpers/replication-action-for-mode';

const pathForAction = (action, replicationMode, clusterMode) => {
  let path;
  if (action === 'reindex' || action === 'recover') {
    path = `sys/replication/${action}`;
  } else {
    path = `sys/replication/${replicationMode}/${clusterMode}/${action}`;
  }
  return path;
};

export default Ember.Route.extend({
  store: Ember.inject.service(),
  model() {
    const store = this.get('store');
    const model = this.modelFor('vault.cluster.replication.mode');

    const replicationMode = this.paramsFor('vault.cluster.replication.mode').replication_mode;
    const clusterMode = model.get(replicationMode).get('modeForUrl');
    const actions = replicationActionForMode([replicationMode, clusterMode]);
    return Ember.RSVP
      .all(
        actions.map(action => {
          return store.findRecord('capabilities', pathForAction(action)).then(capability => {
            model.set(`can${Ember.String.camelize(action)}`, capability.get('canUpdate'));
          });
        })
      )
      .then(() => {
        return model;
      });
  },

  beforeModel() {
    const model = this.modelFor('vault.cluster.replication.mode');
    const replicationMode = this.paramsFor('vault.cluster.replication.mode').replication_mode;
    if (
      model.get(replicationMode).get('replicationDisabled') ||
      model.get(replicationMode).get('replicationUnsupported')
    ) {
      return this.transitionTo('vault.cluster.replication.mode', replicationMode);
    }
  },
});
