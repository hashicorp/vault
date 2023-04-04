/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { create, visitable, collection } from 'ember-cli-page-object';

import { getter } from 'ember-cli-page-object/macros';
import ConfigPKI from 'vault/tests/pages/components/pki/config-pki';

export default create({
  visit: visitable('/vault/settings/secrets/configure/:backend/:section'),
  form: ConfigPKI,
  lastMessage: getter(function () {
    const count = this.flashMessages.length;
    return this.flashMessages.objectAt(count - 1).text;
  }),
  flashMessages: collection('[data-test-flash-message-body]'),
});
