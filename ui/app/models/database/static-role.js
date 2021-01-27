import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { alias } from '@ember/object/computed';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

export default Model.extend({
  backend: attr('string', { readOnly: true }),
  secretPath: lazyCapabilities(apiPath`${'backend'}/static-roles/${'id'}`, 'backend', 'id'),
  canEditRole: computed.or('secretPath.canUpdate', 'secretPath.canCreate'),
  credentialPath: lazyCapabilities(apiPath`${'backend'}/static-creds/${'id'}`, 'backend', 'id'),
  canGenerateCredentials: alias('credentialPath.canRead'),
});
