/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import testHelper from './test-helper';

module('Unit | Adapter | oidc/assignment', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.modelName = 'oidc/assignment';
    this.data = {
      name: 'foo-assignment',
      entity_ids: ['my-entity'],
      group_ids: ['my-group'],
    };
    this.path = '/identity/oidc/assignment/foo-assignment';
  });

  testHelper(test);
});
