import { Base } from '../create';
import { clickable, visitable, create, fillable } from 'ember-cli-page-object';

export default create({
  ...Base,
  path: fillable('[data-test-secret-path]'),
  secretKey: fillable('[data-test-secret-key]'),
  secretValue: fillable('[data-test-secret-value]'),
  save: clickable('[data-test-secret-save]'),
  deleteBtn: clickable('[data-test-secret-delete] button'),
  confirmBtn: clickable('[data-test-confirm-button]'),
  visitEdit: visitable('/vault/secrets/:backend/edit/:id'),
  visitEditRoot: visitable('/vault/secrets/:backend/edit'),
  deleteSecret() {
    return this.deleteBtn().confirmBtn();
  },

  createSecret(path, key, value) {
    return this.path(path).secretKey(key).secretValue(value).save();
  },
});
