import { create, clickable, visitable } from 'ember-cli-page-object';
import fields from '../../../components/form-field';
import flashMessage from '../../../components/flash-message';

export default create({
  ...fields,
  visit: visitable('/vault/settings/auth/configure/:path/:section'),
  flash: flashMessage,
  save: clickable('[data-test-save-config]'),
});
