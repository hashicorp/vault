import { computed } from '@ember/object';
import DS from 'ember-data';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

const { attr } = DS;

export default DS.Model.extend({
  backend: attr('string'),
  der: attr(),
  pem: attr('string'),
  caChain: attr('string'),
  attrList(keys) {
    return expandAttributeMeta(this, keys);
  },

  //urls
  urlsAttrs: computed(function() {
    let keys = ['issuingCertificates', 'crlDistributionPoints', 'ocspServers'];
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
  tidyAttrs: computed(function() {
    let keys = ['tidyCertStore', 'tidyRevocationList', 'safetyBuffer'];
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

  crlAttrs: computed(function() {
    let keys = ['expiry'];
    return this.attrList(keys);
  }),
  //crl
  expiry: attr({
    defaultValue: '72h',
    editType: 'ttl',
  }),
});
