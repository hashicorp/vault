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
      expected: 'vault.cluster.access.identity.entities.index',
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
      expected: 'vault.cluster.access.identity.groups.index',
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

  // verify success toast message is correct
  module('getMessage', function () {
    function makeComponent(owner, mode, model) {
      const componentManager = owner.lookup('component-manager:glimmer');
      const componentClass = owner.factoryFor('component:identity/edit-form').class;
      return componentManager.createComponent(componentClass, { named: { model, mode } });
    }

    test('save uses group name instead of id in success message', function (assert) {
      const model = {
        identityType: 'group',
        itemId: 'some-uuid',
        form: { identityFormType: 'group', data: { name: 'my-group' } },
      };
      const component = makeComponent(this.owner, 'edit', model);
      assert.strictEqual(component.getMessage(model), 'Successfully saved Group: my-group.');
    });

    test('delete uses group name in success message', function (assert) {
      const model = {
        identityType: 'group',
        form: { identityFormType: 'group', data: { name: 'to-delete' } },
      };
      const component = makeComponent(this.owner, 'edit', model);
      assert.strictEqual(component.getMessage(model, true), 'Successfully deleted Group: to-delete.');
    });
  });
});
