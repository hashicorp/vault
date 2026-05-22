/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Component | identity/edit-form', function (hooks) {
  setupTest(hooks);

  const testCases = [
    {
      identityType: 'entity',
      mode: 'create',
      expected: 'vault.cluster.access.identity',
    },
    {
      identityType: 'entity',
      mode: 'edit',
      expected: 'vault.cluster.access.identity.show',
    },
    {
      identityType: 'entity-merge',
      mode: 'merge',
      expected: 'vault.cluster.access.identity',
    },
    {
      identityType: 'entity-alias',
      mode: 'create',
      expected: 'vault.cluster.access.identity.aliases',
    },
    {
      identityType: 'entity-alias',
      mode: 'edit',
      expected: 'vault.cluster.access.identity.aliases.show',
    },
    {
      identityType: 'group',
      mode: 'create',
      expected: 'vault.cluster.access.identity',
    },
    {
      identityType: 'group',
      mode: 'edit',
      expected: 'vault.cluster.access.identity.show',
    },
    {
      identityType: 'group-alias',
      mode: 'create',
      expected: 'vault.cluster.access.identity.aliases',
    },
    {
      identityType: 'group-alias',
      mode: 'edit',
      expected: 'vault.cluster.access.identity.aliases.show',
    },
  ];
  testCases.forEach(function (testCase) {
    test(`it computes cancelLink properly: ${testCase.identityType} ${testCase.mode}`, function (assert) {
      const { identityType, mode } = testCase;
      const named = {
        model: { identityType },
        mode,
      };
      const componentManager = this.owner.lookup('component-manager:glimmer');
      const componentClass = this.owner.factoryFor('component:identity/edit-form').class;
      const component = componentManager.createComponent(componentClass, { named });
      assert.strictEqual(component.cancelLink, testCase.expected, 'cancel link is correct');
    });
  });
});
