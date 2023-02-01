import { attr } from '@ember-data/model';
import { copy } from 'ember-copy';
import { computed } from '@ember/object';
import Certificate from './pki-certificate-sign';

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
  permittedDnsDomains: attr('string', {
    label: 'Permitted DNS domains',
  }),
  ou: attr({
    label: 'OU (OrganizationalUnit)',
    editType: 'string',
  }),
  organization: attr({
    editType: 'string',
  }),
  country: attr({
    editType: 'string',
  }),
  locality: attr({
    editType: 'string',
    label: 'Locality/City',
  }),
  province: attr({
    editType: 'string',
    label: 'Province/State',
  }),
  streetAddress: attr({
    editType: 'string',
  }),
  postalCode: attr({
    editType: 'string',
  }),

  fieldGroups: computed('useCsrValues', function () {
    const options = [
      {
        Options: [
          'altNames',
          'ipSans',
          'ttl',
          'excludeCnFromSans',
          'maxPathLength',
          'permittedDnsDomains',
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
    if (this.useCsrValues === false) {
      groups = groups.concat(options);
    }

    return this.fieldsToAttrs(copy(groups, true));
  }),
});
