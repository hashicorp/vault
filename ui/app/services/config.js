import Service from '@ember/service';

export default Service.extend({
  featureFlags: null,
  setFeatureFlags(flags) {
    this.set('featureFlags', flags);
  },

  get managedNamespaceRoot() {
    if (this.featureFlags && this.featureFlags.includes('MANAGED_NAMESPACE')) {
      return 'admin';
    }
    return null;
  },
});
