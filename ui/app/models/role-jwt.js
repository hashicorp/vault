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
  providerName: computed('authUrl', function() {
    let { hostname } = parseURL(this.authUrl);
    let firstMatch = Object.keys(DOMAIN_STRINGS).find(name => hostname.includes(name));

    return DOMAIN_STRINGS[firstMatch] || null;
  }),
});
