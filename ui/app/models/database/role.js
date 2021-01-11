import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { apiPath } from 'vault/macros/lazy-capabilities';
// import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import attachCapabilities from 'vault/lib/attach-capabilities';

const ModelExport = Model.extend({
  // used for getting appropriate options for backend
  // TODO ARG SEE API DOCS
  // https://www.vaultproject.io/api-docs/secret/databases#create-role
  backend: attr('string', { readOnly: true }),
});

export default attachCapabilities(ModelExport, {
  // ARG TODO: configures a role
  updatePath: apiPath`${'backend'}/roles/${'id'}`,
});
