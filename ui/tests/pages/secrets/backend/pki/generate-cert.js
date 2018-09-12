import { Base } from '../credentials';
import { clickable, text, value, create, fillable, isPresent } from 'ember-cli-page-object';

export default create({
  ...Base,
  title: text('[data-test-title]'),
  commonName: fillable('[data-test-input="commonName"]'),
  commonNameValue: value('[data-test-input="commonName"]'),
  csr: fillable('[data-test-input="csr"]'),
  submit: clickable('[data-test-secret-generate]'),
  back: clickable('[data-test-secret-generate-back]'),
  certificate: text('[data-test-row-value="Certificate"]'),
  toggleOptions: clickable('[data-test-toggle-group]'),
  hasCert: isPresent('[data-test-row-value="Certificate"]'),
  fillInField: fillable('[data-test-field]'),
  issueCert: async function(commonName) {
    await this.commonName(commonName)
      .toggleOptions()
      .fillInField('unit', 'h')
      .submit();
  },

  sign: async function(commonName, csr) {
    return this.csr(csr)
      .commonName(commonName)
      .toggleOptions()
      .fillInField('unit', 'h')
      .submit();
  },
});
