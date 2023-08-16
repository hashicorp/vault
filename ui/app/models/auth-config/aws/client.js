import { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import AuthConfig from '../../auth-config';
import fieldToAttrs from 'vault/utils/field-to-attrs';

export default AuthConfig.extend({
  secretKey: attr('string'),
  accessKey: attr('string'),
  endpoint: attr('string', {
    label: 'EC2 Endpoint',
  }),
  iamEndpoint: attr('string', {
    label: 'IAM Endpoint',
  }),
  stsEndpoint: attr('string', {
    label: 'STS Endpoint',
  }),
  iamServerIdHeaderValue: attr('string', {
    label: 'IAM Server ID Header Value',
  }),

  fieldGroups: computed(function () {
    const groups = [
      { default: ['accessKey', 'secretKey'] },
      { 'AWS Options': ['endpoint', 'iamEndpoint', 'stsEndpoint', 'iamServerIdHeaderValue'] },
    ];

    return fieldToAttrs(this, groups);
  }),
});
