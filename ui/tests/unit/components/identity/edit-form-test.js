/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import EmberObject from '@ember/object';
import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';
import { ROUTES } from 'vault/utils/routes';

module('Unit | Component | identity/edit-form', function (hooks) {
  setupTest(hooks);

  const testCases = [
    {
      identityType: 'entity',
      mode: 'create',
      expected: ROUTES.VAULT_CLUSTER_ACCESS_IDENTITY,
    },
    {
      identityType: 'entity',
      mode: 'edit',
      expected: ROUTES.VAULT_CLUSTER_ACCESS_IDENTITY_SHOW,
    },
    {
      identityType: 'entity-merge',
      mode: 'merge',
      expected: ROUTES.VAULT_CLUSTER_ACCESS_IDENTITY,
    },
    {
      identityType: 'entity-alias',
      mode: 'create',
      expected: ROUTES.VAULT_CLUSTER_ACCESS_IDENTITY_ALIASES_SHOW,
    },
    {
      identityType: 'entity-alias',
      mode: 'edit',
      expected: ROUTES.VAULT_CLUSTER_ACCESS_IDENTITY_ALIASES_SHOW,
    },
    {
      identityType: 'group',
      mode: 'create',
      expected: ROUTES.VAULT_CLUSTER_ACCESS_IDENTITY,
    },
    {
      identityType: 'group',
      mode: 'edit',
      expected: ROUTES.VAULT_CLUSTER_ACCESS_IDENTITY_SHOW,
    },
    {
      identityType: 'group-alias',
      mode: 'create',
      expected: ROUTES.VAULT_CLUSTER_ACCESS_IDENTITY_ALIASES_SHOW,
    },
    {
      identityType: 'group-alias',
      mode: 'edit',
      expected: ROUTES.VAULT_CLUSTER_ACCESS_IDENTITY_ALIASES_SHOW,
    },
  ];
  testCases.forEach(function (testCase) {
    const model = EmberObject.create({
      identityType: testCase.identityType,
      rollbackAttributes: sinon.spy(),
    });
    test(`it computes cancelLink properly: ${testCase.identityType} ${testCase.mode}`, function (assert) {
      const component = this.owner.lookup('component:identity/edit-form');

      component.set('mode', testCase.mode);
      component.set('model', model);
      assert.strictEqual(component.get('cancelLink'), testCase.expected, 'cancel link is correct');
    });
  });
});
