import Model, { attr } from '@ember-data/model';
import { alias } from '@ember/object/computed';
import { computed } from '@ember/object';
import fieldToAttrs, { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

// these arrays define the order in which the fields will be displayed
// see
// https://github.com/hashicorp/vault/blob/main/builtin/logical/ssh/path_roles.go#L542 for list of fields for each key type
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
  'allowedUsersTemplate',
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
  'notBeforeDuration',
  'algorithmSigner',
];

export default Model.extend({
  useOpenAPI: true,
  getHelpUrl: function (backend) {
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
      'Create a list of users who are allowed to use this key (e.g. `admin, dev`, or use `*` to allow all.)',
  }),
  allowedUsersTemplate: attr('boolean', {
    helpText:
      'Specifies that Allowed users can be templated e.g. {{identity.entity.aliases.mount_accessor_xyz.name}}',
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
  algorithmSigner: attr('string', {
    helpText: 'When supplied, this value specifies a signing algorithm for the key',
  }),

  showFields: computed('keyType', function () {
    const keyType = this.keyType;
    const keys = keyType === 'ca' ? CA_FIELDS.slice(0) : OTP_FIELDS.slice(0);
    return expandAttributeMeta(this, keys);
  }),

  fieldGroups: computed('keyType', function () {
    const numRequired = this.keyType === 'otp' ? 3 : 4;
    const fields = this.keyType === 'otp' ? [...OTP_FIELDS] : [...CA_FIELDS];
    const defaultFields = fields.splice(0, numRequired);
    const groups = [
      { default: defaultFields },
      {
        Options: [...fields],
      },
    ];
    return fieldToAttrs(this, groups);
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
