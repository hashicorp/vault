import Ember from 'ember';

const SUPPORTED_SECRET_BACKENDS = ['aws', 'cubbyhole', 'generic', 'kv', 'pki', 'ssh', 'transit'];

export function supportedSecretBackends() {
  return SUPPORTED_SECRET_BACKENDS;
}

export default Ember.Helper.helper(supportedSecretBackends);
