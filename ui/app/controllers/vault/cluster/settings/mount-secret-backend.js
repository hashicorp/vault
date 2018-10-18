import { inject as service } from '@ember/service';
import Controller from '@ember/controller';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';

const SUPPORTED_BACKENDS = supportedSecretBackends();

export default Controller.extend({
  wizard: service(),
  actions: {
    onMountSuccess: function(type, path) {
      let transition;
      if (SUPPORTED_BACKENDS.includes(type)) {
        transition = this.transitionToRoute('vault.cluster.secrets.backend.index', path);
      } else {
        transition = this.transitionToRoute('vault.cluster.secrets.backends');
      }
      return transition.followRedirects().then(() => {
        this.get('wizard').transitionFeatureMachine(this.get('wizard.featureState'), 'CONTINUE', type);
      });
    },
    onConfigError: function(modelId) {
      return this.transitionToRoute('vault.cluster.settings.configure-secret-backend', modelId);
    },
  },
});
