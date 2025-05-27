/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn, findAll, waitFor, click, find } from '@ember/test-helpers';
import sinon from 'sinon';
import hbs from 'htmlbars-inline-precompile';
import Service from '@ember/service';
import { NAMESPACE_PICKER_SELECTORS } from 'vault/tests/helpers/namespace-picker';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

class AuthService extends Service {
  authData = { userRootNamespace: '' };
}

class NamespaceService extends Service {
  accessibleNamespaces = ['parent1', 'parent1/child1'];
  path = 'parent1/child1';

  findNamespacesForUser = {
    perform: () => Promise.resolve(),
  };
}

class StoreService extends Service {
  findRecord(modelType, id) {
    return new Promise((resolve, reject) => {
      if (modelType === 'capabilities' && id === 'sys/namespaces/') {
        resolve(); // Simulate a successful response
      } else {
        reject({ httpStatus: 404, message: 'not found' }); // Simulate an error response
      }
    });
  }
}

function getMockCapabilitiesModel(canList) {
  // Mock for the Capabilities model
  return {
    path: 'sys/namespaces/',
    capabilities: canList ? ['list'] : [],
    get(property) {
      if (property === 'canList') {
        return this.capabilities.includes('list');
      }
      return undefined;
    },
  };
}

module('Integration | Component | namespace-picker', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.owner.register('service:auth', AuthService);
    this.owner.register('service:namespace', NamespaceService);
    this.owner.register('service:store', StoreService);
  });

  test('it focuses the search input field when the component is loaded', async function (assert) {
    await render(hbs`<NamespacePicker />`);
    await click(NAMESPACE_PICKER_SELECTORS.toggle);

    // Verify that the search input field is focused
    const searchInput = find(NAMESPACE_PICKER_SELECTORS.searchInput);
    assert.strictEqual(
      document.activeElement,
      searchInput,
      'The search input field is focused on component load'
    );
  });

  test('it filters namespace options based on search input', async function (assert) {
    await render(hbs`<NamespacePicker/>`);
    await click(NAMESPACE_PICKER_SELECTORS.toggle);

    // Verify all namespaces are displayed initially
    await waitFor(NAMESPACE_PICKER_SELECTORS.link());
    assert.strictEqual(
      findAll(NAMESPACE_PICKER_SELECTORS.link()).length,
      3,
      'All namespaces are displayed initially'
    );

    // Simulate typing into the search input
    await waitFor(NAMESPACE_PICKER_SELECTORS.searchInput);
    await fillIn(NAMESPACE_PICKER_SELECTORS.searchInput, 'child1');

    // Verify that only namespaces matching the search input are displayed
    assert.strictEqual(
      findAll(NAMESPACE_PICKER_SELECTORS.link()).length,
      1,
      'Only matching namespaces are displayed after filtering'
    );

    // Clear the search input
    await fillIn(NAMESPACE_PICKER_SELECTORS.searchInput, '');

    // Verify all namespaces are displayed after clearing the search input
    assert.strictEqual(
      findAll(NAMESPACE_PICKER_SELECTORS.link()).length,
      3,
      'All namespaces are displayed after clearing the search input'
    );
  });

  test('it shows both action buttons when canList is true', async function (assert) {
    const storeStub = this.owner.lookup('service:store');
    sinon.stub(storeStub, 'findRecord').callsFake((modelType, id) => {
      if (modelType === 'capabilities' && id === 'sys/namespaces/') {
        return Promise.resolve(getMockCapabilitiesModel(true));
      }
      return Promise.reject();
    });

    await render(hbs`<NamespacePicker />`);
    await click(NAMESPACE_PICKER_SELECTORS.toggle);

    // Verify that the "Refresh List" button is visible
    assert.dom(GENERAL.button('refresh-namespaces')).exists('Refresh List button is visible');
    assert.dom(GENERAL.button('manage-namespaces')).exists('Manage button is visible');
  });

  test('it hides the refresh button when canList is false', async function (assert) {
    const storeStub = this.owner.lookup('service:store');
    sinon.stub(storeStub, 'findRecord').callsFake((modelType, id) => {
      if (modelType === 'capabilities' && id === 'sys/namespaces/') {
        return Promise.resolve(getMockCapabilitiesModel(false));
      }
      return Promise.reject();
    });

    await render(hbs`<NamespacePicker />`);
    await click(NAMESPACE_PICKER_SELECTORS.toggle);

    // Verify that the buttons are hidden
    assert.dom(GENERAL.button('refresh-namespaces')).doesNotExist('Refresh List button is hidden');
    assert.dom(GENERAL.button('manage-namespaces')).exists('Manage button is hidden');
  });

  test('it hides both action buttons when the capabilities store throws an error', async function (assert) {
    const storeStub = this.owner.lookup('service:store');
    sinon.stub(storeStub, 'findRecord').callsFake(() => {
      return Promise.reject();
    });

    await render(hbs`<NamespacePicker />`);
    await click(NAMESPACE_PICKER_SELECTORS.toggle);

    // Verify that the buttons are hidden
    assert.dom(GENERAL.button('refresh-namespaces')).doesNotExist('Refresh List button is hidden');
    assert.dom(GENERAL.button('manage-namespaces')).doesNotExist('Manage button is hidden');
  });

  test('it updates the namespace list after clicking "Refresh list"', async function (assert) {
    this.owner.lookup('service:namespace').set('hasListPermissions', true);

    const storeStub = this.owner.lookup('service:store');
    sinon.stub(storeStub, 'findRecord').callsFake((modelType, id) => {
      if (modelType === 'capabilities' && id === 'sys/namespaces/') {
        return Promise.resolve(getMockCapabilitiesModel(true)); // Return the mock model
      }
      return Promise.reject();
    });

    await render(hbs`<NamespacePicker />`);
    await click(NAMESPACE_PICKER_SELECTORS.toggle);

    // Dynamically modify the `findNamespacesForUser.perform` method for this test
    const namespaceService = this.owner.lookup('service:namespace');
    namespaceService.set('findNamespacesForUser', {
      perform: () => {
        namespaceService.set('accessibleNamespaces', [
          'parent1',
          'parent1/child1',
          'new-namespace', // Add a new namespace
        ]);
        return Promise.resolve();
      },
    });

    // Verify initial namespaces are displayed
    assert.strictEqual(
      findAll(NAMESPACE_PICKER_SELECTORS.link()).length,
      3,
      'Initially, three namespaces are displayed'
    );

    // Click the "Refresh list" button
    await click(GENERAL.button('refresh-namespaces'));

    // Verify the new namespace is displayed
    assert.strictEqual(
      findAll(NAMESPACE_PICKER_SELECTORS.link()).length,
      4,
      'After refreshing, four namespaces are displayed'
    );

    // Verify the new namespace is specifically shown
    assert
      .dom(NAMESPACE_PICKER_SELECTORS.link('new-namespace'))
      .exists('The new namespace "new-namespace" is displayed after refreshing');
  });
});
