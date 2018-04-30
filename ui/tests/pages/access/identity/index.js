import { create, visitable } from 'ember-cli-page-object';
import flashMessage from 'vault/tests/pages/components/flash-message';

export default create({
  visit: visitable('/vault/access/identity/:item_type'),
  flashMessage,
});
