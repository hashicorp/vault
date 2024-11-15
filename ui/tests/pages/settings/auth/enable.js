/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { create, visitable } from 'ember-cli-page-object';
import flashMessages from '../../components/flash-message';
import { mountBackend } from 'vault/tests/helpers/components/mount-backend-form-helpers';

export default create({
  visit: visitable('/vault/settings/auth/enable'),
  flash: flashMessages,
  enable: async function (type, path) {
    await this.visit();
    await mountBackend(type, path);
  },
});
