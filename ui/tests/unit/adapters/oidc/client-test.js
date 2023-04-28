/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import testHelper from './test-helper';

module('Unit | Adapter | oidc/client', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.modelName = 'oidc/client';
    this.data = {
      name: 'client-1',
      key: 'test-key',
      access_token_ttl: '30m',
      id_token_ttl: '1h',
    };
    this.path = '/identity/oidc/client/client-1';
  });

  testHelper(test);
});
