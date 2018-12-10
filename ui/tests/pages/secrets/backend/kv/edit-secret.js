import { Base } from '../create';
import { isPresent, clickable, visitable, create, fillable } from 'ember-cli-page-object';
import { codeFillable } from 'vault/tests/pages/helpers/codemirror';
export default create({
  ...Base,
  path: fillable('[data-test-secret-path]'),
  secretKey: fillable('[data-test-secret-key]'),
  secretValue: fillable('[data-test-secret-value] textarea'),
  save: clickable('[data-test-secret-save]'),
  deleteBtn: clickable('[data-test-secret-delete] button'),
  confirmBtn: clickable('[data-test-confirm-button]'),
  visitEdit: visitable('/vault/secrets/:backend/edit/:id'),
  visitEditRoot: visitable('/vault/secrets/:backend/edit'),
  toggleJSON: clickable('[data-test-secret-json-toggle]'),
  hasMetadataFields: isPresent('[data-test-metadata-fields]'),
  editor: {
    fillIn: codeFillable('[data-test-component="json-editor"]'),
  },
  deleteSecret() {
    return this.deleteBtn().confirmBtn();
  },
  createSecret: async function(path, key, value) {
    return this.path(path)
      .secretKey(key)
      .secretValue(value)
      .save();
  },
  editSecret: async function(key, value) {
    return this.secretKey(key)
      .secretValue(value)
      .save();
  },
});
