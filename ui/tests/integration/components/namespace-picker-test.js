/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, fillIn, findAll, click, find } from '@ember/test-helpers';
import sinon from 'sinon';
import hbs from 'htmlbars-inline-precompile';
import Service from '@ember/service';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { setupMirage } from 'ember-cli-mirage/test-support';

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

const INITIALIZED_NAMESPACES = ['root', 'parent1', 'parent1/child1'];

module('Integration | Component | namespace-picker', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    const authService = this.owner.lookup('service:auth');
    this.authStub = sinon.stub(authService, 'authData');
    this.authStub.value({ userRootNamespace: '' });

    this.nsService = this.owner.lookup('service:namespace');
    // the path in the namespace service denotes the current namespace context a user is in
    this.nsService.path = 'parent1/child1';
    this.server.get('/sys/internal/ui/namespaces', () => {
      return {
        data: { keys: ['parent1/', 'parent1/child1/'] },
      };
    });

    this.owner.register('service:store', StoreService);
  });

  hooks.afterEach(function () {
    this.authStub.restore();
  });

  test('it focuses the search input field when the component is loaded', async function (assert) {
    await render(hbs`<NamespacePicker />`);
    await click(GENERAL.button('namespace-picker'));

    // Verify that the search input field is focused
    const searchInput = find(GENERAL.inputByAttr('Search namespaces'));
    assert.strictEqual(
      document.activeElement,
      searchInput,
      'The search input field is focused on component load'
    );
  });

  test('it filters namespace options based on search input', async function (assert) {
    await render(hbs`<NamespacePicker/>`);
    await click(GENERAL.button('namespace-picker'));

    // Verify all namespaces are displayed initially which are pre-populated in the NamespaceService
    for (const namespace of INITIALIZED_NAMESPACES) {
      assert.dom(GENERAL.button(namespace)).exists(`Namespace "${namespace}" is displayed initially`);
    }
    // Simulate typing into the search input
    await fillIn(GENERAL.inputByAttr('Search namespaces'), 'child1');

    // Verify that only namespaces matching the search input are displayed
    assert.strictEqual(
      findAll(GENERAL.inputByAttr('Search namespaces')).length,
      1,
      'Only matching namespaces are displayed after filtering'
    );

    // Clear the search input
    await fillIn(GENERAL.inputByAttr('Search namespaces'), '');

    // Verify all namespaces are displayed after clearing the search input
    assert.dom(GENERAL.button('root')).exists('Namespace "root" is displayed');
    assert.dom(GENERAL.button('parent1')).exists('Namespace "parent1" is displayed');
    assert.dom(GENERAL.button('parent1/child1')).exists('Namespace "parent1/child1" is displayed');
    assert.strictEqual(
      findAll(`ul ${GENERAL.button()}`).length,
      3,
      'Three namespaces are displayed after clearing the search input'
    );
  });

  test('it shows both "Manage" and "Refresh list" action buttons when canList is true', async function (assert) {
    const storeStub = this.owner.lookup('service:store');
    sinon.stub(storeStub, 'findRecord').callsFake((modelType, id) => {
      if (modelType === 'capabilities' && id === 'sys/namespaces/') {
        return Promise.resolve(getMockCapabilitiesModel(true));
      }
      return Promise.reject();
    });

    await render(hbs`<NamespacePicker />`);
    await click(GENERAL.button('namespace-picker'));

    // Verify that the "Refresh List" button is visible
    assert.dom(GENERAL.button('Refresh list')).exists('Refresh List button is visible');
    assert.dom(GENERAL.button('Manage')).exists('Manage button is visible');
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
    await click(GENERAL.button('namespace-picker'));

    // Verify that the buttons are hidden
    assert.dom(GENERAL.button('Refresh list')).doesNotExist('Refresh List button is hidden');
    assert.dom(GENERAL.button('Manage')).exists('Manage button is hidden');
  });

  test('it hides both action buttons when the capabilities store throws an error', async function (assert) {
    const storeStub = this.owner.lookup('service:store');
    sinon.stub(storeStub, 'findRecord').callsFake(() => {
      return Promise.reject();
    });

    await render(hbs`<NamespacePicker />`);
    await click(GENERAL.button('namespace-picker'));

    // Verify that the buttons are hidden
    assert.dom(GENERAL.button('Refresh list')).doesNotExist('Refresh List button is hidden');
    assert.dom(GENERAL.button('Manage')).doesNotExist('Manage button is hidden');
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
    await click(GENERAL.button('namespace-picker'));

    // Verify initial namespaces are displayed
    assert.dom(GENERAL.button('parent1')).exists('Namespace "parent1" is displayed');
    assert.dom(GENERAL.button('parent1/child1')).exists('Namespace "parent1/child1" is displayed');
    assert.dom(GENERAL.button('root')).exists('Namespace "root" is displayed');
    assert
      .dom(GENERAL.button('new-namespace'))
      .doesNotExist('Namespace "new-namespace" is not displayed initially');

    // Re-stub request with a new namespace
    this.server.get('/sys/internal/ui/namespaces', () => {
      return {
        data: { keys: ['parent1/', 'parent1/child1/', 'new-namespace/'] },
      };
    });
    // Click the "Refresh list" button
    await click(GENERAL.button('Refresh list'));

    // Verify the new namespace is displayed
    assert
      .dom(GENERAL.button('new-namespace'))
      .exists('Namespace "new-namespace" is displayed after refreshing');
  });

  test("it should display the user's root namespace if it is not true root (an empty string)", async function (assert) {
    this.authStub.value({ userRootNamespace: 'admin' }); // User's root namespace is "admin"
    this.nsService.path = 'admin'; // User is current in the "admin" namespace
    // The user also has access to a child namespace. This additional setup is important because as a fallback
    // the current namespace is displayed in the dropdown if nothing is returned from this endpoint.
    this.server.get('/sys/internal/ui/namespaces', () => {
      return { data: { keys: ['child1/'] } };
    });
    await render(hbs`<NamespacePicker />`);
    assert
      .dom(GENERAL.button('namespace-picker'))
      .hasText('admin', `shows the namespace 'admin' in the toggle component`);
    await click(GENERAL.button('namespace-picker'));
    assert.dom(`li ${GENERAL.button()}`).exists({ count: 2 }, 'namespace picker only contains 2 options');
    assert.dom(GENERAL.button('admin')).exists();
    assert.dom(GENERAL.button('admin/child1')).exists();
  });
});
