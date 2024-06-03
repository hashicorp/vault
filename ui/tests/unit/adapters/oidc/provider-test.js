/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import testHelper from './test-helper';

module('Unit | Adapter | oidc/provider', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.modelName = 'oidc/provider';
    this.data = {
      name: 'foo-provider',
      allowed_client_ids: ['*'],
      scopes_supported: [],
    };
    this.path = '/identity/oidc/provider/foo-provider';
  });

  testHelper(test);
});
