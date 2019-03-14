import { computed } from '@ember/object';
import DS from 'ember-data';
import AuthConfig from '../auth-config';
import fieldToAttrs from 'vault/utils/field-to-attrs';
import { combineFieldGroups } from 'vault/utils/openapi-to-attrs';

const { attr } = DS;

export default AuthConfig.extend({
  useOpenAPI: true,
  oidcDiscoveryUrl: attr('string', {
    label: 'OIDC discovery URL',
    helpText:
      'The OIDC discovery URL, without any .well-known component (base path). Cannot be used with jwt_validation_pubkeys',
  }),

  oidcClientId: attr('string', {
    label: 'OIDC client ID',
  }),

  oidcClientSecret: attr('string', {
    label: 'OIDC client secret',
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

  jwtSupportedAlgs: attr({
    label: 'JWT supported algorithms',
  }),
  boundIssuer: attr('string', {
    helpText: 'The value against which to match the iss claim in a JWT',
  }),
  fieldGroups: computed(function() {
    let type = this.constructor.modelName.split('/')[1].toUpperCase();
    let groups = [
      {
        default: ['oidcDiscoveryUrl', 'defaultRole'],
      },
      {
        [`${type} Options`]: [
          'oidcClientId',
          'oidcClientSecret',
          'oidcDiscoveryCaPem',
          'jwtValidationPubkeys',
          'jwtSupportedAlgs',
          'boundIssuer',
        ],
      },
    ];

    if (this.newFields) {
      groups = combineFieldGroups(groups, this.newFields, []);
    }
    return fieldToAttrs(this, groups);
  }),
});
