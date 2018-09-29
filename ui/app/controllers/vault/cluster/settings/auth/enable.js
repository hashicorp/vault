import { inject as service } from '@ember/service';
import Controller from '@ember/controller';

export default Controller.extend({
  wizard: service(),
  actions: {
    onMountSuccess: function(type) {
      let transition = this.transitionToRoute('vault.cluster.access.methods');
      return transition.followRedirects().then(() => {
        this.get('wizard').transitionFeatureMachine(this.get('wizard.featureState'), 'CONTINUE', type);
      });
    },
    onConfigError: function(modelId) {
      return this.transitionToRoute('vault.cluster.settings.auth.configure', modelId);
    },
  },
});
