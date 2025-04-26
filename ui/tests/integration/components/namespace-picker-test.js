/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn, findAll, waitFor, click, find } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import Service from '@ember/service';
import { NAMESPACE_PICKER_SELECTORS } from 'vault/tests/helpers/namespace-picker';

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

  test('it updates the namespace list after clicking "Refresh list"', async function (assert) {
    // Mock `hasListPermissions`
    this.owner.lookup('service:namespace').set('hasListPermissions', true);

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
    await click(NAMESPACE_PICKER_SELECTORS.refreshList);

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

  test('it displays the "Manage" button when the user has permissions', async function (assert) {
    // Mock `hasListPermissions` to be true
    this.owner.lookup('service:namespace').set('hasListPermissions', true);

    await render(hbs`<NamespacePicker />`);

    await click(NAMESPACE_PICKER_SELECTORS.toggle);

    // Find the "Manage" button
    const manageButton = findAll('a').find((el) => {
      const spans = el.querySelectorAll('span');
      return spans[1]?.textContent.trim() === 'Manage';
    });

    // Verify the "Manage" button is rendered
    assert.ok(manageButton, 'The "Manage" button is displayed');
  });
});
