import { create, clickable, visitable, collection } from 'ember-cli-page-object';
import fields from '../../../components/form-field';
import flashMessage from '../../../components/flash-message';

export default create({
  ...fields,
  tabs: collection('[data-test-auth-section-tab]'),
  visit: visitable('/vault/settings/auth/configure/:path/:section'),
  flash: flashMessage,
  save: clickable('[data-test-save-config]'),
});
