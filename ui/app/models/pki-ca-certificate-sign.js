import { copy } from 'ember-copy';
import { computed } from '@ember/object';
import DS from 'ember-data';
import Certificate from './pki-certificate-sign';

const { attr } = DS;

export default Certificate.extend({
  backend: attr('string', {
    readOnly: true,
  }),
  useCsrValues: attr('boolean', {
    defaultValue: false,
    label: 'Use CSR values',
  }),
  maxPathLength: attr('number', {
    defaultValue: -1,
  }),
  permittedDnsNames: attr('string', {
    label: 'Permitted DNS domains',
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

  fieldGroups: computed('useCsrValues', function() {
    const options = [
      {
        Options: [
          'altNames',
          'ipSans',
          'ttl',
          'excludeCnFromSans',
          'maxPathLength',
          'permittedDnsNames',
          'ou',
          'organization',
          'otherSans',
        ],
      },
      {
        'Address Options': ['country', 'locality', 'province', 'streetAddress', 'postalCode'],
      },
    ];
    let groups = [
      {
        default: ['csr', 'commonName', 'format', 'useCsrValues'],
      },
    ];
    if (this.get('useCsrValues') === false) {
      groups = groups.concat(options);
    }

    return this.fieldsToAttrs(copy(groups, true));
  }),
});
