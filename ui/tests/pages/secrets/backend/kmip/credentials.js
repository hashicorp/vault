/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { create, clickable, visitable } from 'ember-cli-page-object';
import ListView from 'vault/tests/pages/components/list-view';

export default create({
  ...ListView,
  visit: visitable('/vault/secrets-engines/:backend/kmip/scopes/:scope/roles/:role/credentials'),
  visitDetail: visitable(
    '/vault/secrets-engines/:backend/kmip/scopes/:scope/roles/:role/credentials/:serial'
  ),
  create: clickable('[data-test-role-create]'),
  generateCredentialsLink: clickable('[data-test-generate-credentials]'),
  backToRoleLink: clickable('[data-test-back-to-role]'),
});
