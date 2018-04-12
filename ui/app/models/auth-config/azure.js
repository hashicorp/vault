import Ember from 'ember';
import DS from 'ember-data';

import AuthConfig from '../auth-config';
import fieldToAttrs from 'vault/utils/field-to-attrs';

const { attr } = DS;
const { computed } = Ember;

export default AuthConfig.extend({
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
    const groups = [
      { default: ['tenantId', 'resource'] },
      {
        'Azure Options': ['clientId', 'clientSecret'],
      },
    ];
    return fieldToAttrs(this, groups);
  }),
});
