import { alias } from '@ember/object/computed';
import { computed } from '@ember/object';
import DS from 'ember-data';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

const { attr } = DS;

const CREDENTIAL_TYPES = [
  {
    value: 'iam_user',
    displayName: 'IAM User',
  },
  {
    value: 'assumed_role',
    displayName: 'Assumed Role',
  },
  {
    value: 'federation_token',
    displayName: 'Federation Token',
  },
];
export default DS.Model.extend({
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
    defaultValue: 'iam_user',
    possibleValues: CREDENTIAL_TYPES,
  }),
  roleArns: attr({
    editType: 'stringArray',
    label: 'Role ARNs',
  }),
  policyArns: attr({
    editType: 'stringArray',
    label: 'Policy ARNs',
  }),
  policyDocument: attr('string', {
    editType: 'json',
  }),
  fields: computed('credentialType', function() {
    let credentialType = this.credentialType;
    let keysForType = {
      iam_user: ['name', 'credentialType', 'policyArns', 'policyDocument'],
      assumed_role: ['name', 'credentialType', 'roleArns', 'policyDocument'],
      federation_token: ['name', 'credentialType', 'policyDocument'],
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
