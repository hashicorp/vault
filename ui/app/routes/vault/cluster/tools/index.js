import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import { toolsActions } from 'vault/helpers/tools-actions';

export default Route.extend({
  currentCluster: service(),
  beforeModel(transition) {
    const currentCluster = this.get('currentCluster.cluster.name');
    const supportedActions = toolsActions();
    if (transition.targetName === this.routeName) {
      transition.abort();
      return this.replaceWith('vault.cluster.tools.tool', currentCluster, supportedActions[0]);
    }
  },
});
