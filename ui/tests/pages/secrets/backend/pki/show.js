import { Base } from '../show';
import { create, clickable, collection, text, isPresent } from 'ember-cli-page-object';

export default create({
  ...Base,
  rows: collection('data-test-row-label'),
  certificate: text('[data-test-row-value="Certificate"]'),
  hasCert: isPresent('[data-test-row-value="Certificate"]'),
  edit: clickable('[data-test-edit-link]'),
  generateCert: clickable('[data-test-credentials-link]'),
  deleteBtn: clickable('[data-test-role-delete] button'),
  confirmBtn: clickable('[data-test-confirm-button]'),
  deleteRole() {
    return this.deleteBtn().confirmBtn();
  },
});
