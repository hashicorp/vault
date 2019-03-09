import { clickable, collection, fillable, text, selectable, isPresent } from 'ember-cli-page-object';
import { settled } from '@ember/test-helpers';
import fields from './form-field';

export default {
  ...fields,
  scope: '.config-pki-ca',
  text: text('[data-test-text]'),
  title: text('[data-test-title]'),

  hasTitle: isPresent('[data-test-title]'),
  hasError: isPresent('[data-test-error]'),
  hasSignIntermediateForm: isPresent('[data-test-sign-intermediate-form]'),

  replaceCA: clickable('[data-test-go-replace-ca]'),
  replaceCAText: text('[data-test-go-replace-ca]'),
  setSignedIntermediateBtn: clickable('[data-test-go-set-signed-intermediate]'),
  signIntermediateBtn: clickable('[data-test-go-sign-intermediate]'),
  caType: selectable('[data-test-input="caType"]'),
  submit: clickable('[data-test-submit]'),
  back: clickable('[data-test-back-button]'),

  signedIntermediate: fillable('[data-test-signed-intermediate]'),
  downloadLinks: collection('[data-test-ca-download-link]'),
  rows: collection('[data-test-table-row]'),
  rowValues: collection('[data-test-row-value]'),
  csr: text('[data-test-row-value="CSR"]', { normalize: false }),
  csrField: fillable('[data-test-input="csr"]'),
  certificate: text('[data-test-row-value="Certificate"]', { normalize: false }),
  certificateIsPresent: isPresent('[data-test-row-value="Certificate"]'),
  uploadCert: clickable('[data-test-input="uploadPemBundle"]'),
  enterCertAsText: clickable('[data-test-text-toggle]'),
  pemBundle: fillable('[data-test-text-file-textarea="true"]'),
  commonName: fillable('[data-test-input="commonName"]'),

  async generateCA(commonName = 'PKI CA', type = 'root') {
    if (type === 'intermediate') {
      await this.replaceCA()
        .commonName(commonName)
        .caType('intermediate')
        .submit();
      await settled();
      return;
    }
    await this.replaceCA()
      .commonName(commonName)
      .submit();
    await settled();
  },

  async uploadCA(pem) {
    await this.replaceCA()
      .uploadCert()
      .enterCertAsText()
      .pemBundle(pem)
      .submit();
    await settled();
  },

  async signIntermediate(commonName) {
    await this.signIntermediateBtn().commonName(commonName);
    await settled();
  },
};
