import Service from '@ember/service';
import { computed } from '@ember/object';

const FLAGS = {
  vaultCloudNamespace: 'VAULT_CLOUD_ADMIN_NAMESPACE',
};

export default Service.extend({
  featureFlags: null,
  setFeatureFlags(flags) {
    this.set('featureFlags', flags);
  },

  managedNamespaceRoot: computed('featureFlags', function() {
    const flags = this.featureFlags;
    if (flags && flags.includes(FLAGS.vaultCloudNamespace)) {
      return 'admin';
    }
    return null;
  }),
});
