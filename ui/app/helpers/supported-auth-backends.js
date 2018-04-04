import Ember from 'ember';

const SUPPORTED_AUTH_BACKENDS = [
  {
    type: 'token',
    description: 'Token authentication.',
    tokenPath: 'id',
    displayNamePath: 'display_name',
  },
  {
    type: 'userpass',
    description: 'A simple username and password backend.',
    tokenPath: 'client_token',
    displayNamePath: 'metadata.username',
  },
  {
    type: 'LDAP',
    description: 'LDAP authentication.',
    tokenPath: 'client_token',
    displayNamePath: 'metadata.username',
  },
  {
    type: 'Okta',
    description: 'Authenticate with your Okta username and password.',
    tokenPath: 'client_token',
    displayNamePath: 'metadata.username',
  },
  {
    type: 'GitHub',
    description: 'GitHub authentication.',
    tokenPath: 'client_token',
    displayNamePath: ['metadata.org', 'metadata.username'],
  },
];

export function supportedAuthBackends() {
  return SUPPORTED_AUTH_BACKENDS;
}

export default Ember.Helper.helper(supportedAuthBackends);
