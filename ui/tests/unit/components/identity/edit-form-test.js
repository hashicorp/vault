/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';

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

    test('alias create uses the alias name, not the parent entity/group type', function (assert) {
      const model = {
        identityType: 'entity',
        form: { identityFormType: 'alias', data: { name: 'my-alias' } },
      };
      const component = makeComponent(this.owner, 'create', model);
      assert.strictEqual(component.getMessage(model), 'Successfully saved Alias: my-alias.');
    });

    test('alias save uses the alias name', function (assert) {
      const model = {
        identityType: 'group',
        form: { identityFormType: 'alias', data: { name: 'my-alias' } },
      };
      const component = makeComponent(this.owner, 'edit', model);
      assert.strictEqual(component.getMessage(model), 'Successfully saved Alias: my-alias.');
    });

    test('alias delete uses the alias name', function (assert) {
      const model = {
        identityType: 'entity',
        form: { identityFormType: 'alias', data: { name: 'my-alias' } },
      };
      const component = makeComponent(this.owner, 'edit', model);
      assert.strictEqual(component.getMessage(model, true), 'Successfully deleted Alias: my-alias.');
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

  module('deleteNoun', function () {
    function makeComponent(owner, model) {
      const componentManager = owner.lookup('component-manager:glimmer');
      const componentClass = owner.factoryFor('component:identity/edit-form').class;
      return componentManager.createComponent(componentClass, { named: { model, mode: 'edit' } });
    }

    test('uses "alias" for an alias model, not the parent entity/group type', function (assert) {
      const model = { identityType: 'group', form: { identityFormType: 'alias' } };
      const component = makeComponent(this.owner, model);
      assert.strictEqual(component.deleteNoun, 'alias');
    });

    test('uses the identityType for a non-alias model', function (assert) {
      const model = { identityType: 'entity', form: { identityFormType: 'entity' } };
      const component = makeComponent(this.owner, model);
      assert.strictEqual(component.deleteNoun, 'entity');
    });
  });

  module('deleteItem', function (hooks) {
    hooks.beforeEach(function () {
      this.api = this.owner.lookup('service:api');
      this.entityDeleteAliasStub = sinon.stub(this.api.identity, 'entityDeleteAliasById').resolves();
      this.onSave = sinon.stub().resolves();
    });

    function makeComponent(owner, model, onSave) {
      const componentManager = owner.lookup('component-manager:glimmer');
      const componentClass = owner.factoryFor('component:identity/edit-form').class;
      return componentManager.createComponent(componentClass, { named: { model, mode: 'edit', onSave } });
    }

    test('deletes an alias using canonicalId when the model has no id', async function (assert) {
      const model = {
        identityType: 'entity',
        canonicalId: 'alias-id-from-route',
        form: { identityFormType: 'alias', data: { name: 'my-alias' } },
      };
      const component = makeComponent(this.owner, model, this.onSave);
      await component.deleteItem(model);

      assert.true(
        this.entityDeleteAliasStub.calledWith('alias-id-from-route'),
        'calls entityDeleteAliasById with the alias id resolved from canonicalId'
      );
      assert.true(this.onSave.calledOnce, 'onSave is called after a successful delete');
    });

    test('prefers model.id over canonicalId when both are present', async function (assert) {
      const model = {
        identityType: 'entity',
        id: 'preferred-id',
        canonicalId: 'fallback-id',
        form: { identityFormType: 'alias', data: { name: 'my-alias' } },
      };
      const component = makeComponent(this.owner, model, this.onSave);
      await component.deleteItem(model);

      assert.true(this.entityDeleteAliasStub.calledWith('preferred-id'), 'prefers model.id');
    });
  });
});
