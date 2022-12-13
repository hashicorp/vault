import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default Model.extend({
  backend: attr('string'),
  der: attr(),
  pem: attr('string'),
  caChain: attr('string'),
  attrList(keys) {
    return expandAttributeMeta(this, keys);
  },

  //urls
  urlsAttrs: computed(function () {
    const keys = ['issuingCertificates', 'crlDistributionPoints', 'ocspServers'];
    return this.attrList(keys);
  }),
  issuingCertificates: attr({
    editType: 'stringArray',
  }),
  crlDistributionPoints: attr({
    label: 'CRL Distribution Points',
    editType: 'stringArray',
  }),
  ocspServers: attr({
    label: 'OCSP Servers',
    editType: 'stringArray',
  }),

  //tidy
  tidyAttrs: computed(function () {
    const keys = ['tidyCertStore', 'tidyRevocationList', 'safetyBuffer'];
    return this.attrList(keys);
  }),
  tidyCertStore: attr('boolean', {
    defaultValue: false,
    label: 'Tidy the Certificate Store',
  }),
  tidyRevocationList: attr('boolean', {
    defaultValue: false,
    label: 'Tidy the Revocation List (CRL)',
  }),
  safetyBuffer: attr({
    defaultValue: '72h',
    editType: 'ttl',
  }),

  crlAttrs: computed(function () {
    const keys = ['expiry', 'disable'];
    return this.attrList(keys);
  }),
  //crl
  expiry: attr('string', { defaultValue: '72h' }),
  disable: attr('boolean', { defaultValue: false }),
});
