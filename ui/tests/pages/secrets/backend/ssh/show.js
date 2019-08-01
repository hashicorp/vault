import { Base } from '../show';
import { create, clickable, collection, isPresent } from 'ember-cli-page-object';

export default create({
  ...Base,
  rows: collection('data-test-row-label'),
  edit: clickable('[data-test-edit-link]'),
  editIsPresent: isPresent('[data-test-edit-link]'),
  generate: clickable('[data-test-backend-credentials]'),
  generateIsPresent: isPresent('[data-test-backend-credentials]'),
  deleteBtn: clickable('[data-test-confirm-action-trigger]'),
  confirmBtn: clickable('[data-test-confirm-button]'),
  deleteRole() {
    return this.deleteBtn().confirmBtn();
  },
});
