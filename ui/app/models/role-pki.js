import Ember from 'ember';
import DS from 'ember-data';
import { queryRecord } from 'ember-computed-query';

const { attr } = DS;
const { computed, get } = Ember;

export default DS.Model.extend({
  backend: attr('string', {
    readOnly: true,
  }),
  name: attr('string', {
    label: 'Role name',
    fieldValue: 'id',
    readOnly: true,
  }),
  keyType: attr('string', {
    possibleValues: ['rsa', 'ec'],
  }),
  ttl: attr({
    label: 'TTL',
    editType: 'ttl',
  }),
  maxTtl: attr({
    label: 'Max TTL',
    editType: 'ttl',
  }),
  allowLocalhost: attr('boolean', {}),
  allowedDomains: attr('string', {}),
  allowBareDomains: attr('boolean', {}),
  allowSubdomains: attr('boolean', {}),
  allowGlobDomains: attr('boolean', {}),
  allowAnyName: attr('boolean', {}),
  enforceHostnames: attr('boolean', {}),
  allowIpSans: attr('boolean', {
    defaultValue: true,
    label: 'Allow clients to request IP Subject Alternative Names (SANs)',
  }),
  allowedOtherSans: attr({
    editType: 'stringArray',
    label: 'Allowed Other SANs',
  }),
  serverFlag: attr('boolean', {
    defaultValue: true,
  }),
  clientFlag: attr('boolean', {
    defaultValue: true,
  }),
  codeSigningFlag: attr('boolean', {}),
  emailProtectionFlag: attr('boolean', {}),
  keyBits: attr('number', {
    defaultValue: 2048,
  }),
  keyUsage: attr('string', {
    defaultValue: 'DigitalSignature,KeyAgreement,KeyEncipherment',
    editType: 'stringArray',
  }),
  requireCn: attr('boolean', {
    label: 'Require common name',
    defaultValue: true,
  }),
  useCsrCommonName: attr('boolean', {
    label: 'Use CSR common name',
    defaultValue: true,
  }),
  useCsrSans: attr('boolean', {
    defaultValue: true,
    label: 'Use CSR subject alternative names (SANs)',
  }),
  ou: attr({
    label: 'OU (OrganizationalUnit)',
    editType: 'stringArray',
  }),
  organization: attr({
    editType: 'stringArray',
  }),
  country: attr({
    editType: 'stringArray',
  }),
  locality: attr({
    editType: 'stringArray',
    label: 'Locality/City',
  }),
  province: attr({
    editType: 'stringArray',
    label: 'Province/State',
  }),
  streetAddress: attr({
    editType: 'stringArray',
  }),
  postalCode: attr({
    editType: 'stringArray',
  }),
  generateLease: attr('boolean', {}),
  noStore: attr('boolean', {}),
  policyIdentifiers: attr({
    editType: 'stringArray',
  }),
  basicConstraintsValidForNonCA: attr('boolean', {
    label: 'Mark Basic Constraints valid when issuing non-CA certificates.',
  }),

  updatePath: queryRecord(
    'capabilities',
    context => {
      const { backend, id } = context.getProperties('backend', 'id');
      return {
        id: `${backend}/roles/${id}`,
      };
    },
    'id',
    'backend'
  ),
  canDelete: computed.alias('updatePath.canDelete'),
  canEdit: computed.alias('updatePath.canUpdate'),
  canRead: computed.alias('updatePath.canRead'),

  generatePath: queryRecord(
    'capabilities',
    context => {
      const { backend, id } = context.getProperties('backend', 'id');
      return {
        id: `${backend}/issue/${id}`,
      };
    },
    'id',
    'backend'
  ),
  canGenerate: computed.alias('generatePath.canUpdate'),

  signPath: queryRecord(
    'capabilities',
    context => {
      const { backend, id } = context.getProperties('backend', 'id');
      return {
        id: `${backend}/sign/${id}`,
      };
    },
    'id',
    'backend'
  ),
  canSign: computed.alias('signPath.canUpdate'),

  signVerbatimPath: queryRecord(
    'capabilities',
    context => {
      const { backend, id } = context.getProperties('backend', 'id');
      return {
        id: `${backend}/sign-verbatim/${id}`,
      };
    },
    'id',
    'backend'
  ),
  canSignVerbatim: computed.alias('signVerbatimPath.canUpdate'),

  /*
   * this hydrates the map in `fieldGroups` so that it contains
   * the actual field information, not just the name of the field
   */
  fieldsToAttrs(fieldGroups) {
    const attrMap = get(this.constructor, 'attributes');
    return fieldGroups.map(group => {
      const groupKey = Object.keys(group)[0];
      const groupMembers = group[groupKey];
      const fields = groupMembers.map(field => {
        var meta = attrMap.get(field);
        return {
          type: meta.type,
          name: meta.name,
          options: meta.options,
        };
      });
      return { [groupKey]: fields };
    });
  },

  /*
   * returns an array of objects that list attributes so that the form can be programmatically generated
   * the attributes are pulled from the model's attribute hash
   *
   * The keys will be used to label each section of the form.
   * the 'default' key contains fields that are outside of any grouping
   *
   * returns an array of objects:
   *
   * [
   *   {'default': [ { type: 'string', name: 'keyType', options: { label: 'Key Type'}}]},
   *   {'Options': [{ type: 'boolean', name: 'allowAnyName', options: {}}]}
   * ]
   *
   *
   */
  fieldGroups: computed(function() {
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
      { 'Extended Key Usage': ['serverFlag', 'clientFlag', 'codeSigningFlag', 'emailProtectionFlag'] },
      {
        Advanced: ['generateLease', 'noStore', 'basicConstraintsValidForNonCA', 'policyIdentifiers'],
      },
    ];

    return this.fieldsToAttrs(Ember.copy(groups, true));
  }),
});
