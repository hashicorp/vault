import { alias } from '@ember/object/computed';
import { computed } from '@ember/object';
import DS from 'ember-data';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import fieldToAttrs from 'vault/utils/field-to-attrs';

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

  fieldGroups: computed('backend', 'merged', function() {
    const groups = [
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
      let allFields = [];
      for (let group in groups) {
        let fieldName = Object.keys(groups[group])[0];
        allFields.concat(groups[group][fieldName]);
      }
      let otherFields = this.newFields.filter(field => {
        !allFields.includes(field) && !excludedFields.includes(field);
      });
      if (otherFields.length) {
        groups.default.concat(otherFields);
      }
    }

    if (this.newFields) {
      let allFields = [];
      for (let group in groups) {
        let fieldName = Object.keys(groups[group])[0];
        allFields.concat(groups[group][fieldName]);
      }
      let otherFields = this.newFields.filter(field => {
        !allFields.includes(field);
      });
      if (otherFields.length) {
        groups.default.concat(otherFields);
      }
    }

    return fieldToAttrs(this, groups);
  }),
});
