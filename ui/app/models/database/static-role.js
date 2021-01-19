import Model, { attr } from '@ember-data/model';
import { apiPath } from 'vault/macros/lazy-capabilities';
import attachCapabilities from 'vault/lib/attach-capabilities';

const ModelExport = Model.extend({
  backend: attr('string', { readOnly: true }),
});

export default attachCapabilities(ModelExport, {
  updatePath: apiPath`${'backend'}/static-roles/${'id'}`,
});
