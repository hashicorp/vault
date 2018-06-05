import { create, visitable, collection } from 'ember-cli-page-object';

import { getter } from 'ember-cli-page-object/macros';
import ConfigPKI from 'vault/tests/pages/components/config-pki';

export default create({
  visit: visitable('/vault/settings/secrets/configure/:backend/:section'),
  form: ConfigPKI,
  lastMessage: getter(function() {
    const count = this.flashMessages().count;
    return this.flashMessages(count - 1).text;
  }),
  flashMessages: collection({
    itemScope: '[data-test-flash-message-body]',
  }),
});
