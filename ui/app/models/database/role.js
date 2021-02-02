import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { alias } from '@ember/object/computed';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

export default Model.extend({
  backend: attr('string', { readOnly: true }),
  name: attr('string'),
  // CBS TODO: Should we consolidate type and path since they're basically the same thing?
  type: attr('string'),
  path: attr('string', { readOnly: true }),

  secretPath: lazyCapabilities(apiPath`${'backend'}/${'path'}/${'id'}`, 'backend', 'path', 'id'),
  canEditRole: computed.or('secretPath.canUpdate', 'secretPath.canCreate'),
  credentialPath: lazyCapabilities(apiPath`${'backend'}/creds/${'id'}`, 'backend', 'id'),
  canGenerateCredentials: alias('credentialPath.canRead'),
});
