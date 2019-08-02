import { camelize } from '@ember/string';
import { all } from 'rsvp';
import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import { replicationActionForMode } from 'replication/helpers/replication-action-for-mode';

const pathForAction = (action, replicationMode, clusterMode) => {
  let path;
  if (action === 'reindex' || action === 'recover') {
    path = `sys/replication/${action}`;
  } else {
    path = `sys/replication/${replicationMode}/${clusterMode}/${action}`;
  }
  return path;
};

export default Route.extend({
  store: service(),
  model() {
    const store = this.get('store');
    const model = this.modelFor('mode');

    const replicationMode = this.paramsFor('mode').replication_mode;
    const clusterMode = model.get(replicationMode).get('modeForUrl');
    const actions = replicationActionForMode([replicationMode, clusterMode]);
    return all(
      actions.map(action => {
        return store.findRecord('capabilities', pathForAction(action)).then(capability => {
          model.set(`can${camelize(action)}`, capability.get('canUpdate'));
        });
      })
    ).then(() => {
      return model;
    });
  },

  beforeModel() {
    const model = this.modelFor('mode');
    const replicationMode = this.paramsFor('mode').replication_mode;
    if (
      model.get(replicationMode).get('replicationDisabled') ||
      model.get(replicationMode).get('replicationUnsupported')
    ) {
      return this.transitionTo('mode', replicationMode);
    }
  },
});
