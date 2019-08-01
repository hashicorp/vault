import { alias } from '@ember/object/computed';
import { computed } from '@ember/object';
import DS from 'ember-data';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

const { attr } = DS;

// these arrays define the order in which the fields will be displayed
// see
// https://github.com/hashicorp/vault/blob/master/builtin/logical/ssh/path_roles.go#L542 for list of fields for each key type
const OTP_FIELDS = [
  'name',
  'keyType',
  'defaultUser',
  'adminUser',
  'port',
  'allowedUsers',
  'cidrList',
  'excludeCidrList',
];
const CA_FIELDS = [
  'name',
  'keyType',
  'allowUserCertificates',
  'allowHostCertificates',
  'defaultUser',
  'allowedUsers',
  'allowedDomains',
  'ttl',
  'maxTtl',
  'allowedCriticalOptions',
  'defaultCriticalOptions',
  'allowedExtensions',
  'defaultExtensions',
  'allowBareDomains',
  'allowSubdomains',
  'allowUserKeyIds',
  'keyIdFormat',
];

export default DS.Model.extend({
  useOpenAPI: true,
  getHelpUrl: function(backend) {
    return `/v1/${backend}/roles/example?help=1`;
  },
  zeroAddress: attr('boolean', {
    readOnly: true,
  }),
  backend: attr('string', {
    readOnly: true,
  }),
  name: attr('string', {
    label: 'Role Name',
    fieldValue: 'id',
    readOnly: true,
  }),
  keyType: attr('string', {
    possibleValues: ['ca', 'otp'], //overriding the API which also lists 'dynamic' as a type though it is deprecated
  }),
  adminUser: attr('string', {
    helpText: 'Username of the admin user at the remote host',
  }),
  defaultUser: attr('string', {
    helpText: "Username to use when one isn't specified",
  }),
  allowedUsers: attr('string', {
    helpText:
      'Create a whitelist of users that can use this key (e.g. `admin, dev`, or use `*` to allow all)',
  }),
  allowedDomains: attr('string', {
    helpText:
      'List of domains for which a client can request a certificate (e.g. `example.com`, or `*` to allow all)',
  }),
  cidrList: attr('string', {
    helpText: 'List of CIDR blocks for which this role is applicable',
  }),
  excludeCidrList: attr('string', {
    helpText: 'List of CIDR blocks that are not accepted by this role',
  }),
  port: attr('number', {
    helpText: 'Port number for the SSH connection (default is `22`)',
  }),
  allowedCriticalOptions: attr('string', {
    helpText: 'List of critical options that certificates have when signed',
  }),
  defaultCriticalOptions: attr('object', {
    helpText: 'Map of critical options certificates should have if none are provided when signing',
  }),
  allowedExtensions: attr('string', {
    helpText: 'List of extensions that certificates can have when signed',
  }),
  defaultExtensions: attr('object', {
    helpText: 'Map of extensions certificates should have if none are provided when signing',
  }),
  allowUserCertificates: attr('boolean', {
    helpText: 'Specifies if certificates are allowed to be signed for us as a user',
  }),
  allowHostCertificates: attr('boolean', {
    helpText: 'Specifies if certificates are allowed to be signed for us as a host',
  }),
  allowBareDomains: attr('boolean', {
    helpText:
      'Specifies if host certificates that are requested are allowed to use the base domains listed in Allowed Domains',
  }),
  allowSubdomains: attr('boolean', {
    helpText:
      'Specifies if host certificates that are requested are allowed to be subdomains of those listed in Allowed Domains',
  }),
  allowUserKeyIds: attr('boolean', {
    helpText: 'Specifies if users can override the key ID for a signed certificate with the "key_id" field',
  }),
  keyIdFormat: attr('string', {
    helpText: 'When supplied, this value specifies a custom format for the key id of a signed certificate',
  }),

  attrsForKeyType: computed('keyType', function() {
    const keyType = this.get('keyType');
    let keys = keyType === 'ca' ? CA_FIELDS.slice(0) : OTP_FIELDS.slice(0);
    return expandAttributeMeta(this, keys);
  }),

  updatePath: lazyCapabilities(apiPath`${'backend'}/roles/${'id'}`, 'backend', 'id'),
  canDelete: alias('updatePath.canDelete'),
  canEdit: alias('updatePath.canUpdate'),
  canRead: alias('updatePath.canRead'),

  generatePath: lazyCapabilities(apiPath`${'backend'}/creds/${'id'}`, 'backend', 'id'),
  canGenerate: alias('generatePath.canUpdate'),

  signPath: lazyCapabilities(apiPath`${'backend'}/sign/${'id'}`, 'backend', 'id'),
  canSign: alias('signPath.canUpdate'),

  zeroAddressPath: lazyCapabilities(apiPath`${'backend'}/config/zeroaddress`, 'backend'),
  canEditZeroAddress: alias('zeroAddressPath.canUpdate'),
});
