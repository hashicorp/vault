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
];

export default DS.Model.extend({
  useOpenAPI: true,
  getHelpUrl: function(backend) {
    return `/v1/${backend}?help=1`;
  },

  attrsForKeyType: computed('keyType', function() {
    const keyType = this.get('keyType');
    let keys = keyType === 'ca' ? CA_FIELDS.slice(0) : OTP_FIELDS.slice(0);
    return expandAttributeMeta(this, keys);
  }),

  updatePath: lazyCapabilities(apiPath`${'backend'}/transforms/${'id'}`, 'backend', 'id'),
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
