import { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import Certificate from './pki-certificate';

export default Certificate.extend({
  DISPLAY_FIELDS: computed(function () {
    return [
      'csr',
      'certificate',
      'commonName',
      'issueDate',
      'expiryDate',
      'issuingCa',
      'caChain',
      'privateKey',
      'privateKeyType',
      'serialNumber',
    ];
  }),
  addBasicConstraints: attr('boolean', {
    label: 'Add a Basic Constraints extension with CA: true',
    helpText:
      'Only needed as a workaround in some compatibility scenarios with Active Directory Certificate Services',
  }),
  backend: attr('string', {
    readOnly: true,
  }),
  canParse: attr('boolean'),
  caType: attr('string', {
    possibleValues: ['root', 'intermediate'],
    defaultValue: 'root',
    label: 'CA Type',
    readOnly: true,
  }),
  commonName: attr('string'),
  csr: attr('string', {
    editType: 'textarea',
    label: 'CSR',
    masked: true,
  }),
  expiryDate: attr('string', {
    label: 'Expiration date',
  }),
  issueDate: attr('string'),
  keyBits: attr('number', {
    defaultValue: 2048,
  }),
  keyType: attr('string', {
    possibleValues: ['rsa', 'ec', 'ed25519'],
    defaultValue: 'rsa',
  }),
  maxPathLength: attr('number', {
    defaultValue: -1,
  }),
  organization: attr({
    editType: 'stringArray',
  }),
  ou: attr({
    label: 'OU (OrganizationalUnit)',
    editType: 'stringArray',
  }),
  pemBundle: attr('string', {
    label: 'PEM bundle',
    editType: 'file',
  }),
  permittedDnsNames: attr('string', {
    label: 'Permitted DNS domains',
  }),
  privateKeyFormat: attr('string', {
    possibleValues: ['', 'der', 'pem', 'pkcs8'],
    defaultValue: '',
  }),
  type: attr('string', {
    possibleValues: ['internal', 'exported'],
    defaultValue: 'internal',
  }),
  uploadPemBundle: attr('boolean', {
    label: 'Upload PEM bundle',
    readOnly: true,
  }),

  // address attrs
  country: attr({
    editType: 'stringArray',
  }),
  locality: attr({
    editType: 'stringArray',
    label: 'Locality/City',
  }),
  streetAddress: attr({
    editType: 'stringArray',
  }),
  postalCode: attr({
    editType: 'stringArray',
  }),
  province: attr({
    editType: 'stringArray',
    label: 'Province/State',
  }),

  fieldDefinition: computed('caType', 'uploadPemBundle', function () {
    const type = this.caType;
    const isUpload = this.uploadPemBundle;
    let groups = [{ default: ['caType', 'uploadPemBundle'] }];
    if (isUpload) {
      groups[0].default.push('pemBundle');
    } else {
      groups[0].default.push('type', 'commonName');
      if (type === 'root') {
        groups.push({
          Options: [
            'altNames',
            'ipSans',
            'ttl',
            'format',
            'privateKeyFormat',
            'keyType',
            'keyBits',
            'maxPathLength',
            'permittedDnsNames',
            'excludeCnFromSans',
            'ou',
            'organization',
            'otherSans',
          ],
        });
      }
      if (type === 'intermediate') {
        groups.push({
          Options: [
            'altNames',
            'ipSans',
            'format',
            'privateKeyFormat',
            'keyType',
            'keyBits',
            'excludeCnFromSans',
            'addBasicConstraints',
            'ou',
            'organization',
            'otherSans',
          ],
        });
      }
    }
    groups.push({
      'Address Options': ['country', 'locality', 'province', 'streetAddress', 'postalCode'],
    });

    return groups;
  }),
});
