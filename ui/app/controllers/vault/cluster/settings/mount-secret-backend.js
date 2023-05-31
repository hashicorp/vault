import Controller from '@ember/controller';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';

const SUPPORTED_BACKENDS = supportedSecretBackends();

export default Controller.extend({
  actions: {
    onMountSuccess: function (type, path) {
      let transition;
      if (SUPPORTED_BACKENDS.includes(type)) {
        if (type === 'kmip') {
          transition = this.transitionToRoute('vault.cluster.secrets.backend.kmip.scopes', path);
        } else if (type === 'keymgmt') {
          transition = this.transitionToRoute('vault.cluster.secrets.backend.index', path, {
            queryParams: { tab: 'provider' },
          });
        } else {
          transition = this.transitionToRoute('vault.cluster.secrets.backend.index', path);
        }
      } else {
        transition = this.transitionToRoute('vault.cluster.secrets.backends');
      }
      return transition.followRedirects();
    },
  },
});
