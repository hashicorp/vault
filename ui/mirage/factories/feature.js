import { Factory } from 'ember-cli-mirage';

export default Factory.extend({
  feature_flags() {
    return []; // VAULT_CLOUD_ADMIN_NAMESPACE
  },
});
