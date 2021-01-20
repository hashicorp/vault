import Model, { attr } from '@ember-data/model';
import { alias } from '@ember/object/computed';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

export default Model.extend({
  backend: attr('string', { readOnly: true }),
  secretPath: lazyCapabilities(apiPath`${'backend'}/static-roles/${'id'}`, 'backend', 'id'),
  canRead: alias('secretPath.canRead'),
  canEdit: alias('secretPath.canUpdate'),
});
