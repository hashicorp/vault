import { Base } from '../show';
import { create, clickable, collection, isPresent } from 'ember-cli-page-object';
import { code } from 'vault/tests/pages/helpers/codemirror';

export default create({
  ...Base,
  rows: collection('data-test-row-label'),
  toggleJSON: clickable('[data-test-secret-json-toggle]'),
  toggleIsPresent: isPresent('[data-test-secret-json-toggle]'),
  edit: clickable('[data-test-secret-edit]'),
  editIsPresent: isPresent('[data-test-secret-edit]'),
  editor: {
    content: code('[data-test-component="json-editor"]'),
  },
});
