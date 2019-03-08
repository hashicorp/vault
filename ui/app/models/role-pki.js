import { alias } from '@ember/object/computed';
import { computed } from '@ember/object';
import DS from 'ember-data';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import fieldToAttrs from 'vault/utils/field-to-attrs';
import { combineFieldGroups } from 'vault/utils/openapi-to-attrs';
const { attr } = DS;

export default DS.Model.extend({
  backend: attr('string', {
    readOnly: true,
  }),
  name: attr('string', {
    label: 'Role name',
    fieldValue: 'id',
    readOnly: true,
  }),
  useOpenAPI: true,
  getHelpUrl: function(backend) {
    return `/v1/${backend}/roles/example?help=1`;
  },
  updatePath: lazyCapabilities(apiPath`${'backend'}/roles/${'id'}`, 'backend', 'id'),
  canDelete: alias('updatePath.canDelete'),
  canEdit: alias('updatePath.canUpdate'),
  canRead: alias('updatePath.canRead'),

  generatePath: lazyCapabilities(apiPath`${'backend'}/issue/${'id'}`, 'backend', 'id'),
  canGenerate: alias('generatePath.canUpdate'),

  signPath: lazyCapabilities(apiPath`${'backend'}/sign/${'id'}`, 'backend', 'id'),
  canSign: alias('signPath.canUpdate'),

  signVerbatimPath: lazyCapabilities(apiPath`${'backend'}/sign-verbatim/${'id'}`, 'backend', 'id'),
  canSignVerbatim: alias('signVerbatimPath.canUpdate'),

  fieldGroups: computed(function() {
    let groups = [
      { default: ['name', 'keyType'] },
      {
        Options: [
          'keyBits',
          'ttl',
          'maxTtl',
          'allowAnyName',
          'enforceHostnames',
          'allowIpSans',
          'requireCn',
          'useCsrCommonName',
          'useCsrSans',
          'ou',
          'organization',
          'keyUsage',
          'allowedOtherSans',
          'notBeforeDuration',
        ],
      },
      {
        'Address Options': ['country', 'locality', 'province', 'streetAddress', 'postalCode'],
      },
      {
        'Domain Handling': [
          'allowLocalhost',
          'allowBareDomains',
          'allowSubdomains',
          'allowGlobDomains',
          'allowedDomains',
        ],
      },
      {
        'Extended Key Usage': [
          'serverFlag',
          'clientFlag',
          'codeSigningFlag',
          'emailProtectionFlag',
          'extKeyUsageOids',
        ],
      },
      {
        Advanced: ['generateLease', 'noStore', 'basicConstraintsValidForNonCa', 'policyIdentifiers'],
      },
    ];
    let excludedFields = ['extKeyUsage'];
    if (this.newFields) {
      groups = combineFieldGroups(groups, this.newFields, excludedFields);
    }
    return fieldToAttrs(this, groups);
  }),
});
