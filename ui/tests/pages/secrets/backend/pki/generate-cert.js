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
  enableTtl: clickable('[data-test-toggle-input]'),
  hasCert: isPresent('[data-test-row-value="Certificate"]'),
  fillInTime: fillable('[data-test-ttl-value]'),
  fillInField: fillable('[data-test-select="ttl-unit"]'),
  issueCert: async function (commonName) {
    await this.commonName(commonName).toggleOptions().enableTtl().fillInField('h').fillInTime('30').submit();
  },

  sign: async function (commonName, csr) {
    return this.csr(csr)
      .commonName(commonName)
      .toggleOptions()
      .enableTtl()
      .fillInField('h')
      .fillInTime('30')
      .submit();
  },
});
