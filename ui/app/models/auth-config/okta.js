import Ember from 'ember';
import DS from 'ember-data';
import AuthConfig from '../auth-config';
import fieldToAttrs from 'vault/utils/field-to-attrs';

const { attr } = DS;
const { computed } = Ember;

export default AuthConfig.extend({
  orgName: attr('string', {
    label: 'Organization Name',
    helpText: 'Name of the organization to be used in the Okta API',
  }),
  apiToken: attr('string', {
    label: 'API Token',
    helpText:
      'Okta API token. This is required to query Okta for user group membership. If this is not supplied only locally configured groups will be enabled.',
  }),
  baseUrl: attr('string', {
    label: 'Base URL',
    helpText:
      'If set, will be used as the base domain for API requests. Examples are okta.com, oktapreview.com, and okta-emea.com',
  }),
  bypassOktaMfa: attr('boolean', {
    defaultValue: false,
    label: 'Bypass Okta MFA',
    helpText:
      "Useful if Vault's built-in MFA mechanisms. Will also cause certain other statuses to be ignored, such as PASSWORD_EXPIRED",
  }),
  fieldGroups: computed(function() {
    const groups = [
      {
        default: ['orgName'],
      },
      {
        Options: ['apiToken', 'baseUrl', 'bypassOktaMfa'],
      },
    ];
    return fieldToAttrs(this, groups);
  }),
});
