// This model represents the capabilities on a given `path`
// `path` is also the primaryId
// https://www.vaultproject.io/docs/concepts/policies.html#capabilities

import { computed } from '@ember/object';
import DS from 'ember-data';

const { attr } = DS;

const SUDO_PATHS = [
  'sys/seal',
  'sys/replication/performance/primary/secondary-token',
  'sys/replication/dr/primary/secondary-token',
  'sys/replication/reindex',
  'sys/leases/lookup/',
];

const SUDO_PATH_PREFIXES = ['sys/leases/revoke-prefix', 'sys/leases/revoke-force'];

export { SUDO_PATHS, SUDO_PATH_PREFIXES };

const computedCapability = function(capability) {
  return computed('path', 'capabilities', 'capabilities.[]', function() {
    const capabilities = this.get('capabilities');
    const path = this.get('path');
    if (!capabilities) {
      return false;
    }
    if (capabilities.includes('root')) {
      return true;
    }
    if (capabilities.includes('deny')) {
      return false;
    }
    // if the path is sudo protected, they'll need sudo + the appropriate capability
    if (SUDO_PATHS.includes(path) || SUDO_PATH_PREFIXES.find(item => path.startsWith(item))) {
      return capabilities.includes('sudo') && capabilities.includes(capability);
    }
    return capabilities.includes(capability);
  });
};

export default DS.Model.extend({
  path: attr('string'),
  capabilities: attr('array'),
  canSudo: computedCapability('sudo'),
  canRead: computedCapability('read'),
  canCreate: computedCapability('create'),
  canUpdate: computedCapability('update'),
  canDelete: computedCapability('delete'),
  canList: computedCapability('list'),
  allowedParameters: attr(),
  deniedParameters: attr(),
});
