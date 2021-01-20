import Model, { attr } from '@ember-data/model';
import { alias } from '@ember/object/computed';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
// ARG TODO confirm that canRead is the correct policy path
export default Model.extend({
  backend: attr('string', { readOnly: true }),
  secretPath: lazyCapabilities(apiPath`${'backend'}/roles/${'id'}`, 'backend', 'id'),
  canRead: alias('secretPath.canRead'),
  canEdit: alias('secretPath.canUpdate'),
  credentialPath: lazyCapabilities(apiPath`${'backend'}/creds/${'id'}`, 'backend', 'id'),
  canGenerateCredentials: alias('credentialPath.canRead'),
});
