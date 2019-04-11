import { Base } from '../show';
import { create, clickable, collection, isPresent, text } from 'ember-cli-page-object';
import { code } from 'vault/tests/pages/helpers/codemirror';

export default create({
  ...Base,
  breadcrumbs: collection('[data-test-secret-breadcrumb]', {
    text: text(),
  }),
  deleteBtn: clickable('[data-test-secret-delete] button'),
  confirmBtn: clickable('[data-test-confirm-button]'),
  rows: collection('data-test-row-label'),
  toggleJSON: clickable('[data-test-secret-json-toggle]'),
  toggleIsPresent: isPresent('[data-test-secret-json-toggle]'),
  edit: clickable('[data-test-secret-edit]'),
  editIsPresent: isPresent('[data-test-secret-edit]'),
  editor: {
    content: code('[data-test-component="json-editor"]'),
  },
  deleteSecret() {
    return this.deleteBtn().confirmBtn();
  },
});
