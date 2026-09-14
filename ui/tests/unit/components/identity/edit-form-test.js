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
      label: 'entity',
      model: { identityType: 'entity' },
      mode: 'create',
      expected: 'vault.cluster.access.identity.entities.index',
    },
    {
      label: 'entity',
      model: { identityType: 'entity' },
      mode: 'edit',
      expected: 'vault.cluster.access.identity.entities.show',
    },
    {
      label: 'merge',
      model: { form: { identityFormType: 'group' } },
      mode: 'merge',
      expected: 'vault.cluster.access.identity.entities.index',
    },
    {
      label: 'entity-alias',
      model: { identityType: 'entity', form: { identityFormType: 'alias' } },
      mode: 'create',
      expected: 'vault.cluster.access.identity.entities.aliases.index',
    },
    {
      label: 'entity-alias',
      model: { identityType: 'entity', form: { identityFormType: 'alias' } },
      mode: 'edit',
      expected: 'vault.cluster.access.identity.entities.aliases.show',
    },
    {
      label: 'group',
      model: { identityType: 'group' },
      mode: 'create',
      expected: 'vault.cluster.access.identity.groups.index',
    },
    {
      label: 'group',
      model: { identityType: 'group' },
      mode: 'edit',
      expected: 'vault.cluster.access.identity.groups.show',
    },
    {
      label: 'group-alias',
      model: { identityType: 'group', form: { identityFormType: 'alias' } },
      mode: 'create',
      expected: 'vault.cluster.access.identity.groups.aliases.index',
    },
    {
      label: 'group-alias',
      model: { identityType: 'group', form: { identityFormType: 'alias' } },
      mode: 'edit',
      expected: 'vault.cluster.access.identity.groups.aliases.show',
    },
  ];
  testCases.forEach(function (testCase) {
    test(`it computes cancelLink properly: ${testCase.label} ${testCase.mode}`, function (assert) {
      const { mode, model } = testCase;
      const named = {
        model,
        mode,
      };
      const componentManager = this.owner.lookup('component-manager:glimmer');
      const componentClass = this.owner.factoryFor('component:identity/edit-form').class;
      const component = componentManager.createComponent(componentClass, { named });
      assert.strictEqual(component.cancelLink, testCase.expected, 'cancel link is correct');
    });
  });
});
