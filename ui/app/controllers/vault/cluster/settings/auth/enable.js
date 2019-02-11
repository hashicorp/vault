import { inject as service } from '@ember/service';
import Controller from '@ember/controller';

export default Controller.extend({
  wizard: service(),
  actions: {
    onMountSuccess: function(type, path) {
      // We have to remove the trailing '/' from the path to succcessfully redirect with the right params.
      const authPath = path.slice(0, -1);
      let transition = this.transitionToRoute('vault.cluster.settings.auth.configure', authPath);
      return transition.followRedirects();
    },
    onConfigError: function(modelId) {
      return this.transitionToRoute('vault.cluster.settings.auth.configure', modelId);
    },
  },
});
