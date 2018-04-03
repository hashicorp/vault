import Ember from 'ember';
import { toolsActions } from 'vault/helpers/tools-actions';

export default Ember.Route.extend({
  currentCluster: Ember.inject.service(),
  beforeModel(transition) {
    const currentCluster = this.get('currentCluster.cluster.name');
    const supportedActions = toolsActions();
    if (transition.targetName === this.routeName) {
      transition.abort();
      return this.replaceWith('vault.cluster.tools.tool', currentCluster, supportedActions[0]);
    }
  },
});
