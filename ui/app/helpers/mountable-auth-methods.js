import Ember from 'ember';

const MOUNTABLE_AUTH_METHODS = [
  {
    displayName: 'AppRole',
    value: 'approle',
    type: 'approle',
  },
  {
    displayName: 'AWS',
    value: 'aws',
    type: 'aws',
  },
  {
    displayName: 'Azure',
    value: 'azure',
    type: 'azure',
  },
  {
    displayName: 'Google Cloud',
    value: 'gcp',
    type: 'gcp',
  },
  {
    displayName: 'Kubernetes',
    value: 'kubernetes',
    type: 'kubernetes',
  },
  {
    displayName: 'GitHub',
    value: 'github',
    type: 'github',
  },
  {
    displayName: 'LDAP',
    value: 'ldap',
    type: 'ldap',
  },
  {
    displayName: 'Okta',
    value: 'okta',
    type: 'okta',
  },
  {
    displayName: 'RADIUS',
    value: 'radius',
    type: 'radius',
  },
  {
    displayName: 'TLS Certificates',
    value: 'cert',
    type: 'cert',
  },
  {
    displayName: 'Username & Password',
    value: 'userpass',
    type: 'userpass',
  },
];

export function methods() {
  return MOUNTABLE_AUTH_METHODS;
}

export default Ember.Helper.helper(methods);
