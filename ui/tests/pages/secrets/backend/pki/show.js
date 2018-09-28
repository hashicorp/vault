import { Base } from '../show';
import { create, clickable, collection, text, isPresent } from 'ember-cli-page-object';

export default create({
  ...Base,
  rows: collection('data-test-row-label'),
  certificate: text('[data-test-row-value="Certificate"]'),
  hasCert: isPresent('[data-test-row-value="Certificate"]'),
  edit: clickable('[data-test-edit-link]'),
  editIsPresent: isPresent('[data-test-edit-link]'),
  generateCert: clickable('[data-test-credentials-link]'),
  generateCertIsPresent: isPresent('[data-test-credentials-link]'),
  signCert: clickable('[data-test-sign-link]'),
  signCertIsPresent: isPresent('[data-test-sign-link]'),
});
