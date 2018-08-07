import Ember from 'ember';
import DS from 'ember-data';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import fieldToAttrs from 'vault/utils/field-to-attrs';

const { attr } = DS;
const { computed } = Ember;

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
  extKeyUsageOids: attr({
    label: 'Custom extended key usage OIDs',
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

  updatePath: lazyCapabilities(apiPath`${'backend'}/roles/${'id'}`, 'backend', 'id'),
  canDelete: computed.alias('updatePath.canDelete'),
  canEdit: computed.alias('updatePath.canUpdate'),
  canRead: computed.alias('updatePath.canRead'),

  generatePath: lazyCapabilities(apiPath`${'backend'}/issue/${'id'}`, 'backend', 'id'),
  canGenerate: computed.alias('generatePath.canUpdate'),

  signPath: lazyCapabilities(apiPath`${'backend'}/sign/${'id'}`, 'backend', 'id'),
  canSign: computed.alias('signPath.canUpdate'),

  signVerbatimPath: lazyCapabilities(apiPath`${'backend'}/sign-verbatim/${'id'}`, 'backend', 'id'),
  canSignVerbatim: computed.alias('signVerbatimPath.canUpdate'),

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
        Advanced: ['generateLease', 'noStore', 'basicConstraintsValidForNonCA', 'policyIdentifiers'],
      },
    ];

    return fieldToAttrs(this, groups);
  }),
});
