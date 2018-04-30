import { create, text, visitable, collection } from 'ember-cli-page-object';
import flashMessage from 'vault/tests/pages/components/flash-message';

export default create({
  visit: visitable('/vault/access/identity/:item_type'),
  flashMessage,
  items: collection('[data-test-identity-row]', {
    id: text('[data-test-identity-link]'),
  }),
});
