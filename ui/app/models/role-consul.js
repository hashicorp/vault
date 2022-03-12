import Model, { attr } from '@ember-data/model';
import { alias } from '@ember/object/computed';
import { computed } from '@ember/object';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

const TOKEN_TYPES = [
  {
    value: 'service_token',
    displayName: 'Service Token',
  },
  {
    value: 'management',
    displayName: 'Management Token',
  },
];
export default Model.extend({
  backend: attr('string', {
    readOnly: true,
  }),
  name: attr('string', {
    label: 'Role name',
    fieldValue: 'id',
    readOnly: true,
  }),
  useOpenAPI: false,
  // credentialTypes are for backwards compatibility.
  // we use this to populate "credentialType" in
  // the serializer. if there is more than one, the
  // show and edit pages will show a warning
  credentialTypes: attr('array', {
    readOnly: true,
  }),
  credentialType: attr('string', {
    defaultValue: 'service_token',
    possibleValues: TOKEN_TYPES,
  }),
  policies: attr({
    editType: 'stringArray',
    label: 'Policy Names',
  }),
  local: attr({
    editType: 'boolean',
    label: 'Local',
  }),
  fields: computed('credentialType', function() {
    let credentialType = this.credentialType;
    let keysForType = {
      service_token: ['name', 'credentialType', 'policies', 'local'],
      management_token: ['name', 'credentialType']
    };

    return expandAttributeMeta(this, keysForType[credentialType]);
  }),
  updatePath: lazyCapabilities(apiPath`${'backend'}/roles/${'id'}`, 'backend', 'id'),
  canDelete: alias('updatePath.canDelete'),
  canEdit: alias('updatePath.canUpdate'),
  canRead: alias('updatePath.canRead'),

  generatePath: lazyCapabilities(apiPath`${'backend'}/creds/${'id'}`, 'backend', 'id'),
  canGenerate: alias('generatePath.canUpdate'),
});
