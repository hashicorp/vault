import DS from 'ember-data';
import { computed } from '@ember/object';
import parseURL from 'vault/utils/parse-url';
const { attr } = DS;

const DOMAIN_STRINGS = {
  github: 'GitHub',
  gitlab: 'GitLab',
  google: 'Google',
  ping: 'Ping',
  okta: 'Okta',
  auth0: 'Auth0',
};

const PROVIDER_WITH_LOGO = ['GitLab', 'Google', 'Auth0'];

export { DOMAIN_STRINGS, PROVIDER_WITH_LOGO };

export default DS.Model.extend({
  authUrl: attr('string'),

  providerName: computed('authUrl', function() {
    let { hostname } = parseURL(this.authUrl);
    let firstMatch = Object.keys(DOMAIN_STRINGS).find(name => hostname.includes(name));
    return DOMAIN_STRINGS[firstMatch] || null;
  }),

  providerButtonComponent: computed('providerName', function() {
    let { providerName } = this;
    return PROVIDER_WITH_LOGO.includes(providerName) ? `auth-button-${providerName.toLowerCase()}` : null;
  }),
});
