import DS from 'ember-data';
import { computed } from '@ember/object';
import parseURL from 'vault/utils/parse-url';
const { attr } = DS;

const DOMAIN_STRINGS = {
  google: 'Google',
  ping: 'Ping',
  okta: 'Okta',
  auth0: 'Auth0',
};
export default DS.Model.extend({
  authUrl: attr('string'),
  providerMatch: computed('authUrl', function() {
    let { hostname } = parseURL(this.authUrl);
    return Object.keys(DOMAIN_STRINGS).find(name => hostname.includes(name));
  }),

  providerName: computed('providerMatch', function() {
    return DOMAIN_STRINGS[this.providerMatch] || null;
  }),

  providerButtonComponent: computed('providerName', function() {
    let { providerMatch } = this;
    return providerMatch ? `auth-button-${providerMatch}` : null;
  }),
});
