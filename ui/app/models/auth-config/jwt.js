import { computed } from '@ember/object';
import DS from 'ember-data';
import AuthConfig from '../auth-config';
import fieldToAttrs from 'vault/utils/field-to-attrs';

const { attr } = DS;

export default AuthConfig.extend({
  oidcDiscoveryUrl: attr('string', {
    label: 'OIDC discovery URL',
    helpText:
      'The OIDC discovery URL, without any .well-known component (base path). Cannot be used with jwt_validation_pubkeys',
  }),

  oidcDiscoveryCaPem: attr('string', {
    label: 'OIDC discovery CA PEM',
    editType: 'file',
    helpText:
      'The CA certificate or chain of certificates, in PEM format, to use to validate connections to the OIDC Discovery URL. If not set, system certificates are used',
  }),
  jwtValidationPubkeys: attr({
    label: 'JWT validation public keys',
    editType: 'stringArray',
  }),
  boundIssuer: attr('string', {
    helpText: 'The value against which to match the iss claim in a JWT',
  }),
  fieldGroups: computed(function() {
    const groups = [
      {
        default: ['oidcDiscoveryUrl'],
      },
      {
        'JWT Options': ['oidcDiscoveryCaPem', 'jwtValidationPubkeys', 'boundIssuer'],
      },
    ];
    return fieldToAttrs(this, groups);
  }),
});
