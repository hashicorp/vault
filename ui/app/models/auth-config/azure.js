import { computed } from '@ember/object';
import DS from 'ember-data';
import AuthConfig from '../auth-config';
import { combineFieldGroups } from 'vault/utils/openapi-to-attrs';
import fieldToAttrs from 'vault/utils/field-to-attrs';

const { attr } = DS;

export default AuthConfig.extend({
  useOpenAPI: true,
  tenantId: attr('string', {
    label: 'Tenant ID',
    helpText: 'The tenant ID for the Azure Active Directory organization',
  }),
  resource: attr('string', {
    helpText: 'The configured URL for the application registered in Azure Active Directory',
  }),
  clientId: attr('string', {
    label: 'Client ID',
    helpText:
      'The client ID for credentials to query the Azure APIs. Currently read permissions to query compute resources are required.',
  }),
  clientSecret: attr('string', {
    helpText: 'The client secret for credentials to query the Azure APIs',
  }),

  googleCertsEndpoint: attr('string'),

  fieldGroups: computed(function() {
    let groups = [
      { default: ['tenantId', 'resource'] },
      {
        'Azure Options': ['clientId', 'clientSecret'],
      },
    ];
    if (this.newFields) {
      groups = combineFieldGroups(groups, this.newFields, []);
    }

    return fieldToAttrs(this, groups);
  }),
});
