import { Base } from '../credentials';
import { clickable, value, create, fillable, isPresent } from 'ember-cli-page-object';

export default create({
  ...Base,
  userIsPresent: isPresent('[data-test-input="username"]'),
  ipIsPresent: isPresent('[data-test-input="ip"]'),
  user: fillable('[data-test-input="username"]'),
  ip: fillable('[data-test-input="ip"]'),
  warningIsPresent: isPresent('[data-test-warning]'),
  commonNameValue: value('[data-test-input="commonName"]'),
  submit: clickable('[data-test-secret-generate]'),
  back: clickable('[data-test-secret-generate-back]'),
  generateOTP: async function() {
    await this.user('admin')
      .ip('192.168.1.1')
      .submit();
  },
});
